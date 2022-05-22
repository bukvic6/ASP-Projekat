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
	err = json.Unmarshal(pair.Value, pair)
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
