package metrics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Config ...
type Config struct {
	Host        string
	Source      string
	User        string
	APIKey      string
	Environment string
}

// Publisher ...
type Publisher struct {
	Config     Config
	HTTPClient *http.Client
	Logger     *zap.Logger
}

var defaultPublisher = &Publisher{}

// New ...
func New(cfg Config, logger *zap.Logger) (*Publisher, error) {
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives:   true,
			MaxIdleConns:        64,
			MaxIdleConnsPerHost: 64,
			IdleConnTimeout:     time.Second * 30,
		},
		Timeout: time.Second * 5,
	}
	defaultPublisher = &Publisher{
		Config:     cfg,
		HTTPClient: client,
		Logger:     logger,
	}
	return defaultPublisher, nil
}

// Publish ...
func (p *Publisher) Publish(name string, value int) {
	payload := fmt.Sprintf(`[{"name": "%s", "value": %d, "interval": 10, "time": %d}]`, name, value, time.Now().Unix())
	fmt.Println("payload", payload)

	req, err := http.NewRequest(http.MethodPost, p.Config.Host, bytes.NewBufferString(payload))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s:%s", p.Config.User, p.Config.APIKey))

	resp, err := p.HTTPClient.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		fmt.Println(string(body))
		fmt.Println(resp.StatusCode)
		panic(err)
	}
	p.Logger.Sugar().Infof("metric sent %s %d", name, value)
}

// SendMetric ...
func SendMetric(name string, value int) {
	go defaultPublisher.Publish(name, value)
}
