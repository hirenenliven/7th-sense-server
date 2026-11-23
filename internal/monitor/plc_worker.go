package monitor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"7th-sense-server/internal/live"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/opcua"
	"7th-sense-server/internal/service"
)

const (
	// plcRetryMax is the longest wait between two connects of a machine that does not answer.
	plcRetryMax = 5 * time.Minute
	// plcMaxPublishErrors: this many publish errors in a row (no change in between) mean the
	// subscription is dead even though the channel looks open, so the worker reconnects.
	plcMaxPublishErrors = 3
	// plcCloseTimeout bounds closing a session (the controller may be gone).
	plcCloseTimeout = 3 * time.Second
	// plcDefaultPollInterval is used when PLC_POLL_INTERVAL_MS is not positive.
	plcDefaultPollInterval = time.Second
)

// errPlcReload ends a collection to subscribe again with new variables or machine details.
var errPlcReload = errors.New("plc variables reloaded")

// plcWorker reads one machine's PLC variables over OPC UA: connect, subscribe to every PLC node
// id, then for each published change hand it to the live hub first and store it like the CNC
// worker stores snapshot values (machine_variable_values every value, machine_telementry only a
// changed number; text is published, not stored), through the telemetry queue when it is on. Nodes the controller has no room to monitor
// (BadTooManyMonitoredItems) are read every PLC_POLL_INTERVAL_MS instead and handled like
// changes. It reconnects with backoff when the controller cannot be reached, the connection is
// lost, or the health read fails.
type plcWorker struct {
	collector *PlcCollector
	cancel    context.CancelFunc
	reload    chan struct{} // holds one reload request; more collapse into it

	// appliedMachine is what the collector last handed over; only the sync loop goroutine uses it.
	appliedMachine model.MonitoredMachine

	pendingMu      sync.Mutex
	pendingMachine *model.MonitoredMachine

	// The rest belongs to the worker goroutine.
	machine   model.MonitoredMachine
	variables map[uint32]model.VariableConfig // by subscription handle (= variable id)
	stored    map[int]float64                 // last number in machine_variable_values, by variable id
	polled    map[uint32]opcua.Change         // last value read of each node read instead of monitored
}

func newPlcWorker(collector *PlcCollector, machine model.MonitoredMachine, cancel context.CancelFunc) *plcWorker {
	return &plcWorker{collector: collector, cancel: cancel, reload: make(chan struct{}, 1),
		appliedMachine: machine, machine: machine}
}

// requestReload asks the worker to load its variables and subscribe again.
func (plcWorker *plcWorker) requestReload() {
	select {
	case plcWorker.reload <- struct{}{}:
	default:
	}
}

// queueMachineUpdate hands new machine details (name, controller version) to the worker and
// reloads its variables, which depend on the controller version.
func (plcWorker *plcWorker) queueMachineUpdate(machine model.MonitoredMachine) {
	plcWorker.pendingMu.Lock()
	plcWorker.pendingMachine = &machine
	plcWorker.pendingMu.Unlock()
	plcWorker.requestReload()
}

func (plcWorker *plcWorker) applyPendingMachineUpdate() {
	plcWorker.pendingMu.Lock()
	defer plcWorker.pendingMu.Unlock()
	if plcWorker.pendingMachine != nil {
		plcWorker.machine = *plcWorker.pendingMachine
		plcWorker.pendingMachine = nil
	}
}

