package target

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sweetRevenge/src/util"
	"time"
)

func orderItemWithCustomerSelenium(name, phone, itemId, link, socksProxy string) {
	//simpler to do with just tor
	log.Info(fmt.Sprintf("Sending order for (%s, %s, %s) via Selenium",
		name, phone, itemId))

	selenium := util.Connect(link, socksProxy)
	defer selenium.Close()

	selenium.Click("a.fn_fast_order_button")
	time.Sleep(time.Second * 3)
	selenium.Input("#fn_fast_order input.fn_validate_fast_name", name)
	selenium.Input("#fn_fast_order input.fn_validate_fast_phone", phone)
	solveYandexCaptcha(selenium)
	selenium.Click("#fn_fast_order input.fn_fast_order_submit")
	selenium.WaitForRedirect("/order/")

	log.Info("Sent order successfully")
}

func solveYandexCaptcha(selenium *util.Selenium) {
	doc := selenium.GetDocument()
	captcha := doc.Find("#fn_fast_order").Find("div.secret_number").Text()
	if captcha == "" {
		log.Info("No arithmetic captcha found, skipping captcha answer")
		return
	}

	log.Info("Arithmetic captcha found, solving: " + captcha)
	answer := util.SolveArithmeticCaptcha(captcha)
	log.Info("Captcha solution: ", answer)

	selenium.EnterCaptcha(answer)
}
