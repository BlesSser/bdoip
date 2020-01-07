package keys

import (
	. "DOIP/models"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/sirupsen/logrus"
	"strings"
)

//VerifyHandleRecord check pk and signature
func VerifyHandleRecord(hr *HandleRecord) (flag bool, err error) {
	body := Concat(
		[]byte(hr.Author),
		[]byte(hr.Identifier),
		[]byte(hr.Repository),
		[]byte(hr.Timestamp),
		[]byte(hr.Type),
	)
	pk, err := StringToPk(hr.Author)
	sign, err := crypto.ConfigDecodeKey(hr.Signature)
	if err != nil {
		logrus.Error("decode signature String to Bytes error")
	}
	return pk.Verify(body, sign)
}

//VerifyIdentifierWithPk 判断标识中是否和pk对应
//标识：86.5000.470/PKHash_LocalID
func VerifyIdentifierWithPk(id string, pk string) (flag bool) {
	pkFromID := id[strings.Index(id, "/")+1 : strings.Index(id, "_")]
	if len(pk) < 64 {
		logrus.Error("pkString must longer than 64")
		return false
	}
	pkFromPK := pk[10:20]
	return pkFromID == pkFromPK
}

//VerifyResponse check pk and signature
func VerifyResponse(resp *Response) (flag bool, err error) {
	body := Concat(
		[]byte(resp.Attr),
		[]byte(resp.Output),
		[]byte(resp.PKString),
		[]byte(resp.RequestID),
		[]byte(resp.Status),
	)
	pk, err := StringToPk(resp.PKString)
	sign, err := crypto.ConfigDecodeKey(resp.Signature)
	if err != nil {
		logrus.Error("decode signature String to Bytes error")
	}
	return pk.Verify(body, sign)
}

//StringToPk 转换公钥字符串到公钥对象
func StringToPk(pkStr string) (pk crypto.PubKey, err error) {
	pkBytes, err := crypto.ConfigDecodeKey(pkStr)
	if err != nil {
		logrus.Error("decode String to Bytes error")
		return
	}
	pk, err = crypto.UnmarshalPublicKey(pkBytes)
	if err != nil {
		logrus.Error("Unmarshal PK from Bytes error")
		return
	}
	return
}

// Concat an array of bytes.
func Concat(src ...[]byte) []byte {
	siz := 0
	for _, raw := range src {
		siz += len(raw)
	}

	buf := make([]byte, siz)
	pos := 0
	for _, bin := range src {
		pos += copy(buf[pos:], bin)
	}

	return buf
}
