package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type Settle struct {
	ID            int64               `json:"-" xlsx:"-"`
	SettleId      int64               `json:"settle_id,omitempty" gorm:"default:NULL" xlsx:"#"`
	AccountId     int64               `json:"account_id,omitempty" gorm:"default:NULL" xlsx:"-"`
	FromAddress   string              `json:"from_address,omitempty" gorm:"default:NULL" xlsx:"From"`
	ToAddress     string              `json:"to_address,omitempty" gorm:"default:NULL" xlsx:"To"`
	CoinCode      string              `json:"coin_code,omitempty" gorm:"default:NULL" xlsx:"货币"`
	BlockId       string              `json:"block_id,omitempty" gorm:"default:NULL" xlsx:"-"`
	BlockHash     string              `json:"block_hash,omitempty" gorm:"default:NULL" xlsx:"-"`
	TxHash        string              `json:"tx_hash,omitempty" gorm:"default:NULL"`
	SignStr       string              `json:"sign_str,omitempty" gorm:"default:NULL" xlsx:"-"`
	QueryData     string              `json:"query_data,omitempty" gorm:"default:NULL" xlsx:"-"`
	Type          SETTLE_TYPE         `json:"type,omitempty" gorm:"default:NULL" xlsx:"类型;enum:-,充值,提现,赠送,空投,转入,转出,合约化入,合约化出,OTC转入,OTC转出,合约云转入,合约云转出,冻结,销毁,盈利"`
	Status        SETTLE_STATUS       `json:"status,omitempty" gorm:"default:NULL" xlsx:"状态;enum:-,待审核,审核通过,审核拒绝,签名完成,打包中,成功,失败"`
	Vol           decimal.Decimal     `json:"vol,omitempty" gorm:"type:decimal(36,18)" xlsx:"数额"`
	Fee           decimal.Decimal     `json:"fee,omitempty" gorm:"type:decimal(36,18)" xlsx:"手续费"`
	FeeCoinCode   string              `json:"fee_coin_code,omitempty" gorm:"default:NULL" xlsx:"手续费货币"`
	UnfreezeFunds UnfreezeFundsStatus `json:"unfreeze_funds" xlsx:"冻结状态;enum:-,已解冻,冻结中"`
	Auditor       string              `json:"auditor,omitempty" xlsx:"审核"`
	RejectReason  string              `json:"reject_reason,omitempty" xlsx:"拒绝原因"`
	Error         string              `json:"error,omitempty" gorm:"default:NULL" xlsx:"错误"`
	CreatedAt     time.Time           `json:"created_at,omitempty" gorm:"default:NULL" xlsx:"提交时间"`
	UpdatedAt     time.Time           `json:"updated_at,omitempty" gorm:"default:NULL" xlsx:"-"`
	Memo          string              `json:"memo,omitempty" gorm:"default:NULL" xlsx:"memo"`
}

func (*Settle) TableName() string {
	return "settles"
}

func (this *Settle) String() string {
	if nil == this {
		return "nil"
	}
	return fmt.Sprintf(`{ID:%d,SettleId:%d,AccountId:%d,FromAddress:%s,ToAddress:%s,CoinCode:%s,BlockId:%s,BlockHash:%s,TxHash:%s,SignStr:%s,QueryData:%s,Type:%d,Status:%d,Vol:%s,Fee:%s,FeeCoinCode:%s,UnfreezeFunds:%d,this.Auditor:%s,RejectReason:%s,Error:%s,CreatedAt:%s,Updated:%s}`,
		this.ID, this.SettleId, this.AccountId, this.FromAddress,
		this.ToAddress, this.CoinCode, this.BlockId, this.BlockHash,
		this.TxHash, this.SignStr, this.QueryData, this.Type,
		this.Status, this.Vol, this.Fee, this.FeeCoinCode,
		this.UnfreezeFunds, this.Auditor, this.RejectReason,
		this.Error, this.CreatedAt, this.UpdatedAt)
}

func (s *Settle) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, s)
}

func (s Settle) MarshalBinary() (data []byte, err error) {
	return json.Marshal(s)
}

type SETTLE_TYPE int

const (
	_                        SETTLE_TYPE = iota
	SETTLE_TYPE_DEPOSIT                  // 充值
	SETTLE_TYPE_WITHDRAW                 // 提现
	SETTLE_TYPE_BONUS                    // 超发
	SETTLE_TYPE_AIRDROP                  // 空投
	SETTLE_TYPE_TRANSFER_IN              // 内部转账(入账)
	SETTLE_TYPE_TRANSFER_OUT             // 内部转账(出账)
	SETTLE_TYPE_TRANSFER_E2C             // 资金划转(转入)
	SETTLE_TYPE_TRANSFER_C2E             // 资金划转(转出)
	SETTLE_TYPE_OTC_IN                   // OTC转入
	SETTLE_TYPE_OTC_OUT                  // OTC转出
	SETTLE_TYPE_CLOUD_IN                 // 云合约转入
	SETTLE_TYPE_CLOUD_OUT                // 云合约转出
	SETTLE_TYPE_FREEZE                   // 冻结资产
	SETTLE_TYPE_DESTROY                  // 回收超发的资产
	SETTLE_TYPE_EARNING                  // 盈利资产
)

