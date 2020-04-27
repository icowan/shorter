/**
 * @Time : 27/04/2020 11:34 AM
 * @Author : solacowa@gmail.com
 * @File : handler
 * @Software: GoLand
 */

package grpc

import (
	"context"
	"github.com/icowan/shorter/pkg/service"

	kitgrpc "github.com/go-kit/kit/transport/grpc"
	endpoint2 "github.com/icowan/shorter/pkg/endpoint"
	"github.com/icowan/shorter/pkg/grpc/pb"
)

type grpcServer struct {
	get  kitgrpc.Handler
	post kitgrpc.Handler
}

func (g *grpcServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.ServiceResponse, error) {
	_, rep, err := g.get.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceResponse), nil
}

func (g *grpcServer) Post(ctx context.Context, req *pb.PostRequest) (*pb.ServiceResponse, error) {
	_, rep, err := g.post.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ServiceResponse), nil
}

func MakeGRPCHandler(eps endpoint2.Endpoints, opts map[string][]kitgrpc.ServerOption) pb.ShorterServer {
	return &grpcServer{
		get: kitgrpc.NewServer(
			eps.GetEndpoint,
			decodeGetRequest,
			encodeResponse,
			opts["Get"]...,
		),
		post: kitgrpc.NewServer(
			eps.PostEndpoint,
			decodePostRequest,
			encodePostResponse,
			opts["Post"]...,
		),
	}
}

func decodePostRequest(ctx context.Context, r interface{}) (interface{}, error) {
	return endpoint2.PostRequest{URL: r.(*pb.PostRequest).GetDomain()}, nil
}

func decodeGetRequest(_ context.Context, r interface{}) (interface{}, error) {
	return endpoint2.GetRequest{Code: r.(*pb.GetRequest).Code}, nil
}

func encodeResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint2.GetResponse)
	var (
		errStr       string
		err          error
		dataResponse *pb.ResponseData
	)
	if resp.Err != nil {
		errStr = resp.Err.Error()
		err = resp.Err
	}
	if resp.Data != nil {
		data := resp.Data.(*service.Redirect)
		dataResponse = &pb.ResponseData{
			Url:      data.URL,
			Code:     data.Code,
			ShortUri: data.URL,
		}
	}

	return &pb.ServiceResponse{
		Data: dataResponse,
		Err:  errStr,
	}, err
}

func encodePostResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint2.PostResponse)
	var (
		errStr string
		err    error
	)
	if resp.Err != nil {
		errStr = resp.Err.Error()
		err = resp.Err
	}
	return &pb.ServiceResponse{
		Data: &pb.ResponseData{
			Url:      resp.Data.Url,
			Code:     resp.Data.Code,
			ShortUri: resp.Data.ShortUri,
		},
		Err: errStr,
	}, err
}
