/**
 * @Time : 27/04/2020 12:00 PM
 * @Author : solacowa@gmail.com
 * @File : grpc
 * @Software: GoLand
 */

package main

import (
	"context"
	"fmt"
	"github.com/icowan/shorter/pkg/grpc/pb"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8082", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		_ = conn.Close()
	}()

	svc := pb.NewShorterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := svc.Post(ctx, &pb.PostRequest{
		Domain: "https://www.baidu.com",
	})
	if err != nil {
		log.Fatalf("could not put: %v", err)
	}

	fmt.Println("Code", r.Data.Code, "ShortUri", r.Data.ShortUri)

	r, err = svc.Get(ctx, &pb.GetRequest{
		Code: r.Data.Code,
	})
	if err != nil {
		log.Fatalf("could not get: %v", err)
	}
	log.Printf("data: %s", r.GetData())
}
