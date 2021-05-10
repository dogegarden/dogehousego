package dogehousego

import (
	"fmt"
	"testing"
)

func TestConnection_Start(t *testing.T) {
	resp, err := Auth("afc40591-6a37-4ee2-968a-5f3f1761c49f");

	if err != nil {
		fmt.Println(err);
		return;
	}

	con := Connection{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		DebugLog:     false,
	}

	defer func() {
		con.Close();
	}()


	con.OnReady(func(event OnReadyEvent, err error) {

		resp, err := con.GetTopPublicRooms(0);
		if err != nil {
			fmt.Println(err)
			return;
		}
		
		_, err = con.JoinRoomAndGetInfo(resp.Rooms[0].Id);

		fmt.Println(err);
	})


	err = con.Start();

	if err != nil {
		fmt.Println(err)
		t.Fatalf(err.Error());
	}
}
