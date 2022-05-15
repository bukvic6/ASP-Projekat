package main

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	router := mux.NewRouter()
	router.StrictSlash(true)

	server := service{
		data:  map[string][]*Config{},
		data1: map[string][]*Group{},
	}
	//	groupServer := groupService{
	//		data: map[string]*[]Group{},
	//	}

	router.HandleFunc("/config/", server.createPostHandler).Methods("POST")
	router.HandleFunc("/configGroups/", server.createGroupHandler).Methods("POST")
	router.HandleFunc("/config/{version}/{id}/", server.delPostHandler).Methods("DELETE")
	router.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
	router.HandleFunc("/configGroups/", server.getAllGroupHandler).Methods("GET")
	router.HandleFunc("/configGroups/{version}/{id}/", server.getGroupHandler).Methods("GET")
	router.HandleFunc("/config/{version}/{id}/", server.getConfigHandler).Methods("GET")
	//router.HandleFunc("/configGroup/{version}/{id}/", server.delPostGroupHandler).Methods("DELETE")
	router.HandleFunc("/configGroup/{version}/{id}", server.createPutHandler).Methods("PUT")

	srv := &http.Server{Addr: "0.0.0.0:8000", Handler: router}
	go func() {
		log.Println("Server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	<-quit

	log.Println("service shutting down ...")

	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
