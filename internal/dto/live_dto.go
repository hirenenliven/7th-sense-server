package dto

import (
	"net/http"
	"slices"
	"time"

	"7th-sense-server/internal/live"
	"7th-sense-server/internal/response"
)

// Limits of one Socket.IO subscribe message (docs/live-socket-events.md).
const (
	LiveMaxMachines   = 500 // machines per subscribe: the superadmin's whole Monitoring list
	LiveMaxCodes      = 100 // variable codes per subscribe
	LiveMaxCodeLength = 64  // cnc_variables.code is short; longer is not a code
)

// LiveChannels are the channels a live subscription may ask for besides codes.
var LiveChannels = []string{live.ChannelStatus, live.ChannelEvent}

// LiveSubscribeRequest is the Socket.IO "subscribe" message: for these machines send these
// variable codes at most every intervalMs (0 = the default) and these channels. A user id or
// role in the message is ignored: identity always comes from the session.
type LiveSubscribeRequest struct {
	MachineIDs []int64  `json:"machineIds"`
	Codes      []string `json:"codes"`
	Channels   []string `json:"channels"`
	IntervalMs int      `json:"intervalMs"`
}

// Validate checks the machines, codes, channels and interval of a subscribe message.
func (subscribeRequest LiveSubscribeRequest) Validate() error {
	if len(subscribeRequest.MachineIDs) == 0 {
		return response.SelectMachine
	}
	if len(subscribeRequest.MachineIDs) > LiveMaxMachines || len(subscribeRequest.Codes) > LiveMaxCodes {
		return response.NewError(http.StatusBadRequest, "Too many machines or variables in one subscription")
	}
	for _, machineID := range subscribeRequest.MachineIDs {
		if machineID <= 0 {
			return response.SelectMachine
		}
	}
	for _, code := range subscribeRequest.Codes {
		if code == "" || len(code) > LiveMaxCodeLength {
			return response.NewError(http.StatusBadRequest, "Please send valid variable codes")
		}
	}
	for _, channel := range subscribeRequest.Channels {
		if !slices.Contains(LiveChannels, channel) {
			return response.NewError(http.StatusBadRequest, "Please send valid channels")
		}
	}
	if len(subscribeRequest.Codes) == 0 && len(subscribeRequest.Channels) == 0 {
		return response.NewError(http.StatusBadRequest, "Please send variable codes or channels")
	}
	if subscribeRequest.IntervalMs < 0 {
		return response.NewError(http.StatusBadRequest, "Please send a valid interval")
	}
	return nil
}

// Interval is the requested interval (0 = none asked; live.Hub applies the default and minimum).
func (subscribeRequest LiveSubscribeRequest) Interval() time.Duration {
	return time.Duration(subscribeRequest.IntervalMs) * time.Millisecond
}

// LiveUnsubscribeRequest is the Socket.IO "unsubscribe" message: stop every subscription of this
// socket for these machines.
type LiveUnsubscribeRequest struct {
	MachineIDs []int64 `json:"machineIds"`
}

// Validate checks the machines of an unsubscribe message.
func (unsubscribeRequest LiveUnsubscribeRequest) Validate() error {
	if len(unsubscribeRequest.MachineIDs) == 0 || len(unsubscribeRequest.MachineIDs) > LiveMaxMachines {
		return response.SelectMachine
	}
	return nil
}

// LiveSubscribeResponse acknowledges a subscribe: the machines the user may not read (Denied,
// nothing joined for them), the codes the machines do not read (Unknown) and the current values
// of the others. Message is set, and OK false, when the message itself was invalid.
type LiveSubscribeResponse struct {
	OK        bool            `json:"ok"`
	Message   string          `json:"message,omitempty"`
	Denied    []int64         `json:"denied"`
	Unknown   []string        `json:"unknown"`
	Snapshots []live.Snapshot `json:"snapshots"`
}

// LiveAccessChanged is the Socket.IO "access:changed" message: the machines the user may no
// longer read (their subscriptions were dropped) and the ones they may read now.
type LiveAccessChanged struct {
	Removed []int64 `json:"removed"`
	Added   []int64 `json:"added"`
}
