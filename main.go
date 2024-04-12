package main

import (
	"context"
	"flag"
	"time"

	"github.com/kloudlite/iot-devices/devices/common"
	"github.com/kloudlite/iot-devices/devices/hub"
	"github.com/kloudlite/iot-devices/devices/local"
	"github.com/kloudlite/iot-devices/pkg/logging"
	"github.com/kloudlite/iot-devices/utils"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "default", "--mode [local|hub|default]")
	flag.Parse()

	go common.StartPing()

	switch mode {
	case "local":
		if err := onlyLocal(); err != nil {
			println(err.Error())
		}
	case "hub":
		if err := onlyHub(); err != nil {
			println(err.Error())
		}
	default:
		if err := run(); err != nil {
			println(err.Error())
		}
	}
}

func onlyLocal() error {
	l, err := logging.New(&logging.Options{})

	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	println("Connection is unhealthy")
	if err := local.Run(ctx, l); err != nil {
		l.Errorf(err, "Error running local")
		return err
	}

	return nil
}

func onlyHub() error {
	l, err := logging.New(&logging.Options{})
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := hub.Run(ctx, l); err != nil {
		l.Errorf(err, "Error running hub")
		return err
	}

	return nil
}

func run() error {
	l, err := logging.New(&logging.Options{
		Name: "ik-app",
	})
	if err != nil {
		return err
	}

	_, cf := context.WithCancel(context.Background())

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
				l.Infof("Connection status changed to %v", o.IsConnected)
				o.cancel()
			}

			time.Sleep(5 * time.Second)
		}
	}(&obj)

	for {
		ctx, cf2 := context.WithCancel(context.Background())
		obj.cancel = cf2
		if obj.IsConnected {
			if err := hub.Run(ctx, l); err != nil {
				l.Errorf(err, "Error running hub")
				continue
			}
		}

		println("Connection is unhealthy")
		if err := local.Run(ctx, l); err != nil {
			l.Errorf(err, "Error running local")
			continue
		}
	}
}
