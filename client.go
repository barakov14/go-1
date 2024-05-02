package main

import (
	"context"
	"log"

	pb "user.proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewUserServiceClient(conn)

	// Вызов метода AddUser
	user := &pb.User{Id: 1, Name: "John Doe", Email: "john@example.com"}
	addUserResponse, err := c.AddUser(context.Background(), user)
	if err != nil {
		log.Fatalf("could not add user: %v", err)
	}
	log.Printf("User added: %v", addUserResponse)

	// Вызов метода GetUser
	userID := &pb.UserID{Id: 1}
	getUserResponse, err := c.GetUser(context.Background(), userID)
	if err != nil {
		log.Fatalf("could not get user: %v", err)
	}
	log.Printf("User retrieved: %v", getUserResponse)

	// Вызов метода ListUsers
	listUsersStream, err := c.ListUsers(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("could not list users: %v", err)
	}
	for {
		user, err := listUsersStream.Recv()
		if err != nil {
			break
		}
		log.Printf("User: %v", user)
	}
}
