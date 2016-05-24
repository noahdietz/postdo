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
  // db.Exec(`CREATE TABLE todos (id integer primary key, name text, completed boolean,due time);`)
  checkErr(err)
  fmt.Println("DB connection should be live")
}

func InsertTodo(todo Todo) (added Todo) {
  if db == nil {
    err := errors.New("attempted to insert with no DB connection")
    LogError(err)
    return todo
  } else {
    var newId int

    query := fmt.Sprintf("insert into %s values (%d,'%s',%t,'%s') returning id;", TABLE_NAME, todo.Id, todo.Name, todo.Completed, todo.Due)

    err := db.QueryRow(query).Scan(&newId)
    checkErr(err);
    fmt.Println("last inserted id: ", newId)

    return todo
  }
}

func GetTodo(id string) (todo Todo) {
  query := fmt.Sprintf("select * from %s where id=%s;", TABLE_NAME, id)
  err := db.QueryRow(query).Scan(&todo.Id, &todo.Name, &todo.Completed, &todo.Due)
  if err != nil {
    fmt.Println(err)
    todo.Id = -1
  }

  return todo
}

func GetAllTodos() (todos Todos) {
  var temp Todo

  query := fmt.Sprintf("select * from %s;", TABLE_NAME)
  rows, err := db.Query(query)
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