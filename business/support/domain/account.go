package domain

import (
	"business/support/libraries/errors"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	emailRegexp       = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	phoneRegexp       = regexp.MustCompile("^\\+[0-9]+[\\s]+[\\s0-9]+$")
	accountNameRegexp = regexp.MustCompile("^[_0-9A-Za-z\u4e00-\u9fa5]+$")
)

type Account struct {
	ID                         int64          `json:"-" gorm:"primary_key" xlsx:"-"`
	AccountId                  int64          `json:"account_id,omitempty" gorm:"index" xlsx:"#"`
	Email                      EMAIL          `json:"email,omitempty" gorm:"default:NULL;type:varchar(128)" xlsx:"邮箱"`
	Password                   string         `json:"password,omitempty" gorm:"type:varchar(255)" xlsx:"-"`
	Phone                      PHONE          `json:"phone,omitempty" gorm:"default:NULL;type:varchar(32)" xlsx:"手机号"`
	AccountName                string         `json:"account_name,omitempty" gorm:"default:NULL;index;type:varchar(64)" xlsx:"-"`
	Avatar                     string         `json:"avatar,omitempty" gorm:"default:NULL;type:varchar(1024)"`
	AccountType                ACCOUNT_TYPE   `json:"account_type,omitempty" xlsx:"类型;enum:-,邮箱,手机号"`
	OwnerType                  int            `json:"owner_type,omitempty" gorm:"default:1;index" xlsx:"账户;enum:-,普通,内部现货做市,外部现货做市,内部合约做市,外部合约做市,特殊,合约云券商,合约云子账户,成交量结算,OTC,渠道,系统"`
	AccountStatus              ACCOUNT_STATUS `json:"status,omitempty" xlsx:"状态;enum:-,未激活,激活,禁用"`
	Assets                     []UserAssets   `json:"user_assets,omitempty" gorm:"-" xlsx:"-"`
	Settles                    []Settle       `json:"settles,omitempty" gorm:"-" xlsx:"-"`
	RegisterIp                 string         `json:"register_ip,omitempty" gorm:"type:varchar(64)" xlsx:"-"`
	LatestLoginIp              string         `json:"latest_login_ip,omitempty" gorm:"type:varchar(64)" xlsx:"-"`
	LatestLoginAt              *time.Time     `json:"latest_login_at,omitempty" xlsx:"-"`
	AssetPassword              string         `json:"asset_password,omitempty" gorm:"size:64" xlsx:"-" `
	AssetPasswordEffectiveTime int            `json:"asset_password_effective_time" gorm:"default:-2" xlsx:"-"`
	AntiFishingText            string         `json:"anti_fishing_text,omitempty" gorm:"type:varchar(256)" xlsx:"-"`
	GAKey                      string         `json:"ga_key,omitempty" gorm:"default:NULL;type:varchar(64)" xlsx:"-"`
	KYCType                    KYC_TYPE       `json:"kyc_type" gorm:"default:1" xlsx:"-"`
	Nationality                string         `json:"nationality,omitempty" gorm:"type:varchar(244)" xlsx:"-"`
	NationalityLongName        string         `json:"nationality_long_name,omitempty" gorm:"-" xlsx:"-"`
	FirstName                  string         `json:"first_name,omitempty" gorm:"type:varchar(255)" xlsx:"-"`
	LastName                   string         `json:"last_name,omitempty" gorm:"type:varchar(255)" xlsx:"-"`
	IDType                     ID_TYPE        `json:"id_type,omitempty" gorm:"type:int4;default:NULL" xlsx:"-"`
	IDNo                       string         `json:"id_no,omitempty" gorm:"type:varchar(64);default:NULL" xlsx:"-"`
	PassedIDNo                 string         `json:"-" gorm:"type:varchar(64);default:NULL;unique_index" xlsx:"-"`
	IDPhoto1                   string         `json:"id_photo1,omitempty" gorm:"type:varchar(255)" xlsx:"-"`
	IDPhoto2                   string         `json:"id_photo2,omitempty" gorm:"type:varchar(255)" xlsx:"-"`
	IDPhoto3                   string         `json:"id_photo3,omitempty" gorm:"type:varchar(255)" xlsx:"-"`
	KYCStatus                  KYC_STATUS     `json:"kyc_status" gorm:"default:1" xlsx:"KYC;enum:-,未认证,编辑中,已提交,被拒绝,通过"`
	KYCRejectReason            string         `json:"kyc_reject_reason,omitempty" xlsx:"-"`
	KYCAuditor                 string         `json:"kyc_auditor,omitempty" xlsx:"审核"`
	QD                         string         `json:"qd,omitempty" xlsx:"渠道"`
	Remark                     string         `json:"remark,omitempty" gorm:"default:null;type:varchar(1024)" xlsx:"-"`
	CreatedAt                  *time.Time     `json:"created_at,omitempty" xlsx:"注册时间"`
	UpdatedAt                  *time.Time     `json:"updated_at,omitempty" xlsx:"-"`
}

