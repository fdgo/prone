package sequence

import (
	"business/support/libraries/redis"
	"fmt"
	"testing"
)

func Test_GenSequence(t *testing.T) {
	conf := redis.RedisConfig{
		Host: "localhost",
		Port: 6379,
	}
	redisPool := redis.NewPool(&conf)
	Init(redisPool)

	for i := 0; i < 10; i++ {
		id, err := GenSequence()
		if err != nil {
			t.Error(err)
		}
		fmt.Println("Seq:", id)
	}
}
