package domain

type Area struct {
	ID        int       `json:"-" gorm:"primary_key"`
	Language  string    `json:"-"`
	ShortName string    `json:"short_name"`
	LongName  string    `json:"long_name"`
	SubAreas  []SubArea `json:"sub_areas,omitempty" gorm:"-"`
}

type SubArea struct {
	ID         int    `json:"-" gorm:"primary_key"`
	Language   string `json:"-"`
	ParentArea string `json:"-"`
	ShortName  string `json:"short_name"`
	LongName   string `json:"long_name"`
}
