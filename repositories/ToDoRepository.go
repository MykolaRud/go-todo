package repositories

import (
	"database/sql"
	"fmt"
	"log"
	"todoserver/interfaces"
	"todoserver/models"
)

type ToDoRepository struct {
	interfaces.IDbHandler
}

func (repository *ToDoRepository) CreateItem(itemData models.ToDo) (int64, error) {
	Result, err := repository.Execute("INSERT INTO todos (title, is_done, todos.order) VALUES (?, ?, ?)",
		itemData.Title,
		itemData.IsDone,
		itemData.Order,
	)

	if err != nil {
		log.Fatal(err)

		return 0, err
	}

	lastInsertId, err := Result.LastInsertId()
	if err != nil {
		log.Fatal(err)

		return 0, err
	}

	return lastInsertId, nil
}

func (repository *ToDoRepository) UpdateItem(itemData models.ToDo) (sql.Result, error) {
	Result, err := repository.Execute("UPDATE todos SET title = ?, is_done = ?, todos.order = ? WHERE id = ?",
		itemData.Title,
		itemData.IsDone,
		itemData.Order,
		itemData.Id,
	)

	if err != nil {
		log.Fatal("UpdateItem failed: ", err)

		return nil, err
	}

	return Result, nil
}

func (repository *ToDoRepository) DeleteItemById(id int64) (sql.Result, error) {
	Result, err := repository.Execute("DELETE FROM todos WHERE id = ?", id)

	if err != nil {
		log.Fatal("DeleteItemById failed: ", err)

		return nil, err
	}

	return Result, err
}

func (repository *ToDoRepository) SetItemDoneById(id int64) (sql.Result, error) {
	Result, err := repository.Execute("UPDATE todos SET is_done = 1 WHERE id = ?",
		id,
	)

	if err != nil {
		log.Fatal("SetItemDoneById failed: ", err)

		return nil, err
	}

	return Result, nil
}

func (repository *ToDoRepository) GetMaxOrder() int64 {
	Row := repository.QueryRow("SELECT max(todos.order) as cnt FROM todos WHERE 1")

	var cnt int64
	if err := Row.Scan(&cnt); err != nil {
		log.Fatal("GetMaxOrder failed: ", err)
	}

	return cnt
}

func (repository *ToDoRepository) GetAllItems() ([]models.ToDo, error) {
	rows, err := repository.Query("SELECT * FROM todos WHERE 1 ORDER BY todos.order, id")
	if err != nil {
		log.Fatal("GetAllItems failed: ", err)
	}

	todos := []models.ToDo{}

	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var todo models.ToDo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.IsDone, &todo.Order); err != nil {
			return nil, fmt.Errorf("GetAllItems failed %v", err)
		}
		todos = append(todos, todo)
	}

	return todos, err
}
