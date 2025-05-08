package test

import (
	"testing"

	"github.com/jeanhua/jokerhttp/engine"
)

func TestMapReverse(t *testing.T) {
	joker := engine.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	joker.MapReverseProxy("/", "https://imarket.jeanhua.cn")
	joker.Run()
}
