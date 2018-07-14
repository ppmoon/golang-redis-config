package grc_test

import (
	"testing"
	"grc"
	"fmt"
)
//测试grc
func TestNewGrc(t *testing.T) {
	g := grc.NewGrc(".env","localhost:6379","",0)
	go g.WatchFile()
	config := g.GetItem("ack")
	fmt.Println(config)
}
