package tool

import "net/http"

func GetClientIp(r *http.Request) string {
	xRealIp := r.Header.Get("X-Real-Ip")
	if len(xRealIp) > 0 {
		return xRealIp
	}
	xForwardedForIp := r.Header.Get("X-Forwarded-For")
	if len(xForwardedForIp) > 0 {
		return xForwardedForIp
	}
	return r.RemoteAddr
}
