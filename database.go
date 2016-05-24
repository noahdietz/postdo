package main

import (
    "fmt"
    "errors"
    "database/sql"
    _ "github.com/lib/pq"
)

const (
    DB_USER = "todosDB"
    DB_NAME = "todosDB"
    TABLE_NAME = "todos"
)

var db *sql.DB

func InitDb() {
  var err error

  dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable", DB_USER, DB_NAME)
  db, err = sql.Open("postgres", dbinfo)
  checkErr(err)
  fmt.Println("DB connection should be live")
}

func InsertTodo(todo Todo) (added Todo) {
  if db == nil {
    err := errors.New("attempted to insert with no DB connection")
    LogError(err)

    todo.Id = -1
    return todo
  } else {
    err := db.QueryRow("insert into todos values ($1, $2, $3, $4) returning id;",
      todo.Id, todo.Name, todo.Completed, todo.Due).Scan(&todo.Id)
    if err != nil {
      todo.Id = -1
      LogError(err)
    }

    fmt.Println("last inserted id: ", todo.Id)

    return todo
  }
}

func MarkDone(id string) (todo Todo) {
  err := db.QueryRow("update todos set completed=true where id = $1 returning *;", id).Scan(
    &todo.Id, &todo.Name, &todo.Completed, &todo.Due)

  if err != nil {
    todo.Id = -1
    LogError(err)
  }

  return todo
}

func DeleteTodo(id string) (deleted int) {
  err := db.QueryRow("delete from todos where id = $1 returning id;", id).Scan(&deleted)
  if err != nil {
    LogError(err)
    return -1;
  }

  return deleted
}

func GetTodo(id string) (todo Todo) {
  err := db.QueryRow("select * from todos where id = $1;", id).Scan(&todo.Id, &todo.Name, &todo.Completed, &todo.Due)
  if err != nil {
    fmt.Println(err)
    todo.Id = -1
  }

  return todo
}

func GetAllTodos() (todos Todos) {
  var temp Todo

  rows, err := db.Query("select * from todos;")
  checkErr(err)

  for rows.Next() {
    rows.Scan(&temp.Id, &temp.Name, &temp.Completed, &temp.Due)
    todos = append(todos, temp)
  }

  return todos
}

func checkErr(err error) {
    if err != nil {
        LogError(err)
        panic(err)
    }
}