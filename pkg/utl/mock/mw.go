package mock

import "github.com/simonhylander/booker"

// JWT mock
type JWT struct {
	GenerateTokenFn func(booker.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u booker.User) (string, error) {
	return j.GenerateTokenFn(u)
}