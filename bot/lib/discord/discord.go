package discord

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

var (
	Token       = "Bot"
	Discord     = discordgo.Session{}
	handlerList []interface{}
	stopBot     = make(chan bool)
)

type BaseChannel struct {
	Name string
	ID   string
}
type Category struct {
	BaseChannel
}
type Channel struct {
	BaseChannel
	Category
}
type Channels []Channel

var (
	categories []Category
	channels   Channels
)

// Len, Less, Swapを定義
func (c Channels) Len() int {
	return len(c)
}
func (c Channels) Less(i, j int) bool {
	return c[i].ID > c[j].ID
}
func (c Channels) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

func IsOwnMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Author.ID == s.State.User.ID {
		return true
	}
	return false
}

func SendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
	_, err := s.ChannelMessageSend(c.ID, msg)

	log.Println(">>> " + msg)
	if err != nil {
		log.Println("Error sending message: ", err)
	}
}

func GetChannel(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.Channel {
	c, err := s.State.Channel(m.ChannelID) //チャンネル取得
	if err != nil {
		log.Println("Error getting channel: ", err)
		return nil
	}
	return c
}

func getCategoriesChannel(s *discordgo.Session) {
	categories = []Category{}
	channels = Channels{}
	for _, guild := range s.State.Guilds {
		for _, channel := range guild.Channels {
			//fmt.Printf("parrentId: %s, channelId: %s, channelName: %s\r\n", channel.ParentID, channel.ID, channel.Name)
			if channel.Type == discordgo.ChannelTypeGuildCategory {
				categories = append(categories, Category{BaseChannel{channel.Name, channel.ID}})
			}
			if channel.Type == discordgo.ChannelTypeGuildText || channel.Type == discordgo.ChannelTypeGuildVoice {
				tmpC := Channel{
					BaseChannel: BaseChannel{channel.Name, channel.ID},
					Category:    Category{BaseChannel{ID: channel.ParentID}},
				}
				channels = append(channels, tmpC)
			}
		}
	}
	for _, cate := range categories {
		for i, channel := range channels {
			if cate.BaseChannel.ID == channel.Category.ID {
				channels[i].Category.BaseChannel.Name = cate.BaseChannel.Name
			}
		}
	}
}

func GetFixChannel(s *discordgo.Session, category string, channel string) (*discordgo.Channel, error) {
	getCategoriesChannel(s)
	var ctmp Channel
	for _, c := range channels {
		if c.BaseChannel.Name == channel && c.Category.Name == category {
			ctmp = c
		}
	}
	if ctmp.BaseChannel.ID == "" {
		return nil, errors.New("見つかりませんでした。")
	}
	sendChannel, _ := s.State.Channel(ctmp.BaseChannel.ID)
	return sendChannel, nil
}

func SetUpDiscordBot(token string) error {

	Token = Token + " " + token
	Discord, err := discordgo.New()
	if err != nil {
		fmt.Println("Error logging in")
		fmt.Println(err)
		return err
	}
	Discord.Token = Token

	err = Discord.Open()
	if err != nil {
		fmt.Println("openError")
		fmt.Println(err)
		return err
	}
	for _, handler := range handlerList {
		Discord.AddHandler(handler)
	}

	fmt.Println("Listening...")

	<-stopBot
	return nil
}

func AddHandler(onMessageCreate func(s *discordgo.Session, m *discordgo.MessageCreate)) {
	handlerList = append(handlerList, onMessageCreate)
}

func AddHandlerJoinAndLeave(
	onJoinCreate func(s *discordgo.Session, m *discordgo.GuildMemberAdd),
	onLeaveCreate func(s *discordgo.Session, m *discordgo.GuildMemberRemove)) {
	handlerList = append(handlerList, onJoinCreate)
	handlerList = append(handlerList, onLeaveCreate)
}

func GetMemberId(u []*discordgo.User) *discordgo.User {
	return u[len(u)-1]
}

// ComesFromDM returns true if a message comes from a DM channel
func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}