package configstore

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	configId = "config/%s"
	config   = "config/%s/%s"
	configV  = "config/%s/%s"
	all      = "config"

	group         = "configGroup/%s/%s"
	configGroupId = "configGroup/%s"
	allG          = "configGroups"
)

func generateKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id, version), id
}
func configKeyVersion(id string, version string) string {
	return fmt.Sprintf(configV, id, version)

}
func configKey(id string) string {
	return fmt.Sprintf(configId, id)
}

func generateGroupKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(group, id, version), id
}
func configKeyGroupVersion(id string, version string) string {
	return fmt.Sprintf(group, id, version)

}
func configKeyGroup(id string) string {
	return fmt.Sprintf(configGroupId, id)
}
