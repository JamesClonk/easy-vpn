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
				t.Error("Unknown Key Id")
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

func Test_Provider_Vultr_GetAllVMs_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	machines, err := v.GetAllVMs()
	assert.Nil(t, machines)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_GetAllVMs_NoVMs(t *testing.T) {
	server := getTestServer(http.StatusOK, `[]`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	machines, err := v.GetAllVMs()
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, machines)
}

func Test_Provider_Vultr_GetAllVMs_VMs(t *testing.T) {
	server := getTestServer(http.StatusOK,
		`{
			"1":{"SUBID":"1","label":"alpha","OS":"ubuntu","main_ip":"123.456.789.0"},
			"2":{"SUBID":"2","label":"beta","OS":"ubuntu","DCID":"Earth","status":"active"},
			"3":{"SUBID":"3","label":"charlie","OS":"centos"}
		}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	machines, err := v.GetAllVMs()
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, machines) {
		assert.Equal(t, 3, len(machines))
		// machines can be in random order
		for _, vm := range machines {
			switch vm.Id {
			case "1":
				assert.Equal(t, "alpha", vm.Name)
				assert.Equal(t, "ubuntu", vm.OS)
				assert.Equal(t, "123.456.789.0", vm.IP)
			case "2":
				assert.Equal(t, "Earth", vm.Region)
				assert.Equal(t, "active", vm.Status)
			case "3":
				assert.Equal(t, "charlie", vm.Name)
				assert.Equal(t, "centos", vm.OS)
			default:
				t.Error("Unknown VM Id")
			}
		}
	}
}

func Test_Provider_Vultr_CreateVM_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	vmId, err := v.CreateVM("test-vm", "test-os", "test-size", "test-region", "test-key-id")
	assert.Equal(t, "", vmId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_CreateVM_NoVMId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	vmId, err := v.CreateVM("test-vm", "test-os", "test-size", "test-region", "test-key-id")
	assert.Equal(t, "", vmId)
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_CreateVM_VMId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{"SUBID":"test-vm"}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	vmId, err := v.CreateVM("test-vm", "test-os", "test-size", "test-region", "test-key-id")
	if err != nil {
		t.Error(err)
	}
	if assert.NotNil(t, vmId) {
		assert.Equal(t, "test-vm", vmId)
	}
}

func Test_Provider_Vultr_StartVM_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	err := v.StartVM("alpha")
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_StartVM_VMId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{no-response?!}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	err := v.StartVM("alpha")
	if err != nil {
		t.Error(err)
	}
}

func Test_Provider_Vultr_DestroyVM_Error(t *testing.T) {
	server := getTestServer(http.StatusNotAcceptable, `{error-message}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	err := v.DestroyVM("alpha")
	if assert.NotNil(t, err) {
		assert.Equal(t, `{error-message}`, err.Error())
	}
}

func Test_Provider_Vultr_DestroyVM_VMId(t *testing.T) {
	server := getTestServer(http.StatusOK, `{no-response?!}`)
	defer server.Close()

	v := Vultr{Config: testConfig}

	err := v.DestroyVM("alpha")
	if err != nil {
		t.Error(err)
	}
}

func Test_Provider_Vultr_UrlWithApiKey(t *testing.T) {
	v := Vultr{Config: testConfig}
	url := v.urlWithApiKey("http://localhost/base")
	if assert.NotNil(t, url) {
		assert.Equal(t, "http://localhost/base?api_key=xyzabcdefg999", url)
	}
}
