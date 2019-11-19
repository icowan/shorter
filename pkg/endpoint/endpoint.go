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
	URL string `json:"url"`
}

type GetResponse struct {
	Err  error
	Data interface{}
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
		err = s.Post(ctx, req.URL)
		return GetResponse{Err: err}, err
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
