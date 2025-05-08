package test

import (
	"testing"

	"github.com/jeanhua/jokerhttp/engine"
)

func TestMapRedirect(t *testing.T) {
	joker := engine.NewEngine()
	joker.Init()
	joker.SetPort(1314)
	joker.MapRedirect("/baidu", "https://www.baidu.com")
	joker.Run()
}
