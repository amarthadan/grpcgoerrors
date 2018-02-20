package main

import (
	"io"
	"log"
	"net"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/misenko/grpcgoerrors/namesandcolors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

type server struct{}

func (s *server) ListNames(in *google_protobuf.Empty, stream pb.NamesAndColors_ListNamesServer) error {
	names := [5]string{"Jim", "Bob", "Christopher", "Nick", "Grant"}
	for _, name := range names {
		if err := stream.Send(&pb.Name{Name: name}); err != nil {
			return err
		}
	}

	return nil
}

func (s *server) TellColors(stream pb.NamesAndColors_TellColorsServer) error {
	for {
		color, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&google_protobuf.Empty{})
		}
		if err != nil {
			return err
		}
		if len(color.GetColor()) >= 10 {
			return status.Errorf(codes.InvalidArgument, "Length of `Color` cannot be more than 10 characters")
		}
		log.Println(color)
	}
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNamesAndColorsServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
