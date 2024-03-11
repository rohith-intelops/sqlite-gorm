package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id int64 `gorm:"primaryKey;autoIncrement" json:"ID,omitempty"`

	Comments string `json:"comments,omitempty"`

	Likes int `json:"likes,omitempty"`

	Name string `json:"name,omitempty"`
}
