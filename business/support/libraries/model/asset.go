package model

import (
	"business/support/domain"
	"time"

	"business/support/libraries/errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

// SQL expression
type Ex struct {
	expr string
	args []interface{}
}

func Expr(expression string, args ...interface{}) *Ex {
	return &Ex{expr: expression, args: args}
}

// 更新用户资产
func UpdateUserAsset(db *gorm.DB, exprs []Ex, updates map[string]interface{}) (*domain.UserAssets, error) {
	now := time.Now().UTC()
	updates["updated_at"] = &now
	builder := sq.Update("user_assets")
	for i := range exprs {
		builder = builder.Where(exprs[i].expr, exprs[i].args...)
	}
	sqlStr, args, err := builder.SetMap(updates).Suffix("RETURNING *").ToSql()
	if err != nil {
		return nil, err
	}
	var assets []domain.UserAssets
	err = db.Raw(sqlStr, args...).Scan(&assets).Error
	if err != nil {
		return nil, err
	}

	if len(assets) == 0 {
		return nil, nil
	}

	return &assets[0], nil
}

// 增加用户资产(如果不存在则创建)
func AddUserAsset(db *gorm.DB, asset *domain.UserAssets) (*domain.UserAssets, error) {
	if asset.AccountId == 0 {
		return nil, errors.NoAccountId
	}
	if asset.CoinCode == "" {
		return nil, errors.NoCoinCode
	}
	now := time.Now().UTC()
	inserts := map[string]interface{}{
		"account_id":    asset.AccountId,
		"coin_code":     asset.CoinCode,
		"available_vol": asset.AvailableVol,
		"created_at":    &now,
	}
	sqlStr, args, err := sq.Insert("user_assets").SetMap(inserts).
		Suffix("ON CONFLICT (account_id,coin_code) DO UPDATE SET").
		Suffix("available_vol = user_assets.available_vol + ?,", asset.AvailableVol).
		Suffix("updated_at = ?", &now).
		Suffix("RETURNING *").
		ToSql()

	var assets []domain.UserAssets
	err = db.Raw(sqlStr, args...).Scan(&assets).Error
	if err != nil {
		return nil, err
	}

	if len(assets) == 0 {
		return nil, errors.New("unkown error")
	}

	return &assets[0], nil
}

// 减少资产
func DeductUserAsset(db *gorm.DB, accountID int64, coinCode string, vol decimal.Decimal) (*domain.UserAssets, error) {
	if accountID == 0 {
		return nil, errors.NoAccountId
	}
	if coinCode == "" {
		return nil, errors.NoCoinCode
	}
	if !vol.IsPositive() {
		return nil, errors.ParameterError
	}
	querys := []Ex{
		*Expr("coin_code = ?", coinCode),
		*Expr("account_id = ?", accountID),
		*Expr("available_vol >= ?", vol),
	}
	updates := map[string]interface{}{
		"available_vol": sq.Expr("available_vol - ?", vol),
	}
	asset, err := UpdateUserAsset(db, querys, updates)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.InsufficientBalance
	}

	return asset, nil
}

// 冻结用户资产
func FreezeUserAsset(db *gorm.DB, accountID int64, coinCode string, vol decimal.Decimal) (*domain.UserAssets, error) {
	if accountID == 0 || coinCode == "" || !domain.CheckVol(vol) {
		return nil, errors.ParameterError
	}
	querys := []Ex{
		*Expr("coin_code =?", coinCode),
		*Expr("account_id = ?", accountID),
		*Expr("available_vol >= ?", vol),
	}
	updates := map[string]interface{}{
		"available_vol": sq.Expr("available_vol - ?", vol),
		"freeze_vol":    sq.Expr("freeze_vol + ?", vol),
	}
	asset, err := UpdateUserAsset(db, querys, updates)
	if err != nil {
		return nil, err
	}
	if asset == nil {
		return nil, errors.InsufficientBalance
	}

	return asset, nil
}

// 解冻用户资产
func UnfreezeUserAsset(db *gorm.DB, accountID int64, coinCode string, vol decimal.Decimal) (*domain.UserAssets, error) {
	if accountID == 0 || coinCode == "" || !domain.CheckVol(vol) {
		return nil, errors.ParameterError
	}
	querys := []Ex{
		*Expr("account_id = ?", accountID),
		*Expr("coin_code = ?", coinCode),
		*Expr("freeze_vol >= ?", vol),
	}
	updates := map[string]interface{}{
		"freeze_vol":    sq.Expr("freeze_vol - ?", vol),
		"available_vol": gorm.Expr("available_vol + ?", vol),
	}
	asset, err := UpdateUserAsset(db, querys, updates)
	if err != nil {
		db.Rollback()
		return nil, err
	}
	if asset == nil {
		return nil, errors.InsufficientBalance
	}
	return asset, nil
}

func DeductUserFreezeAsset(db *gorm.DB, accountID int64, coinCode string, vol decimal.Decimal) (*domain.UserAssets, error) {
	if accountID == 0 || coinCode == "" || !domain.CheckVol(vol) {
		return nil, errors.ParameterError
	}
	querys := []Ex{
		*Expr("account_id = ?", accountID),
		*Expr("coin_code = ?", coinCode),
		*Expr("freeze_vol >= ?", vol),
	}
	updates := map[string]interface{}{
		"freeze_vol": sq.Expr("freeze_vol - ?", vol),
	}
	asset, err := UpdateUserAsset(db, querys, updates)
	if err != nil {
		db.Rollback()
		return nil, err
	}
	if asset == nil {
		return nil, errors.InsufficientBalance
	}
	return asset, nil
}
