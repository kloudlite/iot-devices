package hub

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/kloudlite/iot-devices/utils"
)

type Dms map[string][]string

func (d *Dms) ToBytes() ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(d); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (d *Dms) FromBytes(b []byte) error {
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	if err := dec.Decode(d); err != nil {
		return err
	}
	return nil
}

func (c *client) resyncDomains() {
	for {
		nd := map[string][]string{}
		nd["ips"] = c.ctx.GetExposedIps()

		for _, v := range c.ctx.GetDomains() {
			if v == "ips" {
				continue
			}

			s, err := utils.GetIps(v)
			if err != nil {
				c.logger.Errorf(err, "Error getting ips for domain %s", v)
			}
			nd[v] = s
		}

		if fmt.Sprintf("%#v", nd) != fmt.Sprintf("%#v", c.domains) {
			c.domains = nd
		}

		time.Sleep(5 * time.Second)
	}

}
