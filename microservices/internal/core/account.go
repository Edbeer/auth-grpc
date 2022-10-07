package models

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Uuid      uuid.UUID `json:"user_id" db:"user_id" validate:"omitempty,uuid"`
	Name      string    `json:"name" db:"name" validate:"required_with,lte=30"`
	Email     string    `json:"email" db:"email" validate:"omitempty,email"`
	Pass      string    `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	Role      string    `json:"role" db:"role" validate:"omitempty,lte=10"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Hash user password with bcrypt
func (u *User) HashPassword() error {
	salt, err := salt()
	if err != nil {
		return err
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Pass + salt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Pass = string(hashedPassword)
	return nil
}

// Compare user password and payload
func (u *User) ComparePasswords(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Pass), []byte(password)); err != nil {
		return err
	}
	return nil
}

// Sanitize user password
func (u *User) SanitizePassword() {
	u.Pass = ""
}

// Prepare user for register
func (u *User) PrepareCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Pass = strings.TrimSpace(u.Pass)

	if err := u.HashPassword(); err != nil {
		return err
	}
	return nil
}

// Prepare user for register
func (u *User) PrepareUpdate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	return nil
}

func salt() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	_, err := r.Read(b)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}