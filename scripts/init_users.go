package scripts

import (
	"context"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	"go.uber.org/zap"
)

func InitUsers(config *serializers.Config, logger *zap.Logger) error {
	connectionConfiguration := config.ConnectionConfiguration
	usersConfiguration := config.UsersConfiguration
	adminCred, err := getUserCreds(connectionConfiguration.TenantID, connectionConfiguration.ClientID, usersConfiguration.AdminEmail, usersConfiguration.AdminPassword)
	if err != nil {
		logger.Error("Unable to create admin creds using username/password", zap.String("User", usersConfiguration.AdminEmail), zap.Error(utils.NormalizeGraphAPIError(err)))
		return err
	}

	client, err := GetAppClient(&connectionConfiguration)
	if err != nil {
		logger.Error("Unable to create client", zap.String("User", usersConfiguration.AdminEmail), zap.Error(utils.NormalizeGraphAPIError(err)))
		return err
	}

	users, err := client.Users().Get(context.Background(), nil)
	if err != nil {
		logger.Error("Unable to get list of users", zap.Error(utils.NormalizeGraphAPIError(err)))
		return err
	}

	if len(users.GetValue()) == 0 {
		logger.Info("No user found on MS Teams")
		return nil
	}

	var userCreds []*serializers.StoredUser
	for _, user := range users.GetValue() {
		email := user.GetMail()
		userID := user.GetId()
		if email != nil && *email != "" && userID != nil && *userID != "" {
			var cred *azidentity.UsernamePasswordCredential
			var err error
			if *email == usersConfiguration.AdminEmail {
				cred = adminCred
			} else {
				cred, err = getUserCreds(connectionConfiguration.TenantID, connectionConfiguration.ClientID, *email, usersConfiguration.UserPassword)
				if err != nil {
					logger.Error("Unable to create creds using username/password", zap.String("User", *email), zap.Error(utils.NormalizeGraphAPIError(err)))
					continue
				}
			}

			token, err := cred.GetToken(context.Background(), policy.TokenRequestOptions{
				Scopes: constants.DefaultOAuthScopes,
			})
			if err != nil {
				logger.Error("Unable to get token", zap.String("User", *email), zap.Error(utils.NormalizeGraphAPIError(err)))
				continue
			}

			userCreds = append(userCreds, &serializers.StoredUser{
				ID:    *userID,
				Email: *email,
				Token: token.Token,
			})
		}
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

func getUserCreds(tenantID, clientID, email, password string) (*azidentity.UsernamePasswordCredential, error) {
	return azidentity.NewUsernamePasswordCredential(tenantID, clientID, email, password,
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
}
