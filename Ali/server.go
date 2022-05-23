package main

import (
	cs "Ali/configstore"
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
)

type configServer struct {
	store *cs.ConfigStore
}

func (cs *configServer) createPostHandler(w http.ResponseWriter, req *http.Request) {
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
	post, err := cs.store.Post(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, post)
}
func (cs *configServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := cs.store.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}
func (cs *configServer) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	allTasks, err := cs.store.GetAllGroups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, allTasks)
}
func (cs *configServer) addConfigVersion(w http.ResponseWriter, req *http.Request) {
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
	id := mux.Vars(req)["id"]
	rt.Id = id
	config, err := cs.store.AddConfigVersion(rt)
	renderJSON(w, config)

}
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := cs.store.GetConf(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, config)
}
func (cs *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := cs.store.Delete(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, config)

}
func (cs *configServer) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	config, err := cs.store.GetConfVersions(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, config)
}
func (cs *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
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
	group, err := cs.store.Group(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}
func (cs *configServer) addConfigGroupVersion(w http.ResponseWriter, req *http.Request) {
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
	id := mux.Vars(req)["id"]
	rt.Id = id
	group, err := cs.store.AddConfigGroupVersion(rt)
	renderJSON(w, group)

}
func (cs *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, err := cs.store.DeleteGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

func (cs *configServer) addConfig(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	_, err := cs.store.DeleteGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	id2 := mux.Vars(req)["id"]
	version2 := mux.Vars(req)["version"]
	rt.Id = id2
	rt.Version = version2

	nova, err := cs.store.Put(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, nova)
}
func (cs *configServer) getGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, err := cs.store.GetGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}
func (cs *configServer) getConfigGroupVersions(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	group, err := cs.store.GetConfGroupVersions(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}
func (cs *configServer) filter(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	group, err := cs.store.FilterGroup(id, version, labels)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
}

/*func (ts *service) createConfigVersionHandler(w http.ResponseWriter, req *http.Request) {
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

	config := ts.data[id]

	config = append(config, rt)
	ts.data[id] = config

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
		}
	}

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

		}
	}
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
		}
	}
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
	for index, v := range task {
		if v.Version == version {
			task = append(task[:index], task[index+1:]...)
			ts.data[id] = task
			break
		}
	}
	renderJSON(w, ts.data)

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

	for index, v := range task {
		if v.Version == version {
			task = append(task[:index], task[index+1:]...)
			gs.data1[id] = task
			break
		}
	}
	renderJSON(w, gs.data1)
*/
//
//for _, v := range task {
//	if v.Version == version {
//		delete(gs.data1, id)
//		renderJSON(w, v)
//	} else {
//		err := errors.New("key not found")
//		http.Error(w, err.Error(), http.StatusNotFound)
//	}
//}
