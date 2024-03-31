package common

import "net/http"

// GetDefaultHttpResponse returns the default http response
func GetDefaultHttpResponse() []byte {
	return []byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 6\r\nConnection: keep-alive\r\nCache-Control: no-cache\r\nCF-Cache-Status: DYNAMIC\r\nServer: cloudflare\r\n\r\nfollow")
}

func GetDefaultHttpHandleFunc() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", "6")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("CF-Cache-Status", "DYNAMIC")
		w.Header().Set("Server", "cloudflare")
		w.Write([]byte("follow"))
	})
}
