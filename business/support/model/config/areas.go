package config

import (
	"encoding/json"
	"fmt"
	"time"
)

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

func GetAreas(language string) ([]Area, error) {
	var (
		areas    []Area
		redisKey = fmt.Sprintf("CONFIG.AREAS.%s", language)
	)
	v, err := redisPool.NewConn().Get(redisKey).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &areas); err == nil {
			return areas, nil
		}
	}
	if err := dbPool.NewConn().Where("language = ?", language).Find(&areas).Error; err != nil {
		return nil, err
	}
	for i := range areas {
		areas[i].SubAreas, _ = GetSubAreas(language, areas[i].ShortName)
	}

	v, err = json.Marshal(areas)
	if err == nil {
		redisPool.NewConn().Set(redisKey, v, time.Hour*2)
	}

	return areas, nil
}

func GetSubAreas(language, parent string) ([]SubArea, error) {
	var subAreas []SubArea
	if err := dbPool.NewConn().Where("language = ?", language).Where("parent_area = ?", parent).Find(&subAreas).Error; err != nil {
		return nil, err
	}
	return subAreas, nil
}

func GetArea(shortName, language string) (*Area, error) {
	var area Area
	if err := dbPool.NewConn().Where("short_name = ?", shortName).Where("language = ?", language).First(&area).Error; err != nil {
		return nil, err
	}
	return &area, nil
}
