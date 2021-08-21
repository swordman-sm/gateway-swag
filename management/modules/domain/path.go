package domain

type Path struct {
	Id                    string    `json:"id"`
	DomainId              string    `json:"domain_id"`
	ReqMethod             string    `json:"req_method"`
	ReqPath               string    `json:"req_path"`
	SearchPath            string    `json:"search_path"`
	ReplacePath           string    `json:"replace_path"`
	CircuitBreakerRequest int       `json:"circuit_breaker_request"`
	CircuitBreakerPercent int       `json:"circuit_breaker_percent"`
	CircuitBreakerTimeout int       `json:"circuit_breaker_timeout"`
	CircuitBreakerMsg     string    `json:"circuit_breaker_msg"`
	CircuitBreakerEnabled bool      `json:"circuit_breaker_enabled"`
	CircuitBreakerForce   bool      `json:"circuit_breaker_force"`
	PrivateProxyEnabled   bool      `json:"private_proxy_enabled"`
	LbType                string    `json:"lb_type"`
	Targets               []*Target `json:"targets"`
	SetTime               string    `json:"set_time"`
}
