package checks

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Edge-Center/edgecenteredgemon-go/edgecenter"
)

type Service[Req any, Resp any] interface {
	Create(ctx context.Context, req *Req) (*CreateResponse, error)
	Get(ctx context.Context, checkID int) (*Resp, error)
	Update(ctx context.Context, checkID int, req *Req) error
	Delete(ctx context.Context, checkID int) error
}

type service[Req any, Resp any] struct {
	r         edgecenter.Requester
	checkType string
}

func NewService[Req any, Resp any](r edgecenter.Requester, checkType string) Service[Req, Resp] {
	return &service[Req, Resp]{
		r:         r,
		checkType: checkType,
	}
}

func (s *service[Req, Resp]) Create(ctx context.Context, req *Req) (*CreateResponse, error) {
	var resp CreateResponse
	u := url.URL{Path: fmt.Sprintf("/rmon/check/%s", s.checkType)}

	if err := s.r.Request(ctx, http.MethodPost, u.String(), req, &resp); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resp, nil
}

func (s *service[Req, Resp]) Get(ctx context.Context, checkID int) (*Resp, error) {
	var resp Resp
	u := url.URL{Path: fmt.Sprintf("/rmon/check/%s/%d", s.checkType, checkID)}

	if err := s.r.Request(ctx, http.MethodGet, u.String(), nil, &resp); err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}

	return &resp, nil
}

func (s *service[Req, Resp]) Update(ctx context.Context, checkID int, req *Req) error {
	u := url.URL{Path: fmt.Sprintf("/rmon/check/%s/%d", s.checkType, checkID)}

	if err := s.r.Request(ctx, http.MethodPut, u.String(), req, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}

func (s *service[Req, Resp]) Delete(ctx context.Context, checkID int) error {
	u := url.URL{Path: fmt.Sprintf("/rmon/check/%s/%d", s.checkType, checkID)}

	if err := s.r.Request(ctx, http.MethodDelete, u.String(), nil, nil); err != nil {
		return fmt.Errorf("request: %w", err)
	}

	return nil
}
