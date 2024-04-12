package local

import (
	"fmt"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/pkg/networkmanager"
)

func (c *client) setRules() error {
	ips := hubs.GetHubs()

	mark := map[string]int{}

	if err := networkmanager.AddExternalDns(); err != nil {
		c.logger.Errorf(err, "error adding external dns")
	}

	defer func() {
		if err := networkmanager.CleanExternalDns(); err != nil {
			c.logger.Errorf(err, "error cleaning external dns")
		}
	}()

	for hub, v := range ips {
		for _, ips := range v.domains {
			for _, ip := range ips {

				met := mark[ip] + 1
				// if err := utils.ExecCmd(fmt.Sprintf("ip route add %s via %s metric %d", ip, hub, met), true); err != nil {
				// 	c.logger.Errorf(err, "error adding ip route")
				// 	continue
				// }

				if err := networkmanager.AddRoute(fmt.Sprintf("%s/32", ip), hub, met); err != nil {
					c.logger.Errorf(err, "error adding route")
					continue
				}

				mark[ip] = met
			}
		}
	}

	defer c.removeRules(ips)

	for {
		select {
		case <-c.ctx.Done():
			return fmt.Errorf("context cancelled")
		default:
			m := hubs.GetHubs()
			if len(m) == 0 {
				c.logger.Infof("No rules to add")
			}

			if !m.compare(ips) {
				c.logger.Infof("Rules changed, updating...")
				return nil
			} else {
				c.logger.Infof("Rules not changed")
			}

			time.Sleep(constants.PingInterval * time.Second)
		}
	}

}

func (c *client) removeRules(ips map[string]hb) {

	mark := map[string]bool{}
	for _, v := range ips {
		for _, ips := range v.domains {
			for _, ip := range ips {
				if mark[ip] {
					continue
				}

				s, err := networkmanager.ListRoutes(fmt.Sprintf("%s/32", ip))
				if err != nil {
					c.logger.Errorf(err, "error listing routes")
					continue
				}

				for _, r := range s {
					if err := networkmanager.DeleteRoute(r.Dst.String(), r.Gw.String()); err != nil {
						c.logger.Errorf(err, "error deleting route")
						continue
					}
				}

				// err := utils.ExecCmd(fmt.Sprintf("ip route delete %s", ip), true)
				// if err != nil {
				// 	c.logger.Errorf(err, "error deleting ip route")
				// 	continue
				// }

				mark[ip] = true
			}
		}
	}

}

func (c *client) ipTableRules() error {
	for {
		if c.ctx.Err() != nil {
			return fmt.Errorf("context cancelled")
		}
		c.setRules()
	}
}
