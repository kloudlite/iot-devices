package hub

import (
	"fmt"
	"net"
	"time"

	"github.com/kloudlite/iot-devices/constants"
)

const udpConnType = "udp"

func (c *client) hc(conn *net.UDPConn) error {
	defer func() {
		time.Sleep(constants.PingInterval * time.Second)
	}()

	d := GetDomains()
	message, err := d.ToBytes()
	if err != nil {
		return err
	}

	if c.ctx.Err() != nil {
		return fmt.Errorf("Context cancelled")
	}

	_, err = conn.Write(message)
	if err != nil {
		return err
	}

	return nil
}

func (c *client) selfBroadcast() error {
	c.logger.Infof("Broadcasting message...")

	broadcastAddr, err := net.ResolveUDPAddr(udpConnType, fmt.Sprintf("%s:%d", constants.BroadcastIP, constants.BroadcastPort))
	if err != nil {
		return err
	}

	localAddr, err := net.ResolveUDPAddr(udpConnType, ":0")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP(udpConnType, localAddr, broadcastAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	for {
		if err := c.hc(conn); err != nil {
			c.logger.Errorf(err, "Error sending broadcast message")
		}
	}
}
