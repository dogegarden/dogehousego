# Saksuka
Golang wrapper for dogehouse API

Done:
- Websocket connection
- Api key authentication
- Event subscription
- Query ws events
- About 50% of mutations

TODO:
- The rest of mutations
- Whole audio module
- Docs

Example usage:

```go
package main

import (
	"fmt"
	"github.com/Explooosion-code/Saksuka"
)

func main() {
	retval, err := Saksuka.Auth("your-api-key")
	if err != nil {
		fmt.Println(err.Error())
		return;
	}

	con := Saksuka.Connection{
		AccessToken:  retval.AccessToken,
		RefreshToken: retval.RefreshToken,
	}

	con.OnNewChatMessage(func(event Saksuka.OnNewMessageEvent, err error) {
		if err != nil {
			fmt.Println("New chat message error! " + err.Error());
			return;
		}

		message := Saksuka.TokensToString(event.Msg.Tokens)

		fmt.Println(message)
	});

	con.OnReady(func(event Saksuka.OnReadyEvent, err error) {

		if err != nil {
			fmt.Println("Error while logging in!: " + err.Error())
		}

		fmt.Println("Logged in as: " + event.User.Username)
		rooms, err := con.GetTopPublicRooms(0)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		info, err := con.JoinRoomAndGetInfo(rooms.Rooms[0].Id)

		if err != nil {
			fmt.Println("Error while joining the room! " + err.Error())
			return
		}

		fmt.Printf("Joined room: %s | User count: %d\n", info.Room.Name, len(info.Users));
	});

	fmt.Println(con.Start()) // This returns an error if present
}
```
