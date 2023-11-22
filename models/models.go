package models

// user table schema
type User struct {
	UserID   uint `gorm:"primaryKey"`
	Name     string
	Email    string `gorm:"unique"`
	Password string
	Notes    []Note `gorm:"foreignKey:OwnerID;references:UserID;constraint:OnDelete:SET NULL;"`
}

// note table schema
type Note struct {
	NoteID  uint32 `gorm:"primaryKey"`
	Note    string
	OwnerID uint
}
