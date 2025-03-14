package domain

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(100);not null"`
	Email     string    `gorm:"type:varchar(200);unique;not null"`
	Password  string    `gorm:"type:varchar(200);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
