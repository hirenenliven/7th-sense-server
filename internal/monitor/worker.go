package monitor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
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
	rateInterval   = 10 * time.Second // how often the reading rate is logged
)

// stopReason says why a connection ended.
type stopReason int

const (
	stopRequested       stopReason = iota // the worker was stopped (collector stop, machine changed)
	stopConnectionStale                   // no usable reading for CNC_STALE_SEC
)

// boundVariable is a configured variable and the snapshot field it is read from.
type boundVariable struct {
	config model.VariableConfig
	field  int
}

// worker collects one machine: it connects, stores every reading in machine_telementry right
// away (no batching), writes mode periods and events, and reconnects with backoff when the connection is lost.
type worker struct {
	collector               *Collector
	machine                 model.MonitoredMachine // owned by the worker goroutine
	appliedMachine          model.MonitoredMachine // owned by the collector's sync loop goroutine
	variableReloadRequested chan struct{}          // holds one request; more requests collapse into it
	cancel                  context.CancelFunc
	done                    chan struct{}

	pendingMu      sync.Mutex
	pendingMachine *model.MonitoredMachine // a new controller version / name, taken on reload

	variables []boundVariable

	mode                   int
	program                string
	alarm                  float64
	emergency              float64
	alarmsKnown            bool
	connectedSince         *time.Time
	lastSeenAt             *time.Time
	lastConnectedReadingAt time.Time // the stale check measures CNC_STALE_SEC from here
	lastChangeAt           *time.Time
	sampleAt               *time.Time
	missingNodes           string
	storedCount            int64
	droppedCount           int64 // telemetry rows a failed insert lost
	processedCount         int64 // readings processReading handled
	reportedError          string
	statusDirty            bool
	statusSavedAt          time.Time // when machine_live_status was last written
	connection             *cnc.Connection

	liveValues map[int]model.LiveValue // cache of current live values by variable_id

	// The counters at the last reading rate log; logReadingRate reports the change since then.
	rateLoggedAt  time.Time
	rateCallbacks int64
	rateProcessed int64
	rateStored    int64
}

// newMachineWorker creates the worker of machine; cancel stops it.
func newMachineWorker(collector *Collector, machine model.MonitoredMachine, cancel context.CancelFunc) *worker {
	return &worker{
		collector:               collector,
		machine:                 machine,
		appliedMachine:          machine,
		variableReloadRequested: make(chan struct{}, 1),
		cancel:                  cancel,
		done:                    make(chan struct{}),
		mode:                    model.ModeOffline, // the collector starts every machine offline
	}
}

// queueMachineUpdate hands the worker machine details that need no reconnect (controller version,
// name) and asks it to reload its variables. Called from the collector's sync loop goroutine;
// the worker takes the details on its own goroutine (applyPendingMachineUpdate).
func (machineWorker *worker) queueMachineUpdate(machine model.MonitoredMachine) {
	machineWorker.appliedMachine = machine
	machineWorker.pendingMu.Lock()
	machineWorker.pendingMachine = &machine
	machineWorker.pendingMu.Unlock()
	machineWorker.requestVariableReload()
}

// requestVariableReload asks the worker to read its variable configuration again (no reconnect).
func (machineWorker *worker) requestVariableReload() {
	select {
	case machineWorker.variableReloadRequested <- struct{}{}:
	default:
	}
}

// runConnectionLoop is the worker goroutine: it connects the machine, collects its readings and,
// when the connection ends, reconnects after a backoff (1 s doubling to 5 min, jittered; reset
// once a connection delivered data) until ctx is cancelled. done is closed when it returns.
func (machineWorker *worker) runConnectionLoop(ctx context.Context) {
	defer close(machineWorker.done)

	// DATA_COLLECTION_STEP_10: Check the machine can be connected and load its variable
	// configuration before every connect attempt.
	if machineWorker.machine.CNCEngine == nil {
		machineWorker.recordEvent(model.EventConfigError, nil,
			"the machine's controller has no cnc_engine; set controllers.cnc_engine and save the machine")
		<-ctx.Done()
		return
	}

	retry := backoff{max: 5 * time.Minute}
	for attempt := 0; ctx.Err() == nil; attempt++ {
		if attempt > 0 {
			sleepUnlessCancelled(ctx, retry.nextWait())
			if ctx.Err() != nil {
				return
			}
			machineWorker.recordEvent(model.EventReconnect, nil, fmt.Sprintf("attempt %d", attempt+1))
		}
		if err := machineWorker.loadVariableConfig(ctx); err != nil {
			slog.Error("cnc: load variable config", "machineId", machineWorker.machine.ID, "error", err)
			continue
		}

		// DATA_COLLECTION_STEP_11: Connect the machine, waiting for a free connect slot.
		connection, id, err := machineWorker.connectWithinLimit(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			// Every callback slot taken is a configuration problem (CNC_SLOTS), not a machine fault.
			eventType := model.EventConnectFailed
			if errors.Is(err, cnc.ErrNoFreeSlot) {
				eventType = model.EventConfigError
			}
			machineWorker.recordEvent(eventType, &id, err.Error())
			continue
		}
		machineWorker.recordEvent(model.EventConnectOK, &id, fmt.Sprintf("connection %d (%s:%d, engine %d)",
			id, machineWorker.machine.LiveIP, machineWorker.machine.Port, *machineWorker.machine.CNCEngine))
		slog.Info("cnc: machine subscribed", "machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name,
			"ip", machineWorker.machine.LiveIP, "port", machineWorker.machine.Port, "engine", *machineWorker.machine.CNCEngine,
			"connectionId", id, "variables", len(machineWorker.variables))

		reason, gotData := machineWorker.collectReadings(ctx, connection)
		// DATA_COLLECTION_STEP_17: Disconnect, record the machine offline, then retry.
		connection.Disconnect()
		machineWorker.recordDisconnection(reason)
		if gotData {
			retry.reset()
		}
	}
}

