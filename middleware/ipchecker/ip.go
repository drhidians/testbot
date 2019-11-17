package ipchecker

import (
	"bytes"
	"net"
	"net/http"
)

//Telegram subnets 149.154.160.0/20 and 91.108.4.0/22.
var f1 = net.ParseIP("149.154.160.0")
var t1 = net.ParseIP("149.154.176.255")

var f2 = net.ParseIP("91.108.4.0")
var t2 = net.ParseIP("91.108.8.255")

//Check checks if incoming request is coming from telegram subnets
func Check(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		IPAddress := r.Header.Get("X-Real-Ip")
		if IPAddress == "" {
			IPAddress = r.Header.Get("X-Forwarded-For")
		}
		if IPAddress == "" {
			IPAddress = r.RemoteAddr
		}

		if !check(IPAddress) {
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

//check checks if ip belongs to telegram subnet
func check(ip string) bool {
	trial := net.ParseIP(ip)
	return ipBetween(f1, t1, trial) || ipBetween(f2, t2, trial)
}

func ipBetween(from net.IP, to net.IP, test net.IP) bool {
	if from == nil || to == nil || test == nil {
		return false
	}

	from16 := from.To16()
	to16 := to.To16()
	test16 := test.To16()
	if from16 == nil || to16 == nil || test16 == nil {
		return false
	}

	if bytes.Compare(test16, from16) >= 0 && bytes.Compare(test16, to16) <= 0 {
		return true
	}
	return false
}
