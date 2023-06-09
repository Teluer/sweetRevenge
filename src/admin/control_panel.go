package admin

import (
	log "github.com/sirupsen/logrus"
	"html/template"
	"net/http"
	"strconv"
	"sweetRevenge/src/config"
	"sweetRevenge/src/db/dao"
	"sweetRevenge/src/rabbitmq"
	"sweetRevenge/src/websites/orsen"
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
            <label for="name">OrderSender Max interval in munites:</label>
            <br>
            <input type="text" id="frequency" name="frequency" value={{.OrdersInterval}}  required>
            <br>
            <label for="email">Enable sending orders: </label>
            <input type="checkbox" id="shouldSend" name="shouldSend" {{if .OrdersEnabled}} checked {{end}}>
            <br>
            <input type="submit" value="Submit">
            <br>
            <label class="response-label" style="color:red;"></label>
        </form>

        <h1>Send Manual OrderSender</h1>
        <form method="POST" id="OrderForm" action="/order">
            <label for="name">Name:</label>
            <br>
            <input type="text" id="name" name="name" value="">
            <br>
            <label for="phone">Phone:</label>
            <br>
            <input type="text" id="phone" name="phone" value="">
            <br>
            <input type="submit" value="Submit">
            <br>
            <label class="response-label" style="color:red;"></label>
        </form>
        <h1>Generate Customer</h1>
        <form method="POST" id="CustomerForm" action="/customer">
            <label>Warning: this marks phone and name as used!</label>
            <br>
            <input type="submit" value="Generate">
            <br>
            <label class="response-label" style="color:red;"></label>
        </form>
        <script>
          listener = function(event) {
			event.preventDefault(); // Prevent the default form submission
		
			var form = event.target;
			var formData = new FormData(form);
		
			fetch(form.action, {
			  method: form.method,
			  body: formData
			})
			.then(function(response) {
			  return response.text();
			})
			.then(function(responseText) {
              document.querySelector('#' + form.id + ' .response-label').textContent = responseText;
			})
			.catch(function(error) {
			  console.error('Error:', error);
			});
		  };
		  document.getElementById("ConfForm").addEventListener("submit", listener);
          document.getElementById("OrderForm").addEventListener("submit", listener);
          document.getElementById("CustomerForm").addEventListener("submit", listener);
		</script>
    </body>
    </html>`

func StartControlPanelServer(cfg *config.OrdersRoutineConfig) {
	log.Infof("Starting Control Panel server at: %s", "localhost:8008/admin")

	http.HandleFunc("/admin", controlPanelHandler(cfg)) // each request calls handler
	http.HandleFunc("/conf", configHandler(cfg))
	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/customer", customerHandler(cfg.OrdersCfg.PhonePrefixes))

	log.Error(http.ListenAndServe("0.0.0.0:8008", nil))
}

func controlPanelHandler(cfg *config.OrdersRoutineConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		log.Info("Accessing control panel from: ", r.RemoteAddr)

		data := &struct {
			OrdersInterval float64
			OrdersEnabled  bool
		}{
			OrdersInterval: float64(cfg.SendOrdersMaxInterval.Nanoseconds()) / float64(time.Minute),
			OrdersEnabled:  cfg.SendOrdersEnabled,
		}

		t, err := template.New("form").Parse(templ)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err = t.Execute(w, data)
		if err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		mu.Unlock()
	}
}

func configHandler(cfg *config.OrdersRoutineConfig) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		//imagine error handling here
		frequencyMinutes, _ := strconv.ParseFloat(r.FormValue("frequency"), 64)
		ordersEnabled := r.FormValue("shouldSend") == "on"

		cfg.SendOrdersMaxInterval = time.Duration(float64(time.Minute) * frequencyMinutes)
		cfg.SendOrdersEnabled = ordersEnabled

		log.Infof("Updated configs: interval=%v, sending=%v", cfg.SendOrdersMaxInterval, cfg.SendOrdersEnabled)
		message := "Configs updated. changes will take effect after the next scheduled order is sent!"
		w.Write([]byte(message))
	}
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	phone := r.FormValue("phone")

	if name == "" && phone == "" {
		w.Write([]byte("Both fields cannot be blank!"))
		return
	}

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

func customerHandler(phonePrefixes string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name, phone := orsen.CreateRandomCustomer(dao.Dao, phonePrefixes)
		w.Write([]byte(name + ", " + phone))
	}
}
