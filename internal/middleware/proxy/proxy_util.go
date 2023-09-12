package proxy

import (
	"io"
	"log"
	"net/http"
)

// https://datatracker.ietf.org/doc/html/rfc2616#section-13.5.1
var deleteHeaders = []string{
	"Connection",
	"Keep-Alive",
	"Proxy-Authenticate",
	"Proxy-Authorization",
	"Te",
	"Trailers",
	"Transfer-Encoding",
	"Upgrade",
	"Proxy-Connection", //???
}

func copyHeaders(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}

func deleteHopByHopHeaders(header http.Header) {
	for _, h := range deleteHeaders {
		header.Del(h)
	}
}

func writeResponse(w http.ResponseWriter, r *http.Response) {
	deleteHopByHopHeaders(r.Header)
	copyHeaders(w.Header(), r.Header)
	w.WriteHeader(r.StatusCode)
	io.Copy(w, r.Body)
}

type Proxy struct {
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, " ", r.Method, " ", r.URL)

	deleteHopByHopHeaders(r.Header)

	response, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		log.Fatal("ServeHTTP error:", err)
	}
	defer response.Body.Close()

	log.Println(r.RemoteAddr, " ", response.Status)

	writeResponse(w, response)
}
