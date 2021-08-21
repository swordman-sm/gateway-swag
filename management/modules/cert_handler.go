package modules

import (
	"crypto/tls"
	"encoding/json"
	"gateway-swag/management/modules/base"
	"gateway-swag/management/modules/domain"
	"github.com/gin-gonic/gin"
	uuid "github.com/iris-contrib/go.uuid"
	"time"
)

//证书管理
func Certs(ctx *gin.Context) {
	datas, err := getAllCertData()
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	var certs []*domain.Cert
	if datas.Count > 0 {
		for _, kv := range datas.Kvs {
			cert := new(domain.Cert)
			err := json.Unmarshal(kv.Value, cert)
			if err == nil {
				certs = append(certs, cert)
			}
		}
		base.Result{Context: ctx}.SucResult(certs)
		return
	}
	base.Result{Context: ctx}.SucResult(make([]string, 0))
}

func PutCert(ctx *gin.Context) {
	certBlock := ctx.PostForm("cert_block")
	certKeyBlock := ctx.PostForm("cert_key_block")
	serName := ctx.PostForm("ser_name")

	if serName == "" || certKeyBlock == "" || certBlock == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	_, err := tls.X509KeyPair([]byte(certBlock), []byte(certKeyBlock))
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	//有接收到certId 就是修改操作， 否则就是新增
	var certId string
	certId = ctx.Param("cert_id")
	if certId == "" {
		certId = uuid.Must(uuid.NewV4()).String()
	}
	cert := new(domain.Cert)
	cert.Id = certId
	cert.SerName = serName
	cert.CertBlock = certBlock
	cert.CertKeyBlock = certKeyBlock
	cert.SetTime = time.Now().Format("2006/1/2 15:04:05")

	certB, err := json.Marshal(cert)
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}

	err = putCertData(cert.Id, string(certB))
	if err != nil {
		base.Result{Context: ctx}.ErrResult(base.SystemError)
		return
	}
	base.Result{Context: ctx}.SucResult(cert)
}

func DelCert(ctx *gin.Context) {
	certId := ctx.Param("cert_id")
	if certId == "" {
		base.Result{Context: ctx}.ErrResult(base.DataParseError)
		return
	}
	deleted := delCertData(certId)
	if deleted {
		base.Result{Context: ctx}.SucResult(nil)
		return
	}
	base.Result{Context: ctx}.ErrResult(base.SystemError)
}
