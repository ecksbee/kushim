package renderables

type ConceptCard struct {
	Source            string
	ID                string
	Namespace         string
	Name              string
	SubstitutionGroup string
	PeriodType        string
	ItemType          string
	BalanceType       string
	PGridMap          map[string]string
	DGridMap          map[string]string
	CGridMap          map[string]string
}
