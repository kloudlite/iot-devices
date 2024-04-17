package k3s

import (
	"fmt"
	"os"
	"sync"

	"github.com/kloudlite/api/pkg/logging"
	"github.com/kloudlite/iot-devices/constants"
	"github.com/kloudlite/iot-devices/templates"
	"github.com/kloudlite/iot-devices/types"

	"github.com/kloudlite/iot-devices/utils"
)

const k3sService = "kloudlite-k3s"

type client struct {
	l logging.Logger
}

func New(ctx types.MainCtx) *client {
	return &client{
		l: ctx.GetLogger(),
	}
}

func (c *client) Start() error {
	return utils.ExecCmd(fmt.Sprintf("systemctl start %s", k3sService), true)
}

func (c *client) Stop() error {
	return utils.ExecCmd(fmt.Sprintf("systemctl stop %s", k3sService), true)
}

func (c *client) Restart() error {
	return utils.ExecCmd(fmt.Sprintf("systemctl restart %s", k3sService), true)
}

func (c *client) Reset() error {
	if err := c.Stop(); err != nil {
		return err
	}

	if err := os.RemoveAll(constants.K3sConfigPath); err != nil {
		return err
	}

	if err := os.RemoveAll(constants.K3sDataPath); err != nil {
		return err
	}

	return c.Start()
}

var (
	mu sync.Mutex
)

func (c *client) write(cf string) error {
	if err := os.WriteFile(constants.K3sConfigPath, []byte(cf), 0644); err != nil {
		return err
	}

	c.Stop()
	return c.Start()
}

func (c *client) UpsertConfig(cf string) error {
	mu.Lock()
	defer mu.Unlock()

	b, err := os.ReadFile(constants.K3sConfigPath)
	if err != nil {
		c.l.Errorf(err, "error reading k3s config")
		return c.write(cf)
	}

	if string(b) == cf {
		return nil
	}

	return c.write(cf)
}

func (c *client) ApplyInstallJob(obj map[string]any) error {

	update := false
	if err := utils.ExecCmd("k3s kubectl get deployments/kl-agent -n kloudlite", true); err != nil {
		c.l.Errorf(err, "error getting kl-agent deployment")
		update = true
	}

	if !update {
		if err := utils.ExecCmd("k3s kubectl get deployments/kl-agent-operator -n kloudlite", true); err != nil {
			c.l.Errorf(err, "error getting kl-agent-operator deployment")
			update = true
		}
	}

	if !update {
		c.l.Infof("kloudlite agent already installed")
		return nil
	}

	utils.ExecCmd("k3s kubectl delete job/helm-job-kloudlite-agent -n kloudlite", true)

	b, err := templates.ParseTemplate(templates.AgentInstallJob, obj)
	if err != nil {
		return err
	}

	if err := os.WriteFile(constants.K3sJobFile, b, 0644); err != nil {
		return err
	}

	if err := utils.ExecCmd(fmt.Sprintf("k3s kubectl apply -f %s", constants.K3sJobFile), true); err != nil {
		return err
	}

	return nil
}
