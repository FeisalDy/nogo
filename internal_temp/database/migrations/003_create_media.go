package migrations

import "gorm.io/gorm"

type Media struct {
	gorm.Model

	UploadBy *uint `json:"upload_by"`
	Uploader *User `json:"uploader" gorm:"foreignKey:UploadBy;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`

	URL         string  `json:"url" gorm:"not null"`
	Type        *string `json:"type"`
	Description *string `json:"description" gorm:"type:text"`
	FileSize    *int    `json:"file_size"`
	MimeType    *string `json:"mime_type"`
}

func Migration003CreateMedia() Migration {
	return Migration{
		ID:          "003_create_media",
		Description: "Create media table with uploader relationship",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&Media{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&Media{})
		},
	}
}
