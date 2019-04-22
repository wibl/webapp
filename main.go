package main

import (
	"log"
	"os"

	"github.com/wibl/webapp/api"
	"github.com/wibl/webapp/storage"

	"net/http"

	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := log.New(os.Stdout, "Log Message ", log.Ldate|log.Ltime|log.Lshortfile)

	stor, err := initializeStorage()
	if err != nil {
		logger.Fatal(err)
	}

	if err = initRPC(stor); err != nil {
		logger.Fatal(err)
	}
}

func initializeStorage() (storage.Storage, error) {
	stor, err := storage.NewDbStorage("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		return nil, err
	}
	return stor, nil
}

func initRPC(stor storage.Storage) error {
	rpcServ := rpc.NewServer()
	rpcServ.RegisterCodec(json.NewCodec(), "application/json")
	rpcServ.RegisterService(&api.GroupService{Stor: stor}, "GS")
	rpcServ.RegisterService(&api.TemplateService{Stor: stor}, "TS")
	http.Handle("/api", rpcServ)
	return http.ListenAndServe(":8080", nil)
}
