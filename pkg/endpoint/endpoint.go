/**
 * @Time : 19/11/2019 10:26 AM
 * @Author : solacowa@gmail.com
 * @File : endpoint
 * @Software: GoLand
 */

package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/icowan/shorter/pkg/service"
	"time"
)

type Endpoints struct {
	GetEndpoint  endpoint.Endpoint
	PostEndpoint endpoint.Endpoint
}

func New(s service.Service, mdw map[string][]endpoint.Middleware) Endpoints {
	eps := Endpoints{
		GetEndpoint:  MakeGetEndpoint(s),
		PostEndpoint: MakePostEndpoint(s),
	}
	for _, m := range mdw["Get"] {
		eps.GetEndpoint = m(eps.GetEndpoint)
	}
	for _, m := range mdw["Post"] {
		eps.PostEndpoint = m(eps.PostEndpoint)
	}
	return eps
}

type GetRequest struct {
	Code string
}

type PostRequest struct {
	URL string `json:"url" validate:"empty=false & format=url"`
}

type GetResponse struct {
	Err  error       `json:"err"`
	Data interface{} `json:"data"`
}

type dataResponse struct {
	Url       string    `json:"url"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"created_at"`
	ShortUri  string    `json:"short_uri"`
}

type PostResponse struct {
	Err  error        `json:"err"`
	Data dataResponse `json:"data"`
}

func MakeGetEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(GetRequest)
		redirect, err := s.Get(ctx, req.Code)
		return GetResponse{Err: err, Data: redirect}, err
	}
}

func MakePostEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(PostRequest)
		res, err := s.Post(ctx, req.URL)
		resp := dataResponse{}
		if err == nil && res != nil {
			resp.Code = res.Code
			resp.CreatedAt = res.CreatedAt
			resp.Url = req.URL
			resp.ShortUri = res.URL
		}
		return PostResponse{Err: err, Data: resp}, err
	}
}

func (r GetResponse) Failed() error {
	return r.Err
}

type Failure interface {
	Failed() error
}

func (e Endpoints) Get(ctx context.Context, code string) (rs interface{}, err error) {
	request := GetRequest{Code: code}
	response, err := e.GetEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetResponse).Data, response.(GetResponse).Err
}

func (e Endpoints) Post(ctx context.Context, uri string) (rs interface{}, err error) {
	request := PostRequest{URL: uri}
	response, err := e.PostEndpoint(ctx, request)
	if err != nil {
		return
	}
	return response.(GetResponse).Data, response.(GetResponse).Err
}
