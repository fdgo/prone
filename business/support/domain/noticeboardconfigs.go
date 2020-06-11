package domain

import (
	"time"
)

type NOTICEBOARD_CONFIG_STATUS int

const (
	NOTICEBOARD_CONFIG_STATUS_UNKOWN  NOTICEBOARD_CONFIG_STATUS = iota
	NOTICEBOARD_CONFIG_STATUS_ENABLE                            // 可以用
	NOTICEBOARD_CONFIG_STATUS_DISABLE                           // 禁用中
)

type NOTICEBOARD_CONFIG_TYPE int

const (
	NOTICEBOARD_CONFIG_TYPE_UNKOWN NOTICEBOARD_CONFIG_TYPE = iota // value --> 0
	NOTICEBOARD_CONFIG_TYPE_IWEB                                  // app内部网页
	NOTICEBOARD_CONFIG_TYPE_PAGE                                  // 页面
	NOTICEBOARD_CONFIG_TYPE_OWEB                                  // app外部网页
)

type NoticeBoardConfig struct {
	ID        int64                     `json:"-" gorm:"primary_key"`
	Status    NOTICEBOARD_CONFIG_STATUS `json:"-"`
	Type      NOTICEBOARD_CONFIG_TYPE   `json:"type"`
	Title     string                    `json:"title"`
	Content   string                    `json:"content"`
	Language  string                    `json:"language"`
	Link      string                    `json:"link"`
	Version   string                    `json:"-"`
	Platform  int32                     `json:"-"`
	CreatedAt *time.Time                `json:"-" gorm:"type(datetime)"`
}

func (this *NoticeBoardConfig) TableName() string {
	return "noticeboard_configs"
}
