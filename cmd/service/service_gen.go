/**
 * @Time : 19/11/2019 10:25 AM
 * @Author : solacowa@gmail.com
 * @File : service_gen
 * @Software: GoLand
 */

package service

import (
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/icowan/shorter/pkg/endpoint"
	"github.com/icowan/shorter/pkg/http"
	"github.com/icowan/shorter/pkg/service"
	"github.com/oklog/oklog/pkg/group"
)

func createService(endpoints endpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	initHttpHandler(endpoints, g)
	return g
}

func defaultHttpOptions(logger log.Logger) map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{"Get": {
		kithttp.ServerErrorEncoder(http.ErrorEncoder),
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerBefore(kithttp.PopulateRequestContext),
	},
		"Post": {
			kithttp.ServerErrorEncoder(http.ErrorEncoder),
			kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
			kithttp.ServerBefore(kithttp.PopulateRequestContext),
		}}
	return options
}

func addDefaultServiceMiddleware(logger log.Logger, mw []service.Middleware) []service.Middleware {
	mw = append(mw, service.LoggingMiddleware(logger))
	return mw
}

func addDefaultEndpointMiddleware(logger log.Logger, duration *prometheus.Summary, mw map[string][]kitendpoint.Middleware) map[string][]kitendpoint.Middleware {

	return mw
}
