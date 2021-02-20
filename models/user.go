package models

import (
	"errors"
	"html"
	"os"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user
type User struct {
	Name       string    `json:"name" bson:"name"`
	Email      string    `json:"email" bson:"email"`
	ProfilePic os.File   `json:"profile_pic" bson:"profile_pic"`
	Password   string    `json:"password" bson:"password"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" bson:"updated_at"`
}

// Hash is used to hash password strings
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword is used to match password with inputted password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// BeforeSave is used to hash password
func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepare the data for saving purposes
func (u *User) Prepare() {
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate is used for validating data for various cases
func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "patch":
		if u.Email != "" {
			if err := checkmail.ValidateFormat(u.Email); err != nil {
				return errors.New("Invalid Email")
			}
		}
		return nil
	case "login":
		if u.Password == "" {
			return errors.New("Password is required")
		}
		if u.Email == "" {
			return errors.New("Password is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
		return nil
	default:
		if u.Name == "" {
			return errors.New("Name is required")
		}
		if u.Password == "" {
			return errors.New("Password is required")
		}
		if u.Email == "" {
			return errors.New("Email is required")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("Invalid email")
		}
		return nil
	}
}