func (*Account) TableName() string {
	return "accounts"
}

type ACCOUNT_STATUS int

const (
	_                                 ACCOUNT_STATUS = iota
	ACCOUNT_STATUS_INACTIVATED                       // 未激活
	ACCOUNT_STATUS_ACTIVATED                         // 激活
	ACCOUNT_STATUS_DISABLE                           // 被禁用
	ACCOUNT_STATUS_PASSWORD_SAFE_WARN                // 登陆密码安全警告
)

type ACCOUNT_TYPE int

const (
	_                  ACCOUNT_TYPE = iota
	ACCOUNT_TYPE_EMAIL              // 邮箱注册
	ACCOUNT_TYPE_PHONE              // 手机号注册
	ACCOUNT_TYPE_CLOUD              // 云合约子账号
)

type OWNER_TYPE int

const (
	_                                  OWNER_TYPE = iota
	OWNER_TYPE_NORMAL                             // 1 普通用户
	OWNER_TYPE_SPOT_INTERNAL_ROBOT                // 2 内部现货做市账户
	OWNER_TYPE_SPOT_EXTERNAL_ROBOT                // 3 外部现货做市账户
	OWNER_TYPE_CONTRACT_INTERNAL_ROBOT            // 4 内部合约做市账户
	OWNER_TYPE_CONTRACT_EXTERNAL_ROBOT            // 5 外部合约做市账户
	OWNER_TYPE_SPECIAL                            // 6 特别交易员账户   （用于多空头寸不平或其他特殊情况交易平帐）
	OWNER_TYPE_CLOUD_PARENT                       // 7 合约云账户
	OWNER_TYPE_CLOUD_SUB                          // 8 合约云子账户
	OWNER_TYPE_VOLUME_SETTLE                      // 9 成交量结算
	OWNER_TYPE_OTC                                // 10 C2C交易账户
	OWNER_TYPE_QD                                 // 11 渠道商户
	OWNER_TYPE_SYSTEM                             // 12 系统账户
	OWNER_TYPE_RISK                               // 13 风控
	OWNER_TYPE_API                                // 14 API用户

	OWNER_TYPE_OTHER = -1
)

type KYC_TYPE int

const (
	_                 KYC_TYPE = iota
	KYC_TYPE_PERSONAL          // 个人认证
	KYC_TYPE_AGENCY            // 机构认证
)

type KYC_STATUS int

const (
	_                        KYC_STATUS = iota
	KYC_STATUS_NOT_CERTIFIED            // 未认证
	KYC_STATUS_EDITING                  // 编辑中
	KYC_STATUS_SUBMITTED                // 已提交
	KYC_STATUS_REJECTED                 // 被拒绝
	KYC_STATUS_APPROVED                 // 认证通过
)

type ID_TYPE int

const (
	_                      ID_TYPE = iota
	ID_TYPE_PASSPORT               // 护照
	ID_TYPE_ID_CARD                // 身份证
	ID_TYPE_DRIVER_LICENSE         // 驾照
)

func (m Account) String() string {
	switch m.AccountType {
	case ACCOUNT_TYPE_EMAIL:
		return string(m.Email)
	case ACCOUNT_TYPE_PHONE:
		return string(m.Phone)
	default:
		return fmt.Sprint(m.AccountId)
	}
}

// Mask 用户名
func (m Account) Mask() string {
	switch m.AccountType {
	case ACCOUNT_TYPE_EMAIL:
		return m.Email.Mask()
	case ACCOUNT_TYPE_PHONE:
		return m.Phone.Mask()
	default:
		return fmt.Sprint(m.AccountId)
	}
}
func (m *Account) UserInfo() *Account {
	user := *m
	user.AssetPassword = ""
	user.Password = ""
	if user.GAKey != "" {
		user.GAKey = "bind"
	} else {
		user.GAKey = "unbound"
	}
	user.KYCAuditor = ""
	user.IDNo = ""
	user.IDPhoto1 = ""
	user.IDPhoto2 = ""
	user.IDPhoto3 = ""
	user.IDType = 0
	user.KYCType = 0
	user.Nationality = ""
	return &user
}

func (m *Account) KYCInfo() *Account {
	return &Account{
		AccountId:       m.AccountId,
		FirstName:       m.FirstName,
		LastName:        m.LastName,
		KYCType:         m.KYCType,
		IDNo:            m.IDNo,
		IDPhoto1:        m.IDPhoto1,
		IDPhoto2:        m.IDPhoto2,
		IDPhoto3:        m.IDPhoto3,
		IDType:          m.IDType,
		Nationality:     m.Nationality,
		KYCStatus:       m.KYCStatus,
		KYCRejectReason: m.KYCRejectReason,
	}
}

