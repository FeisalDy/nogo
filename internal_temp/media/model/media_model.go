package model

import "gorm.io/gorm"

type Media struct {
	gorm.Model
	UploadBy    *uint   `json:"upload_by" gorm:"index"`
	URL         string  `json:"url" gorm:"not null"`
	Type        *string `json:"type"`
	Description *string `json:"description"`
	FileSize    *int    `json:"file_size"`
	MimeType    *string `json:"mime_type"`
}

func (Media) TableName() string {
	return "media"
}
