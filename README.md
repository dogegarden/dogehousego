
<p align="center">
  <img src="https://cdn.discordapp.com/attachments/820450983892222022/820961073980899328/dogegarden-bottom-cropped.png" alt="DogeGarden logo" />
</p>
<p align="center">
  <strong>A wrapper for DogeHouse made in Go 💨</strong>
</p>
<p align="center">
  <a href="https://discord.gg/Nu6KVjJYj6">
    <img src="https://img.shields.io/discord/820442045264691201?style=for-the-badge" alt="discord - users online" />
  </a>
</p>

<h3 align="center">
  <a href="https://dogegarden.net">Website</a>
  <span> · </span>  
  <a href="https://stats.dogegarden.net">Tracker</a>
  <span> · </span>
  <a href="https://discord.gg/Nu6KVjJYj6">Discord</a>
  <span> · </span>
  <a href="https://wiki.dogegarden.net">Documentation</a>
</h3>

---

<p align="center">
Need help or support?<br>
Join our <a href="https://discord.gg/Nu6KVjJYj6">Discord</a>.
</p>

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
	"github.com/dogegarden/dogehousego"
)

func main() {
	retval, err := dogehousego.Auth("your-api-key")
	if err != nil {
		fmt.Println(err.Error())
		return;
	}

	con := dogehousego.Connection{
		AccessToken:  retval.AccessToken,
		RefreshToken: retval.RefreshToken,
	}

	con.OnNewChatMessage(func(event dogehousego.OnNewMessageEvent, err error) {
		if err != nil {
			fmt.Println("New chat message error! " + err.Error());
			return;
		}

		message := dogehousego.TokensToString(event.Msg.Tokens)

		fmt.Println(message)
	});

	con.OnReady(func(event dogehousego.OnReadyEvent, err error) {

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

	fmt.Println(con.Start())
}
```
