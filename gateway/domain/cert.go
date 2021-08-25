package domain

type Cert struct {
	Id           string `json:"id"`
	SerName      string `json:"ser_name"`
	CertBlock    string `json:"cert_block"`
	CertKeyBlock string `json:"cert_key_block"`
	SetTime      string `json:"set_time"`
}
