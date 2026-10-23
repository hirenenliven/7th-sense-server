package monitor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sync"
	"time"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/model"
)

const (
	tickInterval   = time.Second      // the stale check and the status refresh
	statusInterval = 10 * time.Second // machine_live_status refresh while nothing else changes
	writeTimeout   = 5 * time.Second  // one database write
)

// stopReason says why a connection ended.
type stopReason int

const (
	stopStopped stopReason = iota // the worker was stopped (collector stop, machine changed)
	stopStale                     // no usable reading for CNC_STALE_SEC
)

// boundVariable is a configured variable and the snapshot field it is read from.
type boundVariable struct {
	config model.VariableConfig
	field  int
}

// worker collects one machine: it connects, stores every reading in machine_telementry right
// away (no batching), writes mode periods and events, and reconnects with backoff when the connection is lost.
type worker struct {
	collector     *Collector
	machine       model.MonitoredMachine // owned by the worker goroutine
	applied       model.MonitoredMachine // owned by the collector's listener goroutine
	configChanged chan struct{}
	cancel        context.CancelFunc
	done          chan struct{}

	pendingMu      sync.Mutex
	pendingMachine *model.MonitoredMachine // a new controller version / name, taken on reload

	variables []boundVariable

	mode           int
	program        string
	alarm          float64
	emergency      float64
	alarmsKnown    bool
	connectedSince *time.Time
	lastSeenAt     *time.Time
	lastGoodAt     time.Time
	lastChangeAt   *time.Time
	sampleAt       *time.Time
	missingNodes   string
	storedCount    int64
	droppedCount   int64 // telemetry rows a failed insert lost
	reportedError  string
	statusDirty    bool
	statusAt       time.Time
	connection     *cnc.Connection

	liveValues map[int]model.LiveValue // cache of current live values by variable_id
}

func newWorker(collector *Collector, machine model.MonitoredMachine, cancel context.CancelFunc) *worker {
	return &worker{
		collector:     collector,
		machine:       machine,
		applied:       machine,
		configChanged: make(chan struct{}, 1),
		cancel:        cancel,
		done:          make(chan struct{}),
		mode:          model.ModeOffline, // the collector starts every machine offline
	}
}

// updateMachine hands the worker machine details that need no reconnect (controller version,
// name) and reloads its variables. Called from the collector's listener goroutine.
func (machineWorker *worker) updateMachine(machine model.MonitoredMachine) {
	machineWorker.applied = machine
	machineWorker.pendingMu.Lock()
	machineWorker.pendingMachine = &machine
	machineWorker.pendingMu.Unlock()
	machineWorker.reloadConfig()
}

// reloadConfig asks the worker to read its variable configuration again (no reconnect).
func (machineWorker *worker) reloadConfig() {
	select {
	case machineWorker.configChanged <- struct{}{}:
	default:
	}
}

func (machineWorker *worker) run(ctx context.Context) {
	defer close(machineWorker.done)

	if machineWorker.machine.CNCEngine == nil {
		machineWorker.event(model.EventConfigError, nil,
			"the machine's controller has no cnc_engine; set controllers.cnc_engine and save the machine")
		<-ctx.Done()
		return
	}

	retry := backoff{max: 5 * time.Minute}
	for attempt := 0; ctx.Err() == nil; attempt++ {
		if attempt > 0 {
			sleep(ctx, retry.next())
			if ctx.Err() != nil {
				return
			}
			machineWorker.event(model.EventReconnect, nil, fmt.Sprintf("attempt %d", attempt+1))
		}
		if err := machineWorker.loadVariables(ctx); err != nil {
			slog.Error("cnc: load variable config", "machineId", machineWorker.machine.ID, "error", err)
			continue
		}

		connection, id, err := machineWorker.connect(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			eventType := model.EventConnectFailed
			if errors.Is(err, cnc.ErrNoFreeSlot) {
				eventType = model.EventConfigError
			}
			machineWorker.event(eventType, &id, err.Error())
			continue
		}
		machineWorker.event(model.EventConnectOK, &id, fmt.Sprintf("connection %d (%s:%d, engine %d)",
			id, machineWorker.machine.LiveIP, machineWorker.machine.Port, *machineWorker.machine.CNCEngine))

		reason, gotData := machineWorker.collect(ctx, connection)
		connection.Disconnect()
		machineWorker.disconnected(reason)
		if gotData {
			retry.reset()
		}
	}
}

