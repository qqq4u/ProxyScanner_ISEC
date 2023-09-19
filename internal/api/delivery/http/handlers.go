package http

import (
	"ISEC/internal/api/usecase"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type ProxyHandler struct {
	proxyUsecase usecase.ProxyUsecase
}

func NewProxyHandler(uc usecase.ProxyUsecase) *ProxyHandler {
	return &ProxyHandler{
		proxyUsecase: uc,
	}
}

func (h *ProxyHandler) GetAllRequests(w http.ResponseWriter, req *http.Request) {
	requests, err := h.proxyUsecase.GetAllRequests()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(requests)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.Fatalln("write request error", err)
		return
	}
}

func (h *ProxyHandler) GetRequest(w http.ResponseWriter, req *http.Request) {
	idString, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "can't get ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "can't convert ID into int", http.StatusBadRequest)
		return
	}

	request, err := h.proxyUsecase.GetRequest(req.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := json.Marshal(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	if _, err = w.Write(result); err != nil {
		log.Fatalln("write request error", err)
		return
	}
}

func (h *ProxyHandler) RepeatRequest(w http.ResponseWriter, req *http.Request) {
	idString, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "can't get ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "can't convert ID into int", http.StatusBadRequest)
		return
	}

	if err := h.proxyUsecase.RepeatRequest(req.Context(), w, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ProxyHandler) ScanRequest(w http.ResponseWriter, req *http.Request) {
	idString, ok := mux.Vars(req)["id"]
	if !ok {
		http.Error(w, "can't get ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "can't convert ID into int", http.StatusBadRequest)
		return
	}

	result, err := h.proxyUsecase.ScanRequest(req.Context(), w, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write([]byte(result))
	if err != nil {
		log.Fatalln("write request error", err)
		return
	}

}
