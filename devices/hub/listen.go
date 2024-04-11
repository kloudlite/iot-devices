package hub

import (
	"fmt"
	"time"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/utils"
)

var prevDomains = []string{}

func (c *client) listenProxy() error {
	err := caddy.Run(&caddy.Config{})
	if err != nil {
		return err
	}
	defer caddy.Stop()

	domains := utils.GetDomains()
	for {
		if c.ctx.Err() != nil {
			return fmt.Errorf("context cancelled")
		}

		if fmt.Sprintf("%#v", domains) == fmt.Sprintf("%#v", prevDomains) {
			time.Sleep(5 * time.Second)
			continue
		}

		prevDomains = domains

		func() {
			caddyfileConfig := ``

			for _, d := range domains {
				caddyfileConfig += fmt.Sprintf(`
reverse_proxy https://%s {
  transport http {
    tls
    tls_insecure_skip_verify
  }
}
`, d)
			}

			caddyfileConfig = fmt.Sprintf(`
{
  auto_https off
}

:%d {
        %s
}
`, constants.ProxyServerPort, caddyfileConfig)

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

var (
	ipMap = make(map[string]string)
)
