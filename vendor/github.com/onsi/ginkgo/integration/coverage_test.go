package integration_test

import (
	"os"
	"os/exec"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Coverage Specs", func() {
	AfterEach(func() {
		os.RemoveAll("./_fixtures/coverage_fixture/coverage_fixture.coverprofile")
	})

	It("runs coverage analysis in series and in parallel", func() {
		session := startGinkgo("./_fixtures/coverage_fixture", "-cover")
		Eventually(session).Should(gexec.Exit(0))
		output := session.Out.Contents()
		Ω(string(output)).Should(ContainSubstring("coverage: 80.0% of statements"))

		serialCoverProfileOutput, err := exec.Command("go", "tool", "cover", "-func=./_fixtures/coverage_fixture/coverage_fixture.coverprofile").CombinedOutput()
		Ω(err).ShouldNot(HaveOccurred())

		os.RemoveAll("./_fixtures/coverage_fixture/coverage_fixture.coverprofile")

		Eventually(startGinkgo("./_fixtures/coverage_fixture", "-cover", "-nodes=4")).Should(gexec.Exit(0))

		parallelCoverProfileOutput, err := exec.Command("go", "tool", "cover", "-func=./_fixtures/coverage_fixture/coverage_fixture.coverprofile").CombinedOutput()
		Ω(err).ShouldNot(HaveOccurred())

		Ω(parallelCoverProfileOutput).Should(Equal(serialCoverProfileOutput))

		By("handling external packages")
		session = startGinkgo("./_fixtures/coverage_fixture", "-coverpkg=github.com/onsi/ginkgo/integration/_fixtures/coverage_fixture,github.com/onsi/ginkgo/integration/_fixtures/coverage_fixture/external_coverage_fixture")
		Eventually(session).Should(gexec.Exit(0))
		output = session.Out.Contents()
		Ω(output).Should(ContainSubstring("coverage: 71.4% of statements in github.com/onsi/ginkgo/integration/_fixtures/coverage_fixture, github.com/onsi/ginkgo/integration/_fixtures/coverage_fixture/external_coverage_fixture"))

		serialCoverProfileOutput, err = exec.Command("go", "tool", "cover", "-func=./_fixtures/coverage_fixture/coverage_fixture.coverprofile").CombinedOutput()
		Ω(err).ShouldNot(HaveOccurred())

		os.RemoveAll("./_fixtures/coverage_fixture/coverage_fixture.coverprofile")

		Eventually(startGinkgo("./_fixtures/coverage_fixture", "-coverpkg=github.com/onsi/ginkgo/integration/_fixtures/coverage_fixture,github.com/onsi/ginkgo/integration/_fixtures/coverage_fixture/external_coverage_fixture", "-nodes=4")).Should(gexec.Exit(0))

		parallelCoverProfileOutput, err = exec.Command("go", "tool", "cover", "-func=./_fixtures/coverage_fixture/coverage_fixture.coverprofile").CombinedOutput()
		Ω(err).ShouldNot(HaveOccurred())

		Ω(parallelCoverProfileOutput).Should(Equal(serialCoverProfileOutput))
	})

	It("validates coverprofile sets custom profile name", func() {
		session := startGinkgo("./_fixtures/coverage_fixture", "-cover", "-coverprofile=coverage.txt")

		Eventually(session).Should(gexec.Exit(0))

		// Check that the correct file was created
		_, err := os.Stat("./_fixtures/coverage_fixture/coverage.txt")

		Ω(err).ShouldNot(HaveOccurred())

		// Cleanup
		os.RemoveAll("./_fixtures/coverage_fixture/coverage.txt")
	})

	It("Works in recursive mode", func() {
		session := startGinkgo("./_fixtures/combined_coverage_fixture", "-r", "-cover", "-coverprofile=coverage.txt")

		Eventually(session).Should(gexec.Exit(0))

		packages := []string{"first_package", "second_package"}

		for _, p := range packages {
			coverFile := fmt.Sprintf("./_fixtures/combined_coverage_fixture/%s/coverage.txt", p)
			_, err := os.Stat(coverFile)

			Ω(err).ShouldNot(HaveOccurred())

			// Cleanup
			os.RemoveAll(coverFile)
		}
	})

	It("Works in parallel mode", func() {
		session := startGinkgo("./_fixtures/coverage_fixture", "-p", "-cover", "-coverprofile=coverage.txt")

		Eventually(session).Should(gexec.Exit(0))

		coverFile := "./_fixtures/coverage_fixture/coverage.txt"
		_, err := os.Stat(coverFile)

		Ω(err).ShouldNot(HaveOccurred())

		// Cleanup
		os.RemoveAll(coverFile)
	})

	It("Appends coverages if output dir and coverprofile were set", func() {
		session := startGinkgo("./_fixtures/combined_coverage_fixture", "-outputdir=./", "-r", "-cover", "-coverprofile=coverage.txt")

		Eventually(session).Should(gexec.Exit(0))

		_, err := os.Stat("./_fixtures/combined_coverage_fixture/coverage.txt")

		Ω(err).ShouldNot(HaveOccurred())

		// Cleanup
		os.RemoveAll("./_fixtures/combined_coverage_fixture/coverage.txt")
	})

	It("Creates directories in path if they don't exist", func() {
		session := startGinkgo("./_fixtures/combined_coverage_fixture", "-outputdir=./all/profiles/here", "-r", "-cover", "-coverprofile=coverage.txt")

		defer os.RemoveAll("./_fixtures/combined_coverage_fixture/all")
		defer os.RemoveAll("./_fixtures/combined_coverage_fixture/coverage.txt")

		Eventually(session).Should(gexec.Exit(0))

		_, err := os.Stat("./_fixtures/combined_coverage_fixture/all/profiles/here/coverage.txt")

		Ω(err).ShouldNot(HaveOccurred())
	})

	It("Moves coverages if only output dir was set", func() {
		session := startGinkgo("./_fixtures/combined_coverage_fixture", "-outputdir=./", "-r", "-cover")

		Eventually(session).Should(gexec.Exit(0))

		packages := []string{"first_package", "second_package"}

		for _, p := range packages {
			coverFile := fmt.Sprintf("./_fixtures/combined_coverage_fixture/%s.coverprofile", p)

			// Cleanup
			defer func(f string) {
				os.RemoveAll(f)
			}(coverFile)

			defer func(f string) {
				os.RemoveAll(fmt.Sprintf("./_fixtures/combined_coverage_fixture/%s/coverage.txt", f))
			}(p)

			_, err := os.Stat(coverFile)

			Ω(err).ShouldNot(HaveOccurred())
		}
	})
})
