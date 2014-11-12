package main

import (
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitConfig(t *testing.T) {
	initConfig("evhandler")
	actions := viper.GetStringMapString("actions")
	assert.NotEmpty(t, actions)
}
