package main

import (
	"context"
	"flag"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/devices/common"
	"github.com/kloudlite/iot-devices/devices/hub"
	"github.com/kloudlite/iot-devices/devices/local"
	"github.com/kloudlite/iot-devices/types"
	"github.com/kloudlite/iot-devices/utils"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "default", "--mode [local|hub|default]")
	flag.Parse()

	mctx := types.NewMainCtxOrDie(constants.DefaultExposedDomains)

	switch mode {
	case "local":
		if err := onlyLocal(mctx); err != nil {
			println(err.Error())
		}
	case "hub":
		if err := onlyHub(mctx); err != nil {
			println(err.Error())
		}
	default:
		if err := run(mctx); err != nil {
			println(err.Error())
		}
	}
}

func onlyLocal(ctx types.MainCtx) error {

	go common.StartPing(ctx)

	ctx.GetLogger().Infof("Connection is unhealthy")
	if err := local.Run(ctx); err != nil {
		ctx.GetLogger().Errorf(err, "Error running local")
		return err
	}

	return nil
}

func onlyHub(ctx types.MainCtx) error {

	go common.StartPing(ctx)

	if err := hub.Run(ctx); err != nil {
		return err
	}

	return nil
}

func run(ctx types.MainCtx) error {

	go common.StartPing(ctx)

	_, cf := ctx.GetContextWithCancel()

	var obj = struct {
		IsConnected bool
		cancel      context.CancelFunc
	}{
		IsConnected: utils.IsConn(),
		cancel:      cf,
	}

	go func(o *struct {
		IsConnected bool
		cancel      context.CancelFunc
	}) {
		for {
			ic := utils.IsConn()
			if o.IsConnected != ic {
				o.IsConnected = ic
				ctx.GetLogger().Infof("Connection status changed to %v", o.IsConnected)

				o.cancel()
			}

			time.Sleep(5 * time.Second)
		}
	}(&obj)

	for {
		ctx.SetContext(context.Background())
		_, cf := ctx.GetContextWithCancel()

		obj.cancel = cf
		if obj.IsConnected {
			if err := hub.Run(ctx); err != nil {
				ctx.GetLogger().Errorf(err, "Error running hub")
				continue
			}
		}

		println("Connection is unhealthy")
		if err := local.Run(ctx); err != nil {
			ctx.GetLogger().Errorf(err, "Error running local")
			continue
		}
	}
}
