package target

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/firefox"
	"time"
)

func Example() {
	// Start a Selenium WebDriver service (e.g., ChromeDriver, GeckoDriver, etc.)
	// Replace "geckodriver" with the appropriate WebDriver executable for Tor browser
	const port = 4444

	service, err := selenium.NewGeckoDriverService("files/geckodriver.exe", port)
	if err != nil {
		log.Fatalf("Failed to start the WebDriver service: %v", err)
	}
	defer service.Stop()

	// Create capabilities for Tor browser
	caps := selenium.Capabilities{
		"browserName": "firefox",
	}
	//caps.AddProxy(selenium.Proxy{
	//	Type:         selenium.Manual,
	//	SOCKS:        "127.0.0.1:1080",
	//	SOCKSVersion: 5,
	//})

	prefs := map[string]interface{}{}
	prefs["network.proxy.type"] = 1
	prefs["network.proxy.socks"] = "localhost"
	prefs["network.proxy.socks_port"] = 1080

	ffCaps := firefox.Capabilities{
		Prefs:  prefs,
		Binary: "C:\\Projects\\Tor Browser11\\Browser\\firefox.exe",
	}
	//
	//err = ffCaps.SetProfile("files/profile/profile.default")
	err = ffCaps.SetProfile("C:\\Projects\\Tor Browser11\\Browser\\TorBrowser\\Data\\Browser\\profile.default")
	caps.AddFirefox(ffCaps)

	// Start a Tor browser session
	webDriver, err := selenium.NewRemote(caps, fmt.Sprintf("http://localhost:%d", port))
	if err != nil {
		log.Fatalf("Failed to start the Tor browser session: %v", err)
	}
	defer webDriver.Quit()
	webDriver.SetPageLoadTimeout(time.Second * 120)

	connectButton, err := webDriver.FindElement(selenium.ByCSSSelector, "#connectButton")

	jsCode := "arguments[0].removeAttribute(arguments[1]);"
	_, err = webDriver.ExecuteScript(jsCode, []interface{}{connectButton, "hidden"})
	if err != nil {
		log.Fatal(err)
	}
	//jsCode := "arguments[0].click();"
	//_, err = webDriver.ExecuteScript(jsCode, []interface{}{connectButton})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//err = connectButton.Click()

	time.Sleep(time.Second * 10)
	// Navigate to Google
	err = webDriver.Get("https://recaptcha-demo.appspot.com/recaptcha-v2-checkbox-explicit.php")
	if err != nil {
		log.Fatalf("Failed to navigate to Google: %v", err)
	}

	// Find the search input element and enter the search query
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

	// Wait for the search results to load
	//time.Sleep(2 * time.Second)

	//
	//searchResults, err := webDriver.FindElements(selenium.ByCSSSelector, "div.r a")
	//if err != nil {
	//	log.Fatalf("Failed to find the search result links: %v", err)
	//}
	//
	//// Select a random search result from the first page
	//rand.Seed(time.Now().UnixNano())
	//randomIndex := rand.Intn(len(searchResults))
	//randomResult := searchResults[randomIndex]
	//
	//// Click on the random search result
	//err = randomResult.Click()
	//if err != nil {
	//	log.Fatalf("Failed to click on the random search result: %v", err)
	//}

	// Wait for some time to see the result
	time.Sleep(5 * time.Second)
}
