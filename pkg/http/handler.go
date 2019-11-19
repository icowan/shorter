/**
 * @Time : 19/11/2019 10:41 AM
 * @Author : solacowa@gmail.com
 * @File : handler
 * @Software: GoLand
 */

package http

import (
	"context"
	"encoding/json"
	"errors"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/icowan/shorter/pkg/endpoint"
	"github.com/icowan/shorter/pkg/service"
	"net/http"
)

var (
	ErrCodeNotFound = errors.New("code is nil")
)

func NewHTTPHandler(endpoints endpoint.Endpoints, options map[string][]kithttp.ServerOption) http.Handler {
	//m := http.NewServeMux()
	r := mux.NewRouter()
	makeGetHandler(r, endpoints, options["Get"])
	makePostHandler(r, endpoints, options["Post"])
	return r
}

func makeGetHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	m.Handle("/{code}", kithttp.NewServer(
		endpoints.GetEndpoint,
		decodeGetRequest,
		encodeGetResponse,
		options...)).Methods(http.MethodGet)
}
func makePostHandler(m *mux.Router, endpoints endpoint.Endpoints, options []kithttp.ServerOption) {
	m.Handle("/", kithttp.NewServer(
		endpoints.PostEndpoint,
		decodePostRequest,
		encodeGetResponse,
		options...)).Methods(http.MethodGet)
}

func decodePostRequest(_ context.Context, r *http.Request) (interface{}, error) {

	return nil, nil
}

func decodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		return nil, ErrCodeNotFound
	}
	req := endpoint.GetRequest{
		Code: code,
	}
	return req, nil
}

// encodeSayResponse is a transport/http.EncodeResponseFunc that encodes
func encodeGetResponse(ctx context.Context, w http.ResponseWriter, response interface{}) (err error) {
	if f, ok := response.(endpoint.Failure); ok && f.Failed() != nil {
		ErrorEncoder(ctx, f.Failed(), w)
		return nil
	}
	resp := response.(*service.Redirect)
	http.Redirect(w, nil, resp.URL, http.StatusFound)
	return
}
func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.WriteHeader(err2code(err))
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}
func ErrorDecoder(r *http.Response) error {
	var w errorWrapper
	if err := json.NewDecoder(r.Body).Decode(&w); err != nil {
		return err
	}
	return errors.New(w.Error)
}

// This is used to set the http status, see an example here :
// https://github.com/go-kit/kit/blob/master/examples/addsvc/pkg/addtransport/http.go#L133
func err2code(err error) int {
	return http.StatusInternalServerError
}

type errorWrapper struct {
	Error string `json:"error"`
}
