package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"log"
	"net"
	"todoserver/infrastructures"
	"todoserver/models"
	pb "todoserver/pkg/api"
	"todoserver/repositories"
)

type server struct {
	pb.UnimplementedTodoServer
}

func (ser *server) Add(ctx context.Context, in *pb.AddToDoRequest) (*pb.AddToDoResponse, error) {
	log.Printf("Received: %v", in.GetTitle())

	toDoItem := models.ToDo{
		Title:  in.GetTitle(),
		IsDone: in.GetDone(),
		Order:  in.GetOrder(),
	}

	lastId, err := Repo.CreateItem(toDoItem)
	if err != nil {
		log.Fatal("gRPC server Add failed ", err)
	}

	return &pb.AddToDoResponse{Id: int32(lastId)}, nil
}

func (serv *server) List(in *pb.EmptyMessage, stream pb.Todo_ListServer) error {
	items, _ := Repo.GetAllItems()
	for _, item := range items {

		itemGRPS := &pb.ToDoItem{
			Id:    item.Id,
			Title: item.Title,
			Order: item.Order,
			Done:  item.IsDone,
		}

		if err := stream.Send(itemGRPS); err != nil {
			log.Fatal("failed to send item ", itemGRPS)

			return err
		}

	}
	return nil
}

func (serv *server) Delete(ctx context.Context, in *pb.IdRequest) (*pb.ResultBoolResponse, error) {
	_, err := Repo.DeleteItemById(in.GetId())
	if err != nil {
		fmt.Println("couldn't delete item ", in)

		return &pb.ResultBoolResponse{Success: false}, err
	}

	return &pb.ResultBoolResponse{Success: true}, nil
}

const (
	ServerAddr = "localhost:8091"
)

var (
	MySQLConfig = mysql.Config{
		User:   "root",
		Passwd: "diedie11",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "todo",
	}
	Repo *repositories.ToDoRepository
	Conn *sql.DB
)

func main() {
	initDB()
	initRepo()

	lis, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterTodoServer(s, &server{})

	fmt.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initDB() {
	var err error

	Conn, err = sql.Open("mysql", MySQLConfig.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := Conn.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")
}

func initRepo() {
	mysqlHandler := &infrastructures.MySQLHandler{}
	mysqlHandler.Conn = Conn

	Repo = &repositories.ToDoRepository{mysqlHandler}
}
