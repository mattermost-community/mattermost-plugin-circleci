package plugin

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestConfiguration(t *testing.T) {
	t.Run("null configuration", func(t *testing.T) {
		plugin := &Plugin{}
		assert.NotNil(t, plugin.getConfiguration())
	})

	t.Run("changing configuration", func(t *testing.T) {
		plugin := &Plugin{}

		configuration1 := &configuration{WebhooksSecret: "z_8EY-gFuZxIr9TMFyd0TBBN6aGzmGv3", EncryptionKey: "y1xqt8wPCQXIm6Wezx4vfxtAzN_3AhqI"}
		plugin.setConfiguration(configuration1)
		assert.Equal(t, configuration1, plugin.getConfiguration())

		configuration2 := &configuration{WebhooksSecret: "yqXlYzvO2bjqo3aGENy-Gj+D3J5Mg3ua", EncryptionKey: "W3ilr5uKTEt4tYCsf4RWzUP-lZKrFwl5"}
		plugin.setConfiguration(configuration2)
		assert.Equal(t, configuration2, plugin.getConfiguration())
		assert.NotEqual(t, configuration1, plugin.getConfiguration())
		assert.False(t, plugin.getConfiguration() == configuration1)
		assert.True(t, plugin.getConfiguration() == configuration2)
	})

	t.Run("setting same configuration", func(t *testing.T) {
		plugin := &Plugin{}

		configuration1 := &configuration{WebhooksSecret: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", EncryptionKey: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
		plugin.setConfiguration(configuration1)
		assert.Panics(t, func() {
			plugin.setConfiguration(configuration1)
		})
	})

	t.Run("clearing configuration", func(t *testing.T) {
		plugin := &Plugin{}

		configuration1 := &configuration{WebhooksSecret: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", EncryptionKey: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"}
		plugin.setConfiguration(configuration1)
		assert.NotPanics(t, func() {
			plugin.setConfiguration(nil)
		})
		assert.NotNil(t, plugin.getConfiguration())
		assert.NotEqual(t, configuration1, plugin.getConfiguration())
	})
}

func TestOnConfigurationChange(t *testing.T) {
	for name, test := range map[string]struct {
		SetupAPI    func() *plugintest.API
		ShouldError bool
	}{
		"Webhook Secret is not defined": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*plugin.configuration")).Return(nil).Run(func(args mock.Arguments) {
					apiConfiguration := args.Get(0).(*configuration)
					apiConfiguration.EncryptionKey = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
				})
				return api
			},
			ShouldError: true,
		},
		"At Rest Encryption key is not defined": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*plugin.configuration")).Return(nil).Run(func(args mock.Arguments) {
					apiConfiguration := args.Get(0).(*configuration)
					apiConfiguration.WebhooksSecret = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				})
				return api
			},
			ShouldError: true,
		},
		"Webhook secret contains a forwardslash": {
			SetupAPI: func() *plugintest.API {
				api := &plugintest.API{}
				api.On("LoadPluginConfiguration", mock.AnythingOfType("*plugin.configuration")).Return(nil).Run(func(args mock.Arguments) {
					apiConfiguration := args.Get(0).(*configuration)
					apiConfiguration.WebhooksSecret = "aaaaaaaaaaa/aaaaaaaaaaaaaaaaaaaa"
					apiConfiguration.EncryptionKey = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
				})
				return api
			},
			ShouldError: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			api := test.SetupAPI()
			defer api.AssertExpectations(t)

			p := Plugin{}
			p.setConfiguration(&configuration{})
			p.SetAPI(api)

			err := p.OnConfigurationChange()

			if test.ShouldError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
