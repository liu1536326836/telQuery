package config

import (
	"testing"
)

func TestParseConfig(t *testing.T) {
	ParseConfig("./config.toml")

	t.Log("Web:", Conf.Web)
	t.Log("DB:", Conf.DB)
	t.Log("Log:", Conf.Log)
	t.Log("Lib:", Conf.Lib)
	t.Log("Pongo:", Conf.Pongo)
}
