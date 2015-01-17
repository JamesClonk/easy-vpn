package vultr

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

func Test_Provider_Vultr_GetProviderName(t *testing.T) {
	v := Vultr{Config: testConfig}
	if assert.NotNil(t, v) {
		assert.Equal(t, "VULTR", v.GetProviderName())
	}
}

// TODO: add test for GetInstalledSshKeys()
// TODO: add test for InstallNewSshKey()
// TODO: add test for UpdateSshKey()

func Test_Provider_Vultr_GetConfig(t *testing.T) {
	v := Vultr{Config: testConfig}
	cfg := v.GetConfig()
	if assert.NotNil(t, cfg) {
		assert.Equal(t, testConfig, cfg)
	}
}

func Test_Provider_Vultr_UrlWithApiKey(t *testing.T) {
	v := Vultr{Config: testConfig}
	url := v.urlWithApiKey("http://localhost/base")
	if assert.NotNil(t, url) {
		assert.Equal(t, "http://localhost/base?api_key=xyzabcdefg999", url)
	}
}
