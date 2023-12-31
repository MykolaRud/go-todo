package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	pb "todoserver/pkg/api"
)

const (
	gRPCServerAddr = "localhost:8091"
	WebServerAddr  = "localhost:8090"
)

var (
	gRPCConnection *grpc.ClientConn
	gRPCClient     pb.TodoClient
)

func main() {
	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/add", addHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/delete/", deleteHandler)
	http.HandleFunc("/toggle/", toggleHandler)

	fmt.Println("gRPC connection starting")
	initGRPCClient()

	fmt.Println("server starting at ", WebServerAddr)

	err := http.ListenAndServe(WebServerAddr, nil)
	if err != nil {
		log.Fatal("error starting server ", err)
	}

	fmt.Println("server exited")

	gRPCConnection.Close()
}

func initGRPCClient() {
	var err error
	gRPCConnection, err = grpc.Dial(gRPCServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	gRPCClient = pb.NewTodoClient(gRPCConnection)
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func RenderPage(w http.ResponseWriter, Template string, data interface{}) {
	t, err := template.ParseFiles(Template)
	if err != nil {
		log.Fatal("error loading template ", err)
	}
	err = t.Execute(w, data)
	if err != nil {
		log.Fatal("error executing template ", err)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	var Items []pb.ToDoItem

	stream, err := gRPCClient.List(context.Background(), &pb.EmptyMessage{})
	if err != nil {
		log.Fatal("Error getting list form gRPC")
	}

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.List() error %v", gRPCClient, err)
		}

		Items = append(Items, pb.ToDoItem{
			Id:    item.GetId(),
			Title: item.GetTitle(),
			Done:  item.Done,
			Order: item.Order,
		})
	}

	p := make(map[string]interface{})
	p["Items"] = Items

	RenderPage(w, "web/templates/list.html", p)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	RenderPage(w, "web/templates/addform.html", nil)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	title := r.PostFormValue("title")
	order, _ := strconv.Atoi(r.PostFormValue("order"))
	isDoneInt, _ := strconv.Atoi(r.PostFormValue("is_done"))
	isDone := isDoneInt == 1

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := gRPCClient.Add(ctx, &pb.AddToDoRequest{
		Title: title,
		Order: int64(order),
		Done:  isDone,
	})
	if err != nil {
		log.Fatalf("could not add: %v", err)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("new created id is " + strconv.Itoa(int(response.GetId()))))
	w.Write([]byte("<br/><a href='/list'>list</a>"))

}

func editHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/view/"):]
}
func deleteHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/delete/"):]

	idInt, err := strconv.Atoi(id)
	if err != nil {
		w.Write([]byte("Can't convert " + id + " to str"))
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := gRPCClient.Delete(ctx, &pb.IdRequest{
		Id: int64(idInt),
	})
	if err != nil {
		log.Fatalf("could not delete: %v", err)
	}

	if response.GetSuccess() != true {
		w.Write([]byte("failed to delete " + id))
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte("item deleted " + id))
	w.Write([]byte("<br/><a href='/list'>list</a>"))
}

func toggleHandler(w http.ResponseWriter, r *http.Request) {
	urlSegments := strings.Split(r.URL.Path, "/")

	if len(urlSegments) < 3 {
		log.Fatal("not enough arguments")

		return
	}

	id := urlSegments[2]
	done := urlSegments[3]

	idInt, _ := strconv.Atoi(id)
	doneBool := done == "1"

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := gRPCClient.Done(ctx, &pb.DoneToDoRequest{
		Id:   int64(idInt),
		Done: doneBool,
	})
	if err != nil {
		log.Fatalf("could not delete: %v", err)
	}

	if !response.GetSuccess() {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte("couldn't toggle " + id))
		w.Write([]byte("<br/><a href='/list'>list</a>"))
		return
	}

	http.Redirect(w, r, "/list", http.StatusSeeOther)
}
