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
        <form method="POST" id="ConfForm" action="/conf">
            <label class="response-label" style="color:red;"></label>
            <br>
            <label for="name">Order Max interval in munites:</label>
            <br>
            <input type="text" id="frequency" name="frequency" value={{.OrdersInterval}}  required>
            <br>
            <label for="email">Enable sending orders:</label>
            <br>
            <input type="checkbox" id="shouldSend" name="shouldSend" {{if .OrdersEnabled}} checked {{end}}>
            <br>
            <input type="submit" value="Submit">
        </form>

        <h1>Send Manual Order</h1>
        <form method="POST" id="OrderForm" action="/order">
            <label class="response-label" style="color:red;"></label>
            <br>
            <label for="name">Name:</label>
            <br>
            <input type="text" id="name" name="name" value="" required>
            <br>
            <label for="phone">Phone:</label>
            <br>
            <input type="text" id="phone" name="phone" value="" required>
            <br>
            <input type="submit" value="Submit">
        </form>
        <script>
          listener = function(event) {
			event.preventDefault(); // Prevent the default form submission
		
			var form = event.target;
			var formData = new FormData(form);
		
			// Send a POST request to the server
			fetch(form.action, {
			  method: form.method,
			  body: formData
			})
			.then(function(response) {
			  return response.text(); // Extract the response text
			})
			.then(function(responseText) {
			  // Update the label with the response message
              document.querySelector('#' + form.id + ' .response-label').textContent = responseText;
			})
			.catch(function(error) {
			  console.error('Error:', error);
			});
		  };
		  document.getElementById("ConfForm").addEventListener("submit", listener);
          document.getElementById("OrderForm").addEventListener("submit", listener);
		</script>
    </body>
    </html>`

type formData struct {
	OrdersInterval int
	OrdersEnabled  bool
}

func ControlPanel(cfg *config.OrdersRoutineConfig) {
	log.Infof("Starting Control Panel server at: %s", "localhost:8008/admin")

	http.HandleFunc("/admin", controlPanelHandler(cfg)) // each request calls handler
	http.HandleFunc("/conf", configHandler(cfg))
	http.HandleFunc("/order", orderHandler)

	log.Error(http.ListenAndServe("localhost:8000", nil))
}

// handler echoes the Path component of the request URL r.
func controlPanelHandler(cfg *config.OrdersRoutineConfig) func(http.ResponseWriter, *http.Request) {
	// Create a template from the HTML
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		log.Info("Accessing control panel from: ", r.RemoteAddr)

		data := &formData{
			OrdersInterval: int(cfg.SendOrdersMaxInterval.Minutes()),
			OrdersEnabled:  cfg.SendOrdersEnabled,
		}

		t, err := template.New("form").Parse(templ)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, data)
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
		ordersEnabled := r.FormValue("shouldSend") == "on"

		cfg.SendOrdersMaxInterval = time.Minute * time.Duration(frequencyMinutes)
		cfg.SendOrdersEnabled = ordersEnabled

		message := "Configs updated. changes will take effect after the next scheduled order is sent!"
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
		w.Write([]byte("Manual order submitted. Will be sent at scheduled time."))
	}
}