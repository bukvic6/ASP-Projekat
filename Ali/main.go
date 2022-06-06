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
	router.HandleFunc("/config/", countCreateConfig(server.createPostHandler)).Methods("POST")
	router.HandleFunc("/configs/", countGetAll(server.getAllHandler)).Methods("GET")
	router.HandleFunc("/configs/{id}", countConfigVersions(server.getConfigVersionsHandler)).Methods("GET")
	router.HandleFunc("/configs/{id}/{version}", countGetConfig(server.getConfigHandler)).Methods("GET")
	router.HandleFunc("/config/{id}", countAddConfigVersion(server.addConfigVersion)).Methods("POST")
	router.HandleFunc("/config/{id}/{version}", countdelConfigVersion(server.delConfigHandler)).Methods("DELETE")
	router.HandleFunc("/group", counteCreateGroup(server.createGroupHandler)).Methods("POST")
	router.HandleFunc("/group/", countegetAllGroup(server.getAllGroupHandler)).Methods("GET")
	router.HandleFunc("/group/{id}", counteAddGroupVersion(server.addConfigGroupVersion)).Methods("POST")
	router.HandleFunc("/group/{id}", counteGetConfigGroupVersions(server.getConfigGroupVersions)).Methods("GET")
	router.HandleFunc("/group/{id}/{version}", counteGetGroupVersion(server.getGroupVersionsHandler)).Methods("GET")
	router.HandleFunc("/group/{id}/{version}/{labels}", server.filter).Methods("GET")
	router.HandleFunc("/group/{id}/{version}", counteDelgroupHits(server.delGroupHandler)).Methods("DELETE")
	router.HandleFunc("/group/{id}/{version}", counteAddConfigToGroup(server.addConfig)).Methods("PUT")
	router.Path("/metrics").Handler(metricsHandler())

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
