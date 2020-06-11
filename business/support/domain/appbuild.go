package domain

import "time"

type APP_BUILD_STATUS int

const (
	APP_BUILD_STATUS_UNKOWN   APP_BUILD_STATUS = iota
	APP_BUILD_STATUS_ENABLE                    // 安卓
	APP_BUILD_STATUS_DISABLED                  // ios
)

type APP_PLATFORM int

const (
	APP_PLATFORM_UNKOWN APP_PLATFORM = iota
	APP_PLATFORM_ANDOID              // 安卓
	APP_PLATFORM_IOS                 // ios
	APP_PLATFORM_WEB                 // web
)

type APP_UPDATE_TYPE int

const (
	APP_UPDATE_TYPE_UNKOWN APP_UPDATE_TYPE = iota
	APP_UPDATE_TYPE_AUTO                   // 自动检测,但不强制
	APP_UPDATE_TYPE_FORCE                  // 强制
	APP_UPDATE_TYPE_MANUL                  // 手动检测模式
)

type AppBuild struct {
	ID         int64            `json:"-" gorm:"primary_key"`
	Name       string           `json:"name,omitempty"`
	Platform   APP_PLATFORM     `json:"platform,omitempty"`
	Version    string           `json:"version,omitempty"`
	Url        string           `json:"url,omitempty"`
	UpdateType APP_UPDATE_TYPE  `json:"update_type,omitempty"`
	Desc       string           `json:"desc,omitempty"`
	Language   string           `json:"-"`
	Status     APP_BUILD_STATUS `json:"-"`
	CreatedAt  *time.Time       `json:"created_at,omitempty" gorm:"type:timestamp(6);not null"`
}

func (*AppBuild) TableName() string {
	return "appbuilds"
}
