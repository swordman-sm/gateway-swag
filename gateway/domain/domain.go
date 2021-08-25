package domain

type Domain struct {
	Id         string `json:"id"`
	DomainName string `json:"domain_name"`
	DomainUrl  string `json:"domain_url"`
	LbType     string `json:"lb_type"`
	//代理实现
	Targets            []*Target       `json:"targets"`
	BlackIps           map[string]bool `json:"black_ips"`
	RateLimiterNum     float64         `json:"rate_limiter_num"`
	RateLimiterMsg     string          `json:"rate_limiter_msg"`
	RateLimiterEnabled bool            `json:"rate_limiter_enabled"`
	SetTime            string          `json:"set_time"`
}