// connectWithinLimit calls the driver, at most CNC_MAX_CONNECTING machines at a time per
// collector: CNCConnect may block while the DLL connects, so a start with many machines does not
// open every connection at once.
func (machineWorker *worker) connectWithinLimit(ctx context.Context) (*cnc.Connection, int, error) {
	select {
	case machineWorker.collector.connectLimiter <- struct{}{}:
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	}
	defer func() { <-machineWorker.collector.connectLimiter }()
	return machineWorker.collector.driver.Connect(machineWorker.machine)
}

// loadVariableConfig reads the machine's variable configuration and binds each variable to its
// snapshot field. It first takes machine details queued by queueMachineUpdate, because the
// variables depend on the controller version.
func (machineWorker *worker) loadVariableConfig(ctx context.Context) error {
	machineWorker.applyPendingMachineUpdate()

	configs, err := machineWorker.collector.service.GetVariableConfig(ctx, machineWorker.machine, machineWorker.collector.driver.Sources())
	if err != nil {
		return err
	}
	variables := make([]boundVariable, 0, len(configs))
	variableMap := make([]string, 0, len(configs)) // "14:opMode<-MachineMode", for the log
	var unreadable []string                        // active variables without a readable address
	for _, variable := range configs {
		if variable.Address == "" {
			unreadable = append(unreadable, fmt.Sprintf("%d:%s", variable.ID, variable.Code))
			continue
		}
		field, ok := cnc.FieldIndex(variable.Address)
		if !ok {
			machineWorker.recordEvent(model.EventConfigError, nil,
				fmt.Sprintf("variable %d (%s): the DLL snapshot has no field %q", variable.ID, variable.Code, variable.Address))
			continue
		}
		variables = append(variables, boundVariable{config: variable, field: field})
		variableMap = append(variableMap, fmt.Sprintf("%d:%s<-%s", variable.ID, variable.Code, variable.Address))
	}
	machineWorker.variables = variables
	// Clear live values cache on reload - will be reloaded on next reading
	machineWorker.liveValues = nil
	slog.Info("cnc: variables loaded", "machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name,
		"variables", len(variables), "map", variableMap)
	if len(unreadable) > 0 {
		// Logged, not stored as events: it would repeat on every reconnect.
		slog.Warn("cnc: active variables not readable", "machineId", machineWorker.machine.ID,
			"machine", machineWorker.machine.Name, "engine", *machineWorker.machine.CNCEngine,
			"reason", "no address in a source the server can read (only DLL snapshot fields today)",
			"count", len(unreadable), "variables", unreadable)
	}
	return nil
}

// applyPendingMachineUpdate takes the machine details queued by queueMachineUpdate, if any.
func (machineWorker *worker) applyPendingMachineUpdate() {
	machineWorker.pendingMu.Lock()
	defer machineWorker.pendingMu.Unlock()
	if machineWorker.pendingMachine != nil {
		machineWorker.machine = *machineWorker.pendingMachine
		machineWorker.pendingMachine = nil
	}
}

