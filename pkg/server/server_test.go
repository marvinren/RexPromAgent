package server

import (
	"RexPromAgent/pkg/config"
	"RexPromAgent/pkg/log"
	"RexPromAgent/pkg/prometheus/promconfig"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestAlertSync(t *testing.T) {

	logrus.Info("Test the alert sync api.......")

	req, err := http.NewRequest("GET", "/rules/alert/reload", nil)
	if err != nil {
		t.Fatal(err)
	}

	s := NewServer()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(s.SyncAlertRules)

	handler.ServeHTTP(rr, req)

	// check the http status
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// check the generated yaml file
	configPath := viper.GetString("prometheus.alertRulesConfigPath")
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		t.Errorf("yml file read error, maybe the file didn't be generated. error: %s", err)
	}
	t.Log(string(file))
	c := promconfig.PromConfigRules{}
	err = yaml.Unmarshal(file, &c)
	if err != nil {
		t.Errorf("yml file format error , %s", err)
	}
	fmt.Println(len(c.Groups))

}

func TestMain(m *testing.M) {
	// Setup: Initialize
	config.Initialize()
	log.Initialize()

	// Run the test cases
	m.Run()

	// TearDown: remove junk file
	err := os.RemoveAll(viper.GetString("prometheus.alertRulesConfigPath"))
	if err != nil {
		fmt.Printf("delete alert config file error, %s\n", err)
	}
}
