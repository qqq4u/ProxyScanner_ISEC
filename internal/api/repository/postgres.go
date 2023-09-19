package repo

import (
	"ISEC/internal/models"
	"context"
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

const (
	AddRequest     = `INSERT INTO "requests"(scheme, host, path, method, headers, body, params, cookies) VALUES($1,$2,$3,$4,$5,$6,$7,$8);`
	GetRequestById = `SELECT id, scheme, host, path, method, headers, body, params, cookies FROM "requests" where id = $1`
	GetAllRequests = `SELECT id, scheme, host, path, method, headers, body, params, cookies FROM "requests"`
)

type ProxyRepo struct {
	db *sql.DB
}

func NewProxyRepo(db *sql.DB) *ProxyRepo {
	return &ProxyRepo{
		db: db,
	}
}

func (r *ProxyRepo) AddRequest(ctx context.Context, requestData models.Request) error {
	row := r.db.QueryRowContext(ctx, AddRequest, requestData.Scheme, requestData.Host, requestData.Path, requestData.Method, requestData.Headers, requestData.Body, requestData.Params, requestData.Cookies)
	if err := row.Scan(); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	return nil
}
func (r *ProxyRepo) GetRequest(ctx context.Context, requstID int) (models.Request, error) {
	req := models.Request{}
	err := r.db.QueryRowContext(ctx, GetRequestById, requstID).
		Scan(&req.Id, &req.Scheme, &req.Host, &req.Path, &req.Method, &req.Headers, &req.Body, &req.Params, &req.Cookies)
	if err == sql.ErrNoRows {
		return models.Request{}, errors.New("can't find request with that id")
	}
	if err != nil {
		return models.Request{}, err
	}
	return req, nil
}

func (r *ProxyRepo) GetAllRequests() ([]models.Request, error) {
	requests := make([]models.Request, 0, 0)
	rows, err := r.db.Query(GetAllRequests)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		request := models.Request{}
		err := rows.Scan(&request.Id, &request.Scheme, &request.Host, &request.Path, &request.Method, &request.Headers, &request.Body, &request.Params, &request.Cookies)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}
