// Package opcua connects to a controller's OPC UA server with gopcua (github.com/gopcua/opcua,
// ADR 0021) and subscribes to nodes: the protocol side of the PLC reader, like internal/cnc is
// for the DLL. It holds no database or business rules; internal/monitor runs it per machine.
package opcua

import (
	"context"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"sync"
	"time"

	gopcua "github.com/gopcua/opcua"
	"github.com/gopcua/opcua/id"
	"github.com/gopcua/opcua/ua"

	"7th-sense-server/internal/cert"
)

// Security modes of a machine (machines.security_mode, dto.SecurityMode*).
const (
	SecurityModeNone           = "None"
	SecurityModeBasic256Sha256 = "Basic256Sha256"
)

// nodeIDPattern is the node id form the controller understands, e.g. "ns=2;s=/Random/Plc/DB3300.DBX3.0"
// (the same rule as dto.UpdatePlcNodeIDRequest): ua.ParseNodeID alone takes any text as a string id.
var nodeIDPattern = regexp.MustCompile(`^(ns=\d+;)?[sigb]=\S.*$`)

// queueSize is how many changes the controller keeps per node between two publishes; the oldest
// is dropped beyond it (only the newest value matters for the live data).
const queueSize = 10

// Settings say how to reach and log in to one controller.
type Settings struct {
	Endpoint       string // opc.tcp://host:port (Endpoint)
	SecurityMode   string // SecurityModeNone or SecurityModeBasic256Sha256
	Username       string // empty = anonymous
	Password       string
	Certificate    cert.Paths // presented when SecurityMode is not None
	ConnectTimeout time.Duration
}

// Endpoint is the OPC UA address of host and port.
func Endpoint(host string, port int) string {
	return "opc.tcp://" + net.JoinHostPort(host, strconv.Itoa(port))
}

// Item is one node to subscribe to; Handle comes back with each of its changes.
type Item struct {
	Handle uint32
	NodeID string
}

// ItemResult says whether the controller accepted an item: Status "" when it did, else why not
// (an invalid node id, or the controller's status such as BadNodeIdUnknown).
type ItemResult struct {
	Handle uint32
	Status string
}

// Change is one new value of a subscribed item. Good is false when the controller marks the value
// bad (Status names why); Value is then not to be used.
type Change struct {
	Handle uint32
	Value  Value
	Good   bool
	Status string
}

// Connection is an open session with one controller. It never reconnects itself: Lost is closed
// once the connection is gone, and the caller closes it and connects again.
type Connection struct {
	client   *gopcua.Client
	lost     chan struct{}
	lostOnce sync.Once
	// SecurityPolicy is the policy URI of the endpoint used (for the diagnostics).
	SecurityPolicy string
}

// Connect discovers the controller's endpoints, picks the one of settings.SecurityMode and opens
// a session, all within settings.ConnectTimeout.
func Connect(ctx context.Context, settings Settings) (*Connection, error) {
	ctx, cancel := context.WithTimeout(ctx, settings.ConnectTimeout)
	defer cancel()

	endpoints, err := gopcua.GetEndpoints(ctx, settings.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("get endpoints of %s: %w", settings.Endpoint, err)
	}
	policy, mode := securityOf(settings.SecurityMode)
	endpoint, err := gopcua.SelectEndpoint(endpoints, policy, mode)
	if err != nil {
		return nil, fmt.Errorf("no %s endpoint at %s: %w", settings.SecurityMode, settings.Endpoint, err)
	}

	connection := &Connection{lost: make(chan struct{}), SecurityPolicy: endpoint.SecurityPolicyURI}
	options := []gopcua.Option{
		gopcua.ApplicationURI(cert.ApplicationURI),
		gopcua.ApplicationName("7th Sense"),
		gopcua.AutoReconnect(false),
		gopcua.RequestTimeout(settings.ConnectTimeout),
		gopcua.StateChangedFunc(func(state gopcua.ConnState) {
			if state == gopcua.Disconnected || state == gopcua.Closed {
				connection.markLost()
			}
		}),
	}
	if mode != ua.MessageSecurityModeNone {
		options = append(options,
			gopcua.CertificateFile(settings.Certificate.CertificateFile),
			gopcua.PrivateKeyFile(settings.Certificate.PrivateKeyFile))
	}
	// The identity token goes first: SecurityFromEndpoint then sets its policy id.
	authType := ua.UserTokenTypeAnonymous
	if settings.Username != "" {
		authType = ua.UserTokenTypeUserName
		options = append(options, gopcua.AuthUsername(settings.Username, settings.Password))
	} else {
		options = append(options, gopcua.AuthAnonymous())
	}
	options = append(options, gopcua.SecurityFromEndpoint(endpoint, authType))

	client, err := gopcua.NewClient(settings.Endpoint, options...)
	if err != nil {
		return nil, fmt.Errorf("configure client: %w", err)
	}
	if err := client.Connect(ctx); err != nil {
		return nil, fmt.Errorf("connect to %s: %w", settings.Endpoint, err)
	}
	connection.client = client
	return connection, nil
}

// securityOf maps a machine's security mode to the endpoint policy and message mode it needs.
func securityOf(securityMode string) (string, ua.MessageSecurityMode) {
	if securityMode == SecurityModeBasic256Sha256 {
		return ua.SecurityPolicyURIBasic256Sha256, ua.MessageSecurityModeSignAndEncrypt
	}
	return ua.SecurityPolicyURINone, ua.MessageSecurityModeNone
}

