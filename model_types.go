package coffee

type item struct {
	entry *submodelEntry
	text  string
}

type submodelEntry struct {
	submodel Submodel
}

type Section uint8

const (
	SectionHeader Section = iota
	SectionBody
	SectionFooter
)
