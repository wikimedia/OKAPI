package password

import (
	"golang.org/x/crypto/bcrypt"
)

func Hash(clientPassword string) (string, error) {
	password := []byte(clientPassword)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)

	return string(hashedPassword), err
}

func CompareHashWithPassword(clientHash string, clientPassword string) error {
	hash := []byte(clientHash)
	password := []byte(clientPassword)

	return bcrypt.CompareHashAndPassword(hash, password)
}
