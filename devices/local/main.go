package local

import (
	"context"
	"slices"
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

func (h *hubstype) compare(b hubstype) bool {
	if len(*h) != len(b) {
		return false
	}

	for k, v := range *h {
		bv, ok := (b)[k]
		if !ok {
			return false
		}

		for k2, v2 := range v.domains {
			k3, ok := bv.domains[k2]
			if !ok {
				return false
			}

			for _, v3 := range v2 {
				if !slices.Contains(k3, v3) {
					return false
				}
			}
		}

	}

	return true
}

func (h *hubstype) GetHubs() hubstype {
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

	c.logger.Infof("Exiting local")

	return nil
}
