package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"log"
	"os"
)

type MyFormatter struct {
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2022-05-02 15:04:05")
	var newLog string
	newLog = fmt.Sprintf("[%s] [%s] %s\n", timestamp, entry.Level, entry.Message)

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func Initialize() {
	// Set log level
	logrus.SetLevel(logrus.DebugLevel)
	// Add caller's reported
	logrus.SetReportCaller(true)
	// Set the log output
	writerStdout := os.Stdout
	location := viper.GetString("log.location")
	if location != "" {
		writerFile, err := os.OpenFile(location, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.Fatalf("create file %v failed: %v", location, err)
		}
		logrus.SetOutput(io.MultiWriter(writerStdout, writerFile))
	} else {
		logrus.SetOutput(writerStdout)
	}
	// Set the log text format
	logrus.SetFormatter(&MyFormatter{})
}
