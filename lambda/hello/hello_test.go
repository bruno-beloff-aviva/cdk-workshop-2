package main

import (
	"cdk-workshop-2/business"
	"cdk-workshop-2/log"
	"fmt"
	"testing"
)

func MainRunner(sourceIP string, path string) string {
	log.Logger.Debug("sourceIP: ", sourceIP, "path: ", path)
	return business.Hello(sourceIP, path)
}

var tests = []struct {
	sourceIP, path, want string
}{
	{"", "", "Hello Go world!"},
	{"1.2.3.4", "/", "Hello Go world at / from 1.2.3.4"},
	{"1.2.3.4", "/abc", "Hello Go world at /abc from 1.2.3.4"},
}

func TestMain(t *testing.T) {
	log.Init("gohello")

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
