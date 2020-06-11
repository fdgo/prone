package domain

type ErrorConfig struct {
	ID      int32  `json:"-" gorm:"primary_key"`
	ErrCode string `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
	Local   string `json:"local"`
}

func (this *ErrorConfig) TableName() string {
	return "error_configs"
}
