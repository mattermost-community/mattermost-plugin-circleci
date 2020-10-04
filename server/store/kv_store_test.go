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

func TestDefaultProject(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Project")
}

var _ = Describe("Project", func() {
	var (
		pluginAPIMock *plugintest.API
		store         Store
	)

	BeforeEach(func() {
		pluginAPIMock = &plugintest.API{}
		store, _ = NewStore(pluginAPIMock)
		pluginAPIMock.On("LogError", mock.AnythingOfType("string"), mock.AnythingOfType("*model.AppError")).Return(nil)
	})

	It("should save default project", func() {

		project := &ProjectIdentifier{VCSType: "gh", Org: "som-org", Project: "some-project"}
		pluginAPIMock.On("KVSet", "123_default_project", mock.AnythingOfType("[]uint8")).Return(nil)

		err := store.StoreDefaultProject("123", *project)

		Expect(err).To(BeNil())
	})

	It("should return the error in case something goes wrong", func() {

		project := &ProjectIdentifier{VCSType: "gh", Org: "som-org", Project: "some-project"}
		pluginAPIMock.On("KVSet", "123_default_project", mock.AnythingOfType("[]uint8")).Return(&model.AppError{})

		err := store.StoreDefaultProject("123", *project)

		Expect(err).NotTo(BeNil())
	})

	It("should retrieve the saved default project", func() {

		project := &ProjectIdentifier{VCSType: "gh", Org: "som-org", Project: "some-project"}
		projectBytes, _ := json.Marshal(project)
		pluginAPIMock.On("KVGet", "123_default_project").Return(projectBytes, nil)

		savedProject, err := store.GetDefaultProject("123")

		Expect(err).To(BeNil())
		Expect(savedProject).Should(Equal(project))
	})
})
