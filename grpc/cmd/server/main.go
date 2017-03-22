package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	pb "github.com/mycodesmells/golang-examples/grpc/proto/service"
)

func main() {
	addr := ":6000"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to initializa TCP listen: %v", err)
	}
	log.Printf("Listening on %s\n", addr)
	defer lis.Close()

	server := grpc.NewServer()
	pb.RegisterSimpleServerServer(server, NewServer())

	server.Serve(lis)
}

type server struct {
	users map[string]pb.User
}

func NewServer() server {
	return server{
		users: make(map[string]pb.User),
	}
}

func (s server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*empty.Empty, error) {
	log.Println("Creating user...")
	user := req.User

	if user.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	if user.Role == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "role cannot be empty")
	}

	s.users[user.Username] = *user

	log.Println("User created!")
	return &empty.Empty{}, nil
}

func (s server) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	log.Println("Getting user!")

	if req.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}

	u, exists := s.users[req.Username]
	if !exists {
		return nil, grpc.Errorf(codes.NotFound, "user not found")
	}

	log.Println("User found!")
	return &u, nil
}

func (s server) GreetUser(ctx context.Context, req *pb.GreetUserRequest) (*pb.GreetUserResponse, error) {
	log.Println("Greeting user...")
	if req.Username == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "username cannot be empty")
	}
	if req.Greeting == "" {
		return nil, grpc.Errorf(codes.InvalidArgument, "greeting cannot be empty")
	}

	user, err := s.GetUser(ctx, &pb.GetUserRequest{Username: req.Username})
	if err != nil {
		return nil, errors.Wrap(err, "failed to find matching user")
	}

	return &pb.GreetUserResponse{
		Greeting: fmt.Sprintf("%s, %s! You are a great %s!", strings.Title(req.Greeting), user.Username, user.Role),
	}, nil
}