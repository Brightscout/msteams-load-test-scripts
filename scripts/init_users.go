package scripts

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"go.uber.org/zap"
)

func InitUsers(config *serializers.Config, logger *zap.Logger) error {
	var userCreds []*serializers.StoredUser
	for _, userConfig := range config.UsersConfiguration {
		cred, err := azidentity.NewUsernamePasswordCredential(
			config.ConnectionConfiguration.TenantID,
			config.ConnectionConfiguration.ClientID,
			userConfig.Email,
			userConfig.Password,
			&azidentity.UsernamePasswordCredentialOptions{
				ClientOptions: azcore.ClientOptions{
					Retry: policy.RetryOptions{
						MaxRetries:    3,
						RetryDelay:    4 * time.Second,
						MaxRetryDelay: 120 * time.Second,
					},
				},
			},
		)
		if err != nil {
			logger.Error("Unable to create creds using username/password", zap.String("user", userConfig.Email), zap.Error(utils.NormalizeGraphAPIError(err)))
			continue
		}

		token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
			Scopes: []string{"https://graph.microsoft.com/.default"},
		})
		if err != nil {
			logger.Error("Unable to get token", zap.String("user", userConfig.Email), zap.Error(utils.NormalizeGraphAPIError(err)))
			continue
		}

		client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
		if err != nil {
			logger.Error("Unable to create client", zap.String("user", userConfig.Email), zap.Error(utils.NormalizeGraphAPIError(err)))
			continue
		}

		user, err := client.Me().Get(context.Background(), nil)
		if err != nil {
			logger.Error("Unable to get user info", zap.String("user", userConfig.Email), zap.Error(utils.NormalizeGraphAPIError(err)))
			continue
		}

		userCreds = append(userCreds, &serializers.StoredUser{
			ID:    *user.GetId(),
			Email: userConfig.Email,
			Token: token.Token,
		})
	}

	response, err := utils.LoadCreds()
	if err != nil {
		logger.Error("Unable to load creds", zap.Error(err))
		return err
	}

	response.Users = userCreds
	if err := utils.StoreCreds(response); err != nil {
		logger.Error("Unable to store creds", zap.Error(err))
		return err
	}

	return nil
}
