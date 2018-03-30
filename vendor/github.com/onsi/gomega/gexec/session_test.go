package gexec_test

import (
	"os/exec"
	"syscall"
	"time"

	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var command *exec.Cmd
	var session *Session

	var outWriter, errWriter *Buffer

	BeforeEach(func() {
		outWriter = nil
		errWriter = nil
	})

	JustBeforeEach(func() {
		command = exec.Command(fireflyPath)
		var err error
		session, err = Start(command, outWriter, errWriter)
		Ω(err).ShouldNot(HaveOccurred())
	})

	Context("running a command", func() {
		It("should start the process", func() {
			Ω(command.Process).ShouldNot(BeNil())
		})

		It("should wrap the process's stdout and stderr with gbytes buffers", func(done Done) {
			Eventually(session.Out).Should(Say("We've done the impossible, and that makes us mighty"))
			Eventually(session.Err).Should(Say("Ah, curse your sudden but inevitable betrayal!"))
			defer session.Out.CancelDetects()

			select {
			case <-session.Out.Detect("Can we maybe vote on the whole murdering people issue"):
				Eventually(session).Should(Exit(0))
			case <-session.Out.Detect("I swear by my pretty floral bonnet, I will end you."):
				Eventually(session).Should(Exit(1))
			case <-session.Out.Detect("My work's illegal, but at least it's honest."):
				Eventually(session).Should(Exit(2))
			}

			close(done)
		})

		It("should satisfy the gbytes.BufferProvider interface, passing Stdout", func() {
			Eventually(session).Should(Say("We've done the impossible, and that makes us mighty"))
			Eventually(session).Should(Exit())
		})
	})

	Describe("providing the exit code", func() {
		It("should provide the app's exit code", func() {
			Ω(session.ExitCode()).Should(Equal(-1))

			Eventually(session).Should(Exit())
			Ω(session.ExitCode()).Should(BeNumerically(">=", 0))
			Ω(session.ExitCode()).Should(BeNumerically("<", 3))
		})
	})

	Describe("wait", func() {
		It("should wait till the command exits", func() {
			Ω(session.ExitCode()).Should(Equal(-1))
			Ω(session.Wait().ExitCode()).Should(BeNumerically(">=", 0))
			Ω(session.Wait().ExitCode()).Should(BeNumerically("<", 3))
		})
	})

	Describe("exited", func() {
		It("should close when the command exits", func() {
			Eventually(session.Exited).Should(BeClosed())
			Ω(session.ExitCode()).ShouldNot(Equal(-1))
		})
	})

	Describe("kill", func() {
		It("should kill the command and don't wait for it to exit", func() {
			session, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			session.Kill()
			Ω(session).ShouldNot(Exit(), "Should not exit immediately...")
			Eventually(session).Should(Exit(128 + 9))
		})
	})

	Describe("interrupt", func() {
		It("should interrupt the command", func() {
			session, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			session.Interrupt()
			Ω(session).ShouldNot(Exit(), "Should not exit immediately...")
			Eventually(session).Should(Exit(128 + 2))
		})
	})

	Describe("terminate", func() {
		It("should terminate the command", func() {
			session, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			session.Terminate()
			Ω(session).ShouldNot(Exit(), "Should not exit immediately...")
			Eventually(session).Should(Exit(128 + 15))
		})
	})

	Describe("signal", func() {
		It("should send the signal to the command", func() {
			session, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
			Ω(err).ShouldNot(HaveOccurred())

			session.Signal(syscall.SIGABRT)
			Ω(session).ShouldNot(Exit(), "Should not exit immediately...")
			Eventually(session).Should(Exit(128 + 6))
		})
	})

	Context("tracking sessions", func() {
		BeforeEach(func() {
			KillAndWait()
		})

		Describe("kill", func() {
			It("should kill all the started sessions", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Kill()

				Eventually(session1).Should(Exit(128 + 9))
				Eventually(session2).Should(Exit(128 + 9))
				Eventually(session3).Should(Exit(128 + 9))
			})

			It("should not wait for exit", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Kill()
				Ω(session1).ShouldNot(Exit(), "Should not exit immediately...")

				Eventually(session1).Should(Exit(128 + 9))
			})

			It("should not track unstarted sessions", func() {
				_, err := Start(exec.Command("does not exist", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).Should(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Kill()

				Eventually(session2).Should(Exit(128 + 9))
				Eventually(session3).Should(Exit(128 + 9))
			})

		})

		Describe("killAndWait", func() {
			It("should kill all the started sessions and wait for them to finish", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				KillAndWait()
				Ω(session1).Should(Exit(128+9), "Should have exited")
				Ω(session2).Should(Exit(128+9), "Should have exited")
				Ω(session3).Should(Exit(128+9), "Should have exited")
			})
		})

		Describe("terminate", func() {
			It("should terminate all the started sessions", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Terminate()

				Eventually(session1).Should(Exit(128 + 15))
				Eventually(session2).Should(Exit(128 + 15))
				Eventually(session3).Should(Exit(128 + 15))
			})

			It("should not wait for exit", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Terminate()

				Ω(session1).ShouldNot(Exit(), "Should not exit immediately...")
			})
		})

		Describe("terminateAndWait", func() {
			It("should terminate all the started sessions, and wait for them to exit", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				TerminateAndWait()

				Ω(session1).Should(Exit(128+15), "Should have exited")
				Ω(session2).Should(Exit(128+15), "Should have exited")
				Ω(session3).Should(Exit(128+15), "Should have exited")
			})
		})

		Describe("signal", func() {
			It("should signal all the started sessions", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Signal(syscall.SIGABRT)

				Eventually(session1).Should(Exit(128 + 6))
				Eventually(session2).Should(Exit(128 + 6))
				Eventually(session3).Should(Exit(128 + 6))
			})

			It("should not wait", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Signal(syscall.SIGABRT)

				Ω(session1).ShouldNot(Exit(), "Should not exit immediately...")
			})
		})

		Describe("interrupt", func() {
			It("should interrupt all the started sessions, and not wait", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session2, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				session3, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Interrupt()

				Eventually(session1).Should(Exit(128 + 2))
				Eventually(session2).Should(Exit(128 + 2))
				Eventually(session3).Should(Exit(128 + 2))
			})

			It("should not wait", func() {
				session1, err := Start(exec.Command("sleep", "10000000"), GinkgoWriter, GinkgoWriter)
				Ω(err).ShouldNot(HaveOccurred())

				Interrupt()

				Ω(session1).ShouldNot(Exit(), "Should not exit immediately...")
			})
		})
	})

	Context("when the command exits", func() {
		It("should close the buffers", func() {
			Eventually(session).Should(Exit())

			Ω(session.Out.Closed()).Should(BeTrue())
			Ω(session.Err.Closed()).Should(BeTrue())

			Ω(session.Out).Should(Say("We've done the impossible, and that makes us mighty"))
		})

		var So = It

		So("this means that eventually should short circuit", func() {
			t := time.Now()
			failures := InterceptGomegaFailures(func() {
				Eventually(session).Should(Say("blah blah blah blah blah"))
			})
			Ω(time.Since(t)).Should(BeNumerically("<=", 500*time.Millisecond))
			Ω(failures).Should(HaveLen(1))
		})
	})

	Context("when wrapping out and err", func() {
		BeforeEach(func() {
			outWriter = NewBuffer()
			errWriter = NewBuffer()
		})

		It("should route to both the provided writers and the gbytes buffers", func() {
			Eventually(session.Out).Should(Say("We've done the impossible, and that makes us mighty"))
			Eventually(session.Err).Should(Say("Ah, curse your sudden but inevitable betrayal!"))

			Ω(outWriter.Contents()).Should(ContainSubstring("We've done the impossible, and that makes us mighty"))
			Ω(errWriter.Contents()).Should(ContainSubstring("Ah, curse your sudden but inevitable betrayal!"))

			Eventually(session).Should(Exit())

			Ω(outWriter.Contents()).Should(Equal(session.Out.Contents()))
			Ω(errWriter.Contents()).Should(Equal(session.Err.Contents()))
		})
	})

	Describe("when the command fails to start", func() {
		It("should return an error", func() {
			_, err := Start(exec.Command("agklsjdfas"), nil, nil)
			Ω(err).Should(HaveOccurred())
		})
	})
})
