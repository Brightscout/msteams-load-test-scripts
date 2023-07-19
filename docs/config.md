# Load testing configuration

## Connection configuration

### TenantID: *string*
Tenant ID of your Azure Application.

### ClientID: *string*
Client ID of your Azure Application.

### ClientSecret: *string*
Client Secret of your Azure Application.

## User Configurations

### Email: *string*
The email of the MS Teams user.

### Password: *string*
The password of the MS Teams user.

## Channel Configurations

### TeamID: *string*
The team ID to be used for load testing.

### ChannelID: *string*
The channel ID to be used for load testing.

### ChannelDisplayName: *string*
The channel display name to be used for load testing.

**Note**: Either ChannelID or ChannelDisplayName needs to be provided for performing the load testing.

### Type: *string*
The type of new channel to be created. Channel types can be `O` and `P`, which denote open channel and private channel, respectively.

## Post Configurations

### MaxWordsCount: *int*
The maximum number of words in a sentence in a post.

### MaxWordLength: *int*
The maximum length of each word in a post message.

## Load test Configuration

### VirtualUserCount: *int*
The count of virtual users running concurrently and creating posts in the MS Teams channels, and chats.

### Duration: *string*
The duration(in seconds) specifying the total duration of the test run.

### RPS: *boolean*
Set this value to `true` to use the request per second configuration.

### TimeUnit: *string*
Period of time to apply the rate value.

### Executor: *string*
Types of executors to apply for the request rate. Available executors are: `ramping-arrival-rate` and `constant-arrival-rate`.

### Rate: *int*
Number of iterations to start during each time unit period.
