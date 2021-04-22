package main

import (
	"fmt"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", LogHTTP(http.FileServer(http.Dir("."))))
	fmt.Println("Listening localhost:8080")
	http.ListenAndServe(":8080", nil)
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

type LogEntry map[string]interface{}

func Log(le LogEntry) {
	status, ok := le["Status"].(int)
	if ok && status != 404 {
		fmt.Println(le["Status"], le["RequestURI"])
	} else {
		fmt.Println(le["Status"], le["RequestURI"])
	}

}

func LogHTTP(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		w.Header().Add("Cache-Control", "no-store, max-age=0")
		w.Header().Add("Pragma", "no-cache")
		sw := statusWriter{ResponseWriter: w}
		handler.ServeHTTP(&sw, r)
		duration := time.Now().Sub(start)
		Log(LogEntry{
			"Host":       r.Host,
			"RemoteAddr": r.RemoteAddr,
			"Method":     r.Method,
			"RequestURI": r.RequestURI,
			"Proto":      r.Proto,
			"Status":     sw.status,
			"ContentLen": sw.length,
			"UserAgent":  r.Header.Get("User-Agent"),
			"Duration":   duration,
		})
	}
}