type DepositAddress struct {
	ID        int64                  `json:"-"`
	CoinCode  string                 `json:"coin_code,required"`
	Address   string                 `json:"address,required"`
	Status    DEPOSIT_ADDRESS_STATUS `json:"status"`
	CreatedAt *time.Time             `json:"created_at"`
	UpdatedAt *time.Time             `json:"updated_at"`
}

func (*DepositAddress) TableName() string {
	return "deposit_addresss"
}

type DEPOSIT_ADDRESS_STATUS int

const (
	_ DEPOSIT_ADDRESS_STATUS = iota
	DEPOSIT_ADDRESS_FREE
	DEPOSIT_ADDRESS_USED
)

type EMAIL string

func (email EMAIL) MarshalJSON() ([]byte, error) {
	return json.Marshal(strings.ToLower(string(email)))
}

func (email *EMAIL) UnmarshalJSON(e []byte) error {
	if `""` == string(e) {
		return nil
	}
	var eml string
	if err := json.Unmarshal(e, &eml); err != nil {
		return err
	}
	eml = strings.TrimSpace(eml)
	if !CheckEmailFormat(eml) {
		return errors.InvalidEmailFormat
	}
	*email = FormatEmail(eml)
	return nil
}

func CheckEmailFormat(e string) bool {
	if !emailRegexp.MatchString(e) {
		return false
	}
	return true
}

func FormatEmail(e string) EMAIL {
	return EMAIL(strings.TrimSpace(strings.ToLower(e)))
}

type PHONE string

func (phone PHONE) MarshalJSON() ([]byte, error) {
	if "" == string(phone) {
		return []byte(""), nil
	}
	return json.Marshal(strings.ToLower(string(phone)))
}

func (phone *PHONE) UnmarshalJSON(e []byte) error {
	if string(e) == `""` {
		return nil
	}
	var p string
	if err := json.Unmarshal(e, &p); err != nil {
		return err
	}
	p = strings.TrimSpace(p)
	if !CheckPhoneFormat(p) {
		return errors.InvalidPhoneFormat
	}
	*phone = FormatPhone(p)
	return nil
}

func FormatPhone(p string) PHONE {
	arr := strings.Split(
		strings.Replace(
			strings.Trim(p, " "),
			" ",
			"-",
			1,
		),
		"-")
	if len(arr) < 2 {
		return ""
	}
	return PHONE(fmt.Sprintf("%s %s", arr[0], strings.Replace(arr[1], " ", "", -1)))
}

func (phone PHONE) AreaCode() string {
	arr := strings.Split(string(phone), " ")
	if len(arr) == 2 {
		return arr[0]
	} else {
		return ""
	}
}

func (phone PHONE) SubNumber() string {
	arr := strings.Split(string(phone), " ")
	if len(arr) == 2 {
		return arr[1]
	} else {
		return ""
	}
}

func (phone PHONE) FullNumber() string {
	return strings.Replace(string(phone), " ", "", -1)
}

func CheckPhoneFormat(e string) bool {
	if !phoneRegexp.MatchString(e) {
		return false
	}
	return true
}

func CheckAccountName(v string) *errors.Error {
	if len(v) > 18 {
		return errors.AccountNameTooLong
	}
	if !accountNameRegexp.MatchString(v) {
		return errors.AccountNameInvalidSymbol
	}
	return nil
}

type Person struct {
	Name     string `json:"name,omitempty"`
	Password string `json:"password,omitempty"`
	Email    EMAIL  `json:"email,omitempty"`
}

type ACCOUNT_INNOR_TYPE int

const (
	_                      ACCOUNT_INNOR_TYPE = iota
	INNER_TYPE_PLAY                           // 作势
	INNER_TYPE_DEALER                         // 吃单
	INNER_TYPE_DEPTH_MAKER                    // 深度
)

type AccountInner struct {
	ID        int64              `gorm:"primary_key;auto_increment"`
	AccountId int64              `gorm:"unique_index"`
	Type      ACCOUNT_INNOR_TYPE `gorm:"type:int2"`
	CreatedAt time.Time
}

const ASSET_PASSWORD_EFFECTIVE_TIME_PERMANENT = -1

func CheckURL(urlStr string) bool {
	if _, err := url.Parse(urlStr); err != nil {
		return false
	}
	return true
}

func (phone PHONE) Mask() string {
	if len(string(phone)) <= 4 {
		return fmt.Sprintf("*******%s", phone)
	}
	return fmt.Sprintf("*******%s", phone[len(phone)-4:])
}

func (email EMAIL) Mask() string {
	arr := strings.Split(string(email), "@")
	if len(arr) < 2 {
		return string(email)
	}
	if len(arr[0]) <= 4 {
		return fmt.Sprintf("**%s**@%s", arr[0], arr[1])
	}
	return fmt.Sprintf("**%s**@%s", arr[0][2:len(arr[0])-2], arr[1])
}
