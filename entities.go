package dogehousego

import "bytes"

type UUID string;
type NullString string;
type ChatMode string;
type WhisperPrivacySetting string;

const (
	ChatModeDefault      ChatMode = "default"
	ChatModeDisabled     ChatMode = "disabled"
	ChatModeFollowerOnly ChatMode = "follower_only"

	WhisperPrivacyOn WhisperPrivacySetting = "on"
	WhisperPrivacyOff WhisperPrivacySetting = "off"
)

// Make UUID and NullString types nullable
func (c UUID) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(string(c)) == 0 {
		buf.WriteString(`null`)
	} else {
		buf.WriteString(`"` + string(c) + `"`)
	}
	return buf.Bytes(), nil
}

func (c *UUID)UnmarshalJSON(in []byte) error {
	str := string(in)
	if str == `null` {
		*c = ""
		return nil
	}
	res := UUID(str)
	if len(res) >= 2 {
		res = res[1:len(res)-1]
	}
	*c = res
	return nil
}

func (c NullString) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	if len(string(c)) == 0 {
		buf.WriteString(`null`)
	} else {
		buf.WriteString(`"` + string(c) + `"`)
	}
	return buf.Bytes(), nil
}

func (c *NullString)UnmarshalJSON(in []byte) error {
	str := string(in)
	if str == `null` {
		*c = ""
		return nil
	}
	res := NullString(str)
	if len(res) >= 2 {
		res = res[1:len(res)-1]
	}
	*c = res
	return nil
}

type UserPreview struct  {
	Id UUID `json:"id"`
	DisplayName string `json:"displayName"`
	NumFollowers int `json:"numFollowers"`
	AvatarUrl NullString `json:"avatarUrl"`
};

type Room struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	IsPrivate         bool        `json:"isPrivate"`
	NumPeopleInside   int           `json:"numPeopleInside"`
	VoiceServerId     string        `json:"voiceServerId"`
	CreatorId         string        `json:"creatorId"`
	PeoplePreviewList []UserPreview `json:"peoplePreviewList"`
	AutoSpeaker       bool          `json:"autoSpeaker"`
	InsertedAt        string        `json:"inserted_at"`
	ChatMode          ChatMode      `json:"chatMode"`
	ChatThrottle 	  int 			`json:"chatThrottle"`
}

type ScheduledRoom struct {
	RoomId       UUID   `json:"roomId"`
	Description  string `json:"description"`
	ScheduledFor string `json:"scheduledFor"`
	NumAttending int    `json:"numAttending"`
	Name         string `json:"name"`
	Id           UUID   `json:"id"`
	CreatorId    UUID   `json:"creatorId"`
	Creator      User   `json:"creator"`
}

type User struct {
	YouAreFollowing       bool                  `json:"youAreFollowing"`
	Username              string                `json:"username"`
	Online                bool                  `json:"online"`
	NumFollowing          int                   `json:"numFollowing"`
	NumFollowers          int                   `json:"numFollowers"`
	LastOnline            string                `json:"lastOnline"`
	Id                    UUID                  `json:"id"`
	FollowsYou            bool                  `json:"followsYou"`
	BotOwnerId            NullString            `json:"botOwnerId"`
	DisplayName           string                `json:"displayName"`
	CurrentRoomId         UUID                  `json:"currentRoomId"`
	CurrentRoom           Room                  `json:"currentRoom"`
	Bio                   NullString            `json:"bio"`
	AvatarUrl             string                `json:"avatarUrl"`
	BannerUrl             NullString            `json:"bannerUrl"`
	WhisperPrivacySetting WhisperPrivacySetting `json:"whisperPrivacySetting"`
}

type MessageToken struct {
	T string `json:"t"`
	V interface{} `json:"v"`
}

type Message struct {
	Id UUID `json:"id"`
	UserId UUID `json:"userId"`
	AvatarUrl UUID `json:"avatarUrl"`
	Color string `json:"color"`
	DisplayName string `json:"displayName"`
	Tokens []MessageToken `json:"tokens"`
	Username string `json:"username"`
	Deleted bool `json:"deleted"`
	DeleterId UUID `json:"deleterId"`
	SentAt string `json:"sentAt"`
	IsWhisper bool `json:"isWhisper"`
}

