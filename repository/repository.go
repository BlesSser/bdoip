package repository

import (
	"bdoip/msgChannel"
	"bdoip/repository/certificate"
	"bdoip/repository/coder"
	. "bdoip/repository/message"
	"bdoip/repository/operator"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

type Repository struct {
	MsgChannel *msgChannel.MsgChannel
	Coder      *coder.Coder
	Operator   operator.Operator
	Cert       certificate.Certificate
	Logger     logrus.Logger
}

func CreateRepository() (repo Repository, err error) {
	//logger := logrus.Logger{
	//	Out:       nil,
	//	Hooks:     nil,
	//	Formatter: nil,
	//	Level:     0,
	//}
	return
}

func (r *Repository) Run() {

	for {
		requestBox, ok := <-r.MsgChannel.InputChan
		if !ok {
			r.Logger.Info("inputChan is empty,wait")
			time.Sleep(time.Duration(rand.Int63n(500)) * time.Millisecond)
		} else {
			request, err := r.Coder.Decode(requestBox.MessageBytes)
			if err != nil {
				r.Logger.Error("decode request message error: ", err)
				continue
			}
			valid, err := certificate.VerifyByPk(*request)
			if !valid {
				r.Logger.Error("request message signature invalid: ", err)
				continue
			}
			response, err := r.Operator.ProcessMessage(request)
			if err != nil {
				r.Logger.Error("process request message error: ", err)
				continue
			}
			err = r.decorate(response, request.Header.OperationFlag)
			if err != nil {
				r.Logger.Error("decorate response message header error: ", err)
				continue
			}
			respMsgBytes, err := r.Coder.Encode(response)
			if err != nil {
				r.Logger.Error("encode message error: ", err)
				continue
			}
			isEncrypted := false
			if request.Header.OperationFlag.NeedEncrypted {
				respMsgBytes, err = certificate.EncryptByPk(respMsgBytes)
				if err != nil {
					r.Logger.Error("encrypt response message error: ", err)
					continue
				}
				isEncrypted = true
			}
			respBox := &msgChannel.MessageBox{
				IsEncrypted:  isEncrypted,
				MessageBytes: respMsgBytes,
			}
			r.MsgChannel.OutputChan <- respBox
		}
	}
}

func (r *Repository) decorate(resp *Message, reqOpFlag OperationFlag) (err error) {

	return
}

func (r *Repository) packBox(msgBytes []byte) {

}

func (r *Repository) HandleDOIPMessage(msgBox *msgChannel.MessageBox) {

}

func (r *Repository) SetMsgChannel(MsgChannel *msgChannel.MsgChannel) {
	r.MsgChannel = MsgChannel
}
