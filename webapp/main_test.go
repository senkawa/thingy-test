package main

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/tebeka/selenium"
)

var host = os.Getenv("WEBDRIVER_HOST")
var app = os.Getenv("APP_URL")

const (
	StubValidPassword   = "alongpasswordaaaaa"
	StubInvalidPassword = "password"
	StubShortPassword   = "a"
)

func TestLoginWorks(t *testing.T) {
	go runApp()
	RegisterTestingT(t)

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://%s:4444/wd/hub", host))
	Expect(err).NotTo(HaveOccurred())
	defer wd.Quit()

	wd.Get(app + "/")
	out, _ := wd.FindElement(selenium.ByCSSSelector, "#signin")
	outText, _ := out.Text()
	Expect(outText).To(ContainSubstring("Sign in to your account"))

	out, _ = wd.FindElement(selenium.ByCSSSelector, "#enter-password")
	out.SendKeys(StubInvalidPassword)
	out, _ = wd.FindElement(selenium.ByCSSSelector, "#login")
	out.Click()

	out, err = wd.FindElement(selenium.ByCSSSelector, "#error")
	Expect(err).NotTo(HaveOccurred())
	outText, _ = out.Text()
	Expect(outText).To(ContainSubstring("password is insecure"))

	out, _ = wd.FindElement(selenium.ByCSSSelector, "#enter-password")
	out.SendKeys(StubShortPassword)
	out, _ = wd.FindElement(selenium.ByCSSSelector, "#login")
	out.Click()

	out, _ = wd.FindElement(selenium.ByCSSSelector, "#error")
	outText, _ = out.Text()
	Expect(outText).To(ContainSubstring("password must be at least 8 characters"))

	out, err = wd.FindElement(selenium.ByCSSSelector, "#enter-password")
	Expect(err).NotTo(HaveOccurred())
	err = out.SendKeys(StubValidPassword)
	Expect(err).NotTo(HaveOccurred())

	out, err = wd.FindElement(selenium.ByCSSSelector, "#login")
	Expect(err).NotTo(HaveOccurred())
	err = out.Click()
	Expect(err).NotTo(HaveOccurred())

	out, err = wd.FindElement(selenium.ByCSSSelector, "#password")
	Expect(err).NotTo(HaveOccurred())
	outText, err = out.Text()
	Expect(err).NotTo(HaveOccurred())
	Expect(outText).To(Equal(StubValidPassword))
}
