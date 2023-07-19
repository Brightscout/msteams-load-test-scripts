package scripts

import (
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Brightscout/msteams-load-test-scripts/constants"
	"github.com/Brightscout/msteams-load-test-scripts/serializers"
	"github.com/Brightscout/msteams-load-test-scripts/utils"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
)

func GetAppClient(connConfig *serializers.ConnectionConfiguration) (*msgraphsdk.GraphServiceClient, error) {
	cred, err := azidentity.NewClientSecretCredential(
		connConfig.TenantID,
		connConfig.ClientID,
		connConfig.ClientSecret,
		&azidentity.ClientSecretCredentialOptions{
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
		return nil, utils.NormalizeGraphAPIError(err)
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, constants.DefaultOAuthScopes)
	if err != nil {
		return nil, utils.NormalizeGraphAPIError(err)
	}

	return client, nil
}
