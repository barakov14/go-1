package main

import (
	"context"
	"log"
	"net"

	pb "user.proto"

	"google.golang.org/grpc"
)

type userServiceServer struct {
	pb.UnimplementedUserServiceServer
}

func (s *userServiceServer) AddUser(ctx context.Context, in *pb.User) (*pb.User, error) {
	return &pb.User{Id: in.Id, Name: in.Name, Email: in.Email}, nil
}

func (s *userServiceServer) GetUser(ctx context.Context, in *pb.UserID) (*pb.User, error) {
	return &pb.User{Id: in.Id, Name: "John Doe", Email: "john@example.com"}, nil
}

func (s *userServiceServer) ListUsers(in *pb.Empty, stream pb.UserService_ListUsersServer) error {
	for i := 1; i <= 5; i++ {
		user := &pb.User{Id: int32(i), Name: "User" + string(i), Email: "user" + string(i) + "@example.com"}
		if err := stream.Send(user); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &userServiceServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
