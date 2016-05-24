package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "io"
    "io/ioutil"

    "github.com/gorilla/mux"
)

func Index(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Welcome!")
}

func TodoShow(w http.ResponseWriter, r *http.Request) {
    todos := GetAllTodos();

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)

    if err := json.NewEncoder(w).Encode(todos); err != nil {
        panic(err)
    }
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    todoId := vars["todoId"]

    res := GetTodo(todoId)
    if res.Id == -1 {
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        if err := json.NewEncoder(w).Encode(res); err != nil {
            panic(err)
        }
    }
}

func TodoMarkDone(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    todoId := vars["todoId"]

    res := MarkDone(todoId)
    if res.Id == -1 {
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        w.WriteHeader(http.StatusOK)
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")

        if err := json.NewEncoder(w).Encode(res); err != nil {
            panic(err)
        }

    }
}

func TodoDelete(w http.ResponseWriter, r *http.Request) {
    vars:= mux.Vars(r)
    todoId := vars["todoId"]

    res := DeleteTodo(todoId)
    if res == -1 {
        w.WriteHeader(http.StatusInternalServerError)
    } else {
        w.WriteHeader(http.StatusOK)
    }
}

func TodoCreate(w http.ResponseWriter, r *http.Request) {
    var todo Todo
    body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
    if err != nil {
        panic(err)
    }
    if err := r.Body.Close(); err != nil {
        panic(err)
    }
    if err := json.Unmarshal(body, &todo); err != nil {
        w.Header().Set("Content-Type", "application/json; charset=UTF-8")
        w.WriteHeader(422) // unprocessable entity
        if err := json.NewEncoder(w).Encode(err); err != nil {
            panic(err)
        }
    }

    t := InsertTodo(todo)
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    if err := json.NewEncoder(w).Encode(t); err != nil {
        panic(err)
    }
}