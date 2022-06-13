package main

import (
	cs "Ali/configstore"
	"errors"
	"github.com/gorilla/mux"
	"mime"
	"net/http"
	"sort"
	"strings"
)

type configServer struct {
	store *cs.ConfigStore
}

func (cs *configServer) createPostHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	reqKey := req.Header.Get("idempotency-key")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if reqKey == "" {
		renderJSON(w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(reqKey) == true {
		http.Error(w, "You cannot post same request", http.StatusBadRequest)

		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId()

	post, err := cs.store.Post(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, post)
	w.Write([]byte(idempotencyKey))
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
	reqKey := req.Header.Get("idempotency-key")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	if reqKey == "" {
		renderJSON(w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(reqKey) == true {
		http.Error(w, "You cannot post same request", http.StatusBadRequest)

		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId()
	rt, err := decodeBody(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(req)["id"]
	rt.Id = id
	config, err := cs.store.AddConfigVersion(rt)
	if err != nil {
		http.Error(w, "version already exist", http.StatusBadRequest)
	}
	renderJSON(w, config)
	w.Write([]byte(idempotencyKey))

}
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := cs.store.GetConf(id, version)
	if err != nil {
		err := errors.New("not found")
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
		err := errors.New("not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, config)

}
func (cs *configServer) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	config, err := cs.store.GetConfVersions(id)
	if err != nil {
		err := errors.New("not found")

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, config)
}
func (cs *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	reqKey := req.Header.Get("idempotency-key")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}

	rt, err := decodeBodyGroups(req.Body)
	if err != nil || rt.Version == "" || rt.Config == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if reqKey == "" {
		renderJSON(w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(reqKey) == true {
		http.Error(w, "You cannot post same request", http.StatusBadRequest)

		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId()
	group, err := cs.store.Group(rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(w, group)
	w.Write([]byte(idempotencyKey))
}
func (cs *configServer) addConfigGroupVersion(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	reqKey := req.Header.Get("idempotency-key")

	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		err := errors.New("expect application/json Content-Type")
		http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
		return
	}
	if reqKey == "" {
		renderJSON(w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(reqKey) == true {
		http.Error(w, "You cannot post same request", http.StatusBadRequest)
		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId()
	rt, err := decodeBodyGroups(req.Body)
	if err != nil {
		http.Error(w, "incvalid formtat", http.StatusBadRequest)
	}
	id := mux.Vars(req)["id"]
	rt.Id = id
	group, err := cs.store.AddConfigGroupVersion(rt)
	if err != nil {
		http.Error(w, "version alredy exitst", http.StatusBadRequest)
	}
	renderJSON(w, group)
	w.Write([]byte(idempotencyKey))

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
		err := errors.New("expect application/json Content-Type")
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
	group, ok := cs.store.GetGroup(id, version)
	if ok != nil {
		err := errors.New("key not found")
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
	group, err := cs.store.GetGroup(id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	entries := strings.Split(labels, ",")
	//  https://stackoverflow.com/questions/21362950/getting-a-slice-of-keys-from-a-map
	m := make(map[string]string)
	for _, e := range entries {
		parts := strings.Split(e, ":")
		m[parts[0]] = parts[1]
	}

	for i := 0; i < len(group.Config); i++ {
		entries := group.Config[i].Entries
		if len(m) == len(group.Config[i].Entries) {
			check := false
			key := make([]string, 0, len(entries))
			for k := range entries {
				key = append(key, k)
			}

			sort.Strings(key)
			for _, k := range key {
				i, ok := m[k]
				if ok == false {
					check = true
					break
				} else {
					if i != entries[k] {
						check = true
						break
					}
				}

			}
			if check != true {
				renderJSON(w, group.Config[i])
			}
		}
	}
}
