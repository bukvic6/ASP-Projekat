package main

import (
	cs "Ali/configstore"
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

	store, err := cs.New()
	if err != nil {
		log.Fatal(err)
	}
	server := configServer{
		store: store,
	}
	router.HandleFunc("/config/", server.createPostHandler).Methods("POST")

	router.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
	/*	router.HandleFunc("/config", server.createPostHandler).Methods("POST")
		router.HandleFunc("/configGroups", server.createGroupHandler).Methods("POST")
		router.HandleFunc("/config/{id}/{version}", server.delPostHandler).Methods("DELETE")
		router.HandleFunc("/configs/", server.getAllHandler).Methods("GET")
		router.HandleFunc("/configGroups", server.getAllGroupHandler).Methods("GET")
		router.HandleFunc("/configGroups/{version}/{id}", server.getGroupHandler).Methods("GET")
		router.HandleFunc("/config/{id}/{version}", server.getConfigHandler).Methods("GET")
		router.HandleFunc("/configGroups/{version}/{id}", server.delPostGroupHandler).Methods("DELETE")
		router.HandleFunc("/configGroups/{version}/{id}", server.createPutHandler).Methods("PUT")
		router.HandleFunc("/configGroups/{id}", server.createGroupVersionHandler).Methods("PUT")
		router.HandleFunc("/config/{id}", server.createConfigVersionHandler).Methods("PUT")*/

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
