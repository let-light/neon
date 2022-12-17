/*
Copyright © 2022 im-pingo HERE <cczjp89@gmail>

*/
package main

import (
	"fmt"
	"time"

	"github.com/pingopenstack/neon/pkg/module"
	_ "github.com/pingopenstack/neon/src/core"

	_ "github.com/pingopenstack/neon/src/modules/rtsp"
	_ "github.com/pingopenstack/neon/src/modules/webrtc"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	module.Launch()

	for {
		time.Sleep(time.Second * 1)
	}
}
