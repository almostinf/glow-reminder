package main

import (
	"github.com/almostinf/glow-reminder/internal/app"
	"go.uber.org/fx"
)

func main() {
	fx.New(app.CreateApp()).Run()
}
