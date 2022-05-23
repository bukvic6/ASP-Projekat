package configstore

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"os"
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

	sid := configKeyVersion(config.Id, config.Version)

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
func (cs *ConfigStore) AddConfigGroupVersion(group *Group) (*Group, error) {
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
	if err != nil {
		return nil, err
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
