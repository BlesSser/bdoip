package coder

import . "bdoip/repository/message"

type Coder struct {
}

func (c *Coder) Encode(msg *Message) (msgBytes []byte, err error) {

	return
}

func (c *Coder) Decode(msgBytes []byte) (msg *Message, err error) {

	return
}
