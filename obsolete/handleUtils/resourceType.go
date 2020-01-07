package handleUtils

type ResourceMetaData struct {
	ApiAddress  entry `xml:"apiAddress"`
	Method      entry `xml:"method"`
	Sample      entry `xml:"sample"`
	Description entry `xml:"description"`
}
