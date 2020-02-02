package handleSystem

type Record struct {
	Identifier string `json:"identifier"`
	Author     string `json:"author"`
	Repository string `json:"repository"`
	Timestamp  string `json:"timestamp"`
	Type       string `json:"type"`
	Signature  string `json:"signature"`
}
