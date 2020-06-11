package redis

import (
	"business/support/domain"
	"testing"
	"time"
)

func Test_NewPool(t *testing.T) {
	conf := domain.DBConf{
		Host:   "127.0.0.1",
		Port:   6379,
		DBName: "0",
	}
	r := NewPool(&conf)

	key := "UNITTEST"
	expire := time.Second * 10
	conn := r.NewConn()
	if err := conn.Set(key, 1, expire).Err(); err != nil {
		t.Error(err)
	}
	if err := conn.Get(key).Err(); err != nil {
		t.Error(err)
	}
}
