package discordbot

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//Const containing the prefix needed to use bot commands
const prefix string = "+"

//Help contains all the commands available
const help string = "Current commands are:\n\tping\n\tcard <card name>\n\tdice <die sides>\n\tinsult\n\tadvice"

//Const containing string to be sent if decoding fails
const decodingFailed string = "Something wrong happened when decoding data"

//MessageCreate will be called everytime a new message is sent in a channel the bot has access to
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { // Preventing bot from using own commands
		return
	}

	cmd := strings.Split(m.Content, " ") //	Splitting command into string slice

	switch cmd[0] {
	case prefix + "help":
		s.ChannelMessageSend(m.ChannelID, help)
	case prefix + "ping":
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	case prefix + "card":
		if len(cmd) == 1 { // Checks if card name is missing
			log.Println("Missing card name!")
			s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" missing card name!")
		} else {
			s.ChannelMessageSend(m.ChannelID, getCard(cmd))
		}
	case prefix + "dice":
		if len(cmd) == 2 { // Checks if die command was used properly
			s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" rolled "+diceRoll(cmd))
		}
	case prefix + "insult":
		if len(cmd) == 1 { // Checks if insult command was used properly
			s.ChannelMessageSend(m.ChannelID, getInsult())
		}
	case prefix + "advice":
		if len(cmd) == 1 { // Checks if advice command was used properly
			s.ChannelMessageSend(m.ChannelID, getAdvice())
		}
	default:
		return
	}
}