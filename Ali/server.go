package main

import (
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type service struct {
	data  map[string][]*Config `json:"data"`
	data1 map[string][]*Group  `json:"data"`
}

//type groupService struct {
//	data map[string][]*Group
//}

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
	renderJSON(w, ts.data)
}

func (gs *service) createPutHandler(w http.ResponseWriter, req *http.Request) {
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

	val := mux.Vars(req)

	id := val["id"]
	group := gs.data1[id]
	version := val["version"]
	for _, v := range group {
		if v.Version == version {
			v.Config = append(v.Config, *rt)
			renderJSON(w, v.Config)

		} else {
			err := errors.New("key not found")
			http.Error(w, err.Error(), http.StatusNotFound)
		}

	}
	/*	group

		group.Config = append(group.Config, *rt)
		renderJSON(w, group)
	*/
}
func (gs *service) createGroupHandler(w http.ResponseWriter, req *http.Request) {
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
	gs.data1[id] = listConf
	renderJSON(w, gs.data1)
}
func (gs *service) createGroupVersionHandler(w http.ResponseWriter, req *http.Request) {
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
	val := mux.Vars(req)
	id := val["id"]

	group := gs.data1[id]

	rt.Id = id
	group = append(group, rt)
	gs.data1[id] = group
	renderJSON(w, gs.data1)

}

func (gs *service) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := make(map[string][]*Config)
	for k, v := range gs.data {
		allTasks[k] = v
	}

	renderJSON(w, allTasks)
}

func (gs *service) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks := make(map[string][]*Group)
	for k, v := range gs.data1 {
		allTasks[k] = v
	}

	renderJSON(w, allTasks)
}

func (gs *service) getGroupHandler(w http.ResponseWriter, r *http.Request) {
	val := mux.Vars(r)

	id := val["id"]
	//allTasks := make(map[string][]*Group)
	task, ok := gs.data1[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	version := val["version"]
	allTasks := []*Group{}
	for _, v := range task {
		if v.Version == version {
			allTasks = append(allTasks, v)
			renderJSON(w, allTasks)

		} else {
			err := errors.New("key not found")
			http.Error(w, err.Error(), http.StatusNotFound)
		}

	}

	//	version := val["version"]

}

func (gs *service) getConfigHandler(w http.ResponseWriter, r *http.Request) {
	val := mux.Vars(r)

	id := val["id"]
	task, ok := gs.data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	version := val["version"]
	allTasks := []*Config{}
	for _, v := range task {
		if v.Version == version {
			allTasks = append(allTasks, v)
			renderJSON(w, allTasks)

		} else {
			err := errors.New("key not found")
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	} //	version := val["version"]
}

func (ts *service) delPostHandler(w http.ResponseWriter, r *http.Request) {

	val := mux.Vars(r)

	id := val["id"]
	task, ok := ts.data[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	version := val["version"]

	for _, v := range task {
		if v.Version == version {
			delete(ts.data, id)
			renderJSON(w, v)
		} else {
			err := errors.New("key not found")
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
}

func (gs *service) delPostGroupHandler(w http.ResponseWriter, r *http.Request) {
	val := mux.Vars(r)

	id := val["id"]
	task, ok := gs.data1[id]
	if !ok {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	version := val["version"]

	for _, v := range task {
		if v.Version == version {
			delete(gs.data1, id)
			renderJSON(w, v)
		} else {
			err := errors.New("key not found")
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	}
}
