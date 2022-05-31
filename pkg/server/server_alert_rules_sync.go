package server

import (
	"RexPromAgent/pkg/db"
	"RexPromAgent/pkg/prometheus/promconfig"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"strconv"
	"unicode"
)

func (s Server) SyncAlertRules(w http.ResponseWriter, r *http.Request) {
	// Lock the sync lock
	s.syncLocker.Lock()
	defer s.syncLocker.Unlock()
	// Get the prometheus alert rule config file location.
	alertRuleConfigPath := viper.GetString("prometheus.alertRulesConfigPath")
	alertMgrConfigPath := viper.GetString("prometheus.alertMgrConfigPath")
	// Get the alert rules from the database.
	rules := make([]db.AlertRule, 0)
	err := s.db.FetchAlerts(&rules)
	if err != nil {
		logrus.Errorf("err for get alert data: %s", err)
	}

	// Convert the database object to the prometheus's rule config object.
	alertRules := make([]promconfig.AlertRule, 0)
	alertRoutes := make([]promconfig.RouteEle, 0)
	alertReceives := make([]interface{}, 0)

	for _, alert := range rules {
		forDurationTime := alert.Duration
		if isDigit(forDurationTime) {
			forDurationTime = forDurationTime + "m"
		}

		// generate alert rule
		alertRule := promconfig.AlertRule{
			Alert: alert.AlertName,
			For:   forDurationTime,
			Expr:  alert.Expression,
			Labels: map[string]string{
				"severity": alert.AlertLevel,
				"category": alert.AlertType,
				"rule_id":  strconv.FormatInt(alert.RuleId, 10),
			},
			Annotations: promconfig.AlertAnnotations{
				Summary:     alert.Notice,
				Description: alert.Description,
			},
		}
		alertRules = append(alertRules, alertRule)

		// generate alert route & alert receive, if receiver is not nil.
		if alert.Receiver.Valid {
			alertRoute := promconfig.RouteEle{
				map[string]string{"alertname": alert.AlertName},
				alert.AlertName + "-receiver",
			}
			alertRoutes = append(alertRoutes, alertRoute)

			alertReceive := promconfig.ReceiverEmailEle{
				alert.AlertName + "-receiver",
				[]promconfig.ReceiveEmailConfig{
					{alert.Receiver.String},
				},
			}
			alertReceives = append(alertReceives, alertReceive)
		}
	}
	// the root alert rules config.
	alertRuleObj := promconfig.PromConfigRules{
		[]promconfig.AlertGroup{promconfig.AlertGroup{
			"group_auto",
			alertRules,
		}},
	}
	// the root alert manager config.
	alertReceives = append(alertReceives, promconfig.ReceiverWebhookEle{
		"web.hook",
		[]promconfig.ReceiveWebhookConfig{
			{"http://localhost:5001"},
		},
	})
	var alertMgrObj = promconfig.AlertMgrConfig{
		Global: promconfig.GlobalConfig{
			Smtp_smarthost:     "mail.yusys.com.cn:25",
			Smtp_from:          "yucc-aiops-admin@yusys.com.cn",
			Smtp_auth_username: "yucc-aiops-admin",
			Smtp_auth_password: "Aiops123",
		},
		Route: promconfig.RootRoute{
			Group_interval:  "10m",
			Repeat_interval: "1h",
			Group_by:        []string{"alertname"},
			Receiver:        "web.hook",
			Routes:          alertRoutes,
		},
		Receivers: alertReceives,
	}

	// Generate the rule file from the prometheus's rule config object
	err = promconfig.GeneratePromRuleFile(&alertRuleObj, alertRuleConfigPath)
	if err != nil {
		logrus.Errorf("generate the alert rules error, %s", err)
	}
	logrus.Infof("generate the config file %v", alertRuleConfigPath)

	// Generate the alert Manager config file
	err = promconfig.GenerateAlertMgrConfigFile(&alertMgrObj, alertMgrConfigPath)
	if err != nil {
		logrus.Errorf("generate the alert manager error, %s", err)
	}
	logrus.Infof("generate the config file %v", alertMgrConfigPath)

	// Reload the prometheus configuration
	err = promconfig.GetPromManageAPIReload()
	if err != nil {
		logrus.Errorf("reload Prometheus configuration error, %s", err)
	}
	logrus.Info("request the prometheus reload url.")

	// Reload the alertmanager configuration
	err = promconfig.GetAlertMgrManageAPIReload()
	if err != nil {
		logrus.Errorf("reload Alertmanager configuration error, %s", err)
	}
	logrus.Info("request the Alertmanager reload url.")

}

func isDigit(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}
