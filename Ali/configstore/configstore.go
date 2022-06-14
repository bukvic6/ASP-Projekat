package configstore

import (
	tracer "Ali/tracer"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
	"golang.org/x/net/context"
	"os"
	"sort"
)

type ConfigStore struct {
	cli *api.Client
}

func New() (*ConfigStore, error) {
	db := os.Getenv("DB")
	dbport := os.Getenv("DBPORT")

	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%s", db, dbport)
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &ConfigStore{
		cli: client,
	}, nil
}

func (cs *ConfigStore) Post(ctx context.Context, config *Config) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateConfig")
	defer span.Finish()
	kv := cs.cli.KV()

	sid, rid := generateKey(config.Version)
	config.Id = rid

	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return config, nil
}
func (cs *ConfigStore) CheckId(ctx context.Context, reqId string) bool {
	span := tracer.StartSpanFromContext(ctx, "findId")
	defer span.Finish()
	kv := cs.cli.KV()
	k, _, err := kv.Get(reqId, nil)
	if err != nil || k == nil {
		return false
	}

	return true

}

func (cs *ConfigStore) SaveId(ctx context.Context) string {
	span := tracer.StartSpanFromContext(ctx, "SaveRequestId")
	defer span.Finish()
	kv := cs.cli.KV()
	idempotencyId := uuid.New().String()
	p := &api.KVPair{Key: idempotencyId, Value: nil}
	_, err := kv.Put(p, nil)
	if err != nil {
		return "error"
	}
	return idempotencyId
}

func (cs *ConfigStore) GetAll(ctx context.Context) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "Get all")
	defer span.Finish()
	kv := cs.cli.KV()
	data, _, err := kv.List(all, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Config{}
	for _, pair := range data {
		post := &Config{}
		err = json.Unmarshal(pair.Value, post)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}
func (cs *ConfigStore) GetAllGroups(ctx context.Context) ([]*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "Get groups")
	defer span.Finish()
	kv := cs.cli.KV()
	data, _, err := kv.List(allG, nil)
	if err != nil {
		return nil, err
	}

	posts := []*Group{}
	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		posts = append(posts, group)
	}

	return posts, nil
}
func (cs *ConfigStore) AddConfigVersion(ctx context.Context, config *Config) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "AddConfigVersion")
	defer span.Finish()
	kv := cs.cli.KV()
	ctxKey := tracer.ContextWithSpan(ctx, span)
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	sid := configKeyVersion(ctxKey, config.Id, config.Version)
	_, err = cs.GetConf(ctx, config.Id, config.Version)

	if err == nil {
		return nil, errors.New("version already exists! ")
	}

	p := &api.KVPair{Key: sid, Value: data}

	putKey := tracer.StartSpanFromContext(ctxKey, "kv.put")
	_, err = kv.Put(p, nil)
	if err != nil {
		tracer.LogError(span, err)
		return nil, err
	}
	putKey.Finish()
	return config, nil
}
func (cs *ConfigStore) GetConf(ctx context.Context, id string, version string) (*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConfig")
	defer span.Finish()
	kv := cs.cli.KV()

	sid := configKeyVersion(ctx, id, version)
	pair, _, err := kv.Get(sid, nil)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = json.Unmarshal(pair.Value, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func (cs *ConfigStore) Delete(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteConfig")
	defer span.Finish()
	kv := cs.cli.KV()
	_, err := kv.Delete(configKeyVersion(ctx, id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"deleted": id}, nil
}
func (cs *ConfigStore) GetConfVersions(ctx context.Context, id string) ([]*Config, error) {
	span := tracer.StartSpanFromContext(ctx, "GetConfigVersion")
	defer span.Finish()
	kv := cs.cli.KV()
	sid := configKey(ctx, id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	configList := []*Config{}

	for _, pair := range data {
		config := &Config{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configList = append(configList, config)

	}
	return configList, nil

}
func (cs *ConfigStore) Group(ctx context.Context, group *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "CreateGroup")
	defer span.Finish()
	kv := cs.cli.KV()
	sid, rid := generateGroupKey(group.Version)
	group.Id = rid

	for _, v := range group.Config {
		labels := ""
		keys := make([]string, 0, len(v.Entries))
		for k, _ := range v.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			labels += k + ":" + v.Entries[k] + ","
		}
		labels = labels[:len(labels)-1]
		configKeyGroupVersionlabel(ctx, group.Id, group.Version, labels)

	}

	data, err := json.Marshal(group)
	if err != nil {
		return nil, err
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (cs *ConfigStore) AddConfigGroupVersion(ctx context.Context, group *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "AddVersionGroup")
	defer span.Finish()
	kv := cs.cli.KV()
	data, err := json.Marshal(group)

	sid := configKeyGroupVersion(ctx, group.Id, group.Version)
	_, err = cs.GetGroup(ctx, group.Id, group.Version)

	if err == nil {
		return nil, errors.New("version already exists! ")
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}
func (cs *ConfigStore) DeleteGroup(ctx context.Context, id string, version string) (map[string]string, error) {
	span := tracer.StartSpanFromContext(ctx, "DeleteGroup")
	defer span.Finish()
	kv := cs.cli.KV()
	_, err := kv.Delete(configKeyGroupVersion(ctx, id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"deleted": id}, nil
}
func (cs *ConfigStore) GetGroup(ctx context.Context, id string, version string) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "GetGroup")
	defer span.Finish()
	kv := cs.cli.KV()
	ctxKey := tracer.ContextWithSpan(ctx, span)

	sid := configKeyGroupVersion(ctxKey, id, version)
	getKey := tracer.StartSpanFromContext(ctxKey, "kv.get")

	pair, _, err := kv.Get(sid, nil)
	if err != nil || pair == nil {
		return nil, errors.New("not existing")
	}
	getKey.Finish()
	group := &Group{}
	err = json.Unmarshal(pair.Value, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}
func (cs *ConfigStore) GetConfGroupVersions(ctx context.Context, id string) ([]*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "FindConfVersions")
	defer span.Finish()
	kv := cs.cli.KV()
	sid := configKeyGroup(ctx, id)
	data, _, err := kv.List(sid, nil)
	if err != nil {
		return nil, err

	}
	groupList := []*Group{}

	for _, pair := range data {
		group := &Group{}
		err = json.Unmarshal(pair.Value, group)
		if err != nil {
			return nil, err
		}
		groupList = append(groupList, group)

	}
	return groupList, nil

}

func (cs *ConfigStore) Put(ctx context.Context, group *Group) (*Group, error) {
	span := tracer.StartSpanFromContext(ctx, "FindConfVersions")
	defer span.Finish()
	kv := cs.cli.KV()
	data, err := json.Marshal(group)

	sid := configKeyGroupVersion(ctx, group.Id, group.Version)

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}
