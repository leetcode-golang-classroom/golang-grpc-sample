# golang-grpc-sample

This repository is for how to use gRPC as protocol in golang

## how to use?
### install tool
* grpc cli
```shell
go get google.golang.org/grpc
```
* protobuf
```shell
go get github.com/golang/protobuf
```
* protoc-gen-go
```shell
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## define proto file

```proto
syntax = "proto3";

option go_package = "golang-grpc-sample/proto";

package personservice;

service PersonService {
  rpc Create(CreatePersonRequest) returns (PersonProfileResponse);
  rpc Read(SinglePersonRequest) returns (PersonProfileResponse);
  rpc Update(UpdatePersonRequest) returns (SuccessResponse);
  rpc Delete(SinglePersonRequest) returns (SuccessResponse);
}

message CreatePersonRequest {
  string name = 1;
  string email = 2;
  string phoneNumber = 3;
}
message SinglePersonRequest {
  int32 id = 1;
}

message UpdatePersonRequest {
  int32 id = 1;
  string name = 2;
  string email = 3;
  string phoneNumber = 4;
}

message PersonProfileResponse {
  int32 id = 1;
  string name = 2;
  string email = 3;
  string phoneNumber = 4;
}

message SuccessResponse {
  string response = 1;
}
```

## generate client code

```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/*.proto
```

## implement server with rpc interface

```golang
package main

import (
	"context"
	"errors"
	"log"
	"net"

	pb "github.com/leetcode-golang-classroom/golang-grpc-sample/proto"
	"google.golang.org/grpc"
)

type Person struct {
	ID          int32
	Name        string
	Email       string
	PhoneNumber string
}

var nextID int32 = 1
var persons = make(map[int32]Person)

type server struct {
	pb.UnimplementedPersonServiceServer
}

func (s *server) Create(ctx context.Context, in *pb.CreatePersonRequest) (*pb.PersonProfileResponse, error) {
	person := Person{Name: in.GetName(), Email: in.GetEmail(), PhoneNumber: in.GetPhoneNumber()}
	if person.Email == "" || person.Name == "" || person.PhoneNumber == "" {
		return &pb.PersonProfileResponse{}, errors.New("fields missing")
	}
	person.ID = nextID
	persons[person.ID] = person
	nextID = nextID + 1
	return &pb.PersonProfileResponse{
		Id:          person.ID,
		Email:       person.Email,
		Name:        person.Name,
		PhoneNumber: person.PhoneNumber,
	}, nil
}

func (s *server) Read(ctx context.Context, in *pb.SinglePersonRequest) (*pb.PersonProfileResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.PersonProfileResponse{}, errors.New("not found")
	}
	return &pb.PersonProfileResponse{
		Id: person.ID, Name: person.Name, Email: person.Email,
		PhoneNumber: person.PhoneNumber,
	}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdatePersonRequest) (*pb.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.SuccessResponse{Response: "Not found!"}, errors.New("not found")
	}
	person.Name = in.GetName()
	person.Email = in.GetEmail()
	person.PhoneNumber = in.GetPhoneNumber()
	if person.Email == "" || person.Name == "" || person.PhoneNumber == "" {
		return &pb.SuccessResponse{Response: "fields missing!"},
			errors.New("fields missing")
	}
	persons[person.ID] = person
	return &pb.SuccessResponse{Response: "Done!"}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.SinglePersonRequest) (*pb.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.SuccessResponse{Response: "Not found!"}, errors.New("not found")
	}
	delete(persons, id)
	return &pb.SuccessResponse{Response: "Delete!"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	pb.RegisterPersonServiceServer(s, &server{})
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
```
## implement client code with rpc interface

```golang
package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/leetcode-golang-classroom/golang-grpc-sample/proto"
	"google.golang.org/grpc"
)

type Person struct {
	ID          int32
	Name        string
	Email       string
	PhoneNumber string
}

var nextID int32 = 1
var persons = make(map[int32]Person)

type server struct {
	pb.UnimplementedPersonServiceServer
}

func (s *server) Create(ctx context.Context, in *pb.CreatePersonRequest) (*pb.PersonProfileResponse, error) {
	person := Person{Name: in.GetName(), Email: in.GetEmail(), PhoneNumber: in.GetPhoneNumber()}
	if person.Email == "" || person.Name == "" || person.PhoneNumber == "" {
		return &pb.PersonProfileResponse{}, errors.New("fields missing")
	}
	person.ID = nextID
	persons[person.ID] = person
	nextID = nextID + 1
	return &pb.PersonProfileResponse{
		Id:          person.ID,
		Email:       person.Email,
		Name:        person.Name,
		PhoneNumber: person.PhoneNumber,
	}, nil
}

func (s *server) Read(ctx context.Context, in *pb.SinglePersonRequest) (*pb.PersonProfileResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.PersonProfileResponse{}, errors.New("not found")
	}
	return &pb.PersonProfileResponse{
		Id: person.ID, Name: person.Name, Email: person.Email,
		PhoneNumber: person.PhoneNumber,
	}, nil
}

func (s *server) Update(ctx context.Context, in *pb.UpdatePersonRequest) (*pb.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.SuccessResponse{Response: "Not found!"}, errors.New("not found")
	}
	person.Name = in.GetName()
	person.Email = in.GetEmail()
	person.PhoneNumber = in.GetPhoneNumber()
	if person.Email == "" || person.Name == "" || person.PhoneNumber == "" {
		return &pb.SuccessResponse{Response: "fields missing!"},
			errors.New("fields missing")
	}
	persons[person.ID] = person
	return &pb.SuccessResponse{Response: "Done!"}, nil
}

func (s *server) Delete(ctx context.Context, in *pb.SinglePersonRequest) (*pb.SuccessResponse, error) {
	id := in.GetId()
	person := persons[id]
	if person.ID == 0 {
		return &pb.SuccessResponse{Response: "Not found!"}, errors.New("not found")
	}
	delete(persons, id)
	return &pb.SuccessResponse{Response: "Delete!"}, nil
}
func Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", ":8080")
	errCh := make(chan error, 1)
	if err != nil {
		errCh <- fmt.Errorf("failed to listen: %v", err)
	}
	go func() {
		s := grpc.NewServer()
		pb.RegisterPersonServiceServer(s, &server{})
		log.Printf("gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			errCh <- fmt.Errorf("failed to serve: %v", err)
		}
		defer s.GracefulStop()
	}()
	select {
	case err = <-errCh:
		return err
	case <-ctx.Done():
		fmt.Println("server is stopping")
		return lis.Close()
	}
}
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	err := Start(ctx)
	if err != nil {
		log.Println("failed to start app:", err)
	}
}
```

## run server
```shell
make run-server
```
## run client
```shell
make run-client
```