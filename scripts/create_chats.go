package scripts

import (
	"context"
	"errors"

	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
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
		newDM, err := createChatForUsers(client, []string{store.Users[0].ID, store.Users[1].ID})
		if err != nil {
			logger.Error("unable to create the DM", zap.Error(err))
		} else {
			store.DM = &serializers.StoredChat{
				ID: *newDM.GetId(),
			}
		}
	}

	if len(store.Users) >= constants.MinUsersForGM {
		newGM, err := createChatForUsers(client, []string{
			store.Users[0].ID,
			store.Users[1].ID,
			store.Users[2].ID,
		})
		if err != nil {
			logger.Error("unable to create the GM", zap.Error(err))
		} else {
			store.GM = &serializers.StoredChat{
				ID: *newGM.GetId(),
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

func createChatForUsers(client *msgraphsdkgo.GraphServiceClient, usersIDs []string) (models.Chatable, error) {
	chatType := models.GROUP_CHATTYPE
	if len(usersIDs) == 2 {
		chatType = models.ONEONONE_CHATTYPE
	}

	members := make([]models.ConversationMemberable, len(usersIDs))
	for idx, userID := range usersIDs {
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
		return nil, utils.NormalizeGraphAPIError(err)
	}

	return newChat, nil
}
