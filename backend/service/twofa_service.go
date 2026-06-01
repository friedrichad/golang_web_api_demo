package service

import (
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

func Generate2FA(email string) (*otp.Key, error) {
	return totp.Generate(totp.GenerateOpts{
		Issuer:      "MyApp",
		AccountName: email,
	})
}