package digitalocean

import (
	"log"
	"testing"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/stretchr/testify/assert"
)

var testConfig *config.Config

func init() {
	var err error
	testConfig, err = config.LoadConfiguration("../../fixtures/config_test.toml")
	if err != nil {
		log.Println(err)
	}
}

func Test_Provider_Digitalocean_GetProviderName(t *testing.T) {
	d := DO{Config: testConfig}
	if assert.NotNil(t, d) {
		assert.Equal(t, "DigitalOcean", d.GetProviderName())
	}
}

// TODO: add test for GetInstalledSshKeys()
// TODO: add test for InstallNewSshKey()
// TODO: add test for UpdateSshKey()

func Test_Provider_Digitalocean_GetConfig(t *testing.T) {
	d := DO{Config: testConfig}
	cfg := d.GetConfig()
	if assert.NotNil(t, cfg) {
		assert.Equal(t, testConfig, cfg)
	}
}
