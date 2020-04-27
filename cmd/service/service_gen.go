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
	"github.com/go-kit/kit/transport"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/icowan/shorter/pkg/endpoint"
	"github.com/icowan/shorter/pkg/http"
	"github.com/icowan/shorter/pkg/service"
	"github.com/oklog/oklog/pkg/group"
	"golang.org/x/time/rate"
	"time"
)

func createService(endpoints endpoint.Endpoints) (g *group.Group) {
	g = &group.Group{}
	initHttpHandler(endpoints, g)
	initGrpcHandler(endpoints, g)
	return g
}

func defaultHttpOptions(logger log.Logger) map[string][]kithttp.ServerOption {
	options := map[string][]kithttp.ServerOption{"Get": {
		kithttp.ServerErrorEncoder(http.ErrorRedirect),
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

func defaultGrpcOptions(logger log.Logger) map[string][]kitgrpc.ServerOption {
	return map[string][]kitgrpc.ServerOption{
		"Get": {
			kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		},
		"Post": {
			kitgrpc.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		},
	}

}

func addDefaultServiceMiddleware(logger log.Logger, mw []service.Middleware) []service.Middleware {
	mw = append(mw, service.LoggingMiddleware(logger))
	return mw
}

func addDefaultEndpointMiddleware(logger log.Logger, mw map[string][]kitendpoint.Middleware) map[string][]kitendpoint.Middleware {
	mw["Post"] = append(mw["Post"],
		endpoint.LoggingMiddleware(logger),
		endpoint.NewTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), *rateBucketNum)),
	)

	mw["Get"] = append(mw["Get"],
		endpoint.LoggingMiddleware(logger),
		endpoint.NewTokenBucketLimitter(rate.NewLimiter(rate.Every(time.Second*1), *rateBucketNum*100)),
	)

	return mw
}
