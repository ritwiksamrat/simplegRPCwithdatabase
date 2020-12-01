package main

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	_ "github.com/go-sql-driver/mysql"
	"github.com/ritwiksamrat/newkafka/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		panic(err)
	}
	srv := grpc.NewServer()
	proto.RegisterProducerServiceServer(srv, &server{})
	reflection.Register(srv)

	if e := srv.Serve(listener); e != nil {
		panic(err)
	}
}

func (s *server) Producer(ctx context.Context, request *proto.Request) (*proto.Response, error) {

	result := request.GetUsername()
	db, err := sql.Open("mysql", "root:12345@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Database is Connected Successfully!!")
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO testdb (username) values (?);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(result)

	if err != nil {
		fmt.Print(err.Error())
	}
	return &proto.Response{Result: "Success"}, nil
}
