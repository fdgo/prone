package adaptor

import (
	"business/support/domain"
	"business/support/libraries/errors"
	"fmt"
	"time"
)

func UTC2Str() string {
	return time.Now().UTC().Format("2006-01-02 15:04:05.000000")
}

func AddUserAsserts(userAssets *domain.UserAssets) (bool, error) {
	if userAssets.AccountId < 1 || userAssets.AvailableVol.Sign() <= 0 {
		return false, errors.ParameterError
	}
	db := dbPool.NewConn()
	onConflictSql := fmt.Sprintf(
		`ON CONFLICT (account_id,coin_code) 
		DO UPDATE SET available_vol = "user_assets"."available_vol"+'%s',
		updated_at = '%s'`,
		userAssets.AvailableVol.String(),
		UTC2Str())
	dbResult := db.Set("gorm:insert_option", onConflictSql).Create(userAssets)
	if dbResult.Error != nil {
		return false, dbResult.Error
	}
	if dbResult.RowsAffected <= 0 {
		return false, nil
	}
	return true, nil
}
