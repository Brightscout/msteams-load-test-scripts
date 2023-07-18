package constants

// File locations
const (
	ConfigFile    = "config/config.json"
	TempStoreFile = "temp_store.json"
)

// Scripts arguments
const (
	ClearStore     = "clear_store"
	CreateChannels = "create_channels"
	CreateChats    = "create_chats"
	InitUsers      = "init_users"
)

const (
	MinUsersForDM = 2
	MinUsersForGM = 3
)

var DefaultOAuthScopes = []string{"https://graph.microsoft.com/.default"}
