package constants

import "fmt"

const (
	AppName = "kl-iot-devices"

	IotServerEndpoint = "iotnet.dev.kloudlite.io"
	DnsDomain         = "one.one.one.one"

	BroadcastIP   = "255.255.255.255"
	BroadcastPort = 12345

	BroadcastInterval = 4

	ExpireDuration = 10

	ProxyServerPort = 8000

	PingInterval = 5
	PingTimeout  = 3

	K3sConfigPath = "/home/raspberry/runner-config.yml"
	K3sDataPath   = "/var/lib/rancher/k3s/server/db/"

	K3sJobFile = "/tmp/kloudlite-k3s-job.yml"
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
	return fmt.Sprintf("https://%s/device", IotServerEndpoint)
}

func GetIotServerEndpoint() string {
	return fmt.Sprintf("https://%s", IotServerEndpoint)
}

var DefaultExposedDomains = []string{IotServerEndpoint, DnsDomain, "github.com", "kloudlite.github.io"}
var DefaultExposedIps = []string{}
