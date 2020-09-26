package store

import (
	"encoding/json"
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest/mock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConfig(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Config")
}

var _ = Describe("Config", func() {
	var (
		pluginAPIMock *plugintest.API
		store         Store
	)

	BeforeEach(func() {
		pluginAPIMock = &plugintest.API{}
		store, _ = NewStore(pluginAPIMock)
		pluginAPIMock.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("*model.AppError")).Return(nil)
	})

	It("should save config", func() {

		config := &Config{Org: "som-org", Project: "some-project"}
		pluginAPIMock.On("KVSet", "123_circleci_config", mock.AnythingOfType("[]uint8")).Return(nil)

		err := store.SaveConfig("123", *config)

		Expect(err).To(BeNil())
	})

	It("should return the error in case something goes wrong", func() {

		config := &Config{Org: "som-org", Project: "some-project"}
		pluginAPIMock.On("KVSet", "123_circleci_config", mock.AnythingOfType("[]uint8")).Return(&model.AppError{})

		err := store.SaveConfig("123", *config)

		Expect(err).NotTo(BeNil())
	})

	It("should retrieve the saved config", func() {

		config := &Config{Org: "som-org", Project: "some-project"}
		configBytes, _ := json.Marshal(config)
		pluginAPIMock.On("KVGet", "123_circleci_config").Return(configBytes, nil)

		savedConfig, err := store.GetConfig("123")

		Expect(err).To(BeNil())
		Expect(savedConfig).Should(Equal(config))
	})
})