type SETTLE_STATUS int

const (
	_                      SETTLE_STATUS = iota
	SETTLE_STATUS_CREATED                // 1 申请成功(用户提交申请)
	SETTLE_STATUS_PASSED                 // 2 审核通过(运营审核通过)
	SETTLE_STATUS_REJECTED               // 3 审核拒绝(运营审核拒绝)
	SETTLE_STATUS_SIGNED                 // 4 签名完成(生成转账signstr完成)
	SETTLE_STATUS_PENDING                // 5 打包中(待确认链上是否转账成功)
	SETTLE_STATUS_SUCCESS                // 6 成功(转账成功)
	SETTLE_STATUS_FAILED                 // 7 失败(转账失败)
)

type UnfreezeFundsStatus int

const (
	_                       UnfreezeFundsStatus = iota
	UnfreezeFundsStatusDone                     // 已经解冻
	UnfreezeFundsStatusNo                       // 未解冻

)

type WithdrawsQuery struct {
	Status        SETTLE_STATUS       `json:"status,omitempty"`
	SettleId      int64               `json:"settle_id,omitempty"`
	Limit         int                 `json:"limit,omitempty"`
	Offset        int                 `json:"offset,omitempty"`
	CoinCode      string              `json:"coin_code"`
	UnfreezeFunds UnfreezeFundsStatus `json:"unfreeze_funds"`
	MaxVol        decimal.Decimal     `json:"max_vol"`
	MinVol        decimal.Decimal     `json:"min_vol"`
}

type WithdrawsQueryReq struct {
	WithdrawsQuery
	Nonce int64 `json:"nonce,omitempty"`
}

type WithdrawsQueryResp struct {
	Total     int      `json:"total,int"`
	Withdraws []Settle `json:"withdraws,omitempty"`
	Nonce     int64    `json:"nonce,omitempty"`
}

type WithdrawsUpdateReq struct {
	Withdraws []Settle `json:"withdraws,omitempty"`
	Nonce     int64    `json:"nonce,omitempty"`
}

type WithdrawsUpdateResp struct {
	Nonce int64 `json:"nonce,omitempty"`
}

type DepositeAddresssReq struct {
	Limit    int      `json:"limit,omitempty"`
	Offset   int      `json:"offset,omitempty"`
	CoinCode string   `json:"coin_code,omitempty"`
	Start    int64    `json:"start,omitempty"`
	End      int64    `json:"end,omitempty"`
	Nonce    int64    `json:"nonce,omitempty"`
	Ignore   []string `json:"-"`
}

type DepositeAddresssResp struct {
	Total    int      `json:"total,int"`
	Addresss []string `json:"addresss,omitempty"`
	Nonce    int64    `json:"nonce,omitempty"`
}

type RechargeAmount struct {
	AccountId int64           `json:"account_id"`
	CoinCode  string          `json:"coin_code"`
	Address   string          `json:"address"`
	Amount    decimal.Decimal `json:"amount"`
	Nonce     int64           `json:"nonce"`
	Status    bool            `json:"status"`
}

func (r *RechargeAmount) Time() *time.Time {
	t := time.Unix(r.Nonce/1000, r.Nonce%1000*1000)
	return &t
}

type UserWithdrawAddress struct {
	ID        int64                 `json:"-" gorm:"primary_key"`
	Type      WITHDRAW_ADDRESS_TYPE `json:"type" gorm:"default:1"`
	AccountId int64                 `json:"account_id" gorm:"not null"`
	CoinCode  string                `json:"coin_code" gorm:"not null"`
	Address   string                `json:"address,omitempty" gorm:"not null"`
	Remark    string                `json:"remark,omitempty"`
	Memo      string                `json:"memo,omitempty"`
	CreatedAt time.Time             `json:"created_at,omitempty" gorm:"default:NULL"`
	UpdatedAt time.Time             `json:"updated_at,omitempty" gorm:"default:NULL"`
}

type WITHDRAW_ADDRESS_TYPE int

const (
	_ WITHDRAW_ADDRESS_TYPE = iota
	WITHDRAW_ADDRESS_TYPE_CHAIN
	WITHDRAW_ADDRESS_TYPE_ACCOUNT
)
