package interfaces

type AuthInterface interface {
	getUserKey(userId string) string
}
