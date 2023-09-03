package models

import (
	"errors"
	"html"
	"strings"

	"github.com/tantowish/task-5-pbi-btpns-tantowishahhanif/app"
)

type Photo struct {
	ID       uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Title    string     `gorm:"size:100;not null;" json:"title"`
	Caption  string     `gorm:"size:255;not null;" json:"caption"`
	PhotoURL string     `gorm:"size:255;not null;" json:"photo_url"`
	UserID   string     `gorm:"not null" json:"user_id"`
	Author   app.Author `gorm:"author"`
}

func (p *Photo) Initialize() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Caption = html.EscapeString(strings.TrimSpace(p.Caption))
	p.PhotoURL = html.EscapeString(strings.TrimSpace(p.PhotoURL))
}

func (p *Photo) Validate(action string) error {
	switch strings.ToLower(action) {
	case "upload":
		if p.Title == "" {
			return errors.New("required title")
		}
		if p.Caption == "" {
			return errors.New("required caption")
		}
		if p.PhotoURL == "" {
			return errors.New("required photoURL")
		}
		if p.UserID == "" {
			return errors.New("required userID")
		}
		return nil
	case "change":
		if p.Title == "" {
			return errors.New("required title")
		}
		if p.Caption == "" {
			return errors.New("required caption")
		}
		if p.PhotoURL == "" {
			return errors.New("required photoURL")
		}
		if p.UserID == "" {
			return errors.New("required userID")
		}
		return nil
	default:
		return errors.New("invalid action")
	}
}
