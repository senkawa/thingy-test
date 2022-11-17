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

func TestXD(t *testing.T) {
	go runApp()
	RegisterTestingT(t)

	caps := selenium.Capabilities{"browserName": "chrome"}
	wd, err := selenium.NewRemote(caps, fmt.Sprintf("http://%s:4444/wd/hub", host))
	Expect(err).NotTo(HaveOccurred())
	defer wd.Quit()

	wd.Get(app+"/s/aaa")
	out, _ := wd.FindElement(selenium.ByCSSSelector, "body")
	outText, _ := out.Text()
	Expect(outText).To(ContainSubstring("aaa"))

	wd.Get(app+"/s/aaa?script=<script>alert('xss')</script>")
	out, err = wd.FindElement(selenium.ByCSSSelector, "body")
	Expect(err).NotTo(HaveOccurred())
	outText, err = out.Text()
	Expect(err).NotTo(HaveOccurred())
	Expect(outText).To(ContainSubstring("<script>alert('xss')</script>"))
}
