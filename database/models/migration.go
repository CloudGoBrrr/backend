package models

type Migration struct {
	Name string `gorm:"unique"`
}
