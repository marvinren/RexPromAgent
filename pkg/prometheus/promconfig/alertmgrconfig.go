package promconfig

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

type AlertMgrConfig struct {
	Global    GlobalConfig  `yaml:"global"`
	Route     RootRoute     `yaml:"route"`
	Receivers []interface{} `yaml:"receivers"`
}

type GlobalConfig struct {
	Smtp_smarthost     string `yaml:"smtp_smarthost"`
	Smtp_from          string `yaml:"smtp_from"`
	Smtp_auth_username string `yaml:"smtp_auth_username"`
	Smtp_auth_password string `yaml:"smtp_auth_password"`
}

type RootRoute struct {
	Group_interval  string     `yaml:"group_interval"`
	Repeat_interval string     `yaml:"repeat_interval"`
	Group_by        []string   `yaml:"group_by"`
	Receiver        string     `yaml:"receiver"`
	Routes          []RouteEle `yaml:"routes"`
}

type RouteEle struct {
	Match    map[string]string `yaml:"match"`
	Receiver string            `yaml:"receiver"`
}

type ReceiverEmailEle struct {
	Name         string               `yaml:"name"`
	EmailConfigs []ReceiveEmailConfig `yaml:"email_configs"`
}

type ReceiveEmailConfig struct {
	To string `yaml:"to"`
}

type ReceiverWebhookEle struct {
	Name          string                 `yaml:"name"`
	WebhookConfig []ReceiveWebhookConfig `yaml:"webhook_configs"`
}

type ReceiveWebhookConfig struct {
	Url string `yaml:"url"`
}

func GenerateAlertMgrConfigFile(config *AlertMgrConfig, configFilePath string) error {

	data, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("parse the alert config error, %s\n", err)
		return err
	}
	err = ioutil.WriteFile(configFilePath, data, 0777)
	if err != nil {
		log.Printf("write the config file err, %s\n", err)
		return err
	}

	return nil
}

func GetAlertMgrManageAPIReload() error {
	promUrl := viper.GetString("prometheus.alertmgrUrl")
	reloadUrl := promUrl + "/-/reload"
	request, err := http.NewRequest(http.MethodPost, reloadUrl, bytes.NewReader(nil))
	if err != nil {
		logrus.Errorf("build reload (%v) request error: %v", reloadUrl, err)
		return err
	}
	client := http.Client{}
	response, err2 := client.Do(request)
	if err2 != nil {
		logrus.Errorf("request (%v) error: %v", reloadUrl, err2)
		return err2
	}
	if response.StatusCode != http.StatusOK {
		logrus.Errorf("reload request response error: %v", response.StatusCode)
	}
	return nil

}
