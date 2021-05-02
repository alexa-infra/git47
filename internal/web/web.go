package web

import (
	"fmt"
	"path"
	"strings"
)

func joinURL(base string, paths ...string) string {
	p := path.Join(paths...)
	if p == "" {
		return strings.TrimRight(base, "/")
	}
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}
