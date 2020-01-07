package handleUtils

import "encoding/xml"

type root struct {
	MessageID   string  `xml:"messageid"`
	Count       string  `xml:"count"`
	ListRecords Records `xml:"ListRecords"`
}
type MyRespEnvelope struct {
	XMLName xml.Name
	Body    Body
}
type Body struct {
	XMLName             xml.Name
	HandleParseResponse HandleParseResponse `xml:"handleParseResponse"`
}
type HandleParseResponse struct {
	XMLName xml.Name
	Return  Return `xml:"return"`
}
type Return struct {
	Root root `xml:"root"`
}
type entry struct {
	Title string `xml:"title"`
	Value string `xml:"value"`
}
type metaData struct {
	ApiAddress  entry `xml:"apiAddress"`
	Method      entry `xml:"method"`
	Sample      entry `xml:"sample"`
	Description entry `xml:"description"`
}
type header struct {
	Identifier string `xml:"identifier"`
	State      string `xml:"state"`
	Url        string `xml:"url"`
}
type record struct {
	Header   header   `xml:"header"`
	MetaData metaData `xml:"metadata"`
}
type Records struct {
	Record []record `xml:"record"`
}
type Params struct {
	HandleId                 string `xml:"handleId"`
	IsAuthentication         string `xml:"isAuthentication"`
	RelationFlag             string `xml:"relationFlag"`
	DeepRelationFlag         string `xml:"deepRelationFlag"`
	IsReverse                string `xml:"isReverse"`
	UserGroupIds             string `xml:"userGroupIds"`
	OnlyHandleIds            string `xml:"onlyHandleIds"`
	ApplicationHandleId      string `xml:"applicationHandleId"`
	AppHandleIdEncryptBase64 string `xml:"appHandleIdEncryptBase64"`
	AppHandleIdSignBase64    string `xml:"appHandleIdSignBase64"`
	Password                 string `xml:"password"`
}
