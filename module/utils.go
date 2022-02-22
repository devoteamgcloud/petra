package module

import (
	"fmt"
	"path"
)

func modPath(mod Module) string {
	return path.Join(
		fmt.Sprintf("namespace=%s", mod.Namespace),
		fmt.Sprintf("name=%s", mod.Name),
		fmt.Sprintf("provider=%s", mod.Provider),
		fmt.Sprintf("version=%s", mod.Version),
		fmt.Sprintf("%s-%s-%s-%s.tar.gz", mod.Namespace, mod.Name, mod.Provider, mod.Version),
	)
}
