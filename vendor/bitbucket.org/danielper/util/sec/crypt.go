package sec

import "golang.org/x/crypto/bcrypt"

// Hash creates a hash from provided string
func Hash(content string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(content), 14)
	return string(bytes), err
}

// CheckHash verifies if provided content is the same as the hashed
func CheckHash(content, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(content))
	return err == nil
}
