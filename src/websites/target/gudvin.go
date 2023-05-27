package target

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
)

//TODO: call this if captcha is enabled and regular order fails.
func orderItemWithCustomer(name, phone, itemId, link string) {
	PanicIfVpnNotEnabled()

	log.Info(fmt.Sprintf("Sending order for (%s, %s, %s)",
		name, phone, itemId))

	selenium := Connect(link)
	defer selenium.Close()

	selenium.Click("a.fn_fast_order_button")
	time.Sleep(time.Second * 3)
	selenium.Input("input.fn_validate_fast_name", name)
	selenium.Input("input.fn_validate_fast_phone", phone)
	selenium.Click("input.fn_fast_order_submit")

	selenium.SolveCaptcha()

	//todo: check if order success page was opened
	log.Info("Sent order successfully")
}
