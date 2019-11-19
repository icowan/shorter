/**
 * @Time : 19/11/2019 10:25 AM
 * @Author : solacowa@gmail.com
 * @File : service
 * @Software: GoLand
 */

package service

import (
	"flag"
	"fmt"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
	"github.com/icowan/shorter/pkg/endpoint"
	svchttp "github.com/icowan/shorter/pkg/http"
	"github.com/icowan/shorter/pkg/repository/mongodb"
	"github.com/icowan/shorter/pkg/service"
	prometheus2 "github.com/prometheus/client_golang/prometheus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/oklog/oklog/pkg/group"
)

var logger log.Logger

var fs = flag.NewFlagSet("hello", flag.ExitOnError)
var httpAddr = fs.String("http-addr", ":8080", "HTTP listen address")
var mongoAddr = fs.String("mongo-addr", "mongodb://localhost:32768", "mongodb uri, default: mongodb://localhost:27017")
var dbDrive = fs.String("db-drive", "mongo", "db drive type, default: mongo")

func Run() {
	if err := fs.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	if addr := os.Getenv("MONGO_DB"); addr != "" {
		mongoAddr = &addr
	}
	if drive := os.Getenv("DB_DRIVE"); drive != "" {
		dbDrive = &drive
	}

	logger = log.NewLogfmtLogger(log.StdlibWriter{})
	logger = log.With(logger, "caller", log.DefaultCaller)
	logger = level.NewFilter(logger, level.AllowAll())

	repo, err := mongodb.NewMongoRepository(*mongoAddr, "redirect", 60)
	if err != nil {
		_ = level.Error(logger).Log("connect", "db", "err", err.Error())
		//return
	}

	svc := service.New(getServiceMiddleware(logger), repo)
	eps := endpoint.New(svc, getEndpointMiddleware(logger))
	g := createService(eps)
	initCancelInterrupt(g)
	_ = logger.Log("exit", g.Run())

}
func initHttpHandler(endpoints endpoint.Endpoints, g *group.Group) {
	options := defaultHttpOptions(logger)

	httpHandler := svchttp.NewHTTPHandler(endpoints, options)
	httpListener, err := net.Listen("tcp", *httpAddr)
	if err != nil {
		_ = level.Error(logger).Log("transport", "HTTP", "during", "Listen", "err", err)
	}
	g.Add(func() error {
		_ = level.Debug(logger).Log("transport", "HTTP", "addr", *httpAddr)
		return http.Serve(httpListener, accessControl(httpHandler, logger, nil))
	}, func(error) {
		_ = httpListener.Close()
	})

}
func getServiceMiddleware(logger log.Logger) (mw []service.Middleware) {
	mw = []service.Middleware{}
	mw = addDefaultServiceMiddleware(logger, mw)

	return
}
func getEndpointMiddleware(logger log.Logger) (mw map[string][]kitendpoint.Middleware) {
	mw = map[string][]kitendpoint.Middleware{}
	duration := prometheus.NewSummaryFrom(prometheus2.SummaryOpts{
		Help:      "Request duration in seconds.",
		Name:      "request_duration_seconds",
		Namespace: "example",
		Subsystem: "hello",
	}, []string{"method", "success"})
	addDefaultEndpointMiddleware(logger, duration, mw)

	return
}

func initCancelInterrupt(g *group.Group) {
	cancelInterrupt := make(chan struct{})
	g.Add(func() error {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		select {
		case sig := <-c:
			return fmt.Errorf("received signal %s", sig)
		case <-cancelInterrupt:
			return nil
		}
	}, func(error) {
		close(cancelInterrupt)
	})
}

func accessControl(h http.Handler, logger log.Logger, headers map[string]string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for key, val := range headers {
			w.Header().Set(key, val)
		}
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Connection", "keep-alive")

		if r.Method == "OPTIONS" {
			return
		}

		_ = logger.Log("remote-addr", r.RemoteAddr, "uri", r.RequestURI, "method", r.Method, "length", r.ContentLength)
		h.ServeHTTP(w, r)
	})
}
