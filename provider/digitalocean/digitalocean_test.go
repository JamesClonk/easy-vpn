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
		assert.Equal(t, "digitalocean", d.GetProviderName())
	}
}

func Test_Provider_Digitalocean_GetConfig(t *testing.T) {
	d := DO{Config: testConfig}
	cfg := d.GetConfig()
	if assert.NotNil(t, cfg) {
		assert.Equal(t, testConfig, cfg)
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
				t.Error("Unknown Key Id")
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

func Test_Provider_Digitalocean_GetAllVMs_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	d := DO{Config: testConfig}

	machines, err := d.GetAllVMs()
	assert.Nil(t, machines)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Digitalocean_GetAllVMs_NoVMs(t *testing.T) {
	server := getTestServer(http.StatusOK, `{"droplets": [],"links": {},"meta": {"total": 0}}`)
	defer server.Close()

	d := DO{Config: testConfig}

	machines, err := d.GetAllVMs()
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, machines)
}

func Test_Provider_Digitalocean_GetAllVMs_VMs(t *testing.T) {
	server := getTestServer(http.StatusOK, `{"droplets": [{
    "id": 1111,
    "name": "example.com",
    "memory": 512,
    "disk": 20,
    "status": "active",
    "image": {
        "id": 6918990,
        "name": "14.04 x64",
        "distribution": "Ubuntu",
        "slug": "ubuntu-14-04-x64"
    },
    "size_slug": "512mb",
    "networks": {
        "v4": [{
            "ip_address": "104.236.32.111",
            "netmask": "255.255.192.0",
            "gateway": "104.236.0.1",
            "type": "public"
        }],
        "v6": [{}]
    },
    "region": {
        "name": "New York 3",
        "slug": "nyc3"
    }
},{
    "id": 7777,
    "name": "example.com",
    "memory": 512,
    "status": "active",
    "image": {
        "id": 123,
        "distribution": "CentOS",
        "slug": "centos-x64"
    },
    "size_slug": "512mb",
    "networks": {
        "v4": [{
            "ip_address": "104.236.32.999",
            "netmask": "255.255.192.0",
            "gateway": "104.236.0.1",
            "type": "public"
        }]
    },
    "region": {
        "name": "New York 2",
        "slug": "nyc2"
    }
},{
    "id": 9999,
    "name": "test.com",
    "memory": 1024,
    "status": "stopped",
    "image": {
        "slug": "ubuntu-14-10-x64"
    },
    "networks": {"v4": [{"ip_address": "104.236.32.777"}]},
    "region": {
        "name": "New York 1",
        "slug": "nyc1"
    }
}],"links": {},"meta": {"total": 3}}`)
	defer server.Close()

	d := DO{Config: testConfig}

	machines, err := d.GetAllVMs()
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, machines) {
		assert.Equal(t, 3, len(machines))
		// machines can be in random order
		for _, vm := range machines {
			switch vm.Id {
			case "1111":
				assert.Equal(t, "example.com", vm.Name)
				assert.Equal(t, "ubuntu-14-04-x64", vm.OS)
				assert.Equal(t, "nyc3", vm.Region)
				assert.Equal(t, "active", vm.Status)
			case "7777":
				assert.Equal(t, "nyc2", vm.Region)
				assert.Equal(t, "104.236.32.999", vm.IP)
				assert.Equal(t, "centos-x64", vm.OS)
			case "9999":
				assert.Equal(t, "test.com", vm.Name)
				assert.Equal(t, "stopped", vm.Status)
				assert.Equal(t, "104.236.32.777", vm.IP)
			default:
				t.Error("Unknown VM Id")
			}
		}
	}
}
