package app

import (
	"net/http"

	"todo-example/server/pkg/api/v1"
)

type route struct {
	method      string
	path        string
	handlerFunc http.HandlerFunc
}

var routes = []route{
	route{"GET", "/v1/todo/{userID}", v1.GetTodoListHandler},
	route{"POST", "/v1/todo/{userID}", v1.CreateTodoHandler},

	route{"GET", "/v1/todo/{userID}/{id}", v1.GetTodoHandler},
	route{"PUT", "/v1/todo/{userID}/{id}", v1.UpdateTodoHandler},
	route{"DELETE", "/v1/todo/{userID}/{id}", v1.DeleteTodoHandler},
}
