package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/business/hits"
	"fmt"
	"log"
	"testing"

	"github.com/joerdav/zapray"
)

var testLogger *zapray.Logger

func init() {
	var err error
	testLogger, err = zapray.NewProduction()
	if err != nil {
		panic("failed to create logger: " + err.Error())
	}
	testLogger.Info("hello_lambda init!!")

	// level := os.Getenv("LOG_LEVEL") // may also be set by ApplicationLogLevelV2: awslambda.ApplicationLogLevel_INFO
	// fmt.Println("log level: ", level)
}

func MainRunner(sourceIP string, path string) string {
	log.Println("sourceIP: ", sourceIP, "path: ", path)
	hit := hits.NewHits(path)
	return business.Hello(testLogger, sourceIP, hit)
}

var tests = []struct {
	sourceIP, path, want string
}{
	{"", "", "Hello Go world!"},
	{"1.2.3.4", "/", "Hello Go world at / from 1.2.3.4 hits: 0"},
	{"1.2.3.4", "/abc", "Hello Go world at /abc from 1.2.3.4 hits: 0"},
}

func TestMain(t *testing.T) {
	// Test businessFunction
	for _, tt := range tests {
		testName := fmt.Sprintf("%s, %s", tt.sourceIP, tt.path)

		testFunc := func(t *testing.T) {
			ans := MainRunner(tt.sourceIP, tt.path)
			if ans != tt.want {
				t.Errorf("got %s, want %s", ans, tt.want)
			}
		}

		t.Run(testName, testFunc)
	}
}

func TestPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	MainRunner("1.2.3.4", "/panic/")
}
