package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/loggers"
	"business/support/libraries/memcatch"
	"business/support/libraries/utils"
)

func getNoticeboardConfigInDB() ([]*domain.NoticeBoardConfig, error) {
	db := dbCfPool.NewConn()
	var configs []*domain.NoticeBoardConfig
	dbResult := db.Where(`status = ?`, domain.NOTICEBOARD_CONFIG_STATUS_ENABLE).Find(&configs)
	if dbResult.RecordNotFound() {
		return nil, errors.NotFound
	}
	if dbResult.Error != nil {
		return nil, dbResult.Error
	}
	loggers.Debug.Printf("getNoticeboardConfigInDB,configs:%#v", configs)
	if len(configs) <= 0 {
		return nil, errors.NotFound
	}
	return configs, nil
}

func getNoticeboardConfigs() ([]*domain.NoticeBoardConfig, error) {
	configs := memcatch.NoticeboardConfigsCatch.GetConfigs()
	if nil != configs {
		return configs, nil
	}
	var err error
	configs, err = getNoticeboardConfigInDB()
	if nil != err {
		return nil, err
	}
	if nil == configs {
		return nil, errors.NotFound
	}
	memcatch.NoticeboardConfigsCatch.SetConfigs(configs)
	return configs, nil
}

func GetNoticeboardConfigs(language string, version string, platform domain.APP_PLATFORM) ([]*domain.NoticeBoardConfig, error) {
	configs, err := getNoticeboardConfigs()
	if nil != err {
		return nil, err
	}
	// loggers.Debug.Printf("GetNoticeboardConfigs,language:%s,version:%s,dev:%s", language, version, dev)

	confs := make([]*domain.NoticeBoardConfig, 0, len(configs))
	for _, item := range configs {
		if item.Language == language &&
			int32(platform) == int32(platform)&item.Platform {
			if len(version) < 1 ||
				len(item.Version) < 1 ||
				utils.VersionCmp(item.Version, version) >= 0 {
				confs = append(confs, item)
			}
		}
	}
	return confs, nil
}