// run is the worker goroutine: load variables, collect, and on failure wait and try again,
// until ctx is cancelled.
func (plcWorker *plcWorker) run(ctx context.Context) {
	retry := backoff{max: plcRetryMax}
	plcWorker.setStatus(ctx, func(status *opcua.MachineStatus) {
		status.State = opcua.StateConnecting
		status.Endpoint = opcua.Endpoint(plcWorker.machine.LiveIP, plcWorker.machine.Port)
		status.SecurityMode = plcWorker.machine.SecurityMode
	})
	for ctx.Err() == nil {
		plcWorker.applyPendingMachineUpdate()
		count, err := plcWorker.loadVariables(ctx)
		if err == nil && count == 0 {
			plcWorker.setStatus(ctx, func(status *opcua.MachineStatus) {
				status.State, status.ConnectedSince = opcua.StateNoVariables, nil
			})
			plcWorker.log(slog.LevelInfo, "plc: no PLC variable to read, not connecting")
			plcWorker.waitForReload(ctx, 0)
			continue
		}
		if err == nil {
			err = plcWorker.collect(ctx, &retry)
		}
		if ctx.Err() != nil {
			return
		}
		if errors.Is(err, errPlcReload) {
			continue
		}
		wait := retry.nextWait()
		now := time.Now()
		plcWorker.setStatus(ctx, func(status *opcua.MachineStatus) {
			status.State, status.ConnectedSince = opcua.StateFailed, nil
			status.LastError, status.LastErrorAt = err.Error(), &now
		})
		plcWorker.log(slog.LevelWarn, "plc: not reading, retrying", "retryIn", wait.String(), "error", err)
		plcWorker.waitForReload(ctx, wait)
	}
}

// waitForReload waits for a reload request, wait (0 = no limit) or the end of ctx.
func (plcWorker *plcWorker) waitForReload(ctx context.Context, wait time.Duration) {
	var timeout <-chan time.Time
	if wait > 0 {
		timer := time.NewTimer(wait)
		defer timer.Stop()
		timeout = timer.C
	}
	select {
	case <-ctx.Done():
	case <-plcWorker.reload:
	case <-timeout:
	}
}

// loadVariables reads the machine's PLC variables: the active variables whose address the
// collector would read from source plc (service.VariableSources), and returns how many.
func (plcWorker *plcWorker) loadVariables(ctx context.Context) (int, error) {
	sources := service.VariableSources(plcWorker.machine.CNCEngine, true)
	configs, err := plcWorker.collector.service.GetVariableConfig(ctx, plcWorker.machine, sources)
	if err != nil {
		return 0, fmt.Errorf("load variables: %w", err)
	}
	plcWorker.variables = map[uint32]model.VariableConfig{}
	for _, variable := range configs {
		if variable.Source == model.VariableSourcePLC && variable.Address != "" {
			plcWorker.variables[uint32(variable.ID)] = variable
		}
	}
	return len(plcWorker.variables), nil
}

