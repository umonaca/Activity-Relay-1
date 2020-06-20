package api

import (
	"github.com/spf13/viper"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	testConfigPath := "../misc/test/config.yml"
	file, _ := os.Open(testConfigPath)
	defer file.Close()

	viper.SetConfigType("yaml")
	viper.ReadConfig(file)

	m.Run()
}
