package grc_test

import (
	"testing"
	"grc"
	"log"
	"fmt"
	"github.com/go-redis/redis"
)
//测试grc
func TestNewGrc(t *testing.T) {
	go grc.NewGrc(".env")
	client := grc.RedisClient()
	val, err := client.Get("ack").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(val)
}

//测试redis客户端
func TestRedisClient(t *testing.T) {
	client := grc.RedisClient()
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalln("数据库连接失败")
	}
}
//测试redis操作
func TestRedisOperate(t *testing.T) {
	client := grc.RedisClient()
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := client.Get("key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
}
//测试文件监听
func TestWatchFile(t *testing.T) {
	grc.WatchFile(".env")
}