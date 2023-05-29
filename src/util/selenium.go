package util

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"strings"
	"time"
)

type Selenium struct {
	service *selenium.Service
	driver  *selenium.WebDriver
}

func Connect(url, socksProxy string) *Selenium {
	const port = 4444

	service, err := selenium.NewChromeDriverService("chromedriver", port)
	if err != nil {
		log.Errorf("Failed to start the WebDriver service: %v", err)
		panic(err)
	}

	caps := selenium.Capabilities{
		"browserName": "chrome",
	}
	caps.AddProxy(selenium.Proxy{
		Type:         selenium.Manual,
		SOCKS:        socksProxy,
		SOCKSVersion: 5,
	})
	cc := chrome.Capabilities{
		W3C: false,
		Args: []string{
			"--no-sandbox",
			"--headless", // Run Chrome in headless mode (without UI)
		},
	}
	caps.AddChrome(cc)

	// Start a Tor browser session
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		log.Errorf("Failed to start browser session: %v", err)
		panic(err)
	}
	webDriver.SetPageLoadTimeout(time.Minute * 3)
	webDriver.SetImplicitWaitTimeout(time.Minute * 3)

	err = webDriver.Get(url)
	if err != nil {
		log.Errorf("Failed to navigate to Google: %v", err)
		panic(err)
	}

	return &Selenium{service: service, driver: &webDriver}
}

func (s *Selenium) Close() {
	s.service.Stop()
	(*s.driver).Quit()
}

func (s *Selenium) GetDocument() *goquery.Document {
	webDriver := *s.driver
	html, _ := webDriver.PageSource()
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(html))
	return doc
}

func (s *Selenium) Click(selector string) {
	webDriver := *s.driver
	element, err := webDriver.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Errorf("Failed to find element with selector: %q, %v", selector, err)
		panic(err)
	}

	err = element.Click()
	if err != nil {
		log.Errorf("Failed to click element: %v", err)
		panic(err)
	}
}

func (s *Selenium) WaitForRedirect(url string) {
	webDriver := *s.driver

	condition := func(wd selenium.WebDriver) (bool, error) {
		currentUrl, err := wd.CurrentURL()
		if err != nil {
			return false, err
		}
		return strings.Contains(currentUrl, url), nil
	}

	err := webDriver.Wait(condition)
	if err != nil {
		log.Errorf("Wait for redirect failed: %v", err)
		panic(err)
	}
}

func (s *Selenium) Input(selector, value string) {
	webDriver := *s.driver
	element, err := webDriver.FindElement(selenium.ByCSSSelector, selector)
	if err != nil {
		log.Errorf("Failed to find element with selector: %q, %v", selector, err)
		panic(err)
	}

	err = element.SendKeys(value)
	if err != nil {
		log.Errorf("Failed to input: %v", err)
		panic(err)
	}
}

func (s *Selenium) EnterCaptcha(value string) {
	webDriver := *s.driver
	element, err := webDriver.FindElement(selenium.ByCSSSelector, "#fn_fast_order input.form__input_captcha")
	if err != nil {
		log.Errorf("Failed to find captcha field: %v", err)
		panic(err)
	}

	err = element.SendKeys(value)
	if err != nil {
		log.Errorf("Failed to enter captcha: %v", err)
		panic(err)
	}
}

func (s *Selenium) SolveReCaptcha() {
	webDriver := *s.driver

	captcha, err := webDriver.FindElement(selenium.ByCSSSelector, "iframe[title='reCAPTCHA']")
	if err != nil {
		log.Fatalf("Failed to find  captcha: %v", err)
	}
	webDriver.SwitchFrame(captcha)
	checkbox, err := webDriver.FindElement(selenium.ByCSSSelector, "#recaptcha-anchor")
	if err != nil {
		log.Fatalf("Failed to find captcha checkbox: %v", err)
	}
	checkbox.Click()
	webDriver.SwitchFrame(nil)

	time.Sleep(time.Second * 5)
	challenge, err := webDriver.FindElement(selenium.ByCSSSelector, "iframe[title='recaptcha challenge expires in two minutes']") //title="recaptcha challenge expires in two minutes"
	if err != nil {
		log.Fatalf("Failed to find captcha challenge frame: %v", err)
	}
	webDriver.SwitchFrame(challenge)

	button, err := webDriver.FindElement(selenium.ByCSSSelector, "div.help-button-holder")
	if err != nil {
		log.Fatalf("Failed to find the search input element: %v", err)
	}

	button.Click()
}
