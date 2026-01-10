package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a password string and returns a hashed version of it.
// It uses the bcrypt algorithm with a cost factor of 14.
// If the hashing process fails, it returns an error.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword checks if a given password matches a given hash.
// It returns true if the password matches the hash, and false otherwise.
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
