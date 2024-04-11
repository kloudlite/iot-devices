package local

import (
	"fmt"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/utils"
)

func (c *client) setRules() error {
	ips := hubs.GetHubs()

	defer c.removeRules(ips)

	var mark map[string]bool

	for hub, v := range ips {
		for _, ips := range v.domains {
			for _, ip := range ips {

				if mark[ip] {
					continue
				}

				if err := utils.ExecCmd(fmt.Sprintf("ip route add %s via %s", ip, hub), true); err != nil {
					c.logger.Errorf(err, "error adding ip route")
					continue
				}

				mark[ip] = true
			}
		}
	}

	for {

		m := hubs.GetHubs()
		if len(m) == 0 {
			c.logger.Infof("No rules to add")
		}

		if fmt.Sprintf("%#v", m) != fmt.Sprintf("%#v", ips) {
			c.logger.Infof("Rules changed, updating...")
			break
		}

		time.Sleep(constants.PingInterval * time.Second)
	}

	fmt.Println("exiting")
	return nil
}

func (c *client) removeRules(ips map[string]hb) {

	fmt.Println("Removing rules")

	var mark map[string]bool
	for _, v := range ips {
		for _, ips := range v.domains {
			for _, ip := range ips {
				if mark[ip] {
					continue
				}

				err := utils.ExecCmd(fmt.Sprintf("ip route delete %s", ip), true)
				if err != nil {
					c.logger.Errorf(err, "error deleting ip route")
					continue
				}

				mark[ip] = true
			}
		}
	}

}

func (c *client) ipTableRules() {
	for {
		c.setRules()
	}
}
