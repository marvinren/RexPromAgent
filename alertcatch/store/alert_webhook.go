package store

import (
	"RexPromAgent/pkg/db"
	"encoding/json"
	"fmt"
)

func ParsePromAlert(payload []byte) (*db.AlertGroup, error) {
	d := db.AlertGroup{}
	err := json.Unmarshal(payload, &d)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json webhook payload: %s", err)
	}
	return &d, nil
}
