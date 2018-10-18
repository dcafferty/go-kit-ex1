package main

import (
	"context"
	_ "expvar"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/dcafferty/go-kit-ex1/pb"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	grpc "google.golang.org/grpc"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)

	// var svc pb.CounterServer
	c := &countService{}
	var svc endpoint.Endpoint
	svc = makeAddEndpoint(c)

	// requestCount := expvar.NewCounter("request.count")

	svc = loggingMiddlware(logger)(svc)

	errChan := make(chan error)

	// endpoints := Endpoints{
	// 	CounterEndpoint: makeAddEndpoint(&c),
	// }
	//execute grpc server
	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:51251")
		if err != nil {
			errChan <- err
			return
		}
		addHandler := &grpcServer{
			grpctransport.NewServer(
				svc,
				DecodeGRPCAddRequest,
				EncodeGRPCAddResponse,
			),
		}
		gRPCServer := grpc.NewServer()
		pb.RegisterCounterServer(gRPCServer, addHandler)
		errChan <- gRPCServer.Serve(listener)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	select {
	case err := <-errChan:
		panic(err)
		close(errChan)
	}
	close(errChan)
	// http.Handle("/add",
	// 	kth.NewServer(
	// 		svc,
	// 		decodeAddRequest,
	// 		encodeResponse,
	// 		kth.ServerBefore(beforeIDExtractor, beforePATHExtractor),
	// 	),
	// )

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8080"
	// }

	// logger.Log("listening-on", port)
	// if err := http.ListenAndServe(":"+port, nil); err != nil {
	// 	logger.Log("listen.error", err)
	// }
}

// endpoints wrapper
type Endpoints struct {
	CounterEndpoint endpoint.Endpoint
}

func makeAddEndpoint(svc Counter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(addRequest)
		v := svc.Add(req)

		return v, nil
	}
}

func (e Endpoints) Add(ctx context.Context, request interface{}) (result interface{}, err error) {
	// func (e Endpoints) Add(ctx context.Context, in *pb.AddRequest, opts ...grpc.CallOption) (*pb.AddResponse, error) {
	in := request.(addRequest)
	// req := addRequest{int(in.Number)}
	// req := addRequest{int(in.Number)}

	resp, err := e.CounterEndpoint(ctx, &in)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}
	return resp, nil
}
