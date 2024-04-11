package utils

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func IsConn() bool {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(constants.GetHealthyUrl())
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true
	}

	fmt.Println("[ERROR] Connection is unhealthy")

	return false
}

func GenerateWgKeys() ([]byte, []byte, error) {
	key, err := wgtypes.GenerateKey()
	if err != nil {
		return nil, nil, err
	}

	return []byte(key.PublicKey().String()), []byte(key.String()), nil
}

func GetDomains() []string {
	return []string{
		constants.IotServerEndpoint,
	}
}
