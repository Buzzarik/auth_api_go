package otp

import (
	"crypto/rand"
	"fmt"
)

func GenerateOTP() (string, error) {
	const op = "GenerateOTP";
	otp := make([]byte, 2)

	_, err := rand.Read(otp)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err);
	}
	return fmt.Sprintf("%04d", int(otp[0])%10000), nil;
}