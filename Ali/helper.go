package main

import (
	cs "Ali/configstore"
	"Ali/tracer"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func decodeBody(r io.Reader) (*cs.Config, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt cs.Config
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeBodyGroups(ctx context.Context, r io.Reader) (*cs.Group, error) {

	span := tracer.StartSpanFromContext(ctx, "decodeBody")
	defer span.Finish()
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt cs.Group
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func createId() string {
	return uuid.New().String()
}
