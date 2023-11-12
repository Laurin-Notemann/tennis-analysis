package utils

import (
	"time"
)

type MockTokenGenerator struct {
	CallOut           int
	ExpiryDateAccess  time.Duration
	ExpiryDateRefresh time.Duration
}

func (t *MockTokenGenerator) GenerateNewJwtToken(input TokenGenInput) (string, error) {
	gen := ProdTokenGenerator{}
	t.CallOut++
	if input.IsAccessToken {
		time := time.Now().Add(t.ExpiryDateAccess)
		input.ExpiryDate = time
		token, err := gen.GenerateNewJwtToken(input)
		return token, err
	}

	time := time.Now().Add(t.ExpiryDateRefresh)
	input.ExpiryDate = time
	token, err := gen.GenerateNewJwtToken(input)
	return token, err
}
