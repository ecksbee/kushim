package librarian

import (
	"ecksbee.com/telefacts/pkg/hydratables"
	"ecksbee.com/telefacts/pkg/renderables"
)

func pGrid(linkroleURI string, h *hydratables.Hydratable) (renderables.PGrid, []renderables.LabelRole, []renderables.Lang, error) {
	// indentedLabels, labelPacks := getIndentedLabels(linkroleURI, h)
	labelRoles := make([]renderables.LabelRole, 0)
	lang := make([]renderables.Lang, 0)
	return renderables.PGrid{
		// IndentedLabels: indentedLabels,
	}, labelRoles, lang, nil
}
