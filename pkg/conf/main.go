package conf

import (
	"os"

	"github.com/kloudlite/iot-devices/utils"
	"sigs.k8s.io/yaml"
)

type conf struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
}

func (w *conf) parseBytes(b []byte) error {
	return yaml.Unmarshal(b, w)
}

func (w *conf) ToBytes() ([]byte, error) {
	return yaml.Marshal(w)
}

func (w *conf) parseFile(path string) error {
	cp := GetConfigPath()

	if _, err := os.Stat(cp); os.IsNotExist(err) {
		return err
	}

	b, err := os.ReadFile(cp)

	if err != nil {
		return err
	}

	return w.parseBytes(b)
}

func (w *conf) ToFile(path string) error {
	cp := GetConfigPath()

	b, err := w.ToBytes()

	if err != nil {
		return err
	}

	return os.WriteFile(cp, b, 0644)
}

func GetConf() (*conf, error) {
	c := conf{}

	if err := c.parseFile(GetConfigPath()); err == nil {
		return &c, nil
	}

	pub, priv, err := utils.GenerateWgKeys()
	if err != nil {
		return nil, err
	}
	c = conf{
		PrivateKey: string(priv),
		PublicKey:  string(pub),
	}

	if err := c.ToFile(GetConfigPath()); err != nil {
		return nil, err
	}

	return &c, nil
}

func GetConfigPath() string {
	cp, ok := os.LookupEnv("CONFIG_PATH")
	if ok {
		return cp
	}

	return "/home/raspberry/kloudlite-conf.yaml"
}
