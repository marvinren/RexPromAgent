package server

import (
	"RexPromAgent/configgen/promconfig"
	"RexPromAgent/db"
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
	configPath := viper.GetString("prometheus.alertRulesConfigPath")
	// Get the alert rules from the database.
	rules := make([]db.AlertRule, 0)
	err := s.db.FetchAlerts(&rules)
	if err != nil {
		logrus.Errorf("err for get alert data: %s", err)
	}

	// Convert the database object to the prometheus's rule config object.
	alertRules := make([]promconfig.AlertRule, 0)

	for _, alert := range rules {
		forDurationTime := alert.Duration
		if isDigit(forDurationTime) {
			forDurationTime = forDurationTime + "m"
		}

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
	}

	alertRuleObj := promconfig.PromConfigRules{
		[]promconfig.AlertGroup{promconfig.AlertGroup{
			"group_auto",
			alertRules,
		}},
	}
	// Generate the rule file from the prometheus's rule config object
	err = promconfig.GeneratePromRuleFile(&alertRuleObj, configPath)
	if err != nil {
		logrus.Errorf("generate the alert rules error, %s", err)
	}
	logrus.Infof("generate the config file %v", configPath)

	// Reload the prometheus configuration
	err = promconfig.GetPromManageAPIReload()
	if err != nil {
		logrus.Errorf("reload Prometheus configuration error, %s", err)
	}
	logrus.Info("request the prometheus reload url.")

}

func isDigit(str string) bool {
	for _, x := range []rune(str) {
		if !unicode.IsDigit(x) {
			return false
		}
	}
	return true
}
