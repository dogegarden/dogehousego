package dogehousego

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"time"
)

type Token string;

type ListenerHandler struct {
	Data interface{} `json:"data"`
	FetchId string `json:"fetchId"`
	Error error `json:"error"`
}

type Listener struct {
	Opcode string `json:"opcode"`
	Handler func(listenerHandler ListenerHandler) `json:"handler"`
}

type Connection struct {
	Socket *websocket.Conn
	AccessToken string
	RefreshToken string
	Ready bool
	Connected bool
	Listeners []*Listener
	ToRemove chan *Listener
	DebugLog bool
}

func (con *Connection) indexOfListener(listener * Listener) int {
	for k, v := range con.Listeners {
		if listener == v {
			return k
		}
	}
	return -1    //not found.
}

func (con *Connection) removeListener(i int) {
	con.Listeners[i] = con.Listeners[len(con.Listeners)-1]
	con.Listeners = con.Listeners[:len(con.Listeners)-1];
}

func (con *Connection) addListener(opCode string, handler func(listenerHandler ListenerHandler)) func() {
	listener := Listener{
		Opcode:  opCode,
		Handler: handler,
	}

	con.Listeners = append(con.Listeners, &listener);

	return func() { // Delete item from the list
		con.ToRemove <- &listener;
	}
}

func (con *Connection) Start() error  {
	if(con.DebugLog) {
		fmt.Println("Opening websocket connection")
	}
	c, _, err := websocket.DefaultDialer.Dial(WebsocketBaseUrl, nil);

	if err != nil {
		return errors.New("Error connecting to websocket. Error: " + err.Error());
	}

	if(con.DebugLog) {
		fmt.Println("Socket opened")
	}

	con.ToRemove = make(chan *Listener, 100);

	con.Connected = true;

	con.Socket  = c;

	con.Socket.SetReadLimit(4096 * 1024);
	var socketErr error;
	go func() {
		for {
			if !con.Connected {
				break;
			}

			ok := false;
			for !ok { // Delete scheduled listeners
				select {

				case v := <- con.ToRemove:
					index := con.indexOfListener(v);
					if index == -1 {
						continue
					}
					con.removeListener(index);
				default: // No listeners left so break out of the loop
					ok = true;
				}

			}


			_, message, err := con.Socket.ReadMessage();
			if err != nil {
				fmt.Println("Fatal error! Closing socket. Error: " + err.Error());
				socketErr = err;
				con.Connected = false;
				break;
			}

			msg := string(message);

			if msg == "pong" || msg == `"pong""` { // We don't want to do anything with those
				continue;
			}

			if con.DebugLog {
				fmt.Println(msg)
			}

			var obj map[string]interface{};
			json.Unmarshal(message, &obj);

			if val, ok := obj["errors"]; ok {
				fmt.Printf("ERRORS WHEN READING WEBSOCKET DATA: %v\n", val);
				continue;
			}

			op := obj["op"].(string);

			if op == "auth-good" {
				fmt.Println("authorization ok")
				con.Ready = true;
			}

			// Execute listeners
			for _,v := range con.Listeners {
				if v.Opcode == op {
					var fetch string;

					// We need this ugly crap to handle both versions of request
					if val, ok := obj["fetchId"].(string); ok {
						fetch = val;
					} else if val, ok := obj["ref"].(string); ok {
						fetch = val;
					} else {
						fetch = "unset";
					}

					var data interface{};

					// Also 2 versions of api
					if val, ok := obj["d"]; ok {
						data = val;
					} else {
						data = obj["p"];
					}

					var handlerError error = nil;

					// Handle any errors
					if objErr, ok := obj["e"].(map[string]interface{}); ok {
						errString := "";

						for k,v := range objErr {
							errString += fmt.Sprintf("%s %s", k, v);
						}

						handlerError = errors.New(errString);
					}

					// Run as goroutine to not hang main loop
					go v.Handler(ListenerHandler{
						Data:    data,
						FetchId: fetch,
						Error: handlerError,
					})
				}
			}
		}
	}();

	// Loop to send ping messages
	go func() {
		for {
			if !con.Connected {
				break;
			}

			time.Sleep(HeartbeatInterval);
			if err := con.Socket.WriteMessage(websocket.TextMessage, []byte("ping")); err != nil {
				break;
			}
		}
	}()

	// Send authentication data
	var authData = map[string]interface{}{
		"accessToken": con.AccessToken,
		"refreshToken": con.RefreshToken,
		"reconnectToVoice": false,
		"currentRoomId": nil,
		"muted": false,
		"deafened": false,
	}

	con.send("auth", authData);

	// Keep main routine alive while connection is open
	for con.Connected {
		time.Sleep(time.Millisecond)
	}
	return socketErr;
}

