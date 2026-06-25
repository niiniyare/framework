package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	h, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
	fmt.Printf("UPDATE users SET password_hash = '%s' WHERE email = 'admin@acme.com';\n", h)
}
