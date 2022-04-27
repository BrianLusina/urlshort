package models

import "time"

// BaseModel definition same as gorm.Model, but including other common columns
type BaseModel struct {
	Identifier string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