func (con *Connection) send(op string, data interface{}, params ...string) {
	if !con.Connected {
		return
	}

	str, err := json.Marshal(&data);
	if err != nil {
		panic(err);
	}

	additionStr := "";

	if len(params) > 0 {
		additionStr = fmt.Sprintf(`, "fetchId":"%s"`, params[0]);
	}

	raw := fmt.Sprintf(`{"op":"%s", "d":%s %s}`, op, string(str), additionStr);
	err = con.Socket.WriteMessage(websocket.TextMessage, []byte(raw));

	if err != nil {
		fmt.Println("Failed to send ws message. Error: " + err.Error());
		con.Connected = false;
	}
}

func (con *Connection) sendCast(op string, data interface{}, params ...string) {
	if !con.Connected {
		return
	}

	str, err := json.Marshal(&data);
	if err != nil {
		panic(err);
	}

	additionStr := "";

	if len(params) > 0 {
		additionStr = fmt.Sprintf(`, "ref":"%s"`, params[0]);
	}

	raw := fmt.Sprintf(`{"v":"0.2.0", "op":"%s","p":%s %s}`, op, str, additionStr);

	err = con.Socket.WriteMessage(websocket.TextMessage, []byte(raw));

	if err != nil {
		fmt.Println("Failed to send cast ws message. Error: " + err.Error());
		con.Connected = false;
	}
}

func (con *Connection) sendCall(op string, data interface{}, params ...string) (*ListenerHandler, error) {
	if !con.Connected {
		return nil, errors.New("connection is closed");
	}

	doneOpCode := op + ":reply";
	ref := "";

	if len(params) > 0 {
		doneOpCode = params[0];
	} else {
		uid, err := uuid.NewV4();
		if err != nil {
			return nil, err;
		}
		ref = fmt.Sprintf("%s", uid);
	}

	response := make(chan ListenerHandler);


	unsub := con.addListener(doneOpCode, func(listenerHandler ListenerHandler) {
		if !(len(params) > 0) && listenerHandler.FetchId != ref {
			return;
		}

		response <- listenerHandler;
	});

	if ref != "" {
		con.sendCast(op, data, ref);
	} else {
		con.sendCast(op, data);
	}

	select {
		case resp := <- response:
			unsub();
			return &resp, nil;
		case <- time.After(ConnectionTimeout):
			unsub();
			return nil, errors.New("timed out");
	}
}

func (con *Connection) fetch(op string, data interface{}, params ...string) (*ListenerHandler, error) {
	if !con.Connected {
		return nil, errors.New("connection is closed");
	}

	doneOpCode := "fetch_done";
	ref := "";

	if len(params) > 0 {
		doneOpCode = params[0];
	} else {
		uid, err := uuid.NewV4();
		if err != nil {
			return nil, err;
		}
		ref = fmt.Sprintf("%s", uid);
	}

	response := make(chan ListenerHandler);


	unsub := con.addListener(doneOpCode, func(listenerHandler ListenerHandler) {
		if !(len(params) > 0) && listenerHandler.FetchId != ref {
			return;
		}

		response <- listenerHandler;
	});

	if ref != "" {
		con.send(op, data, ref);
	} else {
		con.send(op, data);
	}

	select {
		case resp := <- response:
			unsub();
			return &resp, nil;
		case <- time.After(ConnectionTimeout):
			unsub();
			return nil, errors.New("timed out");
	}
}

func (con *Connection) Close() {
	con.Socket.Close();
	con.Ready = false;
	con.Connected = false;
}