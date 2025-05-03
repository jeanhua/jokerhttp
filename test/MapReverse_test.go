package test

import (
	"testing"

	"github.com/jeanhua/jokerhttp"
)

func TestMapReverse(t *testing.T) {
	joker := jokerhttp.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	joker.MapReverseProxy("/", "https://imarket.jeanhua.cn")
	joker.Run()
}
