package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/containrrr/shoutrrr"
)

func main() {
	if err := _main(); err != nil {
		log.Fatalln(err)
	}
}

func _main() error {
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")
	serverID := os.Getenv("DISCORD_SERVER_ID")
	shoutrrrURL := os.Getenv("SHOUTRRR_URL")

	sender, err := shoutrrr.CreateSender(shoutrrrURL)
	if err != nil {
		return fmt.Errorf("shoutrrr.CreateSender: %w", err)
	}

	discord, err := discordgo.New(fmt.Sprintf("Bot %s", discordToken))
	if err != nil {
		return fmt.Errorf("discordgo.New: %w", err)
	}
	defer discord.Close()

	guild, err := discord.Guild(serverID)

	userIDNickMap := make(map[string]string, len(guild.Members))
	for _, member := range guild.Members {
		userIDNickMap[member.User.ID] = member.Nick
	}
	channelIDNameMap := make(map[string]string, len(guild.Channels))
	for _, channel := range guild.Channels {
		channelIDNameMap[channel.ID] = channel.Name
	}

	channelMembers := make(map[string][]string, len(guild.Channels))
	for _, state := range guild.VoiceStates {
		channelMembers[state.ChannelID] = append(channelMembers[state.ChannelID], state.UserID)
	}

	builder := new(strings.Builder)
	fmt.Fprintf(builder, "Discord (%s) のボイスチャンネルにいる人たちをお知らせします。", guild.Name)
	for channelID, memberIDs := range channelMembers {
		channelName := channelIDNameMap[channelID]
		memberNames := make([]string, 0, len(memberIDs))
		for _, id := range memberIDs {
			memberNames = append(memberNames, userIDNickMap[id])
		}
		fmt.Fprintf(builder, "%s : %s\n", channelName, strings.Join(memberNames, ", "))
	}

	if err := sender.Send(builder.String(), nil); err != nil {
		return fmt.Errorf("sender.Send: %w", err)
	}
}
