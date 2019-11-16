package ipchecker

import (
	"bytes"
	"fmt"
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

		IPAddress := r.RemoteAddr
		fmt.Println(IPAddress)
		if !check(IPAddress) {
			return
		}

		// Token is authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

//check checks if ip belongs to telegram subnet
func check(ip string) bool {
	trial := net.IP(ip)
	return ipBetween(f1, t1, trial) && ipBetween(f2, t2, trial)
}

func ipBetween(from, to, trial net.IP) bool {
	if trial.To4() == nil {
		return false
	}
	if bytes.Compare(trial, from) >= 0 && bytes.Compare(trial, to) <= 0 {
		return true
	}

	return false
}
