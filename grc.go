package grc

import (
	"github.com/go-redis/redis"
	"os"
	"log"
	"bufio"
	"io"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"strings"
)

func init() {
	readEnv()
}
//实例化一个GRC
func NewGrc(file string){
	WatchFile(file)
}
//链接redis的实例
func RedisClient() *redis.Client {
	return redis.NewClient(
		&redis.Options{
			Addr:"localhost:6379",
			Password:"",
			DB:0,
		},
	)
}
//save k-v config to redis
func saveConfig(key string,value string)  {
	client := RedisClient()
	err := client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
//逐行读取文件
func readEnv() {
	file,err := os.Open(".env")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	br := bufio.NewReader(file)
	for {
		line,_,err := br.ReadLine()
		if err == io.EOF {
			break
		}
		if 0 == len(line) || string(line) == "\r\n" {
			continue
		}

		k,v := parse(string(line))
		saveConfig(k,v)
	}
}
//解析env文件
func parse(line string) (string,string) {
	line = strings.Replace(line, " ", "", -1)
	line = strings.Replace(line, "\n", "", -1)
	result := strings.Split(line, "=")
	//fmt.Println(result[0],result[1])
	return result[0],result[1]
}
//监听文件变化
func WatchFile(filePath string) {
	watcher,err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("ERROR",err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Printf("EVENT! %#v\n", event)
				readEnv()
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
	if err := watcher.Add(filePath); err != nil {
		fmt.Println("ERROR",err)
	}
	<- done
}