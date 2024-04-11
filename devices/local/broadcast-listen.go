package local

import (
	"fmt"
	"net"
	"time"

	"github.com/kloudlite/iot-devices/constants"
)

const (
	udpConnType = "udp"
)

func (c *client) listenBroadcast() error {

	c.logger.Infof("Listening for broadcast messages...")

	listenAddr, err := net.ResolveUDPAddr(udpConnType, fmt.Sprintf(":%d", constants.BroadcastPort))
	if err != nil {
		return err
	}

	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP(udpConnType, listenAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		if c.ctx.Err() != nil {
			return fmt.Errorf("Context cancelled")
		}

		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			c.logger.Errorf(err, "Error reading from UDP connection")
			continue
		}

		if addr.String() != localAddr.String() {
			if constants.IsDebug() {
				c.logger.Infof("Received message from %s: %s", addr, string(buffer[:n]))
			}

			hubs[addr.IP.String()] = hub{
				lastPing: time.Now(),
			}
		}
	}
}
