package vultr

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
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

func getTestServer(code int, body string) *httptest.Server {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, body)
	}))
	baseUrl = ts.URL
	return ts
}

func Test_Provider_Vultr_GetProviderName(t *testing.T) {
	v := Vultr{Config: testConfig}
	if assert.NotNil(t, v) {
		assert.Equal(t, "vultr", v.GetProviderName())
	}
}

func Test_Provider_Vultr_GetConfig(t *testing.T) {
	v := Vultr{Config: testConfig}
	cfg := v.GetConfig()
	if assert.NotNil(t, cfg) {
		assert.Equal(t, testConfig, cfg)
	}
}

func Test_Provider_Vultr_GetInstalledSshKeys_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keys, err := v.GetInstalledSshKeys()
	assert.Nil(t, keys)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_GetInstalledSshKeys_NoKeys(t *testing.T) {
	server := getTestServer(http.StatusOK, `[]`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keys, err := v.GetInstalledSshKeys()
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, keys)
}

func Test_Provider_Vultr_GetInstalledSshKeys_Keys(t *testing.T) {
	server := getTestServer(http.StatusOK,
		`{
			"one":{"SSHKEYID":"1","name":"alpha","ssh_key":"aaaa","date_created":null},
			"two":{"SSHKEYID":"2","name":"beta","ssh_key":"bbbb","date_created":null},
			"three":{"SSHKEYID":"3","name":"charlie","ssh_key":"cccc"}
		}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keys, err := v.GetInstalledSshKeys()
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, keys) {
		assert.Equal(t, 3, len(keys))
		// keys can be in random order
		for _, key := range keys {
			switch key.Id {
			case "1":
				assert.Equal(t, "alpha", key.Name)
			case "2":
				assert.Equal(t, "beta", key.Name)
			case "3":
				assert.Equal(t, "cccc", key.Key)
			default:
				t.Fail()
			}
		}
	}
}

func Test_Provider_Vultr_InstallNewSshKey_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keyId, err := v.InstallNewSshKey("delta", "ddddd")
	assert.Equal(t, "", keyId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_InstallNewSshKey_NoKeyId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keyId, err := v.InstallNewSshKey("delta", "ddddd")
	assert.Equal(t, "", keyId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_InstallNewSshKey_KeyId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{"SSHKEYID":"d1d2d3d4"}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keyId, err := v.InstallNewSshKey("delta", "ddddd")
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, keyId) {
		assert.Equal(t, "d1d2d3d4", keyId)
	}
}

func Test_Provider_Vultr_UpdateSshKey_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keyId, err := v.UpdateSshKey("o1", "omega", "oooo")
	assert.Equal(t, "", keyId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_UpdateSshKey_KeyId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{no-response?!}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	keyId, err := v.UpdateSshKey("o1", "omega", "oooo")
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, keyId) {
		assert.Equal(t, "o1", keyId)
	}
}

// TODO: add test for GetAllVMs()

func Test_Provider_Vultr_UrlWithApiKey(t *testing.T) {
	v := Vultr{Config: testConfig}
	url := v.urlWithApiKey("http://localhost/base")
	if assert.NotNil(t, url) {
		assert.Equal(t, "http://localhost/base?api_key=xyzabcdefg999", url)
	}
}
