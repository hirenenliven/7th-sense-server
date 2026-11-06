package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"7th-sense-server/internal/dto"
)

const (
	portProbeTimeout = 4 * time.Second // Windows retries a refused connect for about 2 s
	pingTimeout      = 3 * time.Second
	dnsTimeout       = 2 * time.Second
)

// How a TCP connect ended.
type portResult int

const (
	portOpen portResult = iota
	portRefused
	portTimedOut
	portUnreachable
	portOtherError
)

// Windows socket error numbers syscall does not name (WSAECONNREFUSED, WSAENETUNREACH,
// WSAEHOSTUNREACH); the Unix ones have names.
const (
	windowsConnectionRefused  = 10061
	windowsNetworkUnreachable = 10051
	windowsHostUnreachable    = 10065
)

// probeNetwork tests host:port from this server, the place the collector connects from. A TCP
// connect is the only proof about the port; for the IP, any answer counts, because many plant
// networks block ping and a refused connect also shows the host is there.
func probeNetwork(ctx context.Context, host string, port int) dto.MachineNetworkResponse {
	address := net.JoinHostPort(host, strconv.Itoa(port))
	result := dto.MachineNetworkResponse{CheckedAt: time.Now(), Address: address}

	ip, err := resolveIP(ctx, host)
	if err != nil {
		result.IP = dto.DiagnosticCheck{Key: "ip", Status: dto.DiagnosticFailed,
			Message: fmt.Sprintf("%q is not an IP address and could not be resolved: %v", host, err),
			Hint:    "Correct the Live IP on the machine form."}
		result.Port = dto.DiagnosticCheck{Key: "port", Status: dto.DiagnosticSkipped, Message: "Not tested: the IP is not usable."}
		return result
	}

	var replied bool
	var outcome portResult
	var latency time.Duration
	var dialError error
	var group sync.WaitGroup
	group.Go(func() { replied = pingHost(ctx, ip) })
	group.Go(func() { outcome, latency, dialError = dialPort(ctx, net.JoinHostPort(ip, strconv.Itoa(port))) })
	group.Wait()

	result.IP = ipCheck(ip, replied, outcome)
	result.Port = portCheck(port, outcome, dialError)
	if outcome == portOpen {
		milliseconds := latency.Milliseconds()
		result.LatencyMs = &milliseconds
	}
	return result
}

// resolveIP returns host when it is an IP address, else the first address DNS gives. ping only
// ever receives this, so a stored value can never be taken as a command line option.
func resolveIP(ctx context.Context, host string) (string, error) {
	if ip := net.ParseIP(host); ip != nil {
		return ip.String(), nil
	}
	lookup, cancel := context.WithTimeout(ctx, dnsTimeout)
	defer cancel()
	addresses, err := net.DefaultResolver.LookupHost(lookup, host)
	if err != nil {
		return "", err
	}
	return addresses[0], nil
}

// pingHost reports whether ip answered one ping. Replies are recognised by their TTL: Windows ping
// also exits 0 when a router answers "destination host unreachable", and its text is translated.
func pingHost(ctx context.Context, ip string) bool {
	probe, cancel := context.WithTimeout(ctx, pingTimeout)
	defer cancel()
	arguments := []string{"-c", "1", "-W", "2", ip}
	if runtime.GOOS == "windows" {
		arguments = []string{"-n", "1", "-w", "2000", ip}
	}
	output, err := exec.CommandContext(probe, "ping", arguments...).CombinedOutput()
	return err == nil && strings.Contains(strings.ToLower(string(output)), "ttl=")
}

// dialPort opens and closes one TCP connection. It sends no data, so the controller's protocol
// sessions and the collector's connection are not disturbed.
func dialPort(ctx context.Context, address string) (portResult, time.Duration, error) {
	dialer := net.Dialer{Timeout: portProbeTimeout}
	started := time.Now()
	connection, err := dialer.DialContext(ctx, "tcp", address)
	latency := time.Since(started)
	if err != nil {
		return classifyDialError(err), latency, err
	}
	_ = connection.Close()
	return portOpen, latency, nil
}

func classifyDialError(err error) portResult {
	var timeout interface{ Timeout() bool }
	if errors.As(err, &timeout) && timeout.Timeout() {
		return portTimedOut
	}
	var errno syscall.Errno
	if errors.As(err, &errno) {
		switch {
		case errno == syscall.ECONNREFUSED || errno == windowsConnectionRefused:
			return portRefused
		case errno == syscall.ENETUNREACH || errno == syscall.EHOSTUNREACH ||
			errno == windowsNetworkUnreachable || errno == windowsHostUnreachable:
			return portUnreachable
		}
	}
	return portOtherError
}

func ipCheck(ip string, replied bool, outcome portResult) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: "ip", Status: dto.DiagnosticOK}
	switch {
	case replied:
		check.Message = fmt.Sprintf("%s replied to ping from this server.", ip)
	case outcome == portOpen || outcome == portRefused:
		check.Message = fmt.Sprintf("%s is reachable: it answered on TCP. It did not reply to ping, which is often blocked.", ip)
	case outcome == portTimedOut:
		// Silence on both says little more than "no answer": a firewall looks the same as a machine that is off
		check.Status = dto.DiagnosticFailed
		check.Message = fmt.Sprintf("%s did not reply to ping or TCP from this server.", ip)
		check.Hint = "Check the machine is powered on, the IP is right, the plant network or VPN is up and no firewall blocks this server."
	default:
		check.Status = dto.DiagnosticFailed
		check.Message = fmt.Sprintf("This server has no route to %s.", ip)
		check.Hint = "Check the server's network, VPN and gateway reach the machine's network."
	}
	return check
}

func portCheck(port int, outcome portResult, dialError error) dto.DiagnosticCheck {
	check := dto.DiagnosticCheck{Key: "port", Status: dto.DiagnosticFailed}
	switch outcome {
	case portOpen:
		check.Status = dto.DiagnosticOK
		check.Message = fmt.Sprintf("Port %d accepted a TCP connection.", port)
		check.Hint = "This proves the port only; the Controller check shows whether the controller answers on it."
	case portRefused:
		check.Message = fmt.Sprintf("The machine refused port %d: it is on, but nothing listens on that port.", port)
		check.Hint = "Check the port on the machine form (Fanuc FOCAS is usually 8193, Siemens OPC UA 4840) and that the service is enabled on the controller."
	case portTimedOut:
		check.Message = fmt.Sprintf("No answer on port %d within %s.", port, portProbeTimeout)
		check.Hint = "A firewall may be dropping it, or the machine is off."
	case portUnreachable:
		check.Message = fmt.Sprintf("Port %d could not be tried: this server has no route to the machine.", port)
	default:
		check.Message = fmt.Sprintf("Port %d could not be tested: %v", port, dialError)
	}
	return check
}