// collect connects, subscribes and stores changes until the connection is lost, the health read
// fails, a reload is asked for (errPlcReload) or ctx ends. retry is reset once subscribed.
func (plcWorker *plcWorker) collect(ctx context.Context, retry *backoff) error {
	machine := plcWorker.machine
	settings := opcua.Settings{
		Endpoint:       opcua.Endpoint(machine.LiveIP, machine.Port),
		SecurityMode:   machine.SecurityMode,
		Username:       machine.Username,
		Password:       machine.Password,
		Certificate:    plcWorker.collector.certificate,
		ConnectTimeout: plcWorker.collector.config.PLC.ConnectTimeout,
		SessionTimeout: plcWorker.collector.config.PLC.SessionTimeout,
	}
	plcWorker.setStatus(ctx, func(status *opcua.MachineStatus) {
		status.State, status.Endpoint, status.SecurityMode = opcua.StateConnecting, settings.Endpoint, settings.SecurityMode
	})
	connection, err := opcua.Connect(ctx, settings)
	if err != nil {
		return err
	}
	defer func() {
		closeCtx, cancel := context.WithTimeout(context.Background(), plcCloseTimeout)
		defer cancel()
		connection.Close(closeCtx)
	}()

	sessionCtx, stopSession := context.WithCancel(ctx)
	defer stopSession()
	plc := plcWorker.collector.config.PLC
	results, changes, publishErrors, err := connection.Subscribe(sessionCtx, plc.PublishInterval, plc.SamplingInterval, plcWorker.items())
	if err != nil {
		return err
	}
	plcWorker.recordSubscribeResults(ctx, results, connection.SecurityPolicy)
	plcWorker.loadStoredValues(ctx)
	retry.reset()

	// nodes refused for lack of room are read instead; nil pollTick = nothing to read
	pollHandles, pollNodeIDs := plcWorker.noRoomNodes(results)
	var pollTick <-chan time.Time
	if len(pollHandles) > 0 {
		pollInterval := plc.PollInterval
		if pollInterval <= 0 {
			pollInterval = plcDefaultPollInterval
		}
		poller := time.NewTicker(pollInterval)
		defer poller.Stop()
		pollTick = poller.C
		plcWorker.polled = map[uint32]opcua.Change{}
		if err := plcWorker.poll(ctx, connection, pollHandles, pollNodeIDs); err != nil {
			return err
		}
	}

	health := time.NewTicker(plc.HealthInterval)
	defer health.Stop()
	failedPublishes := 0
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-plcWorker.reload:
			plcWorker.log(slog.LevelInfo, "plc: variables changed, subscribing again")
			return errPlcReload
		case <-connection.Lost():
			return errors.New("connection lost")
		case batch := <-changes:
			failedPublishes = 0
			plcWorker.handleChanges(ctx, batch)
		case <-pollTick:
			if err := plcWorker.poll(ctx, connection, pollHandles, pollNodeIDs); err != nil {
				return err
			}
		case publishErr := <-publishErrors:
			failedPublishes++
			plcWorker.log(slog.LevelWarn, "plc: publish failed", "inARow", failedPublishes, "error", publishErr)
			if failedPublishes >= plcMaxPublishErrors {
				return fmt.Errorf("publish failed %d times: %w", failedPublishes, publishErr)
			}
		case <-health.C:
			checkCtx, cancel := context.WithTimeout(ctx, plc.ConnectTimeout)
			err := connection.CheckAlive(checkCtx)
			cancel()
			if err != nil {
				return err
			}
		}
	}
}

// items lists the subscription's nodes: one per variable, its id as the handle.
func (plcWorker *plcWorker) items() []opcua.Item {
	items := make([]opcua.Item, 0, len(plcWorker.variables))
	for handle, variable := range plcWorker.variables {
		items = append(items, opcua.Item{Handle: handle, NodeID: variable.Address})
	}
	return items
}

// noRoomNodes lists the handles and node ids the controller refused for lack of room only.
func (plcWorker *plcWorker) noRoomNodes(results []opcua.ItemResult) (handles []uint32, nodeIDs []string) {
	for _, result := range results {
		if result.NoRoom {
			handles = append(handles, result.Handle)
			nodeIDs = append(nodeIDs, plcWorker.variables[result.Handle].Address)
		}
	}
	return handles, nodeIDs
}

// poll reads nodeIDs in one request and handles the values that changed since the last read
// like a publish. A failed read ends the session, like a failed health read.
func (plcWorker *plcWorker) poll(ctx context.Context, connection *opcua.Connection, handles []uint32, nodeIDs []string) error {
	readCtx, cancel := context.WithTimeout(ctx, plcWorker.collector.config.PLC.ConnectTimeout)
	defer cancel()
	results, err := connection.ReadValues(readCtx, nodeIDs)
	if err != nil {
		return err
	}
	if batch := polledChanges(handles, results, plcWorker.polled); len(batch) > 0 {
		plcWorker.handleChanges(ctx, batch)
	}
	return nil
}

// polledChanges turns read results (results[i] is handles[i]) into the changes since last, the
// previous read of each handle, and moves last on: an unchanged value or status is left out, so
// reading sends no more than a subscription would.
func polledChanges(handles []uint32, results []opcua.ReadResult, last map[uint32]opcua.Change) []opcua.Change {
	var changes []opcua.Change
	for index, result := range results {
		if index >= len(handles) {
			break
		}
		change := opcua.Change{Handle: handles[index], Status: result.Status}
		if result.Status == "" {
			change.Value, change.Good = result.Value, true
		}
		if previous, found := last[change.Handle]; found && previous == change {
			continue
		}
		last[change.Handle] = change
		changes = append(changes, change)
	}
	return changes
}

