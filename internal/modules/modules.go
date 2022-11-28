package modules

import (
	"fmt"
	"os"
	"path"

	"github.com/devoteamgcloud/petra/internal/storage"
)

type module struct {
	Namespace string `yaml:"namespace"`
	Name      string `yaml:"name"`
	Provider  string `yaml:"provider"`
	Version   string `yaml:"version"`
}

var backend *storage.GCSBackend

func Init() error {
	var err error
	backend, err = storage.InitGCSBackend(os.Getenv("GCS_BUCKET"))
	if err != nil {
		return fmt.Errorf("failed to setup storage: %v", err)
	}
	return nil
}

func modPath(mod module) string {
	return path.Join(
		mod.Namespace,
		mod.Name,
		mod.Provider,
		mod.Version,
		fmt.Sprintf("%s-%s-%s-%s.tar.gz", mod.Namespace, mod.Name, mod.Provider, mod.Version),
	)
}

func modPathPartial(mod module) string {
	return path.Join(
		mod.Namespace,
		mod.Name,
		mod.Provider,
	)
}
