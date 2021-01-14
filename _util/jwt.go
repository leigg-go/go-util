package _util

import (
	"github.com/pkg/errors"
	"gopkg.in/jose.v1/crypto"
	"gopkg.in/jose.v1/jws"
	"gopkg.in/jose.v1/jwt"
	"time"
)

const (
	JwtKey_UID = "UID"
)

type JwtHelper struct {
	secret       []byte
	validator    *jwt.Validator
	cryptoMethod crypto.SigningMethod
	timeout      time.Duration
}

var JWT *JwtHelper

var AdminBgClaims = jws.Claims{
	//"aud": "http://adminbg.com",
	"iss": "https://github.com/chaseSpace/adminbg",
	//"iat": time.Now(), timestamp of issue at
}

func InitJWT(timeout, notValidBefore time.Duration, secret string) {
	JWT = &JwtHelper{
		validator:    jws.NewValidator(AdminBgClaims, timeout, notValidBefore, nil),
		cryptoMethod: crypto.SigningMethodHS256,
		secret:       []byte(secret),
		timeout:      timeout,
	}
}

func (j *JwtHelper) GenToken(claims jws.Claims) ([]byte, error) {
	claims.SetExpiration(time.Now().Add(j.timeout))

	token := jws.NewJWT(claims, j.cryptoMethod)
	b, err := token.Serialize(j.secret)
	if err != nil {
		return nil, errors.Wrap(err, "jwt")
	}
	return b, nil
}

func (j *JwtHelper) Verify(token []byte) (jwt.Claims, error) {
	JWT, err := jws.ParseJWT(token)
	if err != nil {
		return nil, errors.Wrap(err, "jwt")
	}
	err = JWT.Validate(j.secret, j.cryptoMethod, j.validator)
	if err != nil {
		return nil, errors.Wrap(err, "jwt")
	}
	return JWT.Claims(), nil
}
