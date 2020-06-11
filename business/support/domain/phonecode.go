package domain

type PhoneCode struct {
	Codes string `json:"codes"`
}

func (this *PhoneCode) TableName() string {
	return "phone_codes"
}
