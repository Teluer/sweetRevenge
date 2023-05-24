package admin

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
	"sweetRevenge/src/config"
	"sweetRevenge/src/rabbitmq"
	"sync"
	"time"
)

var mu sync.Mutex
var templ = `<!DOCTYPE html>
    <html>
    <head>
        <title>SweetRevenge control panel</title>
    </head>
    <body>
        <h1>Update Configs</h1>
        <form method="POST" action="/conf">
            <label for="name">Order Max interval in munites:</label>
            <input type="text" id="frequency" name="frequency" value={{.SendOrdersMaxInterval}}  required>
            <br>
            <label for="email">Enable sending orders:</label>
            <input type="checkbox" id="shouldSend" name="shouldSend" value="shouldSend" checked={{.SendOrdersEnabled}}>
            <br>
            <input type="submit" value="Submit">
        </form>

        <h1>Send Manual Order</h1>
        <form method="POST" action="/order">
            <label for="name">Name:</label>
            <input type="text" id="name" name="name" required>
            <br>
            <label for="phone">Enable sending orders:</label>
            <input type="text" id="phone" name="phone" required>
            <br>
            <input type="submit" value="Submit">
        </form>
    </body>
    </html>`

func ControlPanel(cfg *config.OrdersRoutineConfig) {
	log.Infof("Starting Control Panel server at: %s", "localhost:8008/admin")

	http.HandleFunc("/admin", controlPanelHandler(cfg)) // each request calls handler
	http.HandleFunc("/conf", configHandler(cfg))
	http.HandleFunc("/order", orderHandler)
	
	log.Error(http.ListenAndServe("localhost:8008", nil))
}

// handler echoes the Path component of the request URL r.
func controlPanelHandler(cfg *config.OrdersRoutineConfig) func(http.ResponseWriter, *http.Request) {
	// Create a template from the HTML
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()

		t, err := template.New("form").Parse(templ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = t.Execute(w, *cfg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mu.Unlock()
	}
}

func configHandler(cfg *config.OrdersRoutineConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//imagine error handling here
		frequencyMinutes, _ := strconv.Atoi(r.FormValue("frequency"))
		ordersEnabled, _ := strconv.ParseBool(r.FormValue("shouldSend"))

		cfg.SendOrdersMaxInterval = time.Minute * time.Duration(frequencyMinutes)
		cfg.SendOrdersEnabled = ordersEnabled

		message := "Config changes will take effect after the next scheduled order is sent!"
		w.Write([]byte(message))
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	order := rabbitmq.ManualOrder{
		Name:  name,
		Phone: phone,
	}

	err := rabbitmq.Publish(&order)

	if err != nil {
		w.Write([]byte("Manual order submission failed! See logs for more info."))
	} else {
		w.Write([]byte("Manual order will be sent after the next scheduled order!"))
	}
}
