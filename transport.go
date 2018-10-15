package main

import (
	"context"
	"fmt"

	pb "github.com/dcafferty/go-kit-ex1/pb"
	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type addRequest struct {
	V int `json:"v"`
}

func (r addRequest) String() string {
	return fmt.Sprintf("%d", r.V)
}

type addResponse struct {
	V   int `json:"v"`
	Err string
}

func (r addResponse) String() string {
	return fmt.Sprintf("%d", r.V)
}

func makeAddEndpoint(svc Counter) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pb.AddRequest)
		v, err := svc.Add(ctx, &req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

//Encode and Decode Counter Request
func EncodeGRPCAddRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(addRequest)
	return &pb.AddRequest{int32(req.V)}, nil
}

func DecodeGRPCAddRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.AddRequest)
	return addRequest{int(req.Number)}, nil
}

//Encode and Decode Counter Request
func EncodeGRPCAddResponse(_ context.Context, r interface{}) (interface{}, error) {
	res := r.(addResponse)
	return &pb.AddRequest{int32(res.V)}, nil
}

func DecodeGRPCAddResponse(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.AddResponse)
	return addRequest{int(req.Number)}, nil
}

type grpcServer struct {
	Counter grpctransport.Handler
}

func (s *grpcServer) Add(ctx context.Context, r *pb.AddRequest) (*pb.AddResponse, error) {
	_, resp, err := s.Counter.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AddResponse), nil
}

// func decodeAddRequest(ctx context.Context, r *http.Request) (interface{}, error) {
// 	var req addRequest
// 	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
// 		return nil, err
// 	}
// 	return req, nil
// }

// func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
// 	return json.NewEncoder(w).Encode(response)
// }
