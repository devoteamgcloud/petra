package module

import (
	"fmt"
	"path"
)

func modPath(mod Module) string {
	return path.Join(
		mod.Namespace,
		mod.Name,
		mod.Provider,
		mod.Version,
		fmt.Sprintf("%s-%s-%s-%s.tar.gz", mod.Namespace, mod.Name, mod.Provider, mod.Version),
	)
}

func modPathPartial(mod Module) string {
	return path.Join(
		mod.Namespace,
		mod.Name,
		mod.Provider,
	)
}
