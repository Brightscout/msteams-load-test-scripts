package serializers

import (
	"errors"
	"fmt"
	"strings"
)

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

func (c *Config) IsConnectionConfigurationValid() error {
	config := c.ConnectionConfiguration
	config.TenantID = strings.TrimSpace(config.TenantID)
	config.ClientID = strings.TrimSpace(config.ClientID)
	config.ClientSecret = strings.TrimSpace(config.ClientSecret)

	if config.TenantID == "" {
		return errors.New("tenantID should not be empty")
	}

	if config.ClientID == "" {
		return errors.New("clientID should not be empty")
	}

	if config.ClientSecret == "" {
		return errors.New("clientSecret should not be empty")
	}

	return nil
}

func (c *Config) IsUsersConfigurationValid() error {
	for idx, user := range c.UsersConfiguration {
		user.Email = strings.TrimSpace(user.Email)
		user.Password = strings.TrimSpace(user.Password)

		if user.Email == "" {
			return fmt.Errorf("user email should not be empty. index: %d", idx)
		}

		if user.Password == "" {
			return fmt.Errorf("user password should not be empty. index: %d", idx)
		}
	}

	return nil
}

func (c *Config) IsChannelsConfigurationValid() error {
	for idx, channel := range c.ChannelsConfiguration {
		channel.TeamID = strings.TrimSpace(channel.TeamID)
		channel.Type = strings.TrimSpace(channel.Type)
		channel.ChannelID = strings.TrimSpace(channel.ChannelID)
		channel.ChannelDisplayName = strings.TrimSpace(channel.ChannelDisplayName)

		if channel.TeamID == "" {
			return fmt.Errorf("team ID should not be empty. index: %d", idx)
		}

		if channel.ChannelID == "" {
			if channel.ChannelDisplayName == "" {
				return fmt.Errorf("channel display name should not be empty. index: %d", idx)
			}

			if channel.Type == "" {
				return fmt.Errorf("channel type should not be empty. index: %d", idx)
			}

			if !strings.EqualFold(channel.Type, "P") && !strings.EqualFold(channel.Type, "O") {
				return fmt.Errorf("invalid channel type. allowed values are %q and %q. index: %d", "O", "P", idx)
			}
		}
	}

	return nil
}
