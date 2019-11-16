package ipchecker

import (
	"net"
	"testing"
)

func TestIPBetween(t *testing.T) {
	HandleIpBetween(t, "0.0.0.0", "255.255.255.255", "128.128.128.128", true)
	HandleIpBetween(t, "0.0.0.0", "128.128.128.128", "255.255.255.255", false)
	HandleIpBetween(t, "74.50.153.0", "74.50.153.4", "74.50.153.0", true)
	HandleIpBetween(t, "74.50.153.0", "74.50.153.4", "74.50.153.4", true)
	HandleIpBetween(t, "74.50.153.0", "74.50.153.4", "74.50.153.5", false)

}

func HandleIpBetween(t *testing.T, from string, to string, test string, assert bool) {
	res := ipBetween(net.ParseIP(from), net.ParseIP(to), net.ParseIP(test))
	if res != assert {
		t.Errorf("Assertion (have: %v should be: %v) failed on range %s-%s with test %s", res, assert, from, to, test)
	}
}
