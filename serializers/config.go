package serializers

type Config struct {
	ConnectionConfiguration ConnectionConfiguration
	UsersConfiguration      []UsersConfiguration
	ChannelsConfiguration   []ChannelsConfiguration
}

type ConnectionConfiguration struct {
	TenantID     string
	ClientID     string
	ClientSecret string
}

type UsersConfiguration struct {
	Email    string
	Password string
}

type ChannelsConfiguration struct {
	TeamID             string
	ChannelID          string
	ChannelDisplayName string
	Type               string
}
