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
	HandleIpBetween(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "74.50.153.4", "74.50.153.2", false)
	HandleIpBetween(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:7334", true)
	HandleIpBetween(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:7350", true)
	HandleIpBetween(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", true)
	HandleIpBetween(t, "2001:0db8:85a3:0000:0000:8a2e:0370:7334", "2001:0db8:85a3:0000:0000:8a2e:0370:8334", "2001:0db8:85a3:0000:0000:8a2e:0370:8335", false)
	HandleIpBetween(t, "::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.127", false)
	HandleIpBetween(t, "::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.128", true)
	HandleIpBetween(t, "::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.129", true)
	HandleIpBetween(t, "::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.250", true)
	HandleIpBetween(t, "::ffff:192.0.2.128", "::ffff:192.0.2.250", "::ffff:192.0.2.251", false)
	HandleIpBetween(t, "::ffff:192.0.2.128", "::ffff:192.0.2.250", "192.0.2.130", true)
	HandleIpBetween(t, "192.0.2.128", "192.0.2.250", "::ffff:192.0.2.130", true)
	HandleIpBetween(t, "idonotparse", "192.0.2.250", "::ffff:192.0.2.130", false)

}

func HandleIpBetween(t *testing.T, from string, to string, test string, assert bool) {
	res := ipBetween(net.ParseIP(from), net.ParseIP(to), net.ParseIP(test))
	if res != assert {
		t.Errorf("Assertion (have: %v should be: %v) failed on range %s-%s with test %s", res, assert, from, to, test)
	}
}
