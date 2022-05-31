package promconfig

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"testing"
)

func TestAlertMgrConfigGenerate(t *testing.T) {

	var config = &AlertMgrConfig{
		Global: GlobalConfig{
			"mail.yusys.com.cn",
			"yusys_aiops@yusys.com",
			"aiops",
			"aiops",
		},
		Route: RootRoute{
			"5m",
			"1h",
			[]string{"alertname"},
			"renzq2",
			[]RouteEle{
				{
					map[string]string{"alertname": "dev1"},
					"renzq2",
				},
			},
		},
		Receivers: []interface{}{
			ReceiverEmailEle{"renzq2",
				ReceiveEmailConfig{
					[]string{"renzq2@yusys.com.cn"},
				},
			},
			ReceiverWebhookEle{
				"webhook",
				ReceiveWebhookConfig{
					"http://localhost:5001",
				},
			},
		},
	}

	err := GenerateAlertMgrConfigFile(config, "./alertmgr.yml")
	if err != nil {
		t.Errorf("generate the alertmgr file error: %v", err)
	}

	ymldata, err := ioutil.ReadFile("./alertmgr.yml")
	if err != nil {
		t.Errorf("read result yaml file, error: %v", err)
	}

	result := AlertMgrConfig{}
	err = yaml.Unmarshal(ymldata, &result)
	if err != nil {
		t.Errorf("unmarshal the yaml string error: %v", err)
	}

	if len(result.Receivers) <= 0 {
		t.Errorf("there is no receiver.")
	}

}

func TestAlertMgrManageAPIReload(t *testing.T) {

	err := GetAlertMgrManageAPIReload()
	if err != nil {
		t.Errorf("reload alertmgr config error: %v", err)
	}
}
