package quota_test

import (
	"github.com/cloudfoundry/cli/cf/api/quotas/fakes"
	"github.com/cloudfoundry/cli/cf/errors"
	"github.com/cloudfoundry/cli/cf/models"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/quota"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("quotas command", func() {
	var (
		ui                  *testterm.FakeUI
		quotaRepo           *fakes.FakeQuotaRepository
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		quotaRepo = &fakes.FakeQuotaRepository{}
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	runCommand := func() bool {
		cmd := NewListQuotas(ui, testconfig.NewRepositoryWithDefaults(), quotaRepo)
		return testcmd.RunCommand(cmd, []string{}, requirementsFactory)
	}

	Describe("requirements", func() {
		It("requires the user to be logged in", func() {
			requirementsFactory.LoginSuccess = false
			Expect(runCommand()).ToNot(HavePassedRequirements())
		})
	})

	Context("when quotas exist", func() {
		BeforeEach(func() {
			quotaRepo.FindAllReturns([]models.QuotaFields{
				models.QuotaFields{
					Name:                    "quota-name",
					MemoryLimit:             1024,
					InstanceMemoryLimit:     512,
					RoutesLimit:             111,
					ServicesLimit:           222,
					NonBasicServicesAllowed: true,
				},
				models.QuotaFields{
					Name:                    "quota-non-basic-not-allowed",
					MemoryLimit:             434,
					InstanceMemoryLimit:     3,
					RoutesLimit:             1,
					ServicesLimit:           2,
					NonBasicServicesAllowed: false,
				},
			}, nil)
		})

		It("lists quotas", func() {
			Expect(Expect(runCommand()).To(HavePassedRequirements())).To(HavePassedRequirements())
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting quotas as", "my-user"},
				[]string{"OK"},
				[]string{"name", "total memory limit", "instance memory limit", "routes", "service instances", "paid service plans"},
				[]string{"quota-name", "1G", "512M", "111", "222", "allowed"},
				[]string{"quota-non-basic-not-allowed", "434M", "3M", "1", "2", "disallowed"},
			))
		})
	})

	Context("when an error occurs fetching quotas", func() {
		BeforeEach(func() {
			quotaRepo.FindAllReturns([]models.QuotaFields{}, errors.New("I haz a borken!"))
		})

		It("prints an error", func() {
			Expect(runCommand()).To(HavePassedRequirements())
			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Getting quotas as", "my-user"},
				[]string{"FAILED"},
			))
		})
	})

})
