package domain

type RequestListen struct {
	DomainUrl  string `json:"domain_url"`
	ListenPath string `json:"listen_path"`
}
