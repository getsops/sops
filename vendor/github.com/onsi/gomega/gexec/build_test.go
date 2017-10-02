package gexec_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var packagePath = "./_fixture/firefly"

var _ = Describe(".Build", func() {
	Context("when there have been previous calls to Build", func() {
		BeforeEach(func() {
			_, err := gexec.Build(packagePath)
			Ω(err).ShouldNot(HaveOccurred())
		})

		It("compiles the specified package", func() {
			compiledPath, err := gexec.Build(packagePath)
			Ω(err).ShouldNot(HaveOccurred())
			Ω(compiledPath).Should(BeAnExistingFile())
		})

		Context("and CleanupBuildArtifacts has been called", func() {
			BeforeEach(func() {
				gexec.CleanupBuildArtifacts()
			})

			It("compiles the specified package", func() {
				var err error
				fireflyPath, err = gexec.Build(packagePath)
				Ω(err).ShouldNot(HaveOccurred())
				Ω(fireflyPath).Should(BeAnExistingFile())
			})
		})
	})
})

var _ = Describe(".BuildWithEnvironment", func() {
	var err error
	env := []string{
		"GOOS=linux",
		"GOARCH=amd64",
	}

	It("compiles the specified package with the specified env vars", func() {
		compiledPath, err := gexec.BuildWithEnvironment(packagePath, env)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(compiledPath).Should(BeAnExistingFile())
	})

	It("returns the environment to a good state", func() {
		_, err = gexec.BuildWithEnvironment(packagePath, env)
		Ω(err).ShouldNot(HaveOccurred())
		Ω(os.Environ()).ShouldNot(ContainElement("GOOS=linux"))
	})
})
