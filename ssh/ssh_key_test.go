package ssh

import (
	"log"
	"os/user"
	"testing"

	"github.com/JamesClonk/easy-vpn/config"
	"github.com/JamesClonk/easy-vpn/provider"
	"github.com/JamesClonk/easy-vpn/test"
	"github.com/stretchr/testify/assert"
)

var cfg *config.Config

func init() {
	var err error
	cfg, err = config.LoadConfiguration("../fixtures/config_test.toml")
	if err != nil {
		log.Println(err)
	}
}

func Test_Main_GetEasyVpnKeyId(t *testing.T) {
	mockedProvider1 := test.MockProvider{
		Config: cfg,
		Keys: []provider.SshKey{
			provider.SshKey{},
			provider.SshKey{
				Name: "easy-vpn",
				Id:   "mockId",
			},
			provider.SshKey{
				Name: "mockName",
			},
		},
	}

	keyId1 := GetEasyVpnKeyId(mockedProvider1, "easy-vpn")
	if assert.NotNil(t, keyId1) {
		assert.Equal(t, "mockId:easy-vpn:this would be a public key!\n;)\n", keyId1)
	}

	mockedProvider2 := test.MockProvider{Config: cfg}
	keyId2 := GetEasyVpnKeyId(mockedProvider2, "easy-vpn")
	if assert.NotNil(t, keyId2) {
		assert.Equal(t, "easy-vpn:this would be a public key!\n;)\n", keyId2)
	}
}

func Test_Main_ReadKeyFile(t *testing.T) {
	if assert.NotNil(t, cfg) {
		assert.Equal(t, "vultr", cfg.Provider)
	}

	pubkey := string(readKeyFile(cfg.PublicKeyFile))
	if assert.NotNil(t, pubkey) {
		assert.Equal(t, "this would be a public key!\n;)\n", pubkey)
	}

	privkey := string(readKeyFile(cfg.PrivateKeyFile))
	if assert.NotNil(t, privkey) {
		assert.Equal(t, "this would be a private key!\n;)\n", privkey)
	}
}

func Test_Main_SanitizeFilename(t *testing.T) {
	filename := sanitizeFilename("~/test/123.txt")
	if assert.NotNil(t, filename) {
		usr, _ := user.Current()
		home := usr.HomeDir
		assert.Equal(t, home+"/test/123.txt", filename)
	}
}
