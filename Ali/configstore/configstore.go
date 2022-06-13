package configstore

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
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

func (cs *ConfigStore) Post(config *Config) (*Config, error) {
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
func (cs *ConfigStore) CheckId(reqId string) bool {
	kv := cs.cli.KV()
	k, _, err := kv.Get(reqId, nil)
	if err != nil || k == nil {
		return false
	}

	return true

}

func (cs *ConfigStore) SaveId() string {
	kv := cs.cli.KV()
	idempotencyId := uuid.New().String()
	p := &api.KVPair{Key: idempotencyId, Value: nil}
	_, err := kv.Put(p, nil)
	if err != nil {
		return "error"
	}
	return idempotencyId
}

func (cs *ConfigStore) GetAll() ([]*Config, error) {
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
func (cs *ConfigStore) GetAllGroups() ([]*Group, error) {
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
func (cs *ConfigStore) AddConfigVersion(config *Config) (*Config, error) {
	kv := cs.cli.KV()
	data, err := json.Marshal(config)
	if err != nil {
		return nil, err
	}

	sid := configKeyVersion(config.Id, config.Version)
	_, err = cs.GetConf(config.Id, config.Version)

	if err == nil {
		return nil, errors.New("version already exists! ")
	}

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return config, nil
}
func (cs *ConfigStore) GetConf(id string, version string) (*Config, error) {
	kv := cs.cli.KV()

	sid := configKeyVersion(id, version)
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
func (cs *ConfigStore) Delete(id string, version string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(configKeyVersion(id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"deleted": id}, nil
}
func (cs *ConfigStore) GetConfVersions(id string) ([]*Config, error) {
	kv := cs.cli.KV()
	sid := configKey(id)
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
func (cs *ConfigStore) Group(group *Group) (*Group, error) {
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
		configKeyGroupVersionlabel(group.Id, group.Version, labels)

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

	/*	for _, v := range group.Config {
		labels := ""
		keys := make([]string, 0, len(v.Entries))
		for k, _ := range v.Entries {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			labels += k + ":" + v.Entries[k] + ";"
		}
		labels = labels[:len(labels)-1]
		sid, rid := generateGroupKey(group.Version, labels)
		group.Id = rid
		data, err := json.Marshal(group)
		if err != nil {
			return nil, err
		}
		p := &api.KVPair{Key: sid, Value: data}
		_, err = kv.Put(p, nil)
		if err != nil {
			return nil, err
		}*/

}

func (cs *ConfigStore) AddConfigGroupVersion(group *Group) (*Group, error) {
	kv := cs.cli.KV()
	data, err := json.Marshal(group)

	sid := configKeyGroupVersion(group.Id, group.Version)
	_, err = cs.GetGroup(group.Id, group.Version)

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
func (cs *ConfigStore) DeleteGroup(id string, version string) (map[string]string, error) {
	kv := cs.cli.KV()
	_, err := kv.Delete(configKeyGroupVersion(id, version), nil)
	if err != nil {
		return nil, err
	}
	return map[string]string{"deleted": id}, nil
}
func (cs *ConfigStore) GetGroup(id string, version string) (*Group, error) {
	kv := cs.cli.KV()

	sid := configKeyGroupVersion(id, version)
	pair, _, err := kv.Get(sid, nil)
	if err != nil || pair == nil {
		return nil, errors.New("not existing")
	}
	group := &Group{}
	err = json.Unmarshal(pair.Value, group)
	if err != nil {
		return nil, err
	}
	return group, nil
}
func (cs *ConfigStore) GetConfGroupVersions(id string) ([]*Group, error) {
	kv := cs.cli.KV()
	sid := configKeyGroup(id)
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

/*func (cs *ConfigStore) Filter(id string, version string, label string) ([]*ConfigG, error) {
	kv := cs.cli.KV()
	data, _, err := kv.List(filter(id, version, label), nil)
	if err != nil {
		return nil, err
	}

	configs := []*ConfigG{}
	for _, pair := range data {
		config := &ConfigG{}
		err = json.Unmarshal(pair.Value, config)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}*/

func (cs *ConfigStore) Put(group *Group) (*Group, error) {
	kv := cs.cli.KV()
	data, err := json.Marshal(group)

	sid := configKeyGroupVersion(group.Id, group.Version)

	p := &api.KVPair{Key: sid, Value: data}
	_, err = kv.Put(p, nil)
	if err != nil {
		return nil, err
	}
	return group, nil
}
