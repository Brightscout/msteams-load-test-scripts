package serializers

type Store struct {
	Users    []*StoredUser
	Channels []*StoredChannel
	DM       *StoredChat
	GM       *StoredChat
}

type StoredUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type StoredChannel struct {
	ID     string `json:"id"`
	TeamID string `json:"team_id"`
}

type StoredChat struct {
	ID string `json:"id"`
}