// Lost is closed once the connection is gone (secure channel error or closed).
func (connection *Connection) Lost() <-chan struct{} {
	return connection.lost
}

func (connection *Connection) markLost() {
	connection.lostOnce.Do(func() { close(connection.lost) })
}

// Close ends the session and the connection.
func (connection *Connection) Close(ctx context.Context) {
	if connection.client != nil {
		_ = connection.client.Close(ctx)
	}
	connection.markLost()
}

// CheckAlive reads the server's current time: an error means the connection is dead even though
// no channel error came (subscriptions are silent while no value changes).
func (connection *Connection) CheckAlive(ctx context.Context) error {
	response, err := connection.client.Read(ctx, &ua.ReadRequest{
		NodesToRead: []*ua.ReadValueID{{NodeID: ua.NewNumericNodeID(0, id.Server_ServerStatus_CurrentTime), AttributeID: ua.AttributeIDValue}},
	})
	if err != nil {
		return fmt.Errorf("read server time: %w", err)
	}
	if len(response.Results) == 0 || isBad(response.Results[0].Status) {
		return errors.New("read server time: no good value")
	}
	return nil
}

// Subscribe creates one subscription publishing every publishInterval and a monitored item per
// item, sampled every samplingInterval. It returns each item's result (an item refused does not
// stop the others) and the channel its changes arrive on until ctx ends or the connection is
// lost. A publish error is sent as a Change-less error on errs.
func (connection *Connection) Subscribe(ctx context.Context, publishInterval, samplingInterval time.Duration, items []Item) (
	results []ItemResult, changes <-chan []Change, errs <-chan error, err error) {
	notifications := make(chan *gopcua.PublishNotificationData, 64)
	subscription, err := connection.client.Subscribe(ctx, &gopcua.SubscriptionParameters{Interval: publishInterval}, notifications)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create subscription: %w", err)
	}

	results = make([]ItemResult, 0, len(items))
	requests := make([]*ua.MonitoredItemCreateRequest, 0, len(items))
	for _, item := range items {
		if !nodeIDPattern.MatchString(item.NodeID) {
			results = append(results, ItemResult{Handle: item.Handle, Status: "invalid node id: expected ns=<n>;s=<path>"})
			continue
		}
		nodeID, parseErr := ua.ParseNodeID(item.NodeID)
		if parseErr != nil {
			results = append(results, ItemResult{Handle: item.Handle, Status: "invalid node id: " + parseErr.Error()})
			continue
		}
		request := gopcua.NewMonitoredItemCreateRequestWithDefaults(nodeID, ua.AttributeIDValue, item.Handle)
		request.RequestedParameters.SamplingInterval = float64(samplingInterval / time.Millisecond)
		request.RequestedParameters.QueueSize = queueSize
		requests = append(requests, request)
	}
	if len(requests) > 0 {
		response, monitorErr := subscription.Monitor(ctx, ua.TimestampsToReturnBoth, requests...)
		if monitorErr != nil {
			return nil, nil, nil, fmt.Errorf("monitor items: %w", monitorErr)
		}
		for index, request := range requests {
			result := ItemResult{Handle: request.RequestedParameters.ClientHandle}
			if index >= len(response.Results) {
				result.Status = "no result from the controller"
			} else if status := response.Results[index].StatusCode; status != ua.StatusOK {
				result.Status = statusName(status)
			}
			results = append(results, result)
		}
	}

	changeChannel := make(chan []Change, 16)
	errorChannel := make(chan error, 4)
	go forwardNotifications(ctx, connection.lost, notifications, changeChannel, errorChannel)
	return results, changeChannel, errorChannel, nil
}

// forwardNotifications turns gopcua's publish notifications into Change batches (one per
// publish) and errors, until ctx ends or the connection is lost.
func forwardNotifications(ctx context.Context, lost <-chan struct{}, notifications <-chan *gopcua.PublishNotificationData,
	changes chan<- []Change, errs chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-lost:
			return
		case notification := <-notifications:
			if notification == nil {
				continue
			}
			if notification.Error != nil {
				select {
				case errs <- notification.Error:
				default: // the reader is behind; a newer error will follow
				}
				continue
			}
			dataChange, ok := notification.Value.(*ua.DataChangeNotification)
			if !ok {
				continue
			}
			batch := make([]Change, 0, len(dataChange.MonitoredItems))
			for _, item := range dataChange.MonitoredItems {
				batch = append(batch, changeOf(item))
			}
			select {
			case changes <- batch:
			case <-ctx.Done():
				return
			case <-lost:
				return
			}
		}
	}
}

// changeOf converts one monitored item notification.
func changeOf(item *ua.MonitoredItemNotification) Change {
	change := Change{Handle: item.ClientHandle}
	if item.Value == nil {
		change.Status = "no value"
		return change
	}
	if isBad(item.Value.Status) {
		change.Status = statusName(item.Value.Status)
		return change
	}
	value, ok := ValueOf(item.Value.Value)
	if !ok {
		change.Status = "unsupported value type"
		return change
	}
	change.Value, change.Good = value, true
	return change
}

// isBad reports whether status has severity Bad (the top two bits 10).
func isBad(status ua.StatusCode) bool {
	return uint32(status)>>30 == 2
}

// statusName is the controller's status code in words, e.g. "BadNodeIdUnknown".
func statusName(status ua.StatusCode) string {
	if description, ok := ua.StatusCodes[status]; ok {
		return description.Name
	}
	return fmt.Sprintf("status 0x%08X", uint32(status))
}
