package librarian

import (
	"ecksbee.com/telefacts/pkg/renderables"
)

func stringify(e *renderables.Entity) string {
	if e == nil {
		return ""
	}
	return e.Scheme + "/" + e.CharData
}
