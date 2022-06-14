package main

import (
	cs "Ali/configstore"
	"Ali/tracer"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/net/context"
	"io"
	"mime"
	"net/http"
	"sort"
	"strings"
)

const (
	name = "config_service"
)

type configServer struct {
	store  *cs.ConfigStore
	tracer opentracing.Tracer
	closer io.Closer
}

func NewCOnfigServer() (*configServer, error) {
	store, err := cs.New()
	if err != nil {
		return nil, err
	}

	tracer, closer := tracer.Init(name)
	opentracing.SetGlobalTracer(tracer)
	return &configServer{
		store:  store,
		tracer: tracer,
		closer: closer,
	}, nil
}
func (c *configServer) GetTracer() opentracing.Tracer {
	return c.tracer
}

func (c *configServer) GetCloser() io.Closer {
	return c.closer
}

func (c *configServer) CloseTracer() error {
	return c.closer.Close()
}

func (cs *configServer) createPostHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("createConfigHandler", cs.tracer, req)
	defer span.Finish()

	ctx := tracer.ContextWithSpan(context.Background(), span)
	span.LogFields(
		tracer.LogString("Handler", fmt.Sprintf("Handling greate config at %s\n", req.URL.Path)),
	)
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

	rt, err := decodeBody(ctx, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if reqKey == "" {
		renderJSON(ctx, w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(ctx, reqKey) == true {
		http.Error(w, "OK", http.StatusBadRequest)

		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId(ctx)

	post, err := cs.store.Post(ctx, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, post)
	w.Write([]byte(idempotencyKey))
}
func (cs *configServer) getAllHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("GetAllHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all configs at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	allTasks, err := cs.store.GetAll(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, allTasks)
}
func (cs *configServer) getAllGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getAllGroupHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get all groups at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	allTasks, err := cs.store.GetAllGroups(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, allTasks)
}
func (cs *configServer) addConfigVersion(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("addConfigVersionHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling add config version at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
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
		renderJSON(ctx, w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(ctx, reqKey) == true {
		http.Error(w, "OK", http.StatusBadRequest)

		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId(ctx)
	rt, err := decodeBody(ctx, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := mux.Vars(req)["id"]
	rt.Id = id
	config, err := cs.store.AddConfigVersion(ctx, rt)
	if err != nil {
		http.Error(w, "version already exist", http.StatusBadRequest)
	}
	renderJSON(ctx, w, config)
	w.Write([]byte(idempotencyKey))

}
func (cs *configServer) getConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling getConfig handler at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := cs.store.GetConf(ctx, id, version)
	if err != nil {
		err := errors.New("not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, config)
}
func (cs *configServer) delConfigHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("delConfigHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del Config Handler at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	config, err := cs.store.Delete(ctx, id, version)
	if err != nil {
		err := errors.New("not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, config)

}
func (cs *configServer) getConfigVersionsHandler(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]

	span := tracer.StartSpanFromRequest("GetConfVersionHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get config version handler at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	config, err := cs.store.GetConfVersions(ctx, id)
	if err != nil {
		err := errors.New("not found")

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, config)
}
func (cs *configServer) createGroupHandler(w http.ResponseWriter, req *http.Request) {

	span := tracer.StartSpanFromRequest("greateGroupHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling create Group handler at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
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

	rt, err := decodeBodyGroups(ctx, req.Body)
	if err != nil || rt.Version == "" || rt.Config == nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	if reqKey == "" {
		renderJSON(ctx, w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(ctx, reqKey) == true {
		http.Error(w, "OK", http.StatusBadRequest)

		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId(ctx)
	group, err := cs.store.Group(ctx, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, group)
	w.Write([]byte(idempotencyKey))
}
func (cs *configServer) addConfigGroupVersion(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("addConfigVersion", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling add config version handler at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
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
		renderJSON(ctx, w, "Idempotency-key is missing")
		return
	}
	if cs.store.CheckId(ctx, reqKey) == true {
		http.Error(w, "OK", http.StatusBadRequest)
		return
	}

	idempotencyKey := ""
	idempotencyKey = cs.store.SaveId(ctx)
	rt, err := decodeBodyGroups(ctx, req.Body)
	if err != nil {
		http.Error(w, "incvalid formtat", http.StatusBadRequest)
	}
	id := mux.Vars(req)["id"]
	rt.Id = id
	group, err := cs.store.AddConfigGroupVersion(ctx, rt)
	if err != nil {
		http.Error(w, "version alredy exitst", http.StatusBadRequest)
	}
	renderJSON(ctx, w, group)
	w.Write([]byte(idempotencyKey))

}
func (cs *configServer) delGroupHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("delGroupHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling del group handlere at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, err := cs.store.DeleteGroup(ctx, id, version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, group)
}

func (cs *configServer) addConfig(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("addConfigToGroup", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling add config to group at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	_, err := cs.store.DeleteGroup(ctx, id, version)
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

	rt, err := decodeBodyGroups(ctx, req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id2 := mux.Vars(req)["id"]
	version2 := mux.Vars(req)["version"]
	rt.Id = id2
	rt.Version = version2

	nova, err := cs.store.Put(ctx, rt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, nova)
}
func (cs *configServer) getGroupVersionsHandler(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getGroupVersion", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get group version handler at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	group, ok := cs.store.GetGroup(ctx, id, version)
	if ok != nil {
		err := errors.New("key not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, group)
}
func (cs *configServer) getConfigGroupVersions(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("getConfigGroupVersionsHandler", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling get Config Group Versions at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	group, err := cs.store.GetConfGroupVersions(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	renderJSON(ctx, w, group)
}

func (cs *configServer) filter(w http.ResponseWriter, req *http.Request) {
	span := tracer.StartSpanFromRequest("filterGroupBy version", cs.tracer, req)
	defer span.Finish()
	span.LogFields(
		tracer.LogString("handler", fmt.Sprintf("handling cfilter group at %s\n", req.URL.Path)),
	)
	ctx := tracer.ContextWithSpan(context.Background(), span)
	id := mux.Vars(req)["id"]
	version := mux.Vars(req)["version"]
	labels := mux.Vars(req)["labels"]
	group, err := cs.store.GetGroup(ctx, id, version)
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
				renderJSON(ctx, w, group.Config[i])
			}
		}
	}
}
