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
	PGridHashes       []string
	DGridHashes       []string
	CGridHashes       []string
}
