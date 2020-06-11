package adaptor

//func getUserLevelConfigsInDB() ([]*domain.UserLevelConfig, error) {
//	db := dbCfPool.NewConn()
//	var configs []*domain.UserLevelConfig
//	dbResult := db.Where(`status = ?`, domain.USER_LEVEL_CONFIG_STATUS_ENABLED).
//		Find(&configs)
//	if dbResult.RecordNotFound() {
//		return nil, errors.NotFound
//	}
//	if dbResult.Error != nil {
//		return nil, dbResult.Error
//	}
//	return configs, nil
//}

//func GetUserLevelConfigs() ([]*domain.UserLevelConfig, error) {
//	configs := memcatch.UserLevelConfigCatch.GetConfigs()
//	if nil != configs {
//		return configs, nil
//	}
//	var err error
//	configs, err = getUserLevelConfigsInDB()
//	if nil != err {
//		return nil, err
//	}
//	if nil == configs || len(configs) < 1 {
//		return nil, errors.NotFound
//	}
//	memcatch.UserLevelConfigCatch.SetConfigs(configs)
//	return configs, nil
//}

//func QueryUserCapacity(
//	accountId int64,
//	capacityType domain.USER_CAPACITY_TYPE,
//	coinCode string) (*domain.UserLevelConfig, error) {
//	userPoints, err := GetAllUserPoint(accountId)
//	if nil != err {
//		return nil, err
//	}
//	configs, err := GetUserLevelConfigs()
//	if nil != err {
//		return nil, err
//	}
//	var userLevel *domain.UserLevelConfig = nil
//	for _, point := range userPoints {
//		for _, level := range configs {
//			if (int32(level.CapacityType) & int32(capacityType)) != int32(capacityType) {
//				continue
//			}
//			if int32(point.PointType) != (int32(point.PointType)&level.PointType) ||
//				(len(level.CapacityKey) > 0 && level.CapacityKey != coinCode) {
//				continue
//			}
//			if level.PointCount < 1 || level.PointCount > point.Count.IntPart() {
//				continue
//			}
//			if nil == userLevel {
//				userLevel = level
//			} else if userLevel.Level < level.Level {
//				userLevel = level
//			}
//		}
//	}
//	return userLevel, nil
//}
