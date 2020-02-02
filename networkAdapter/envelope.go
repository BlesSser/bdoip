package networkAdapter

type MsgWrapper struct {
	Envelope Envelope
	Message  []byte
}

type Envelope struct {
	IsCompressed   bool
	IsEncrypted    bool
	IsTruncated    bool
	IsRequested    bool
	MajorVersion   uint8
	MinVersion     uint8
	RequestId      int32
	SequenceNumber int32
	MessageLength  int32
}
