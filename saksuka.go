package Saksuka

import "time"

// Some constants
const VERSION = "0.0.0"
const HttpBaseUrl = "https://api.dogehouse.tv";
const WebsocketBaseUrl = "wss://api.dogehouse.tv/socket";
const ConnectionTimeout = time.Second * 15;
const HeartbeatInterval = 8 * time.Second;