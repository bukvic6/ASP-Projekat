package configstore

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	config  = "config/%s"
	configV = "config/%s/%s"
	all     = "config"
)

func generateKey() (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id), id
}
func configKeyVerion(id string, version string) string {
	return fmt.Sprintf(configV, id, version)

}
func constructKey(id string) string {
	return fmt.Sprintf(config, id)
}
