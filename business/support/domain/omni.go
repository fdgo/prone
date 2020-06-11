package domain

import (
	"github.com/shopspring/decimal"
)

// OmniTx 交易结构
type OmniTx struct {
	TxID             string          `json:"txid,omitempty"`             // the hash of the transaction of this offer
	SendingAddress   string          `json:"sendingaddress,omitempty"`   // the Bitcoin address of the sender
	ReferenceAddress string          `json:"referenceaddress,omitempty"` // the Bitcoin address of the sender
	Confirmations    int             `json:"confirmations"`              // the number of transaction confirmations
	IsMine           bool            `json:"ismine"`                     // whether the transaction involves an address in the wallet
	Fee              decimal.Decimal `json:"fee,omitempty"`              // the transaction fee in bitcoins
	BlockTime        int64           `json:"blocktime,omitempty"`        // the timestamp of the block that contains the transaction
	Valid            bool            `json:"valid"`                      // whether the transaction is valid
	PositionInBlock  int64           `json:"positioninblock,omitempty"`  // the position (index) of the transaction within the block
	Version          int             `json:"version"`                    // the transaction version
	TypeInt          int             `json:"type_int"`                   // the transaction type as number
	Type             string          `json:"type"`                       // the transaction type as string
	Amount           decimal.Decimal `json:"amount,omitempty"`
	PropertyID       int             `json:"propertyid,omitempty"`
	PropertyName     string          `json:"propertyname,omitempty"`
	Block            int64           `json:"block,omitempty"`
	BlockHash        string          `json:"blockhash,omitempty"`
}

// OmniBalance 余额
type OmniBalance struct {
	Balance    decimal.Decimal `json:"balance"`              // the available balance of the address
	Reserved   decimal.Decimal `json:"reserved"`             // the amount reserved by sell offers and accepts
	Frozen     decimal.Decimal `json:"frozen,omitempty"`     // the amount frozen by the issuer (applies to managed properties only)
	PropertyID int             `json:"propertyid,omitempty"` // the property identifier
}

// OmniAllBalance 全部余额
type OmniAllBalance []OmniBalance

// OmniProperty 代币信息
type OmniProperty struct {
	PropertyID      int             `json:"propertyid"`      // the identifier
	Name            string          `json:"name"`            // the name of the tokens
	Category        string          `json:"category"`        // the category used for the tokens
	Subcategory     string          `json:"subcategory"`     // the subcategory used for the tokens
	Information     string          `json:"data"`            // additional information or a description
	URL             string          `json:"url"`             // an URI, for example pointing to a website
	Divisible       bool            `json:"divisible"`       // whether the tokens are divisible
	Issuer          ETHADDRESS      `json:"issuer"`          // the Bitcoin address of the issuer on record
	Hash            string          `json:"creationtxid"`    // the hex-encoded creation transaction hash
	FixeDissuance   bool            `json:"fixedissuance"`   // whether the token supply is fixed
	ManageDissuance bool            `json:"managedissuance"` // whether the token supply is managed by the issuer
	FreezingEnabled bool            `json:"freezingenabled"` // whether freezing is enabled for the property (managed properties only)
	TotalTokens     decimal.Decimal `json:"totaltokens"`     // the total number of tokens in existence
}
