package utils

import (
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	// "github.com/adrg/xdg"
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

func ExecCmd(cmdString string, verbose bool) error {
	r := csv.NewReader(strings.NewReader(cmdString))
	r.Comma = ' '
	cmdArr, err := r.Read()
	if err != nil {
		return err
	}

	cmd := exec.Command(cmdArr[0], cmdArr[1:]...)
	if verbose {
		fmt.Println("[#] " + strings.Join(cmdArr, " "))
		cmd.Stdout = os.Stdout
	}

	cmd.Stderr = os.Stderr
	// s.Start()
	err = cmd.Run()
	// s.Stop()
	return err
}

func GetIps(domain string) ([]string, error) {
	d := []string{}

	ips, err := net.LookupIP(domain)
	if err != nil {
		return d, err
	}

	for _, ip := range ips {

		if ip.To4() != nil {
			d = append(d, ip.String())
		}
	}

	return d, nil

}

// func GenerateRandomString(length int) string {
// 	const charset = "abcdefghijklmnopqrstuvwxyz" +
// 		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
//
// 	b := make([]byte, length)
// 	for i := range b {
// 		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
// 	}
// 	return string(b)
// }
