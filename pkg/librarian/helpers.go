package librarian

import (
	"encoding/hex"
	"hash/fnv"

	"ecksbee.com/telefacts/pkg/renderables"
)

func hash(schemedEntity string, linkroleURI string) string {
	h := fnv.New128a()
	h.Write([]byte(schemedEntity + linkroleURI))
	return hex.EncodeToString(h.Sum([]byte{}))
}
func stringify(e *renderables.Entity) string {
	if e == nil {
		return ""
	}
	return e.Scheme + "/" + e.CharData
}
