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
	ToRemove []*Listener
	DebugLog bool
}

func (con *Connection) indexOfListener(listener * Listener) (int) {
	for k, v := range con.Listeners {
		if listener == v {
			return k
		}
	}
	return -1    //not found.
}

func (con *Connection) removeListener(i int, j int) {
	con.Listeners[i] = con.Listeners[len(con.Listeners)-1]
	con.Listeners = con.Listeners[:len(con.Listeners)-1];

	con.ToRemove[j] = con.ToRemove[len(con.ToRemove)-1]
	con.ToRemove = con.ToRemove[:len(con.ToRemove)-1];
}

func (con *Connection) addListener(opCode string, handler func(listenerHandler ListenerHandler)) func() {
	listener := Listener{
		Opcode:  opCode,
		Handler: handler,
	}

	con.Listeners = append(con.Listeners, &listener);

	return func() { // Delete item from the list
		con.ToRemove = append(con.ToRemove, &listener);
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

	con.Connected = true;

	con.Socket  = c;

	con.Socket.SetReadLimit(4096 * 1024);

	go func() {
		for {
			if !con.Connected {
				break;
			}

			for k, v := range con.ToRemove {
				index := con.indexOfListener(v);
				if index == -1 {
					continue
				}
				con.removeListener(index, k);
			}

			_, message, err := con.Socket.ReadMessage();
			if err != nil {
				fmt.Println("Fatal error! Closing socket. Error: " + err.Error());
				con.Connected = false;
				break;
			}

			msg := string(message);

			if msg == "pong" || msg == `"pong""` {
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

			for _,v := range con.Listeners {
				if v.Opcode == op {
					var fetch string;

					if val, ok := obj["fetchId"].(string); ok {
						fetch = val;
					} else if val, ok := obj["ref"].(string); ok {
						fetch = val;
					} else {
						fetch = "unset";
					}

					var data interface{};

					if val, ok := obj["d"]; ok {
						data = val;
					} else {
						data = obj["p"];
					}

					var handlerError error = nil;

					if objErr, ok := obj["e"].(map[string]interface{}); ok {
						fmt.Println("TEST")
						errString := "";

						for k,v := range objErr {
							errString += fmt.Sprintf("%s %s", k, v);
						}

						handlerError = errors.New(errString);
					}

					go v.Handler(ListenerHandler{
						Data:    data,
						FetchId: fetch,
						Error: handlerError,
					})
				}
			}
		}
	}();

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

	var authData = map[string]interface{}{
		"accessToken": con.AccessToken,
		"refreshToken": con.RefreshToken,
		"reconnectToVoice": false,
		"currentRoomId": nil,
		"muted": false,
		"deafened": false,
	}

	con.send("auth", authData);

	for con.Connected {
		time.Sleep(time.Millisecond)
	}
	return nil;
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