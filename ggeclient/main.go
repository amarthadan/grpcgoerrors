package main

import (
	"context"
	"io"
	"log"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/misenko/grpcgoerrors/namesandcolors"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewNamesAndColorsClient(conn)

	streamIn, err := c.ListNames(context.Background(), &google_protobuf.Empty{})
	if err != nil {
		log.Fatalf("%v.ListNames(_) = _, %v", c, err)
	}
	for {
		name, err := streamIn.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListNames(_) = _, %v", c, err)
		}
		log.Println(name)
	}

	streamOut, err := c.TellColors(context.Background())
	if err != nil {
		log.Fatalf("%v.TellColors(_) = _, %v", c, err)
	}
	colors := [5]string{"blue", "white", "yellow", "AuroMetalSaurus", "green"}
	for _, color := range colors {
		if err := streamOut.Send(&pb.Color{Color: color}); err != nil {
			log.Fatalf("%v.Send(%v) = %v", streamOut, color, err)
		}
	}
	_, err = streamOut.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", streamOut, err, nil)
	}
}