// recordSubscribeResults marks the worker connected, keeps each node's status for the
// diagnostics, tells the hub which codes it reads, and logs the nodes the controller refused and
// the ones it had no room for (read instead, see poll).
func (plcWorker *plcWorker) recordSubscribeResults(ctx context.Context, results []opcua.ItemResult, securityPolicy string) {
	now := time.Now()
	codes := make([]string, 0, len(results))
	var refused, readInstead []string
	nodes := make(map[int]opcua.NodeStatus, len(results))
	for _, result := range results {
		variable := plcWorker.variables[result.Handle]
		status := result.Status
		if result.NoRoom {
			status = ""
			readInstead = append(readInstead, fmt.Sprintf("%d:%s %s", variable.ID, variable.Code, variable.Address))
		}
		nodes[variable.ID] = opcua.NodeStatus{NodeID: variable.Address, Status: status, At: now}
		if status != "" {
			refused = append(refused, fmt.Sprintf("%d:%s %s (%s)", variable.ID, variable.Code, variable.Address, result.Status))
			continue
		}
		codes = append(codes, variable.Code)
	}
	plcWorker.collector.hub.SetMachineCodes(plcWorker.machine.ID, live.ReaderPLC, codes)
	plcWorker.setStatus(ctx, func(status *opcua.MachineStatus) {
		status.State, status.ConnectedSince, status.SecurityPolicy = opcua.StateConnected, &now, securityPolicy
		status.Nodes = nodes
	})
	plcWorker.log(slog.LevelInfo, "plc: connected and subscribed", "endpoint", opcua.Endpoint(plcWorker.machine.LiveIP, plcWorker.machine.Port),
		"securityPolicy", securityPolicy, "variables", len(results), "subscribed", len(codes)-len(readInstead), "readInstead", len(readInstead))
	if len(refused) > 0 {
		plcWorker.log(slog.LevelWarn, "plc: nodes refused by the controller", "count", len(refused), "nodes", refused)
	}
	if len(readInstead) > 0 {
		plcWorker.log(slog.LevelWarn, "plc: controller has no room to monitor these nodes, reading them instead",
			"count", len(readInstead), "interval", plcWorker.collector.config.PLC.PollInterval.String(), "nodes", readInstead)
	}
}

// loadStoredValues reads the last stored numbers, so the first change of each variable after a
// restart is stored in machine_telementry only when it differs.
func (plcWorker *plcWorker) loadStoredValues(ctx context.Context) {
	variableIDs := make([]int, 0, len(plcWorker.variables))
	for _, variable := range plcWorker.variables {
		variableIDs = append(variableIDs, variable.ID)
	}
	plcWorker.stored = map[int]float64{}
	readCtx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	values, err := plcWorker.collector.service.GetLiveValues(readCtx, plcWorker.machine.ID, variableIDs)
	if err != nil {
		plcWorker.log(slog.LevelError, "plc: load stored values", "error", err)
		return
	}
	for variableID, value := range values {
		if value.Value != nil {
			plcWorker.stored[variableID] = *value.Value
		}
	}
}