// collectReadings handles readings until the worker stops or the connection goes stale. gotData
// reports whether the controller sent a usable reading.
func (machineWorker *worker) collectReadings(ctx context.Context, connection *cnc.Connection) (reason stopReason, gotData bool) {
	machineWorker.connection = connection
	machineWorker.alarmsKnown = false
	machineWorker.lastConnectedReadingAt = time.Now()
	machineWorker.reportedError = ""
	machineWorker.startReadingRateWindow(time.Now(), 0) // a new connection's callbacks count from 0

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	var lastReading time.Time

	// DATA_COLLECTION_STEP_14: Wait for DLL readings, variable reloads and the 1 s tick on the
	// worker goroutine; only this goroutine touches the worker's state.
	for {
		select {
		case <-ctx.Done():
			return stopRequested, gotData
		case <-machineWorker.variableReloadRequested:
			if err := machineWorker.loadVariableConfig(ctx); err != nil {
				slog.Error("cnc: reload variable config", "machineId", machineWorker.machine.ID, "error", err)
			}
		case <-connection.Signal():
			machineWorker.recordDLLError(connection.Stats())
			// Signals collapse, so take the newest reading and skip one already processed.
			reading, ok := connection.Latest()
			if !ok || !reading.ReceivedAt.After(lastReading) {
				continue
			}
			lastReading = reading.ReceivedAt
			machineWorker.processedCount++
			if machineWorker.processReading(reading) {
				gotData = true
			}
		case now := <-ticker.C:
			// DATA_COLLECTION_STEP_16: Drop a stale connection (no connected reading for
			// CNC_STALE_SEC) and refresh machine_live_status when it changed or every 10 s.
			if now.Sub(machineWorker.lastConnectedReadingAt) > machineWorker.collector.config.Stale {
				return stopConnectionStale, gotData
			}
			if machineWorker.statusDirty || now.Sub(machineWorker.statusSavedAt) >= statusInterval {
				machineWorker.saveLiveStatus(now)
			}
			if now.Sub(machineWorker.rateLoggedAt) >= rateInterval {
				machineWorker.logReadingRate(now, connection.Stats())
			}
		}
	}
}

// processReading stores one reading: every numeric value is upserted into
// machine_variable_values, only the values that changed are inserted into machine_telementry,
// and when something changed the alarm events and the mode are worked out. Text variables are
// not stored. It returns whether the controller was connected.
func (machineWorker *worker) processReading(reading cnc.Reading) bool {
	// DATA_COLLECTION_STEP_15: Store the reading (change detection against the cached live values).
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
			slog.Warn("cnc: controller data stopped (the DLL reports it lost the controller)",
				"machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name)
		}
		machineWorker.recordModeChange(now, model.ModeOffline)
		return false
	}
	machineWorker.lastConnectedReadingAt = now
	if machineWorker.connectedSince == nil {
		machineWorker.connectedSince = &now
		machineWorker.statusDirty = true
		slog.Info("cnc: controller data started", "machineId", machineWorker.machine.ID,
			"machine", machineWorker.machine.Name, "missingNodes", reading.MissingNodes)
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
			if machineWorker.collector.config.LogValues {
				machineWorker.logValueChange(variable.config, currentLive.Value, value.Number, now)
			}
		}

		numbers[variable.config.ID] = value.Number
	}

	// Alarms and mode are derived from the stored values, so they are rechecked only when a
	// value changed.
	if len(liveValueRows) > 0 {
		machineWorker.upsertLiveValues(liveValueRows, now)
	}
	if len(telemetryRows) > 0 {
		machineWorker.insertTelemetry(telemetryRows, now)
		machineWorker.recordAlarmChanges(now, numbers, alarmMessage)
		machineWorker.recordModeChange(now, modeFor(true, numbers))
	}
	return true
}

// logValueChange logs one changed variable value (CNC_LOG_VALUES=true). previous is nil for
// the first value the machine sends for the variable.
func (machineWorker *worker) logValueChange(variable model.VariableConfig, previous *float64, value float64, readAt time.Time) {
	slog.Info("cnc: variable value changed", "machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name,
		"variableId", variable.ID, "variable", variable.Code, "address", variable.Address,
		"previous", previous, "value", value, "readAt", readAt.Format(time.RFC3339Nano))
}

// startReadingRateWindow starts counting the reading rate from now; callbacks is the
// connection's callback count at that moment.
func (machineWorker *worker) startReadingRateWindow(now time.Time, callbacks int64) {
	machineWorker.rateLoggedAt = now
	machineWorker.rateCallbacks = callbacks
	machineWorker.rateProcessed = machineWorker.processedCount
	machineWorker.rateStored = machineWorker.storedCount
}

