package local

import (
	"fmt"
	"strings"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/kloudlite/iot-devices/constants"

	_ "github.com/caddyserver/caddy/v2/modules/standard"
)

var (
	currentTargetIndex int32 = 0
)

var prevIps = []string{}

func (c *client) listenProxy() error {
	err := caddy.Run(&caddy.Config{})
	if err != nil {
		return err
	}

	defer caddy.Stop()

	for {
		if c.ctx.Err() != nil {
			return fmt.Errorf("context cancelled")
		}

		hbs := hubs.GetHubs()

		func() {

			if len(hbs) == 0 {
				c.logger.Infof("No hubs found")
				return
			}

			if fmt.Sprintf("%#v", hbs) == fmt.Sprintf("%#v", prevIps) {
				return
			}

			caddyfileConfig := fmt.Sprintf(`
* {
  reverse_proxy %s
}
			`, strings.Join(func() []string {

				var targets []string

				for _, h := range hbs {
					targets = append(targets, fmt.Sprintf("http://%s:%d", h, constants.ProxyServerPort))
				}

				return targets

			}(), " "))

			fmt.Println(caddyfileConfig)

			adapter := caddyconfig.GetAdapter("caddyfile")
			if adapter == nil {
				c.logger.Errorf(nil, "Caddyfile adapter not found")
				return
			}
			caddyConfig, _, err := adapter.Adapt([]byte(caddyfileConfig), nil)
			if err != nil {
				c.logger.Errorf(err, "Failed to adapt Caddyfile")
				return
			}

			if err := caddy.Load(caddyConfig, true); err != nil {
				c.logger.Errorf(err, "Failed to load Caddy config")
			}
		}()

		time.Sleep(5 * time.Second)
	}
}
