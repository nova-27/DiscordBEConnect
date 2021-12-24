package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

func BotLogin(token string) *discordgo.Session {
	bot, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
	}

	bot.AddHandler(messageCreate)
	bot.Identify.Intents = discordgo.IntentsGuildMessages
	err = bot.Open()
	if err != nil {
		log.Print("Botのログインに失敗しました")
		log.Fatalln(err)
	}

	log.Println("Botが起動しました")
	return bot
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.ChannelID != CHANNEL {
		return
	}

	senderName := m.Author.Username
	if m.Member.Nick != "" {
		senderName = m.Member.Nick
	}

	log.Printf("[Discord]  <%s> %s", senderName, m.Message.Content)
	if mcwsclient != nil {
		mcwsclient.Exec(
			fmt.Sprintf(
				"tellraw @a {\"rawtext\":[{\"text\":\"§9[Discord]§r<%s> %s\"}]}",
				senderName,
				m.Message.Content),
			nil)
	}
}
