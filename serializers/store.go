package serializers

type Store struct {
	Users    []*StoredUser
	Channels []*StoredChannel
}

type StoredUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}

type StoredChannel struct {
	ID string `json:"id"`
}
