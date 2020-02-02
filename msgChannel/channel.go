package msgChannel

type MsgChannel struct {
	InputChan  chan *MessageBox
	OutputChan chan *MessageBox
}

type MessageBox struct {
	RequestID    int32
	IsEncrypted  bool
	MessageBytes []byte
}

func CreateMsgChannel(inLen int, outLen int) *MsgChannel {
	mc := &MsgChannel{
		InputChan:  make(chan *MessageBox, inLen),
		OutputChan: make(chan *MessageBox, outLen),
	}
	return mc
}

func (mc *MsgChannel) PutInMessageBytes(msgBox *MessageBox) {
	mc.InputChan <- msgBox
}

func (mc *MsgChannel) PutOutMessageBytes(msgBox *MessageBox) {
	mc.OutputChan <- msgBox
}
