package types

import (
	"context"
	"encoding/json"
	"slices"

	"github.com/kloudlite/api/pkg/logging"
	"github.com/kloudlite/iot-devices/constants"
)

type Response struct {
	AccountName       string   `json:"accountName"`
	CreationTime      string   `json:"creationTime"`
	DeploymentName    string   `json:"deploymentName"`
	DisplayName       string   `json:"displayName"`
	ID                string   `json:"id"`
	IP                string   `json:"ip"`
	MarkedForDeletion string   `json:"markedForDeletion"`
	Name              string   `json:"name"`
	PodCIDR           string   `json:"podCIDR"`
	ProjectName       string   `json:"projectName"`
	PublicKey         string   `json:"publicKey"`
	RecordVersion     int      `json:"recordVersion"`
	ServiceCIDR       string   `json:"serviceCIDR"`
	UpdateTime        string   `json:"updateTime"`
	Version           string   `json:"version"`
	ExposedDomains    []string `json:"exposedDomains"`
	ExposedIPs        []string `json:"exposedIps"`
	ClusterToken      string   `json:"clusterToken"`

	Reset bool `json:"reset"`
}

func (c *Response) FromJson(data []byte) error {
	return json.Unmarshal(data, c)
}

type MainCtx interface {
	GetDomains() []string
	UpdateDomains([]string)
	GetDevice() (*Response, error)
	UpdateDevice(*Response)

	GetLogger() logging.Logger
	GetContext() context.Context
	SetContext(context.Context)

	GetContextWithCancel() (context.Context, context.CancelFunc)

	SetExposedIps([]string)
	GetExposedIps() []string
}

type mainCtx struct {
	ips     []string
	domains []string
	device  *Response
	logger  logging.Logger
	ctx     context.Context
}

func (m *mainCtx) GetExposedIps() []string {
	return m.ips
}

func (m *mainCtx) SetExposedIps(ips []string) {
	m.ips = ips

	for _, ip := range constants.DefaultExposedIps {
		if !slices.Contains(m.ips, ip) {
			m.ips = append(m.ips, ip)
		}
	}
}

func (m *mainCtx) GetDomains() []string {
	return m.domains
}

func (m *mainCtx) UpdateDomains(domains []string) {
	m.domains = domains
	for _, d := range constants.DefaultExposedDomains {
		if !slices.Contains(m.domains, d) {
			m.domains = append(m.domains, d)
		}
	}

}

func (m *mainCtx) GetDevice() (*Response, error) {
	return m.device, nil
}

func (m *mainCtx) UpdateDevice(device *Response) {
	m.device = device
}

func (m *mainCtx) GetLogger() logging.Logger {
	return m.logger
}

func (m *mainCtx) GetContext() context.Context {
	return m.ctx
}
func (m *mainCtx) SetContext(ctx context.Context) {
	m.ctx = ctx
}

func (m *mainCtx) GetContextWithCancel() (context.Context, context.CancelFunc) {
	ctx, cf := context.WithCancel(m.ctx)
	m.ctx = ctx

	return ctx, cf
}

func NewMainCtx(domains []string) (MainCtx, error) {
	l, err := logging.New(&logging.Options{
		Name: constants.AppName,
	})
	if err != nil {
		return nil, err
	}

	return &mainCtx{
		domains: domains,
		logger:  l,
	}, nil

}

func NewMainCtxOrDie(domains []string) MainCtx {

	l, err := logging.New(&logging.Options{
		Name: constants.AppName,
		Dev:  true,
	})
	if err != nil {
		panic(err)
	}

	return &mainCtx{
		domains: domains,
		logger:  l,
		ctx:     context.Background(),
	}
}
