package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/onrik/logrus/filename"
	log "github.com/sirupsen/logrus"

	"todo-example/server/pkg/common"
	"todo-example/server/pkg/config"
	"todo-example/server/pkg/db"
	"todo-example/server/pkg/filter"
)

type App struct {
	router *mux.Router
	port   string
}

// Initialize creates a mux router and
// applies configuration.
func (a *App) Initialize(conf config.Config) {
	var err error
	var logLevel log.Level
	if logLevel, err = log.ParseLevel(conf.Server.Loglevel); err != nil {
		log.Fatal(err)
	}
	log.SetLevel(logLevel)
	log.AddHook(filename.NewHook())
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	dbStr := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		conf.Mysql.Username,
		conf.Mysql.Password,
		conf.Mysql.Endpoint,
		conf.Mysql.Db)
	db.DBInit(dbStr)
	a.router = mux.NewRouter()
	a.port = conf.Server.Port
}

func (a *App) Run() {
	a.router.NotFoundHandler = filter.Wrap(
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			common.RespondErr(w, 404, "no such route!")
		}), filter.Init(), filter.Logger())

	for _, route := range routes {
		handler := filter.Wrap(route.handlerFunc,
			filter.Init(),
			filter.Logger(),
			filter.CheckSchema(),
			filter.CheckToken(),
		)
		a.router.Methods(route.method).Path(route.path).Handler(handler)
	}
        a.router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))
	log.Info("start server on *:" + a.port)
	log.Fatal(http.ListenAndServe(":"+a.port, a.router))
}
