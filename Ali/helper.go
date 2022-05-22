package main

import (
	cs "Ali/configstore"
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
func decodeConfig(r io.Reader) (*cs.ConfigG, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var rt cs.ConfigG
	if err := dec.Decode(&rt); err != nil {
		return nil, err
	}
	return &rt, nil
}

func decodeBodyGroups(r io.Reader) (*cs.Group, error) {
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
