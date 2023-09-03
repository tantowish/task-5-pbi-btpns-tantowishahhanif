package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	"github.com/tantowish/task-5-vix-btpns-tantowishahhanif/helpers/hash"
)

type User struct {
	ID        string    `gorm:"primary_key;unique" json:"id"`
	Username  string    `gorm:"size:50;not null;unique" json:"username"`
	Email     string    `gorm:"size:50;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Photos    Photo     `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"photos"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) HashPassword() error {
	hashedPassword, err := hash.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) error {
	return hash.VerifyPassword(u.Password, password)
}

func (u *User) Initialize() {
	u.ID = uuid.New().String()
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if u.Email == "" {
			return errors.New("required email")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		return nil
	case "register":
		if u.ID == "" {
			return errors.New("required id")
		}
		if u.Username == "" {
			return errors.New("required username")
		}
		if u.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		if len(u.Password) < 6 {
			return errors.New("password is too short")
		}
		return nil
	case "update":
		if u.ID == "" {
			return errors.New("required id")
		}
		if u.Username == "" {
			return errors.New("required username")
		}
		if u.Email == "" {
			return errors.New("required email")
		}
		if err := checkmail.ValidateFormat(u.Email); err != nil {
			return errors.New("invalid email")
		}
		if u.Password == "" {
			return errors.New("required password")
		}
		if len(u.Password) < 6 {
			return errors.New("password is too short")
		}
		return nil
	default:
		return errors.New("invalid Action")
	}
}
