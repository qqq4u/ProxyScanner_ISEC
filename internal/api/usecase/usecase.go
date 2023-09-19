package usecase

import (
	repo "ISEC/internal/api/repository"
	"ISEC/internal/middleware/proxy"
	"ISEC/internal/middleware/utils"
	"ISEC/internal/models"
	"context"
	"net/http"
)

type ProxyUsecase struct {
	repo  repo.ProxyRepo
	proxy *proxy.Proxy
}

func NewProxyUsecase(repo repo.ProxyRepo, proxy *proxy.Proxy) *ProxyUsecase {
	return &ProxyUsecase{
		repo:  repo,
		proxy: proxy,
	}
}
func (uc *ProxyUsecase) RepeatRequest(ctx context.Context, w http.ResponseWriter, requestID int) error {
	request, err := uc.GetRequest(ctx, requestID)
	if err != nil {
		return err
	}

	HTTPrequest, err := utils.GetRequestFromStruct(request)
	if err != nil {
		return err
	}

	uc.proxy.ServeHTTP(w, &HTTPrequest)

	return nil
}
func (uc *ProxyUsecase) GetRequest(ctx context.Context, requsetID int) (models.Request, error) {
	return uc.repo.GetRequest(ctx, requsetID)
}
func (uc *ProxyUsecase) GetAllRequests() ([]models.Request, error) {
	return uc.repo.GetAllRequests()
}

func (uc *ProxyUsecase) ScanRequest(ctx context.Context, w http.ResponseWriter, requestID int) (string, error) {
	request, err := uc.GetRequest(ctx, requestID)
	if err != nil {
		return "", err
	}
	_, err = utils.GetRequestFromStruct(request)
	if err != nil {
		return "", err
	}

	return "", nil
	//IMPLEMENT ME
}
