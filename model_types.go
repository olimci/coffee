package coffee

type item struct {
	entry *submodelEntry
	text  string
	opts  logOptions
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
