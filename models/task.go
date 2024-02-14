package models

import "time"

// Task モデルはデータベースの tasks テーブルを表します
type Task struct {
	ID          uint           `gorm:"primary_key"`
	Task        string         `gorm:"size:255"`
	IsCompleted bool           `gorm:"default:false"`
	CreatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt   *time.Time     `gorm:"index"`
}