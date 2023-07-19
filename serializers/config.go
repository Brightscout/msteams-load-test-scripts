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
	if c.ConnectionConfiguration.TenantID == "" {
		return errors.New("tenantID should not be empty")
	}

	if c.ConnectionConfiguration.ClientID == "" {
		return errors.New("clientID should not be empty")
	}

	if c.ConnectionConfiguration.ClientSecret == "" {
		return errors.New("clientSecret should not be empty")
	}

	config := c.ConnectionConfiguration
	config.TenantID = strings.TrimSpace(config.TenantID)
	config.ClientID = strings.TrimSpace(config.ClientID)
	config.ClientSecret = strings.TrimSpace(config.ClientSecret)

	return nil
}

func (c *Config) IsUsersConfigurationValid() error {
	for idx, user := range c.UsersConfiguration {
		if user.Email == "" {
			return fmt.Errorf("%s. index: %d", "user email should not be empty", idx)
		}

		if user.Password == "" {
			return fmt.Errorf("%s. index: %d", "user password should not be empty", idx)
		}

		user.Email = strings.TrimSpace(user.Email)
		user.Password = strings.TrimSpace(user.Password)
	}

	return nil
}

func (c *Config) IsChannelsConfigurationValid() error {
	for idx, channel := range c.ChannelsConfiguration {
		if channel.TeamID == "" {
			return fmt.Errorf("%s. index: %d", "team ID should not be empty", idx)
		}

		if channel.ChannelID == "" {
			if channel.ChannelDisplayName == "" {
				return fmt.Errorf("%s. index: %d", "channel display name should not be empty", idx)
			}

			if channel.Type == "" {
				return fmt.Errorf("%s. index: %d", "channel type should not be empty", idx)
			}

			if !strings.EqualFold(channel.Type, "P") && !strings.EqualFold(channel.Type, "O") {
				return fmt.Errorf("%s. index: %d", `invalid channel type. allowed values are "P" and "O"`, idx)
			}
		}

		channel.TeamID = strings.TrimSpace(channel.TeamID)
		channel.Type = strings.TrimSpace(channel.Type)
		channel.ChannelID = strings.TrimSpace(channel.ChannelID)
		channel.ChannelDisplayName = strings.TrimSpace(channel.ChannelDisplayName)
	}

	return nil
}
