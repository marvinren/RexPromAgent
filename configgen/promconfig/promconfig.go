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

type PromConfigRules struct {
	Groups []AlertGroup `yaml:"groups"`
}

type AlertGroup struct {
	Name  string      `yaml:"name" json:"name"`
	Rules []AlertRule `yaml:"rules" json:"rules"`
}

type AlertRule struct {
	Alert       string            `yaml:"alert" json:"alert"`
	Expr        string            `yaml:"expr" json:"expr"`
	For         string            `yaml:"for" json:"for"`
	Labels      map[string]string `yaml:"labels" json:"labels"`
	Annotations AlertAnnotations  `yaml:"annotations" json:"annotations"`
}

type AlertAnnotations struct {
	Summary     string `yaml:"summary" json:"summary"`
	Description string `yaml:"description" json:"description"`
}

func GeneratePromRuleFile(config *PromConfigRules, config_file_path string) error {

	data, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("parse the alert rule error, %s\n", err)
		return err
	}
	err = ioutil.WriteFile(config_file_path, data, 0777)
	if err != nil {
		log.Printf("write the config file err, %s\n", err)
		return err
	}

	return nil
}

func GetPromManageAPIReload() error {
	promUrl := viper.GetString("prometheus.prometheusUrl")
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