type BaseUser struct {
	Username string `json:"username"`
	Online bool `json:"online"`
	LastOnline string `json:"lastOnline"`
	Id string `json:"id"`
	Bio string `json:"bio"`
	DisplayName string `json:"displayName"`
	AvatarUrl string `json:"avatarUrl"`
	BannerUrl string `json:"bannerUrl"`
	NumFollowing int `json:"numFollowing"`
	NumFollowers int `json:"numFollowers"`
	CurrentRoom *Room `json:"currentRoom"`
	BotOwnerId NullString `json:"botOwnerId"`
}

type PaginatedBaseUsers struct {
	Users []BaseUser `json:"users"`
	NextCursor *int `json:"nextCursor"`
}

type RoomPermissions struct {
	AskedToSpeak bool `json:"askedToSpeak"`
	IsSpeaker bool `json:"isSpeaker"`
	IsMod bool `json:"isMod"`
}

type UserWithFollowInfo struct {
	Username string `json:"username"`
	Online bool `json:"online"`
	LastOnline string `json:"lastOnline"`
	Id string `json:"id"`
	Bio string `json:"bio"`
	DisplayName string `json:"displayName"`
	AvatarUrl string `json:"avatarUrl"`
	BannerUrl string `json:"bannerUrl"`
	NumFollowing int `json:"numFollowing"`
	NumFollowers int `json:"numFollowers"`
	CurrentRoom *Room `json:"currentRoom"`
	BotOwnerId NullString `json:"botOwnerId"`
	FollowsYou *bool `json:"followsYou"`
	YouAreFollowing *bool `json:"youAreFollowing"`
	IBlockedThem *bool `json:"iBlockedThem"`
}

type RoomUser struct {
	Username string `json:"username"`
	Online bool `json:"online"`
	LastOnline string `json:"lastOnline"`
	Id string `json:"id"`
	Bio string `json:"bio"`
	DisplayName string `json:"displayName"`
	AvatarUrl string `json:"avatarUrl"`
	BannerUrl string `json:"bannerUrl"`
	NumFollowing int `json:"numFollowing"`
	NumFollowers int `json:"numFollowers"`
	CurrentRoom *Room `json:"currentRoom"`
	BotOwnerId NullString `json:"botOwnerId"`
	FollowsYou *bool `json:"followsYou"`
	YouAreFollowing *bool `json:"youAreFollowing"`
	IBlockedThem *bool `json:"iBlockedThem"`
	RoomPermissions *RoomPermissions `json:"roomPermissions"`
}

type BooleanMap map[UUID]bool

type CurrentRoom struct {
	Id                string        `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	IsPrivate         bool        `json:"isPrivate"`
	NumPeopleInside   int           `json:"numPeopleInside"`
	VoiceServerId     string        `json:"voiceServerId"`
	CreatorId         string        `json:"creatorId"`
	PeoplePreviewList []UserPreview `json:"peoplePreviewList"`
	InsertedAt        string        `json:"inserted_at"`
	ChatMode          ChatMode      `json:"chatMode"`
	ChatThrottle 	  int 			`json:"chatThrottle"`
	Users []RoomUser `json:"users"`
	MuteMap BooleanMap `json:"muteMap"`
	DeafMap BooleanMap `json:"deafMap"`
	ActiveSpeakerMap BooleanMap `json:"activeSpeakerMap"`
	AutoSpeaker bool `json:"autoSpeaker"`
}

type RoomUpdateArgs struct {
	Name *string `json:"name,omitempty"`
	Privacy *string `json:"privacy,omitempty"`
	ChatThrottle *int `json:"chatThrottle,omitempty"`
	Description *string `json:"description,omitempty"`
	AutoSpeaker *bool `json:"autoSpeaker,omitempty"`
	ChatMode *ChatMode `json:"chatMode,omitempty"`
}