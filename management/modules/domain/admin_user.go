package domain

//管理员结构体
type AdminUser struct {
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
}
