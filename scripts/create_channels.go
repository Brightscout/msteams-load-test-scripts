package scripts

import (
	"context"
	"fmt"
	"strings"

	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/teams"
	"go.uber.org/zap"
)

func CreateChannels(config *serializers.Config, logger *zap.Logger) error {
	client, err := GetAppClient(&config.ConnectionConfiguration)
	if err != nil {
		logger.Error("unable to get client", zap.Error(err))
		return err
	}

	store, err := utils.LoadCreds()
	if err != nil {
		return err
	}

	// Create a set of unique team IDs
	teamIDs := make(map[string]bool)
	for _, channelConfig := range config.ChannelsConfiguration {
		teamIDs[channelConfig.TeamID] = true
	}

	teamMembersRequestBody := teams.NewItemMembersAddPostRequestBody()
	// Create an array of conversation members to add in teams and channels
	var members []models.ConversationMemberable
	for idx, user := range store.Users {
		conversationMember := models.NewAadUserConversationMember()
		roles := []string{}
		if idx == 0 {
			roles = append(roles, "owner")
		}
		odataType := "#microsoft.graph.aadUserConversationMember"
		conversationMember.SetOdataType(&odataType)
		conversationMember.SetRoles(roles)
		additionalData := map[string]interface{}{
			"user@odata.bind": fmt.Sprintf("https://graph.microsoft.com/v1.0/users('%s')", user.ID),
		}
		conversationMember.SetAdditionalData(additionalData)
		members = append(members, conversationMember)
	}

	teamMembersRequestBody.SetValues(members)
	// Add team members in the teams
	for teamID := range teamIDs {
		// TODO: Add batching of requests
		if _, err := client.Teams().ByTeamId(teamID).Members().Add().Post(context.Background(), teamMembersRequestBody, nil); err != nil {
			logger.Error("error in adding team members", zap.String("TeamID", teamID), zap.Error(utils.NormalizeGraphAPIError(err)))
		}
	}

	var storedChannels []*serializers.StoredChannel
	for _, channelConfig := range config.ChannelsConfiguration {
		channel, newCreated, err := getOrCreateChannel(client, &channelConfig, members)
		if err != nil {
			logger.Error("error in getting or creating channel", zap.String("Channel", channelConfig.ChannelDisplayName), zap.Error(err))
			continue
		}

		storedChannels = append(storedChannels, &serializers.StoredChannel{
			ID:     *channel.GetId(),
			TeamID: channelConfig.TeamID,
		})

		if !newCreated {
			// TODO: Check if channel members are already present in the channel
			// requestFilter := fmt.Sprintf("email eq '%s'", user.Email)
			// requestParameters := &teams.ItemChannelsItemMembersRequestBuilderGetQueryParameters{
			// 	Filter: &requestFilter,
			// }

			// configuration := &teams.ItemChannelsItemMembersRequestBuilderGetRequestConfiguration{
			// 	QueryParameters: requestParameters,
			// }
			// client.Teams().ByTeamId(channelConfig.TeamID).Channels().ByChannelId(*channel.GetId()).Members().Get(context.Background(), configuration)

			channelMembersRequestBody := teams.NewItemChannelsItemMembersAddPostRequestBody()
			channelMembersRequestBody.SetValues(members)
			if _, err := client.Teams().ByTeamId(channelConfig.TeamID).Channels().ByChannelId(*channel.GetId()).Members().Add().Post(
				context.Background(), channelMembersRequestBody, nil,
			); err != nil {
				logger.Error("error in adding channel members", zap.String("ChannelName", *channel.GetDisplayName()), zap.Error(utils.NormalizeGraphAPIError(err)))
			}
		}
	}

	store.Channels = storedChannels
	if err = utils.StoreCreds(store); err != nil {
		logger.Error("Unable to store creds", zap.Error(err))
		return err
	}

	return nil
}

func getOrCreateChannel(client *msgraphsdkgo.GraphServiceClient, channelConfig *serializers.ChannelsConfiguration, members []models.ConversationMemberable) (channel models.Channelable, newCreated bool, err error) {
	teamID, channelID, channelDisplayName := channelConfig.TeamID, channelConfig.ChannelID, channelConfig.ChannelDisplayName
	if channelID != "" {
		channel, err = client.Teams().ByTeamId(teamID).Channels().ByChannelId(channelID).Get(context.Background(), nil)
		if err != nil {
			return nil, false, utils.NormalizeGraphAPIError(err)
		}

		return channel, false, nil
	}

	requestFilter := fmt.Sprintf("displayName eq '%s'", channelDisplayName)
	requestParameters := &teams.ItemChannelsRequestBuilderGetQueryParameters{
		Filter: &requestFilter,
	}

	configuration := &teams.ItemChannelsRequestBuilderGetRequestConfiguration{
		QueryParameters: requestParameters,
	}

	result, err := client.Teams().ByTeamId(teamID).Channels().Get(context.Background(), configuration)
	if err != nil {
		return nil, false, utils.NormalizeGraphAPIError(err)
	}

	channels := result.GetValue()
	if len(channels) > 0 {
		return channels[0], false, nil
	}

	channelType := models.STANDARD_CHANNELMEMBERSHIPTYPE
	if strings.EqualFold(channelConfig.Type, "P") {
		channelType = models.PRIVATE_CHANNELMEMBERSHIPTYPE
	}

	newChannel := models.NewChannel()
	newChannel.SetDisplayName(&channelDisplayName)
	newChannel.SetMembershipType(&channelType)

	newChannel.SetMembers(members)
	channel, err = client.Teams().ByTeamId(teamID).Channels().Post(context.Background(), newChannel, nil)
	if err != nil {
		return nil, false, utils.NormalizeGraphAPIError(err)
	}

	return channel, true, nil
}