// connect calls the driver, at most CNC_MAX_CONNECTING machines at a time per collector.
func (machineWorker *worker) connect(ctx context.Context) (*cnc.Connection, int, error) {
	select {
	case machineWorker.collector.connecting <- struct{}{}:
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
	defer func() { <-machineWorker.collector.connecting }()
	return machineWorker.collector.driver.Connect(machineWorker.machine)
}

// loadVariables reads the machine's variable configuration and binds each variable to its
// snapshot field.
func (machineWorker *worker) loadVariables(ctx context.Context) error {
	machineWorker.pendingMu.Lock()
	if machineWorker.pendingMachine != nil {
		machineWorker.machine = *machineWorker.pendingMachine
		machineWorker.pendingMachine = nil
	}
	machineWorker.pendingMu.Unlock()

	configs, err := machineWorker.collector.service.GetVariableConfig(ctx, machineWorker.machine, machineWorker.collector.driver.Sources())
	if err != nil {
		return err
	}
	variables := make([]boundVariable, 0, len(configs))
	for _, variable := range configs {
		field, ok := cnc.FieldIndex(variable.Address)
		if !ok {
			machineWorker.event(model.EventConfigError, nil,
				fmt.Sprintf("variable %d (%s): the DLL snapshot has no field %q", variable.ID, variable.Code, variable.Address))
			continue
		}
		variables = append(variables, boundVariable{config: variable, field: field})
	}
	machineWorker.variables = variables
	// Clear live values cache on reload - will be reloaded on next reading
	machineWorker.liveValues = nil
	slog.Info("cnc: variables loaded", "machineId", machineWorker.machine.ID, "variables", len(variables))
	return nil
}

// collect handles readings until the worker stops or the connection goes stale. gotData
// reports whether the controller sent a usable reading.
func (machineWorker *worker) collect(ctx context.Context, connection *cnc.Connection) (reason stopReason, gotData bool) {
	machineWorker.connection = connection
	machineWorker.alarmsKnown = false
	machineWorker.lastGoodAt = time.Now()
	machineWorker.reportedError = ""

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	var lastReading time.Time

	for {
		select {
		case <-ctx.Done():
			return stopStopped, gotData
		case <-machineWorker.configChanged:
			if err := machineWorker.loadVariables(ctx); err != nil {
				slog.Error("cnc: reload variable config", "machineId", machineWorker.machine.ID, "error", err)
			}
		case <-connection.Signal():
			machineWorker.reportDLLError(connection.Stats())
			reading, ok := connection.Latest()
			if !ok || !reading.ReceivedAt.After(lastReading) {
				continue
			}
			lastReading = reading.ReceivedAt
			if machineWorker.handleReading(reading) {
				gotData = true
			}
		case now := <-ticker.C:
			if now.Sub(machineWorker.lastGoodAt) > machineWorker.collector.config.Stale {
				return stopStale, gotData
			}
			if machineWorker.statusDirty || now.Sub(machineWorker.statusAt) >= statusInterval {
				machineWorker.saveStatus(now)
			}
		}
	}
}

// handleReading stores one reading in machine_telementry (one row per variable) and works out
// the mode and events. It returns whether the controller was connected.
func (machineWorker *worker) handleReading(reading cnc.Reading) bool {
	now := reading.ReceivedAt
	machineWorker.lastSeenAt = &now
	if !reading.SampleAt.IsZero() {
		sampleAt := reading.SampleAt
		machineWorker.sampleAt = &sampleAt
	}
	if reading.MissingNodes != machineWorker.missingNodes {
		machineWorker.missingNodes = reading.MissingNodes
		machineWorker.statusDirty = true
	}
	if !reading.Connected {
		if machineWorker.connectedSince != nil {
			machineWorker.connectedSince = nil
			machineWorker.statusDirty = true
		}
		machineWorker.setMode(now, model.ModeOffline)
		return false
	}
	machineWorker.lastGoodAt = now
	if machineWorker.connectedSince == nil {
		machineWorker.connectedSince = &now
		machineWorker.statusDirty = true
	}

	// Load live values if not already cached
	if machineWorker.liveValues == nil {
		variableIDs := make([]int, len(machineWorker.variables))
		for i, variable := range machineWorker.variables {
			variableIDs[i] = variable.config.ID
		}
		ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
		defer cancel()
		liveValues, err := machineWorker.collector.service.GetLiveValues(ctx, machineWorker.machine.ID, variableIDs)
		if err != nil {
			slog.Error("cnc: load live values", "machineId", machineWorker.machine.ID, "error", err)
			// Continue without live values - we'll still upsert them
		} else {
			machineWorker.liveValues = liveValues
		}
	}

	numbers := make(map[int]float64, len(machineWorker.variables))
	alarmMessage := ""
	telemetryRows := make([]model.TelemetryValue, 0, len(machineWorker.variables))
	liveValueRows := make([]model.LiveValue, 0, len(machineWorker.variables))
	for _, variable := range machineWorker.variables {
		value := reading.Values[variable.field]
		currentLive := machineWorker.liveValues[variable.config.ID]

		// Only process numeric variables
		if variable.config.IsText() {
			// Text variables are not stored - skip
			switch {
			case variable.config.ID == variableProgName:
				machineWorker.program = value.Text
			case variable.config.ID == variableAlarmMsg:
				alarmMessage = value.Text
			}
			continue
		}

		// Check if value changed
		valueChanged := false
		if currentLive.Value == nil || *currentLive.Value != value.Number {
			valueChanged = true
		}

		// Update live values
		liveRow := machineWorker.liveValueRow(variable.config, value, now)
		liveValueRows = append(liveValueRows, liveRow)

		// Insert into time-series if value changed
		if valueChanged {
			telemetryRows = append(telemetryRows, machineWorker.telemetryRow(variable.config, value, now))
		}

		numbers[variable.config.ID] = value.Number
	}

	// Upsert live values and insert changed telemetry
	if len(liveValueRows) > 0 {
		machineWorker.upsertLiveValues(liveValueRows, now)
	}
	if len(telemetryRows) > 0 {
		machineWorker.insertTelemetry(telemetryRows, now)
		machineWorker.checkAlarms(now, numbers, alarmMessage)
		machineWorker.setMode(now, modeFor(true, numbers))
	}
	return true
}

// liveValueRow is the machine_variable_values row of value, read at readAt.
// Only numeric values are cached here (text variables are not stored in this table).
func (machineWorker *worker) liveValueRow(variable model.VariableConfig, value cnc.Value, readAt time.Time) model.LiveValue {
	row := model.LiveValue{MachineID: machineWorker.machine.ID, VariableID: variable.ID, UpdatedAt: readAt}
	row.Value = &value.Number
	return row
}

// telemetryRow is the machine_telementry row of value, read at readAt.
func (machineWorker *worker) telemetryRow(variable model.VariableConfig, value cnc.Value, readAt time.Time) model.TelemetryValue {
	row := model.TelemetryValue{MachineID: machineWorker.machine.ID, VariableID: variable.ID, RecordedAt: readAt}
	row.Value = &value.Number
	return row
}

// upsertLiveValues upserts the live values for all variables.
func (machineWorker *worker) upsertLiveValues(rows []model.LiveValue, readAt time.Time) {
	if len(rows) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	if err := machineWorker.collector.service.UpsertLiveValues(ctx, rows); err != nil {
		machineWorker.reportedError = err.Error()
		machineWorker.statusDirty = true
		slog.Error("cnc: upsert live values", "machineId", machineWorker.machine.ID, "rows", len(rows), "error", err)
	} else {
		// Update cache
		for _, row := range rows {
			machineWorker.liveValues[row.VariableID] = row
		}
		machineWorker.lastChangeAt = &readAt
		machineWorker.statusDirty = true
	}
}

// insertTelemetry stores the rows of one reading. When the insert fails the rows are lost
// (counted in machine_live_status.dropped_count); the next reading is stored again.
func (machineWorker *worker) insertTelemetry(rows []model.TelemetryValue, readAt time.Time) {
	if len(rows) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	if err := machineWorker.collector.service.InsertTelemetry(ctx, rows); err != nil {
		machineWorker.droppedCount += int64(len(rows))
		machineWorker.reportedError = err.Error()
		machineWorker.statusDirty = true
		slog.Error("cnc: insert telemetry", "machineId", machineWorker.machine.ID, "rows", len(rows), "error", err)
		return
	}
	machineWorker.storedCount += int64(len(rows))
	machineWorker.lastChangeAt = &readAt
}

// checkAlarms writes alarm_on / alarm_off and emergency_on / emergency_off when they change.
func (machineWorker *worker) checkAlarms(now time.Time, numbers map[int]float64, alarmMessage string) {
	alarm, emergency := numbers[variableAlarm], numbers[variableEmergency]
	previousAlarm, previousEmergency := machineWorker.alarm, machineWorker.emergency
	if !machineWorker.alarmsKnown {
		previousAlarm, previousEmergency = 0, 0
		machineWorker.alarmsKnown = true
	}
	if alarm != previousAlarm {
		if previousAlarm != 0 {
			code := int(previousAlarm)
			machineWorker.eventAt(now, model.EventAlarmOff, &code, "")
		}
		if alarm != 0 {
			code := int(alarm)
			machineWorker.eventAt(now, model.EventAlarmOn, &code, alarmMessage)
		}
	}
	if (emergency != 0) != (previousEmergency != 0) {
		eventType := model.EventEmergencyOff
		if emergency != 0 {
			eventType = model.EventEmergencyOn
		}
		machineWorker.eventAt(now, eventType, nil, "")
	}
	machineWorker.alarm, machineWorker.emergency = alarm, emergency
}

// setMode opens a new mode period when the mode changed.
func (machineWorker *worker) setMode(at time.Time, mode int) {
	if mode == machineWorker.mode {
		return
	}
	machineWorker.mode = mode
	machineWorker.statusDirty = true
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	err := machineWorker.collector.service.ChangeMode(ctx, model.ModeChange{
		MachineID: machineWorker.machine.ID, ModeID: mode, ProgramName: machineWorker.program, At: at,
	})
	if err != nil {
		slog.Error("cnc: change mode", "machineId", machineWorker.machine.ID, "mode", mode, "error", err)
	}
}

// reportDLLError writes a dll_error event when the DLL reports a new error message.
func (machineWorker *worker) reportDLLError(stats cnc.Stats) {
	if stats.LastError == "" || stats.LastError == machineWorker.reportedError {
		return
	}
	machineWorker.reportedError = stats.LastError
	machineWorker.statusDirty = true
	machineWorker.event(model.EventDLLError, nil, stats.LastError)
}

// saveStatus writes the machine's machine_live_status row.
func (machineWorker *worker) saveStatus(now time.Time) {
	status := model.MachineLiveStatus{
		MachineID:      machineWorker.machine.ID,
		Shard:          machineWorker.collector.shard,
		IsConnected:    machineWorker.connectedSince != nil,
		ModeID:         machineWorker.mode,
		ConnectedSince: machineWorker.connectedSince,
		LastSeenAt:     machineWorker.lastSeenAt,
		LastChangeAt:   machineWorker.lastChangeAt,
		SampleAt:       machineWorker.sampleAt,
		MissingNodes:   machineWorker.missingNodes,
		StoredCount:    machineWorker.storedCount,
		DroppedCount:   machineWorker.droppedCount,
	}
	if machineWorker.connection != nil {
		stats := machineWorker.connection.Stats()
		status.CallbackCount, status.ErrorCount, status.LastError = stats.Callbacks, stats.Errors, stats.LastError
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	if err := machineWorker.collector.service.SaveStatus(ctx, status); err != nil {
		slog.Error("cnc: save live status", "machineId", machineWorker.machine.ID, "error", err)
	}
	machineWorker.statusDirty = false
	machineWorker.statusAt = now
}

// disconnected records the end of a connection: an event, the Offline mode and the status.
func (machineWorker *worker) disconnected(reason stopReason) {
	now := time.Now()
	if reason == stopStale {
		machineWorker.event(model.EventStale, nil, fmt.Sprintf("no usable reading for %s", machineWorker.collector.config.Stale))
	}
	machineWorker.event(model.EventDisconnected, nil, "")
	machineWorker.setMode(now, model.ModeOffline)
	machineWorker.connectedSince = nil
	machineWorker.saveStatus(now)
	machineWorker.connection = nil
}

func (machineWorker *worker) event(eventType string, code *int, message string) {
	machineWorker.eventAt(time.Now(), eventType, code, message)
}

// eventAt stores a machine event and logs it.
func (machineWorker *worker) eventAt(at time.Time, eventType string, code *int, message string) {
	machineID := machineWorker.machine.ID
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout)
	defer cancel()
	err := machineWorker.collector.service.AddEvent(ctx, model.MachineEvent{
		MachineID: &machineID, Shard: machineWorker.collector.shard, Type: eventType, Code: code, Message: message, At: at,
	})
	if err != nil {
		slog.Error("cnc: store event", "machineId", machineID, "event", eventType, "error", err)
	}

	attributes := []any{"machineId", machineID, "machine", machineWorker.machine.Name, "event", eventType}
	if code != nil {
		attributes = append(attributes, "code", *code)
	}
	if message != "" {
		attributes = append(attributes, "message", message)
	}
	switch eventType {
	case model.EventConnectFailed, model.EventConfigError, model.EventDLLError, model.EventStale:
		slog.Warn("cnc machine event", attributes...)
	default:
		slog.Info("cnc machine event", attributes...)
	}
}

// backoff gives waits that double from 1 s up to max, with ±20 % jitter so machines that
// failed together do not retry together.
type backoff struct {
	current time.Duration
	max     time.Duration
}

func (retry *backoff) next() time.Duration {
	if retry.current == 0 {
		retry.current = time.Second
	} else {
		retry.current = min(retry.current*2, retry.max)
	}
	jitter := time.Duration((rand.Float64()*0.4 - 0.2) * float64(retry.current))
	return retry.current + jitter
}

func (retry *backoff) reset() {
	retry.current = 0
}

// sleep waits for duration or until ctx ends.
func sleep(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
