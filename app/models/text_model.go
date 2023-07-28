package models

import (
	"time"
)

type TextTracking struct {
	Id            int64     `json:"-"`
	TextId        string    `json:"id"`
	TextData      string    `json:"data"`
	UploadDate    time.Time `json:"uploadDate"`
	DownloadCount int64     `json:"downloadCount"`
	DownloadLimit int64     `json:"downloadLimit"`
	ExpireTime    time.Time `json:"expireTime"`
	IsDeleted     bool      `json:"-"`
}
