package main

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
	"todoserver/infrastructures"
	"todoserver/models"
	"todoserver/repositories"
)

func main() {
	fmt.Println("Hello hell")

	conn := initDB()

	mysqlHandler := &infrastructures.MySQLHandler{}
	mysqlHandler.Conn = conn

	Repo := &repositories.ToDoRepository{mysqlHandler}

	item := models.ToDo{
		Title:  "first " + time.Now().Format("01-02-2006 15:04:05"),
		IsDone: false,
		Order:  1 + Repo.GetMaxOrder(),
	}

	lastID, err := Repo.CreateItem(item)
	if err != nil {
		log.Fatal("Could not add item")

		return
	}
	item.Id = lastID
	fmt.Println("added item with ID ", lastID)

	item.IsDone = true
	_, err = Repo.UpdateItem(item)
	if err != nil {
		log.Fatal("Could not update item")

		return
	}
	fmt.Println("updated item with ID ", item.Id)

	//_, err = Repo.DeleteItemById(item.Id - 1)
	//if err != nil {
	//	log.Fatal("Could not delete item")
	//
	//	return
	//}
	//fmt.Println("deleted item with ID ", item.Id-1)

	items, err := Repo.GetAllItems()

	fmt.Println(items)

	fmt.Println("OK")

}

func initDB() *sql.DB {
	var err error
	var db *sql.DB

	cfg := mysql.Config{
		User:   "root",
		Passwd: "diedie11",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "todo",
	}

	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return db
}
