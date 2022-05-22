package configstore

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	config = "config/%s"
)

func generateKey() (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id), id
}
func constructKey(id string) string {
	return fmt.Sprintf(config, id)
}
