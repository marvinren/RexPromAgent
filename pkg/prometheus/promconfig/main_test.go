package promconfig

import (
	"RexPromAgent/pkg/config"
	"RexPromAgent/pkg/log"
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	log.Initialize()
	config.Initialize()
	m.Run()
	// delete file prom.yml
	_, err := os.Stat("./prom.yml")
	if err == nil {
		err := os.Remove("./prom.yml")
		if err != nil {
			fmt.Printf("remove file error, %v \n", err)
		}
	}

	// delete file alertmgr.yml
	_, err = os.Stat("./alertmgr.yml")
	if err == nil {
		err = os.Remove("./alertmgr.yml")
		if err != nil {
			fmt.Printf("remove file error, %v \n", err)
		}
	}

}
