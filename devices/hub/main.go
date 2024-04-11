package hub

import (
	"context"

	"github.com/kloudlite/iot-devices/pkg/logging"
)

type client struct {
	ctx    context.Context
	logger logging.Logger
}

func Run(ctx context.Context, logger logging.Logger) error {
	c := &client{
		ctx:    ctx,
		logger: logger,
	}

	c.logger.Infof("Starting hub")

	go c.resyncDomains()

	if err := c.selfBroadcast(); err != nil {
		c.logger.Errorf(err, "Error broadcasting message")

		// TODO: handle error
		return err
	}
	return nil
}
