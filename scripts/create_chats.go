package scripts

import (
	"context"
	"errors"

	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"go.uber.org/zap"
)

func CreateChats(config *serializers.Config, logger *zap.Logger) error {
	store, err := utils.LoadCreds()
	if err != nil {
		return err
	}

	if len(store.Users) == 0 {
		return errors.New("no user initialized")
	}

	client, err := GetAppClient(&config.ConnectionConfiguration)
	if err != nil {
		logger.Error("unable to get client", zap.Error(err))
		return err
	}

	if len(store.Users) >= constants.MinUsersForDM {
		newDMID, err := getOrCreateChatForUsers(client, []string{store.Users[0].ID, store.Users[1].ID})
		if err != nil {
			logger.Error("unable to create the DM", zap.Error(err))
		} else {
			store.DM = &serializers.StoredChat{
				ID: newDMID,
			}
		}
	}

	if len(store.Users) >= constants.MinUsersForGM {
		newGMID, err := getOrCreateChatForUsers(client, []string{
			store.Users[0].ID,
			store.Users[1].ID,
			store.Users[2].ID,
		})
		if err != nil {
			logger.Error("unable to create the GM", zap.Error(err))
		} else {
			store.GM = &serializers.StoredChat{
				ID: newGMID,
			}
		}
	}

	if store.DM != nil || store.GM != nil {
		if err := utils.StoreCreds(store); err != nil {
			return err
		}
	}

	return nil
}

func getOrCreateChatForUsers(client *msgraphsdkgo.GraphServiceClient, userIDs []string) (string, error) {
	if len(userIDs) == 2 {
		return createChatForUsers(client, userIDs, models.ONEONONE_CHATTYPE)
	}

	chatID, err := getChatForUsers(client, userIDs)
	if err == nil && chatID != "" {
		return chatID, nil
	}

	return createChatForUsers(client, userIDs, models.GROUP_CHATTYPE)
}

func getChatForUsers(client *msgraphsdkgo.GraphServiceClient, userIDs []string) (string, error) {
	requestParameters := &users.ItemChatsRequestBuilderGetQueryParameters{
		Select: []string{"members", "id"},
		Expand: []string{"members"},
	}
	configuration := &users.ItemChatsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	res, err := client.Users().ByUserId(userIDs[0]).Chats().Get(context.Background(), configuration)
	if err != nil {
		return "", utils.NormalizeGraphAPIError(err)
	}

	for _, c := range res.GetValue() {
		if len(c.GetMembers()) == len(userIDs) {
			matches := map[string]bool{}
			for _, m := range c.GetMembers() {
				for _, u := range userIDs {
					userID, err2 := m.GetBackingStore().Get("userId")
					if err2 == nil && userID != nil && *(userID.(*string)) == u {
						matches[u] = true
						break
					}
				}
			}
			if len(matches) == len(userIDs) {
				return *c.GetId(), nil
			}
		}
	}

	return "", nil
}

func createChatForUsers(client *msgraphsdkgo.GraphServiceClient, userIDs []string, chatType models.ChatType) (string, error) {
	members := make([]models.ConversationMemberable, len(userIDs))
	for idx, userID := range userIDs {
		conversationMember := models.NewConversationMember()
		odataType := "#microsoft.graph.aadUserConversationMember"
		conversationMember.SetOdataType(&odataType)
		conversationMember.SetAdditionalData(map[string]interface{}{
			"user@odata.bind": "https://graph.microsoft.com/v1.0/users('" + userID + "')",
		})
		conversationMember.SetRoles([]string{"owner"})
		members[idx] = conversationMember
	}

	chatRequestBody := models.NewChat()
	chatRequestBody.SetChatType(&chatType)
	chatRequestBody.SetMembers(members)
	newChat, err := client.Chats().Post(context.Background(), chatRequestBody, nil)
	if err != nil {
		return "", utils.NormalizeGraphAPIError(err)
	}

	return *newChat.GetId(), nil
}
