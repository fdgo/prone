package domain

type APP_SWITCH_STATUS int

const (
	APP_SWITCH_STATUS_UNKOWN APP_SWITCH_STATUS = iota
	APP_SWITCH_STATUS_OFF                      // 关闭开关,统统返回实体
	APP_SWITCH_STATUS_ON                       // 打开开关,检测版本,根据版本返回实体
)

type AppSwitch struct {
	ID         int64             `json:"-" gorm:"primary_key"`
	Platform   APP_PLATFORM      `json:"platform,omitempty"`
	MinVersion string            `json:"min_version,omitempty"`
	MaxVersion string            `json:"max_version,omitempty"`
	Status     APP_SWITCH_STATUS `json:"-"`
}

func (*AppSwitch) TableName() string {
	return "appswitchs"
}

type APP_SKIN_TYPE int

const (
	APP_SKIN_TYPE_UNKOWN APP_SKIN_TYPE = iota
	APP_SKIN_TYPE_REAL                 // 实体
	APP_SKIN_TYPE_FAKE                 // 马甲
)

type APP_SKIN_STATUS int

const (
	APP_SKIN_STATUS_UNKOWN APP_SKIN_STATUS = iota
	APP_SKIN_STATUS_ENABLED
	APP_SKIN_STATUS_DISABLED
)

type AppSkin struct {
	ID         int64           `json:"-" gorm:"primary_key"`
	Skin       APP_SKIN_TYPE   `json:"skin"`
	Ridge      string          `json:"ridge"`
	Platform   APP_PLATFORM    `json:"platform"`
	MinVersion string          `json:"min_version"`
	MaxVersion string          `json:"max_version"`
	Status     APP_SKIN_STATUS `json:"-"`
}

func (*AppSkin) TableName() string {
	return "appskins"
}

type AppLocal struct {
	ID         int64           `json:"-" gorm:"primary_key"`
	Local      string          `json:"local"`
	Language   string          `json:"language"`
	Platform   APP_PLATFORM    `json:"platform"`
	MinVersion string          `json:"min_version"`
	MaxVersion string          `json:"max_version"`
	Status     APP_SKIN_STATUS `json:"-"`
}

func (*AppLocal) TableName() string {
	return "applocals"
}

type AppSkinPkg struct {
	Swtichs []*AppSwitch
	Skins   []*AppSkin
	Locals  []*AppLocal
}
