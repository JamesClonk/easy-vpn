package digitalocean

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

func Test_Provider_Digitalocean_GetProviderName(t *testing.T) {
	d := DO{Config: testConfig}
	if assert.NotNil(t, d) {
		assert.Equal(t, "DigitalOcean", d.GetProviderName())
	}
}

func Test_Provider_Digitalocean_GetInstalledSshKeys_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	d := DO{Config: testConfig}

	keys, err := d.GetInstalledSshKeys()
	assert.Nil(t, keys)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Digitalocean_GetInstalledSshKeys_NoKeys(t *testing.T) {
	server := getTestServer(http.StatusOK,
		`{"ssh_keys": [], "links": {}, "meta": {"total": 0}}`)
	defer server.Close()

	d := DO{Config: testConfig}

	keys, err := d.GetInstalledSshKeys()
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, keys)
}

func Test_Provider_Digitalocean_GetInstalledSshKeys_Keys(t *testing.T) {
	server := getTestServer(http.StatusOK,
		`{"ssh_keys": [
			{"id": 1, "fingerprint": "who cares", "public_key": "aaaa", "name": "alpha"},
			{"id": 2, "fingerprint": "who cares", "public_key": "bbbb", "name": "beta"},
			{"id": 3, "fingerprint": "who cares", "public_key": "cccc", "name": "charlie"}
			], "links": {}, "meta": {"total": 3}
		}`)
	defer server.Close()

	d := DO{Config: testConfig}

	keys, err := d.GetInstalledSshKeys()
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, keys) {
		assert.Equal(t, 3, len(keys))
		// keys can be in random order
		for _, key := range keys {
			switch key.Id {
			case "1":
				assert.Equal(t, "aaaa", key.Key)
			case "2":
				assert.Equal(t, "bbbb", key.Key)
			case "3":
				assert.Equal(t, "charlie", key.Name)
			default:
				t.Fail()
			}
		}
	}
}

func Test_Provider_Digitalocean_InstallNewSshKey_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	d := DO{Config: testConfig}

	keyId, err := d.InstallNewSshKey("delta", "ddddd")
	assert.Equal(t, "", keyId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Digitalocean_InstallNewSshKey_NoKeyId(t *testing.T) {
	server := getTestServer(http.StatusCreated, `{error-message}`)
	defer server.Close()

	d := DO{Config: testConfig}

	keyId, err := d.InstallNewSshKey("delta", "ddddd")
	assert.Equal(t, "", keyId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Digitalocean_InstallNewSshKey_KeyId(t *testing.T) {
	server := getTestServer(http.StatusCreated,
		`{"ssh_key":{"id":4,"public_key":"ddddd","name":"My SSH Public Key"}}`)
	defer server.Close()

	d := DO{Config: testConfig}

	keyId, err := d.InstallNewSshKey("delta", "ddddd")
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, keyId) {
		assert.Equal(t, "4", keyId)
	}
}

func Test_Provider_Digitalocean_DeleteSshKey_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	d := DO{Config: testConfig}

	err := d.deleteSshKey("123")
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Digitalocean_DeleteSshKey_Deleted(t *testing.T) {
	server := getTestServer(http.StatusNoContent, `{no-response?!}`)
	defer server.Close()

	d := DO{Config: testConfig}

	err := d.deleteSshKey("123")
	if err != nil {
		t.Error(err)
	}
}

func Test_Provider_Digitalocean_GetConfig(t *testing.T) {
	d := DO{Config: testConfig}
	cfg := d.GetConfig()
	if assert.NotNil(t, cfg) {
		assert.Equal(t, testConfig, cfg)
	}
}
