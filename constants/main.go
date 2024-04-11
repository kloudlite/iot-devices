package constants

import "fmt"

const (
	IotServerEndpoint = "iotnet.dev.kloudlite.io"

	BroadcastIP   = "255.255.255.255"
	BroadcastPort = 12345

	BroadcastInterval = 4

	ExpireDuration = 10

	ProxyServerPort = 8000

	PingInterval = 5
	PingTimeout  = 3
)

var (
	Debug = "false"
)

func IsDebug() bool {
	return Debug == "true"
}

func GetHealthyUrl() string {
	return fmt.Sprintf("https://%s/healthy", IotServerEndpoint)
}

func GetPingUrl() string {
	return fmt.Sprintf("https://%s/healthy", IotServerEndpoint)
}

func GetIotServerEndpoint() string {
	return fmt.Sprintf("https://%s", IotServerEndpoint)
}
