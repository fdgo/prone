package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"business/support/libraries/memcatch"
	"business/support/libraries/utils"
)

func getAppBuilsInDB() ([]*domain.AppBuild, error) {
	db := dbCfPool.NewConn()
	var builds []*domain.AppBuild
	dbResult := db.Where(`status = ?`, domain.APP_BUILD_STATUS_ENABLE).
		Find(&builds)
	if dbResult.RecordNotFound() {
		return nil, nil
	}
	if dbResult.Error != nil {
		return nil, errors.SysError
	}
	return builds, nil
}

func GetAppBuilds() ([]*domain.AppBuild, error) {
	appBuilds := memcatch.AppbuildsCatch.GetAppBuilds()
	if nil != appBuilds {
		return appBuilds, nil
	}
	var err error
	appBuilds, err = getAppBuilsInDB()
	if nil != err {
		return nil, err
	}
	if nil == appBuilds {
		return nil, errors.NotFound
	}
	memcatch.AppbuildsCatch.SetAppBuilds(appBuilds)
	return appBuilds, nil
}

func GetAppPlatformBuilds(platform domain.APP_PLATFORM, language string) ([]*domain.AppBuild, error) {
	appBuilds, _ := GetAppBuilds()
	var arr []*domain.AppBuild
	if platform == domain.APP_PLATFORM_UNKOWN {
		for _, item := range appBuilds {
			if item.Language == language {
				arr = append(arr, item)
			}
		}
	} else {
		for _, item := range appBuilds {
			if item.Platform == platform && item.Language == language {
				arr = append(arr, item)
			}
		}
	}
	return arr, nil
}

func QueryAppUpdate(
	platform domain.APP_PLATFORM,
	currentVersion string,
	language string) (*domain.AppBuild, error) {
	appBuilds, err := GetAppBuilds()
	if nil != err {
		return nil, err
	}
	var newestBuild *domain.AppBuild = nil
	for _, item := range appBuilds {
		if item.Platform != platform || item.Language != language {
			continue
		}
		if utils.VersionCmp(item.Version, currentVersion) <= 0 {
			continue
		}
		if nil == newestBuild {
			newestBuild = item
		} else if utils.VersionCmp(item.Version, newestBuild.Version) > 0 {
			newestBuild = item
		}
	}
	return newestBuild, nil
}
