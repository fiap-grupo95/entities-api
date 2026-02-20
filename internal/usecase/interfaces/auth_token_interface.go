package interfaces

type AuthTokenInterface interface {
	GenerateToken(subject string) (string, error)
}
