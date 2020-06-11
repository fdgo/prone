package adaptor

import (
	"business/support/domain"
	"encoding/json"
	"fmt"
)

func GetAreas(language string) ([]domain.Area, error) {
	var (
		areas []domain.Area
		key   = fmt.Sprintf("%s.%s", domain.KAreas, language)
	)
	v, err := redisPool.NewConn().Get(key).Bytes()
	if err == nil {
		if err := json.Unmarshal(v, &areas); err == nil {
			return areas, nil
		}
	}
	if err := dbCfPool.NewConn().Where("language = ?", language).Find(&areas).Error; err != nil {
		return nil, err
	}
	for i := range areas {
		areas[i].SubAreas, _ = GetSubAreas(language, areas[i].ShortName)
	}

	v, err = json.Marshal(areas)
	if err == nil {
		redisPool.NewConn().Set(key, v, domain.KAreas.Timeout)
	}

	return areas, nil
}

func GetSubAreas(language, parent string) ([]domain.SubArea, error) {
	var subAreas []domain.SubArea
	if err := dbCfPool.NewConn().Where("language = ?", language).Where("parent_area = ?", parent).Find(&subAreas).Error; err != nil {
		return nil, err
	}
	return subAreas, nil
}
