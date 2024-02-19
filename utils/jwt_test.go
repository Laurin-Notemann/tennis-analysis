package utils

import (
	"regexp"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type TestError struct {
	IsError       bool
	ExpectedError error
}

type TestGenerateJwtInput struct {
	name  string
	error TestError
	input TokenGenInput
}

var tokenGenerator = ProdTokenGenerator{}

func TestGenerateJwt(t *testing.T) {
	testInput := []TestGenerateJwtInput{
		{
			name: "successful valid jwt token",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: TokenGenInput{
				UserId:        uuid.New(),
				Username:      "tim",
				Email:         "tim@test",
				ExpiryDate:    time.Now().Add(5 * time.Minute),
				SigningKey:    "Random",
				IsAccessToken: true,
			},
		},
		{
			name: "successful valid jwt token",
			error: TestError{
				IsError:       false,
				ExpectedError: nil,
			},
			input: TokenGenInput{
				UserId:        uuid.New(),
				Username:      "tim",
				Email:         "tim@test",
				ExpiryDate:    time.Now(),
				SigningKey:    "Random",
				IsAccessToken: true,
			},
		},
	}

	for _, data := range testInput {
		t.Run("token: "+data.name, func(t *testing.T) {
			jwt, err := tokenGenerator.GenerateNewJwtToken(data.input)

			if data.error.IsError {
			} else {
				if assert.NoError(t, err, "Problem with generating a jwt") {
					assert.Regexp(t, regexp.MustCompile("^[A-Za-z0-9_-]{2,}(?:\\.[A-Za-z0-9_-]{2,}){2}$"), jwt)
				}
			}
		})
	}
}
