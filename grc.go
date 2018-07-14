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
type grconfig struct {
	filePath string
	addr string
	password string
	db int
}
func (g *grconfig) init() {
	g.readEnv()
}
//实例化一个GRC
func NewGrc(file string,addr string,password string,db int) *grconfig {
	return &grconfig{
		filePath:file,
		addr:addr,
		password:password,
		db:db,
	}
}
//save k-v config to redis
func (g *grconfig) saveConfig(key string,value string)  {
	client := redis.NewClient(
		&redis.Options{
			Addr:g.addr,
			Password:g.password,
			DB:g.db,
		},
	)
	err := client.Set(key, value, 0).Err()
	if err != nil {
		panic(err)
	}
}
//读取配置字段
func (g *grconfig) GetItem(key string) string{
	client := redis.NewClient(
		&redis.Options{
			Addr:g.addr,
			Password:g.password,
			DB:g.db,
		},
	)
	val, err := client.Get(key).Result()
	if err != nil {
		panic(err)
	}
	return val
}
//逐行读取文件
func (g *grconfig) readEnv() {
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
		g.saveConfig(k,v)
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
func (g *grconfig) WatchFile() {
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
				g.readEnv()
				// watch for errors
			case err := <-watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()
	if err := watcher.Add(g.filePath); err != nil {
		fmt.Println("ERROR",err)
	}
	<- done
}