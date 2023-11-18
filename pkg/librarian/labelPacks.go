package librarian

import (
	"ecksbee.com/telefacts/pkg/renderables"
)

func reduce(labelPacks []renderables.LabelPack) *renderables.LabelPack {
	if len(labelPacks) <= 0 {
		return nil
	}
	if len(labelPacks) == 1 {
		return &labelPacks[0]
	}
	ret := make(renderables.LabelPack)
	for i := 0; i < len(labelPacks); i++ {
		item := labelPacks[i]
		for labelRole, langPack := range item {
			ret[labelRole] = make(renderables.LanguagePack)
			for lang, chardata := range langPack {
				ret[labelRole][lang] = chardata
			}
		}
	}
	return &ret
}

func destruct(labelPack renderables.LabelPack) ([]renderables.LabelRole, []renderables.Lang) {
	labelRoles := make([]renderables.LabelRole, 0, 20)
	langs := make([]renderables.Lang, 0, 8)
	for labelRole, langPack := range labelPack {
		labelRoles = append(labelRoles, labelRole)
		for lang := range langPack {
			langs = append(langs, lang)
		}
	}
	return labelRoles, langs
}
