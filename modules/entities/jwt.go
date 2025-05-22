package entities

import "github.com/golang-jwt/jwt/v4"

type (
	JwtClaim struct {
		Id         uint
		Username   string
		Hospital   string
		HospitalID uint
		jwt.RegisteredClaims
	}

	Jwtpassport struct {
		Id         uint
		Username   string
		Hospital   string
		HospitalID uint
	}
)
