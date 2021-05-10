package dogehousego

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"time"
)

//Events
func (con *Connection) OnNewChatMessage(handler func(event OnNewMessageEvent, err error)) {
	con.addListener("new_chat_msg", func(listenerHandler ListenerHandler) {
		var invokeArg OnNewMessageEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnNewRoomDetails(handler func(event OnNewRoomDetailsEvent, err error)) {
	con.addListener("new_room_details", func(listenerHandler ListenerHandler) {
		var invokeArg OnNewRoomDetailsEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnUserJoinRoom(handler func(event OnUserJoinRoomEvent, err error)) {
	con.addListener("new_user_join_room", func(listenerHandler ListenerHandler) {
		var invokeArg OnUserJoinRoomEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnUserLeaveRoom(handler func(event OnUserLeftRoomEvent, err error)) {
	con.addListener("user_left_room", func(listenerHandler ListenerHandler) {
		var invokeArg OnUserLeftRoomEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnInvitationToRoom(handler func(event OnInvitationToRoomEvent, err error)) {
	con.addListener("invitation_to_room", func(listenerHandler ListenerHandler) {
		var invokeArg OnInvitationToRoomEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnHandRaised(handler func(event OnHandRaisedEvent, err error)) {
	con.addListener("hand_raised", func(listenerHandler ListenerHandler) {
		var invokeArg OnHandRaisedEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnSpeakerAdded(handler func(event OnSpeakerAddedEvent, err error)) {
	con.addListener("speaker_added", func(listenerHandler ListenerHandler) {
		var invokeArg OnSpeakerAddedEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnSpeakerRemoved(handler func(event OnSpeakerRemovedEvent, err error)) {
	con.addListener("speaker_removed", func(listenerHandler ListenerHandler) {
		var invokeArg OnSpeakerRemovedEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) OnReady(handler func(event OnReadyEvent, err error)) {
	con.addListener("auth-good", func(listenerHandler ListenerHandler) {
		var invokeArg OnReadyEvent;
		err := mapstructure.Decode(listenerHandler.Data, &invokeArg);
		go handler(invokeArg, err);
	})
}

func (con *Connection) Search(query string) (SearchResult, error){
	resp, err := con.fetch("search", map[string]interface{}{
		"query": query,
	});
	var retval SearchResult;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) GetMyScheduledRoomsAboutToStart(query string) (GetMyScheduledRoomsAboutToStartResponse, error){
	resp, err := con.fetch("get_my_scheduled_rooms_about_to_start", map[string]interface{}{
		"query": query,
	});
	var retval GetMyScheduledRoomsAboutToStartResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) JoinRoomAndGetInfo(roomId string) (JoinRoomAndGetInfoResponse, error) {
	resp, err := con.fetch("join_room_and_get_info", map[string]interface{}{
		"roomId": roomId,
	});

	var retval JoinRoomAndGetInfoResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}


func (con *Connection) GetInviteList(cursor int) (GetFollowListResponse, error){
	resp, err := con.fetch("get_invite_list", map[string]interface{}{
		"cursor": cursor,
	});
	var retval GetFollowListResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) GetFollowList(cursor int) (GetFollowListResponse, error){
	resp, err := con.fetch("get_follow_list", map[string]interface{}{
		"cursor": cursor,
	});
	var retval GetFollowListResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) GetBlockedFromRoomUsers(cursor int) (GetBlockedFromRoomUsersResponse, error){
	resp, err := con.fetch("get_blocked_from_room_users", map[string]interface{}{
		"cursor": cursor,
	});
	var retval GetBlockedFromRoomUsersResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) GetMyFollowing(cursor int) (GetMyFollowingResponse, error){
	resp, err := con.fetch("get_my_following", map[string]interface{}{
		"cursor": cursor,
	});
	var retval GetMyFollowingResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) GetTopPublicRooms(cursor int) (GetTopPublicRoomsResponse, error){
	resp, err := con.fetch("get_top_public_rooms", map[string]interface{}{
		"cursor": cursor,
	});
	fmt.Println()
	var retval GetTopPublicRoomsResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) GetUserProfile(idOrUsername string) (*UserWithFollowInfo, error){
	resp, err := con.fetch("get_user_profile", map[string]interface{}{
		"userId": idOrUsername,
	});
	var retval UserWithFollowInfo;

	if err != nil {
		return nil, err;
	}

	if resp.Data == nil {
		return nil, errors.New("could not find specified user");
	}

	if val, ok := resp.Data.(map[string]interface{})["error"]; ok {
		fmt.Println(val)
		return nil, errors.New(val.(string));
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return &retval, err;
}

func (con *Connection) GetScheduledRooms(cursor string, getOnlyMyScheduledRooms bool) (GetScheduledRoomsResponse, error){
	resp, err := con.fetch("get_scheduled_rooms", map[string]interface{}{
		"cursor": cursor,
		"getOnlyMyScheduledRooms": getOnlyMyScheduledRooms,
	});
	var retval GetScheduledRoomsResponse;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

// TODO: Not implemented (and exported) as i do not have a good idea to add Partial<User> in golang
// The only "good" way i can think of is to just pass map[string]interface and pass it to the ws
// but that will make it really easy to make mistakes.
// I would use struct just like in RoomUpdate but there is no documentation on fields that you
// can accualy change, and I can't read elixir code to find out myself
func (con *Connection) userUpdate() {
	return;
}

func (con *Connection) UserBlock(userId string) error {
	_, err := con.sendCall("user:block", map[string]interface{}{
		"userId": userId,
	});

	return err;
}

func (con *Connection) UserUnblock(userId string) error {
	_, err := con.sendCall("user:unblock", map[string]interface{}{
		"userId": userId,
	});

	return err;
}

// Initialization of RoomUpdateArgs is a pain in the ass but i don't have any ideas how to make is work
// and not make it even more confusing than it is
func (con *Connection) RoomUpdate(args RoomUpdateArgs) error {
	_, err := con.sendCall("room:update", args);

	return err;
}

func (con *Connection) RoomBan(userId string, shouldBanIp bool) {
	con.sendCast("room:update", map[string]interface{}{
		"userId": userId,
		"shouldBanIp": shouldBanIp,
	});
}

func (con *Connection) SetDeaf(deaf bool) error {
	_, err := con.sendCall("room:deafen", map[string]interface{}{
		"deafened": deaf,
	});

	return err;
}

func (con *Connection) UserCreateBot(username string) (*UserCreateBotResponse, error) {
	resp, err := con.sendCall("user:create_bot", map[string]interface{}{
		"username": username,
	});

	var retval UserCreateBotResponse;

	if err != nil {
		return nil, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	if retval.Error != nil {
		return nil, errors.New(*retval.Error);
	}

	return &retval, err;
}

func (con *Connection) Ban(username string, reason string) {
	con.send("ban", map[string]interface{}{
		"username": username,
		"reason": reason,
	});
}

func (con *Connection) DeleteScheduledRoom(id string) error {
	_, err := con.fetch("delete_scheduled_room", map[string]interface{}{
		"id": id,
	});

	return err;
}

func (con *Connection) CreateRoomFromScheduledRoom(id, name, description string) (Room, error) {
	resp, err := con.fetch("create_room_from_scheduled_room", map[string]interface{}{
		"id": id,
		"name": name,
		"description": description,
	});

	var retval Room;

	if err != nil {
		return retval, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return retval, err;
}

func (con *Connection) CreateScheduledRoom(name, description string, scheduledFor time.Time) (*CreateScheduledRoomResponse, error) {
	resp, err := con.fetch("schedule_room", map[string]interface{}{
		"name": name,
		"description": description,
		"scheduledFor": scheduledFor.Format(time.RFC3339),
	});

	var retval CreateScheduledRoomResponse;

	if err != nil {
		return nil, err;
	}

	if val, ok := resp.Data.(map[string]interface{})["error"]; ok {
		fmt.Println(val)
		return nil,errors.New(val.(string));
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return &retval, err;
}

func (con *Connection) EditScheduledRoom(id, name, description string, scheduledFor time.Time) (*EditScheduledRoomResponse, error) {
	resp, err := con.fetch("edit_scheduled_room", map[string]interface{}{
		"id": id,
		"data": map[string]interface{} {
			"name": name,
			"description": description,
			"scheduledFor": scheduledFor.Format(time.RFC3339),
		},
	});

	var retval EditScheduledRoomResponse;

	if err != nil {
		return nil, err;
	}

	if val, ok := resp.Data.(map[string]interface{})["error"]; ok {
		fmt.Println(val)
		return nil,errors.New(val.(string));
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return &retval, err;
}

func (con *Connection) AskToSpeak() {
	con.send("ask_to_speak", map[string]interface{}{});
}

func (con *Connection) InviteToRoom(userId string) {
	con.send("invite_to_room", map[string]interface{}{
		"userId": userId,
	});
}

func (con *Connection) SpeakingChange(value bool) {
	con.send("invite_to_room", map[string]interface{}{
		"value": value,
	});
}

func (con *Connection) UnbanFromRoom(userId string) error {
	_, err := con.fetch("unban_from_room", map[string]interface{}{
		"userId ": userId,
	});

	return err;
}

func (con *Connection) Follow(userId string) error {
	_, err := con.fetch("follow", map[string]interface{}{
		"userId ": userId,
	});

	return err;
}

func (con *Connection) SendRoomChatMessage(message string, whisperedTo ...UUID) {
	con.send("send_room_chat_msg", map[string]interface{}{
		"tokens ": StringToToken(message),
		"whisperedTo": whisperedTo,
	});
}

func (con *Connection) ChangeModStatus(userId string, value bool) {
	con.send("change_mod_status", map[string]interface{}{
		"userId ": userId,
		"value": value,
	});
}

func (con *Connection) ChangeRoomCreator(userId string) {
	con.send("change_room_creator", map[string]interface{}{
		"userId ": userId,
	});
}

func (con *Connection) AddSpeaker(userId string) {
	con.send("add_speaker", map[string]interface{}{
		"userId ": userId,
	});
}

func (con *Connection) DeleteRoomChatMessage(userId, messageId string) {
	con.send("delete_room_chat_message", map[string]interface{}{
		"userId ": userId,
		"messageId": messageId,
	});
}

func (con *Connection) UnbanFromRoomChat(userId string) {
	con.send("unban_from_room_chat", map[string]interface{}{
		"userId ": userId,
	});
}

func (con *Connection) BanFromRoomChat(userId string) {
	con.send("ban_from_room_chat", map[string]interface{}{
		"userId ": userId,
	});
}

func (con *Connection) SetListener(userId string) {
	con.send("set_listener", map[string]interface{}{
		"userId ": userId,
	});
}

func (con *Connection) SetMute(isMuted bool) error {
	_, err := con.fetch("mute", map[string]interface{}{
		"value": isMuted,
	});

	return err;
}

func (con *Connection) LeaveRoom() (*LeaveRoomResponse, error) {
	resp, err := con.fetch("leave_room", map[string]interface{}{});

	var retval LeaveRoomResponse;

	if err != nil {
		return nil, err;
	}

	err = mapstructure.Decode(resp.Data, &retval);

	return &retval, err;
}