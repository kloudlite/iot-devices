package hub

import "time"

var ()

func (c *client) setRules() error {
	return nil
}

func (c *client) handleIpTableRules() {
	for {
		if err := c.setRules(); err != nil {
			c.logger.Errorf(err, "Error setting rules")
		}

		time.Sleep(1 * time.Second)
	}
}
