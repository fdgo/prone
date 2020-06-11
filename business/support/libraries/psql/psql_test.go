package psql

import (
	"business/support/domain"
	"testing"
)

func Test_NewClient(t *testing.T) {
	conf := domain.DBConf{
		Host:     "192.168.5.141",
		Username: "root",
		Password: "Wans198059",
		DBName:   "exchange",
	}
	pool := NewClient(&conf, true)
	if err := pool.NewConn().Error; err != nil {
		t.Error(err)
	}
}
