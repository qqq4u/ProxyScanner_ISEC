package proxy

import (
	repo "ISEC/internal/api/repository"
	"ISEC/internal/middleware/utils"
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
	ProxyRepo *repo.ProxyRepo
}

func ProcessRequest(w http.ResponseWriter, r *http.Request, client http.Client) error {
	res, err := client.Do(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	log.Println(r.RemoteAddr, " ", res.Status)

	writeResponse(w, res)

	return nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.RemoteAddr, " ", r.Method, " ", r.URL)

	httpClient := http.Client{}

	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	deleteHopByHopHeaders(r.Header)
	r.RequestURI = ""

	requestStruct, err := utils.ConvertRequestToStruct(r)
	if err != nil {
		log.Fatal("convert request error:", err)
	}
	if err = p.ProxyRepo.AddRequest(r.Context(), requestStruct); err != nil {
		log.Fatal("convert request error:", err)
	}

	if err := ProcessRequest(w, r, httpClient); err != nil {
		log.Fatal("Error during processing request:", err)
	}
}
