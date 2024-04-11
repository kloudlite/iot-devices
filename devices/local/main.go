package local

import (
	"context"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/devices/hub"
	"github.com/kloudlite/iot-devices/pkg/logging"
)

type hb struct {
	lastPing *time.Time
	domains  hub.Dms
}

type hubstype map[string]hb

func (h *hubstype) cleanup() {
	for k, v := range *h {
		if v.lastPing != nil && time.Since(*v.lastPing) > constants.ExpireDuration*time.Second {
			delete(*h, k)
		}
	}
}

func (h *hubstype) GetHubs() map[string]hb {
	h.cleanup()

	d := map[string]hb{}
	for k, v := range *h {
		v.lastPing = nil
		d[k] = v
	}

	return d
}

var hubs = hubstype{}

type client struct {
	logger logging.Logger
	ctx    context.Context
}

func Run(ctx context.Context, logger logging.Logger) error {
	c := &client{
		logger: logger,
		ctx:    ctx,
	}

	c.logger.Infof("Starting local")

	go func() {
		c.ipTableRules()
	}()

	if err := c.listenBroadcast(); err != nil {
		return err
	}

	return nil
}
