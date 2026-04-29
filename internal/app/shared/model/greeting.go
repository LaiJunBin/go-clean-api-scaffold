package model

import "gorm.io/gorm"

type Greeting struct {
	gorm.Model
	Name    string
	Message string
}
