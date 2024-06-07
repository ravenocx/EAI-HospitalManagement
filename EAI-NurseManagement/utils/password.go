package utils

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func NormalizePassword(p string) []byte {
	return []byte(p)
}

func GeneratePassword(p string) string {
	bytePwd := NormalizePassword(p)

	saltPassword := os.Getenv("BCRYPT_SALT")

	if saltPassword == "" {
		saltPassword = "10"
	}

	salt, err := strconv.Atoi(saltPassword)
	if err != nil {
		return err.Error()
	}

	hash, err := bcrypt.GenerateFromPassword(bytePwd, salt)
	if err != nil {
		return err.Error()
	}

	return string(hash)
}

func ComparePasswords(hashedPwd, inputPwd string) error {
	byteHash := NormalizePassword(hashedPwd)
	byteInput := NormalizePassword(inputPwd)

	if err := bcrypt.CompareHashAndPassword(byteHash, byteInput); err != nil {
		return err
	}

	return nil
}
