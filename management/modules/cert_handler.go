package modules

import (
	"crypto/tls"
	"encoding/json"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"time"
)

type Cert struct {
	Id           string `json:"id"`
	SerName      string `json:"ser_name"`
	CertBlock    string `json:"cert_block"`
	CertKeyBlock string `json:"cert_key_block"`
	SetTime      string `json:"set_time"`
}

//证书管理
func Certs(ctx *gin.Context) {
	datas, err := getAllCertData()
	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	var certs []*Cert
	if datas.Count > 0 {
		for _, kv := range datas.Kvs {
			cert := new(Cert)
			err := json.Unmarshal(kv.Value, cert)
			if err == nil {
				certs = append(certs, cert)
			}
		}
		resultCtx{ctx}.SucResult(certs)
		return
	}
	resultCtx{ctx}.SucResult(make([]string, 0))
}

func PutCert(ctx *gin.Context) {
	certBlock := ctx.PostForm("cert_block")
	certKeyBlock := ctx.PostForm("cert_key_block")
	serName := ctx.PostForm("ser_name")

	if serName == "" || certKeyBlock == "" || certBlock == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	_, err := tls.X509KeyPair([]byte(certBlock), []byte(certKeyBlock))
	if err != nil {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	//有接收到certId 就是修改操作， 否则就是新增
	var certId string
	certId = ctx.Param("cert_id")
	if certId == "" {
		certId = uuid.Must(uuid.NewV4()).String()
	}
	cert := new(Cert)
	cert.Id = certId
	cert.SerName = serName
	cert.CertBlock = certBlock
	cert.CertKeyBlock = certKeyBlock
	cert.SetTime = time.Now().Format("2006/1/2 15:04:05")

	certB, err := json.Marshal(cert)
	if err != nil {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}

	err = putCertData(cert.Id, string(certB))
	if err != nil {
		resultCtx{ctx}.ErrResult(SystemError)
		return
	}
	resultCtx{ctx}.SucResult(cert)
}

func DelCert(ctx *gin.Context) {
	certId := ctx.Param("cert_id")
	if certId == "" {
		resultCtx{ctx}.ErrResult(DataParseError)
		return
	}
	deleted := delCertData(certId)
	if deleted {
		resultCtx{ctx}.SucResult(nil)
		return
	}
	resultCtx{ctx}.ErrResult(SystemError)
}
