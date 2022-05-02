package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type service struct {
	data map[string][]*Config
}
type groupService struct {
	data map[string][]*Group
}

func (ts *service) createPostHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	rt.Entries["id"] = id

	var listConf []*Config

	listConf = append(listConf, rt)
	ts.data[id] = listConf
	renderJSON(w, listConf)
}
func (gs *groupService) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("Expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodyGroups(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := createId()
	rt.Id = id
	var listConf []*Group

	listConf = append(listConf, rt)
	gs.data[id] = listConf
	renderJSON(w, listConf)
}

func (ts *service) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := [][]*Config{}
	for _, v := range ts.data {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}
func (gs *groupService) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := [][]*Group{}
	for _, v := range gs.data {
		allTasks = append(allTasks, v)
	}

	renderJSON(w, allTasks)
}

func (ts *service) delPostHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := ts.data[id]; ok {
		delete(ts.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}
func (gs *groupService) delPostGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if v, ok := gs.data[id]; ok {
		delete(gs.data, id)
		renderJSON(w, v)
	} else {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

//NOTE: pokusajte odraditi prosirenje konfiguracione grupe slicno ovome
//u sustini to je funkcija prvo brisanja pa onda opet dodavanja

//func updateMovie(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "application/json")
//	params := mux.Vars(r)
//	for index, item := range movies {
//		if item.ID == params["id"] {
//			movies = append(movies[:index], movies[index+1:]...)
//			var movie Movie
//			_ = json.NewDecoder(r.Body).Decode(&movie)
//			movie.ID = params["id"]
//			movies = append(movies, movie)
//			json.NewEncoder(w).Encode(movie)
//			return
//		}
//	}
//}