// handleChanges takes one publish: the screens get the values first, then every number is
// upserted and the changed ones are added to the time series.
func (plcWorker *plcWorker) handleChanges(ctx context.Context, batch []opcua.Change) {
	now := time.Now()
	published := make(map[string]any, len(batch))
	liveRows := make([]model.LiveValue, 0, len(batch))
	telemetryRows := make([]model.TelemetryValue, 0, len(batch))
	var changeLog []string
	nodeStatuses := map[int]string{}
	for _, change := range batch {
		variable, found := plcWorker.variables[change.Handle]
		if !found {
			continue
		}
		nodeStatuses[variable.ID] = change.Status
		if !change.Good {
			continue
		}
		if change.Value.IsText {
			published[variable.Code] = change.Value.Text
			continue
		}
		number := change.Value.Number
		published[variable.Code] = number
		if variable.IsText() {
			continue // a text variable is published, never stored (as the CNC worker does)
		}
		liveRows = append(liveRows, model.LiveValue{MachineID: plcWorker.machine.ID, VariableID: variable.ID, Value: &number, UpdatedAt: now})
		if previous, found := plcWorker.stored[variable.ID]; !found || previous != number {
			telemetryRows = append(telemetryRows, model.TelemetryValue{MachineID: plcWorker.machine.ID, VariableID: variable.ID, Value: &number, RecordedAt: now})
			changeLog = append(changeLog, fmt.Sprintf("%d:%s %s -> %s", variable.ID, variable.Code, previousText(plcWorker.stored, variable.ID),
				strconv.FormatFloat(number, 'g', -1, 64)))
		}
	}

	plcWorker.collector.hub.PublishValues(plcWorker.machine.ID, now, published)
	plcWorker.store(ctx, now, liveRows, telemetryRows)
	plcWorker.setStatus(ctx, func(status *opcua.MachineStatus) {
		status.LastChangeAt = &now
		status.ChangeCount += int64(len(batch))
		for variableID, nodeStatus := range nodeStatuses {
			node := status.Nodes[variableID]
			node.Status, node.At = nodeStatus, now
			status.Nodes[variableID] = node
		}
	})
	if len(changeLog) > 0 {
		writeLog(plcWorker.collector.config.Logs, logTypeVariable, slog.LevelInfo, &plcWorker.machine,
			"plc: variable values changed", "count", len(changeLog), "changes", changeLog)
	}
}

// previousText is the stored number of variableID for the change log, "none" before the first.
func previousText(stored map[int]float64, variableID int) string {
	if value, found := stored[variableID]; found {
		return strconv.FormatFloat(value, 'g', -1, 64)
	}
	return "none"
}

// store hands liveRows and telemetryRows, read at readAt, to the telemetry queue
// (RABBITMQ_ENABLED), or upserts and inserts them directly when the queue is off or cannot take
// them; the stored numbers move on only when the rows are queued or written, so a failed insert
// is retried with the next change.
func (plcWorker *plcWorker) store(ctx context.Context, readAt time.Time, liveRows []model.LiveValue, telemetryRows []model.TelemetryValue) {
	if len(liveRows) == 0 {
		return
	}
	telemetryQueue := plcWorker.collector.queue
	if telemetryQueue != nil &&
		telemetryQueue.Enqueue(model.NewTelemetryReading(plcWorker.machine.ID, readAt, model.TelemetrySourcePLC, liveRows, telemetryRows)) {
		plcWorker.moveStoredOn(liveRows)
		return
	}
	writeCtx, cancel := context.WithTimeout(ctx, writeTimeout)
	defer cancel()
	if err := plcWorker.collector.service.UpsertLiveValues(writeCtx, liveRows); err != nil {
		plcWorker.log(slog.LevelError, "plc: upsert live values", "rows", len(liveRows), "error", err)
		return
	}
	if len(telemetryRows) > 0 {
		if err := plcWorker.collector.service.InsertTelemetry(writeCtx, telemetryRows); err != nil {
			plcWorker.log(slog.LevelError, "plc: insert telemetry", "rows", len(telemetryRows), "error", err)
			return
		}
	}
	plcWorker.moveStoredOn(liveRows)
}

// moveStoredOn takes the numbers of liveRows as the stored ones (change detection).
func (plcWorker *plcWorker) moveStoredOn(liveRows []model.LiveValue) {
	for _, row := range liveRows {
		plcWorker.stored[row.VariableID] = *row.Value
	}
}

// setStatus updates the machine's diagnostics status, unless the worker is stopping (a newer
// worker of the same machine may already report).
func (plcWorker *plcWorker) setStatus(ctx context.Context, change func(status *opcua.MachineStatus)) {
	if ctx.Err() != nil {
		return
	}
	plcWorker.collector.board.Update(plcWorker.machine.ID, change)
}

// log writes a PLC line about this machine (CNC_LOG_PLC).
func (plcWorker *plcWorker) log(level slog.Level, message string, attributes ...any) {
	writeLog(plcWorker.collector.config.Logs, logTypePLC, level, &plcWorker.machine, message, attributes...)
}
