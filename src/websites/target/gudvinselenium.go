package target

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"time"
)

type Selenium struct {
	service *selenium.Service
	driver  *selenium.WebDriver
}

func Connect(url string) Selenium {
	const port = 4444

	service, err := selenium.NewGeckoDriverService("files/geckodriver.exe", port)
	if err != nil {
		log.Errorf("Failed to start the WebDriver service: %v", err)
		panic(err)
	}

	caps := selenium.Capabilities{
		"browserName": "firefox",
	}
	ffCaps := firefox.Capabilities{}
	err = ffCaps.SetProfile("files/a4r2akqa.SeleniumFF")
	caps.AddFirefox(ffCaps)

	// Start a Tor browser session
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		log.Errorf("Failed to start the Tor browser session: %v", err)
		panic(err)
	}
	webDriver.SetPageLoadTimeout(time.Second * 60)

	err = webDriver.Get(url)
	if err != nil {
		log.Errorf("Failed to navigate to Google: %v", err)
		panic(err)
	}

	return Selenium{service: service, driver: &webDriver}
}

func (s *Selenium) Close() {
	s.service.Stop()
	(*s.driver).Quit()
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

func (s *Selenium) SolveCaptcha() {
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
