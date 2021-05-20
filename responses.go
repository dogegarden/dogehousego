package dogehousego

type OnNewMessageEvent struct {
	Msg Message `json:"msg"`
	UserId UUID `json:"userId"`
}

type OnNewRoomDetailsEvent struct {
	RoomId UUID `json:"roomId"`
	Name string `json:"name"`
	ChatThrottle int `json:"chatThrottle"`
	IsPrivate bool `json:"isPrivate"`
	Description string `json:"description"`
}

type OnUserJoinRoomEvent struct {
	User User `json:"user"`
	MuteMap BooleanMap `json:"muteMap"`
	DeafMap BooleanMap `json:"deafMap"`
	RoomId UUID `json:"roomId"`
}

type OnUserLeftRoomEvent struct {
	RoomId UUID `json:"roomId"`
	UserId UUID `json:"userId"`
}

type OnInvitationToRoomEvent struct {
	Type string `json:"type"`
	Username string `json:"username"`
	DisplayName string `json:"displayName"`
	AvatarUrl string `json:"avatarUrl"`
	BannerUrl string `json:"bannerUrl"`
	RoomName string `json:"roomName"`
	RoomId UUID `json:"roomId"`
}

type OnHandRaisedEvent struct {
	UserId UUID `json:"userId"`
}

type OnSpeakerAddedEvent struct {
	UserId UUID `json:"userId"`
	MuteMap BooleanMap `json:"muteMap"`
	DeafMap BooleanMap `json:"deafMap"`
}

type OnSpeakerRemovedEvent OnSpeakerAddedEvent;

type OnReadyEvent struct {
	User User `json:"user"`
}

type SearchResult struct {
	Items []interface{} `json:"items"`
	Rooms []Room `json:"rooms"`
	Users []User `json:"users"`
}

type GetMyScheduledRoomsAboutToStartResponse struct {
	ScheduledRooms []ScheduledRoom `json:"scheduledRooms"`
}

type JoinRoomAndGetInfoResponse struct {
	Room Room `json:"room"`
	Users []RoomUser `json:"users"`
	MuteMap BooleanMap `json:"muteMap"`
	DeafMap BooleanMap `json:"deafMap"`
	RoomId UUID `json:"roomId"`
	ActiveSpeakerMap BooleanMap `json:"activeSpeakerMap"`
}

type GetInviteListResponse struct {
	Users []User `json:"users"`
	NextCursor *int `json:"nextCursor"`
}

type GetFollowListResponse struct {
	Users []UserWithFollowInfo `json:"users"`
	NextCursor *int            `json:"nextCursor"`
}

type GetBlockedFromRoomUsersResponse GetInviteListResponse;

type GetMyFollowingResponse GetFollowListResponse;

type GetTopPublicRoomsResponse struct {
	Rooms []Room `json:"rooms"`
	NextCursor *int `json:"nextCursor"`
}

type GetUserProfileResponse struct {
	Rooms []Room `json:"rooms"`
}

type GetScheduledRoomsResponse struct {
	NextCursor NullString `json:"nextCursor"`
	ScheduledRooms []ScheduledRoom `json:"scheduledRooms"`
}

type UserCreateBotResponse struct {
	ApiKey *string `json:"apiKey"`
	IsUsernameToken *bool `json:"isUsernameToken"`
	Error *string `json:"error"`
}

type CreateScheduledRoomResponse struct {
	ScheduledRoom ScheduledRoom `json:"scheduledRoom"`
}

type EditScheduledRoomResponse CreateScheduledRoomResponse

type LeaveRoomResponse struct {
	RoomId UUID `json:"roomId"`
}

type CreateRoomResponse struct {
	Room Room `json:"room"`
}