// logReadingRate logs how many readings the DLL sent since the last log, how many the worker
// processed and how many values it stored. skipped readings arrived while the worker was still
// busy with an earlier one: only the newest reading is processed (Connection.Signal collapses),
// so a growing skipped count means the worker (mostly its database writes) is slower than the DLL.
func (machineWorker *worker) logReadingRate(now time.Time, stats cnc.Stats) {
	seconds := now.Sub(machineWorker.rateLoggedAt).Seconds()
	callbacks := stats.Callbacks - machineWorker.rateCallbacks
	processed := machineWorker.processedCount - machineWorker.rateProcessed
	stored := machineWorker.storedCount - machineWorker.rateStored
	slog.Info("cnc: reading rate", "machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name,
		"seconds", roundToTenth(seconds), "callbacks", callbacks, "callbacksPerSecond", roundToTenth(float64(callbacks)/seconds),
		"processed", processed, "skipped", callbacks-processed,
		"storedValues", stored, "storedValuesPerSecond", roundToTenth(float64(stored)/seconds))
	machineWorker.startReadingRateWindow(now, stats.Callbacks)
}

// roundToTenth rounds value to one decimal for logs.
func roundToTenth(value float64) float64 {
	return math.Round(value*10) / 10
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

// recordAlarmChanges writes alarm_on / alarm_off and emergency_on / emergency_off when they change.
func (machineWorker *worker) recordAlarmChanges(now time.Time, numbers map[int]float64, alarmMessage string) {
	alarm, emergency := numbers[variableAlarm], numbers[variableEmergency]
	previousAlarm, previousEmergency := machineWorker.alarm, machineWorker.emergency
	if !machineWorker.alarmsKnown {
		previousAlarm, previousEmergency = 0, 0
		machineWorker.alarmsKnown = true
	}
	if alarm != previousAlarm {
		if previousAlarm != 0 {
			code := int(previousAlarm)
			machineWorker.recordEventAt(now, model.EventAlarmOff, &code, "")
		}
		if alarm != 0 {
			code := int(alarm)
			machineWorker.recordEventAt(now, model.EventAlarmOn, &code, alarmMessage)
		}
	}
	if (emergency != 0) != (previousEmergency != 0) {
		eventType := model.EventEmergencyOff
		if emergency != 0 {
			eventType = model.EventEmergencyOn
		}
		machineWorker.recordEventAt(now, eventType, nil, "")
	}
	machineWorker.alarm, machineWorker.emergency = alarm, emergency
}

// recordModeChange opens a new mode period when the mode changed.
func (machineWorker *worker) recordModeChange(at time.Time, mode int) {
	if mode == machineWorker.mode {
		return
	}
	slog.Info("cnc: mode changed", "machineId", machineWorker.machine.ID, "machine", machineWorker.machine.Name,
		"from", modeName(machineWorker.mode), "to", modeName(mode), "program", machineWorker.program)
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

// recordDLLError writes a dll_error event when the DLL reports a new error message.
func (machineWorker *worker) recordDLLError(stats cnc.Stats) {
	if stats.LastError == "" || stats.LastError == machineWorker.reportedError {
		return
	}
	machineWorker.reportedError = stats.LastError
	machineWorker.statusDirty = true
	machineWorker.recordEvent(model.EventDLLError, nil, stats.LastError)
}

// saveLiveStatus writes the machine's machine_live_status row.
func (machineWorker *worker) saveLiveStatus(now time.Time) {
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
	machineWorker.statusSavedAt = now
}

// recordDisconnection records the end of a connection: a stale event when it went stale, a
// disconnected event, the Offline mode and the live status.
func (machineWorker *worker) recordDisconnection(reason stopReason) {
	now := time.Now()
	if reason == stopConnectionStale {
		machineWorker.recordEvent(model.EventStale, nil, fmt.Sprintf("no usable reading for %s", machineWorker.collector.config.Stale))
	}
	machineWorker.recordEvent(model.EventDisconnected, nil, "")
	machineWorker.recordModeChange(now, model.ModeOffline)
	machineWorker.connectedSince = nil
	machineWorker.saveLiveStatus(now)
	machineWorker.connection = nil
}

// recordEvent stores and logs a machine event at the current time.
func (machineWorker *worker) recordEvent(eventType string, code *int, message string) {
	machineWorker.recordEventAt(time.Now(), eventType, code, message)
}

// recordEventAt stores a machine event and logs it.
func (machineWorker *worker) recordEventAt(at time.Time, eventType string, code *int, message string) {
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

// nextWait returns the next wait: 1 s first, then double the previous one up to max, ±20 %.
func (retry *backoff) nextWait() time.Duration {
	if retry.current == 0 {
		retry.current = time.Second
	} else {
		retry.current = min(retry.current*2, retry.max)
	}
	jitter := time.Duration((rand.Float64()*0.4 - 0.2) * float64(retry.current))
	return retry.current + jitter
}

// reset makes the next wait 1 s again.
func (retry *backoff) reset() {
	retry.current = 0
}

// sleepUnlessCancelled waits for duration or until ctx ends.
func sleepUnlessCancelled(ctx context.Context, duration time.Duration) {
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-ctx.Done():
	case <-timer.C:
	}
}
