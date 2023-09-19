package api

import (
	"ISEC/internal/models"
	"context"
	"net/http"
)

type ProxyUsecase interface {
	GetRequest(context.Context, int) (models.Request, error)
	GetAllRequests() ([]models.Request, error)
	RepeatRequest(context.Context, http.ResponseWriter, int) error
	ScanRequest(context.Context, http.ResponseWriter, int) (string, error)
}

type ProxyRepo interface {
	AddRequest(context.Context, models.Request) error
	GetRequest(context.Context, int) (models.Request, error)
	GetAllRequests() ([]models.Request, error)
}
