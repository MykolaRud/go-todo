package infrastructures

import (
	"database/sql"
	"fmt"
)

type MySQLHandler struct {
	Conn *sql.DB
}

func (handler *MySQLHandler) Execute(statement string, args ...any) (sql.Result, error) {

	return handler.Conn.Exec(statement, args...)
}

func (handler *MySQLHandler) Query(statement string, args ...any) (*sql.Rows, error) {
	rows, err := handler.Conn.Query(statement, args...)

	if err != nil {
		fmt.Println(err)
		return rows, err
	}

	return rows, nil
}

func (handler *MySQLHandler) QueryRow(statement string, args ...any) *sql.Row {
	row := handler.Conn.QueryRow(statement, args...)

	return row
}

//type SqliteRow struct {
//	Rows *sql.Rows
//}
//
//func (r SqliteRow) Scan(dest ...interface{}) error {
//	err := r.Rows.Scan(dest...)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

//func (r SqliteRow) Next() bool {
//	return r.Rows.Next()
//}
