package local

import (
	"fmt"
	"net"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/devices/hub"
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
		select {
		case <-c.ctx.Done():
			return fmt.Errorf("Context cancelled")
		default:

			// Set a timeout on the read operation
			if err := conn.SetReadDeadline(time.Now().Add(constants.PingInterval * time.Second)); err != nil {
				c.logger.Errorf(err, "Error setting read deadline")
				continue
			}

			n, addr, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				c.logger.Errorf(err, "Error reading from UDP connection")
				continue
			}

			if addr.String() != localAddr.String() {

				var dm hub.Dms

				if err := dm.FromBytes(buffer[:n]); err != nil {
					c.logger.Errorf(err, "Error decoding message")
					continue
				}

				if constants.IsDebug() {
					c.logger.Infof("Received message from %s: %s", addr, string(buffer[:n]))
				}

				d := time.Now()
				hubs[addr.IP.String()] = hb{
					lastPing: &d,
					domains:  dm,
				}
			}
		}
	}
}
