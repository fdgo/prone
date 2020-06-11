package psql

import (
	"business/support/domain"
	"fmt"
	"strings"
	"time"

	"errors"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Client struct {
	DB *gorm.DB
}

// NewClient Create a PostgresSQL client, retry if failed
func NewClient(conf *domain.DBConf, isDebug bool) *Client {
	var err error
	var db *gorm.DB

	if conf.Host == "" || conf.Username == "" || conf.DBName == "" {
		panic(errors.New("DBConf Error"))
	}

	for {
		db, err = gorm.Open("postgres", fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s",
			conf.Host,
			conf.Port,
			conf.Username,
			conf.DBName,
			conf.Password,
		))
		if err != nil {
			fmt.Println("Failed to connect database", err, conf)
			time.Sleep(2 * time.Second)
			fmt.Println("Retrying to connect to postgres...")
		} else {
			fmt.Println("Connected to Postgres Server:", conf.Host, conf.DBName)
			break
		}
	}

	if isDebug {
		db.LogMode(true)
	}
	if conf.MaxIdleConns == 0 {
		conf.MaxIdleConns = 4
	}
	if conf.MaxOpenConns == 0 {
		conf.MaxOpenConns = 8
	}
	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)

	return &Client{db}
}

// NewConn Fetch a db connection from pool
func (self *Client) NewConn() *gorm.DB {
	db := self.DB.New()
	db.BlockGlobalUpdate(true)
	return db
}

func IsDuplicateKeyError(err error) bool {
	if nil == err {
		return false
	}
	// duplicate key
	if strings.Contains(err.Error(), "duplicate key") {
		return true
	}
	return false
}
