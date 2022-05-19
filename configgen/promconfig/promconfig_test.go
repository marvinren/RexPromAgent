package promconfig

import (
	"testing"
)

func TestWriteConfig(t *testing.T) {

	rule1 := AlertRule{
		Alert: "test",
		For:   "1m",
		Expr:  "(1 - avg by (instance)(irate(node_cpu_seconds_total{mode=\"idle\"}[5m]))) > 0.4",
		Labels: map[string]string{
			"severity": "critical",
			"category": "prealarm",
		},
		Annotations: AlertAnnotations{
			Summary:     "pre alerm {{$labels.instance}} cpu utils > {{$value}}",
			Description: "pre alerm {{$labels.instance}} cpu utils > {{$value}}",
		},
	}

	group1 := AlertGroup{
		"group1",
		[]AlertRule{rule1},
	}

	rules := PromConfigRules{
		[]AlertGroup{group1},
	}

	err := GeneratePromRuleFile(&rules, "./test.yml")
	if err != nil {
		t.Fatal(err)
	}
}
