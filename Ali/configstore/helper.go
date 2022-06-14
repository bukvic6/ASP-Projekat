package configstore

import (
	"Ali/tracer"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/net/context"
)

const (
	configId = "config/%s"
	config   = "config/%s/%s"
	configV  = "config/%s/%s"
	all      = "config"

	grouplabel    = "group/%s/%s/%s"
	group         = "group/%s/%s"
	configGroupId = "group/%s"
	allG          = "group"
)

func generateKey(version string) (string, string) {
	id := uuid.New().String()
	return fmt.Sprintf(config, id, version), id
}
func configKeyVersion(ctx context.Context, id string, version string) string {
	span := tracer.StartSpanFromContext(ctx, "constructKeyVersion")
	defer span.Finish()
	return fmt.Sprintf(configV, id, version)

}
func configKey(ctx context.Context, id string) string {
	span := tracer.StartSpanFromContext(ctx, "ConstructConfigKey")
	defer span.Finish()
	return fmt.Sprintf(configId, id)
}

func generateGroupKey(version string) (string, string) {

	id := uuid.New().String()
	return fmt.Sprintf(group, id, version), id
}
func configKeyGroupVersion(ctx context.Context, id string, version string) string {
	span := tracer.StartSpanFromContext(ctx, "ConstructKeyGroupVersion")
	defer span.Finish()
	return fmt.Sprintf(group, id, version)

}
func configKeyGroupVersionlabel(ctx context.Context, id string, version string, labels string) string {
	span := tracer.StartSpanFromContext(ctx, "ConstructConfigKey")
	defer span.Finish()
	return fmt.Sprintf(grouplabel, id, version, labels)

}
func configKeyGroup(ctx context.Context, id string) string {
	span := tracer.StartSpanFromContext(ctx, "configKeyGroup")
	defer span.Finish()
	return fmt.Sprintf(configGroupId, id)
}
