package websocketserver

import (
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "localhost:8080", "http service address")

func checkOrigin(r *http.Request) bool {
	// origin := r.Header.Get("Origin")
	// origin == "http://localhost:3000" || origin == "http://localhost:4001" // TODO: env var this
	return true
}

var upgrader = websocket.Upgrader{
	CheckOrigin: checkOrigin,
} // use default options

func Echo(w http.ResponseWriter, r *http.Request) {
	log.SetFlags(0)
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("WEBSOCKET SERVER: upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("WEBSOCKET SERVER: read:", err)
			break
		}
		log.Printf("WEBSOCKET SERVER: received: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("ws://" + r.Host + "/echo")
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

// func main() {
// 	flag.Parse()
// 	log.SetFlags(0)
// 	http.HandleFunc("/echo", Echo)
// 	http.HandleFunc("/", Home)
// 	log.Fatal(http.ListenAndServe(*addr, nil))
// }

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<script>  
window.addEventListener("load", function(evt) {
    var output = document.getElementById("output");
    var input = document.getElementById("input");
    var ws;
    var print = function(message) {
        var d = document.createElement("div");
        d.textContent = message;
        output.appendChild(d);
        output.scroll(0, output.scrollHeight);
    };
    document.getElementById("open").onclick = function(evt) {
        if (ws) {
            return false;
        }
		console.log("hello")
		console.log("{{.}}")
        ws = new WebSocket("{{.}}");
        ws.onopen = function(evt) {
            print("OPEN");
        }
        ws.onclose = function(evt) {
            print("CLOSE");
            ws = null;
        }
        ws.onmessage = function(evt) {
            print("RESPONSE: " + evt.data);
        }
        ws.onerror = function(evt) {
            print("ERROR: " + evt.data);
        }
        return false;
    };
    document.getElementById("send").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        print("SEND: " + input.value);
        ws.send(input.value);
        return false;
    };
    document.getElementById("close").onclick = function(evt) {
        if (!ws) {
            return false;
        }
        ws.close();
        return false;
    };
});
</script>
</head>
<body>
<table>
<tr><td valign="top" width="50%">
<p>Click "Open" to create a connection to the server, 
"Send" to send a message to the server and "Close" to close the connection. 
You can change the message and send multiple times.
<p>
<form>
<button id="open">Open</button>
<button id="close">Close</button>
<p><input id="input" type="text" value="Hello world!">
<button id="send">Send</button>
</form>
</td><td valign="top" width="50%">
<div id="output" style="max-height: 70vh;overflow-y: scroll;"></div>
</td></tr></table>
</body>
</html>
`))