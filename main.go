package main

import (
	"context"
	"errors"
	_ "expvar"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/dcafferty/go-kit-ex1/pb"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics/expvar"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/juju/ratelimit"
	grpc "google.golang.org/grpc"
)

func main() {
	logger := log.NewLogfmtLogger(os.Stdout)

	var svc pb.CounterClient
	c := &grpcServer{}
	svc = Endpoints{
		CounterEndpoint: makeAddEndpoint(c)}

	limit := ratelimit.NewBucket(2*time.Second, 1)

	requestCount := expvar.NewCounter("request.count")

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
		handler := &grpcServer{
			grpctransport.NewServer(
				svc,
				DecodeGRPCAddRequest,
				EncodeGRPCAddResponse,
			),
		}
		gRPCServer := grpc.NewServer()
		pb.RegisterCounterServer(gRPCServer, svc)
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

func (e Endpoints) Add(ctx context.Context, in *pb.AddRequest, opts ...grpc.CallOption) (*pb.AddResponse, error) {
	req := addRequest{int(in.Number)}

	resp, err := e.CounterEndpoint(ctx, req)
	if err != nil {
		return nil, err
	}
	addResponse := resp.(pb.AddResponse)
	if addResponse.Err != "" {
		return &pb.AddResponse{int32(0), addResponse.Err}, errors.New(addResponse.Err)
	}
	return &addResponse, nil
}
