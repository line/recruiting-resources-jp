package v1

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"

	"todo-example/server/pkg/common"
	"todo-example/server/pkg/db"
)

func PutCORSHeaders(w http.ResponseWriter) {
	w.Header().Set( "Access-Control-Allow-Origin", "*" )
	w.Header().Set( "Access-Control-Allow-Credentials", "true" )
	w.Header().Set( "Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization" )
	w.Header().Set( "Access-Control-Allow-Methods","GET, POST, PUT, DELETE, OPTIONS" )
}

func isOptionMethod(r *http.Request) bool {
	if r.Method == "OPTIONS" {
		return true
	}else{
		return false
	}
}

func GetTodoListHandler(w http.ResponseWriter, r *http.Request) {
	PutCORSHeaders(w)
	if isOptionMethod(r) {
		w.WriteHeader(204)
		return
	}
	logger := r.Context().Value("logger").(*common.Log)
	dbClient := r.Context().Value("dbClient").(db.DBClient)
	vars := mux.Vars(r)
	logger.Debug(log.Fields{"user_id": vars["userID"]}, "Try to get todo from db")
	list, err := dbClient.GetTodoList(
		map[string]interface{}{"user_id": vars["userID"]}, -1, -1)
	if err != nil {
		logger.Error(log.Fields{"error": err.Error()}, "Failed to get todo list")
		common.RespondErr(w, 500, "Failed to get list!")
		return
	}
	resp := TodoListResponse{List: []TodoResponse{}}
	for _, t := range list {
		resp.List = append(resp.List, generateTodoResponse(t))
	}
	common.RespondJSON(w, 200, resp)
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	PutCORSHeaders(w)
	if isOptionMethod(r) {
		w.WriteHeader(204)
		return
	}
	logger := r.Context().Value("logger").(*common.Log)
	dbClient := r.Context().Value("dbClient").(db.DBClient)
	var opt TodoCreateRequest
	json.NewDecoder(r.Body).Decode(&opt)
	defer r.Body.Close()
	vars := mux.Vars(r)
	todo := db.Todo{
		Title:       opt.Title,
		UserID:      vars["userID"],
		Description: opt.Description,
	}
	id, err := dbClient.CreateTodo(todo)
	if err != nil {
		logger.Error(log.Fields{"error": err.Error()}, "Failed to create new todo item")
		common.RespondErr(w, 500, "Failed to create!")
	} else {
		logger.Info(log.Fields{"id": id}, "Created new todo")
		resp := map[string]string{
			"id": id,
		}
		common.RespondJSON(w, 201, resp)
	}
}

func GetTodoHandler(w http.ResponseWriter, r *http.Request) {
	PutCORSHeaders(w)
	if isOptionMethod(r) {
		w.WriteHeader(204)
		return
	}
	logger := r.Context().Value("logger").(*common.Log)
	dbClient := r.Context().Value("dbClient").(db.DBClient)
	vars := mux.Vars(r)
	logger.Debug(log.Fields{"id": vars["id"]}, "Try to get todo from db")
	t, err := dbClient.GetTodo(vars["id"])
	if err != nil {
		logger.Error(log.Fields{"error": err.Error()}, "Failed to get todo")
		common.RespondErr(w, 404, "No such item")
		return
	}
	common.RespondJSON(w, 200, generateTodoResponse(t))
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	PutCORSHeaders(w)
	if isOptionMethod(r) {
		w.WriteHeader(204)
		return
	}
	logger := r.Context().Value("logger").(*common.Log)
	dbClient := r.Context().Value("dbClient").(db.DBClient)
	vars := mux.Vars(r)
	logger.Debug(log.Fields{"id": vars["id"]}, "Try to get todo from db")
	err := dbClient.DeleteTodo(vars["id"])
	if err != nil {
		logger.Error(log.Fields{"error": err.Error()}, "Failed to delete todo")
		common.RespondErr(w, 404, "No such item")
		return
	}
	resp := map[string]string{
		"id": vars["id"],
	}
	common.RespondJSON(w, 200, resp)
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	PutCORSHeaders(w)
	if isOptionMethod(r) {
		w.WriteHeader(204)
		return
	}
	logger := r.Context().Value("logger").(*common.Log)
	dbClient := r.Context().Value("dbClient").(db.DBClient)
	vars := mux.Vars(r)
	var opt TodoUpdateRequest
	json.NewDecoder(r.Body).Decode(&opt)
	defer r.Body.Close()
	logger.Debug(log.Fields{"id": vars["id"]}, "Try to get todo from db")
	t, err := dbClient.GetTodo(vars["id"])
	if err != nil {
		logger.Error(log.Fields{"error": err.Error()}, "Failed to get todo")
		common.RespondErr(w, 404, "No such item")
		return
	}
	if opt.Title != nil {
		t.Title = *opt.Title
	}
	if opt.Description != nil {
		t.Description = *opt.Description
	}
	if opt.Finished != nil {
		t.Finished = *opt.Finished
	}
	err = dbClient.UpdateTodo(t)
	if err != nil {
		logger.Error(log.Fields{"error": err.Error()}, "Failed to update todo item")
		common.RespondErr(w, 500, "Failed to update!")
	} else {
		logger.Info(log.Fields{"id": vars["id"]}, "updated todo")
		resp := map[string]string{
			"id": vars["id"],
		}
		common.RespondJSON(w, 200, resp)
	}
}

func generateTodoResponse(t db.Todo) TodoResponse {
	return TodoResponse{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Finished:    t.Finished,
		CreatedAt:   JsonTime{Time: t.CreatedAt},
		UpdatedAt:   JsonTime{Time: t.UpdatedAt},
	}
}
