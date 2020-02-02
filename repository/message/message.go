package message

type Message struct {
	Header     Header
	Body       []byte
	Credential Credential
}

type Header struct {
	HeaderLength   int32
	OperationFlag  OperationFlag
	CacheValidTime uint32
	ExpiredTime    uint32
	BodyLength     int32
	Operation      Operation
}

type Operation struct {
	Identifier string `json:"identifier"`
	Operation  string `json:"operation"`
	Response   string `json:"response"`
}

type OperationFlag struct {
	NeedAuthorize bool
	NeedCertify   bool
	NeedEncrypted bool
	KeepConnect   bool
}

type Credential struct {
	CredentialLength int32
	Version          int8
	AttributeLength  int32
	Attribute        CredAttr
	SignedDataLength int32
	SignedData       []byte
}

type CredAttr struct {
	Signer    string `json:"signer"`
	Algorithm string `json:"algorithm"`
}
