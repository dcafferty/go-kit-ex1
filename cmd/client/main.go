package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	pb "github.com/dcafferty/go-kit-ex1/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	address = "localhost:51251"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()
	c := pb.NewCounterClient(conn)

	// Get NetDevice in Network Container
	n := &pb.AddRequest{
		Number: int32(3),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 9*time.Second)
	defer cancel()
	r1, err := c.Add(ctx, n)
	s, ok := status.FromError(err)
	if !ok {
		log.Fatal("not ok")
	}

	switch s.Code() {
	case codes.OK:

	case codes.DeadlineExceeded:
		log.Fatalf("DEADLINE & RETRY: err (%s)", err.Error())
	case codes.NotFound:
		log.Fatalf("NOTFOUND: err (%s)", err.Error())
	default:

	}
	b1, _ := json.MarshalIndent(r1, "", "  ")
	fmt.Printf("\n#############\n%v\n#############\n", string(b1))

}
