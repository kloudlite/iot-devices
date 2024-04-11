package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/pkg/conf"
	// "github.com/kloudlite/iot-devices/pkg/logging"
)

func StartPing() {

	// l, err := logging.New(&logging.Options{})
	// if err != nil {
	// 	panic(err)
	// }

	for {
		if err := ping(); err != nil {
			// l.Errorf(err, "sending ping to server")
		}

		time.Sleep(constants.PingInterval * time.Second)
	}
}

func ping() error {
	client := &http.Client{
		Timeout: constants.PingTimeout * time.Second,
	}

	c, err := conf.GetConf()
	if err != nil {
		return err
	}

	var data = struct {
		PublicKey string `json:"public_key"`
	}{
		PublicKey: c.PublicKey,
	}

	dataBytes, err := json.Marshal(data)

	if err != nil {
		return err
	}

	resp, err := client.Post(constants.GetPingUrl(), "application/json", bytes.NewBuffer(dataBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return fmt.Errorf("status code: %d", resp.StatusCode)
}
