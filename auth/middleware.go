package auth

import (
	"encoding/json"
	"strings"

	"github.com/golang-jwt/jwt"
)

type ArctirClaims struct {
	jwt.StandardClaims
	Audience multiString `json:"aud,omitempty"`
	Email    string      `json:"email"`
	Groups   []string    `json:"groups"`
}

// multiString is a temporary necessity to allow receiving the claim
// aud from the token with either string or array. There is a PR opened
// on jwt-go that once merged, this code can be removed and we can use only
// slice of strings for the audience
type multiString string

func (ms *multiString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 {
		switch data[0] {
		case '"':
			var s string
			if err := json.Unmarshal(data, &s); err != nil {
				return err
			}
			*ms = multiString(s)
		case '[':
			var s []string
			if err := json.Unmarshal(data, &s); err != nil {
				return err
			}
			*ms = multiString(strings.Join(s, ","))
		}
	}
	return nil
}
