package message

type STATUS string
type OPERATION string
type DOType string

type DO struct {
	ID        string    `json:"id"`
	Type      DOType    `json:"type"`
	Attr      string    `json:"attr"`
	Elements  []Element `json:"elements"`
	Signature string    `json:"signature"` //非必须，对于特殊的DOI，如：0.Type/DOIPServiceInfo等，需要签名。或，对于DO不存在本地的情况？
}

type Element struct {
	ID     string `json:"id"`
	Length int    `json:"length"`
	Type   string `json:"type"`
	Attr   string `json:"attr"`
	Data   []byte `json:"data"`
}

//Request 是面向DO的，但是发送给DORepository，每个Request中带有DO标识
type Request struct {
	TargetID       string    `json:"target_id"`
	RequestID      string    `json:"request_id"`
	ClientID       string    `json:"client_id"`
	OperationID    OPERATION `json:"operation_id"`
	Attr           string    `json:"attr"`
	Authentication string    `json:"authentication"`
	Input          string    `json:"input"`
}

type Response struct {
	RequestID string `json:"request_id"`
	Status    STATUS `json:"status"`
	Attr      string `json:"attr"`
	Output    string `json:"output"`
	PKString  string `json:"pk_string"`
	Signature string `json:"signature"`
}
