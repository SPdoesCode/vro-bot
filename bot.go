package main

/*
 *
 * 	Hourly vro bot by SP649/SPdoesCode
 *
 */

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"vro-bot/util"

	"github.com/bwmarrin/discordgo"
)

var cmds = []*discordgo.ApplicationCommand{ // slash commands init
	{
		Name:        "getvro",
		Description: "send a randome vro image",
	},
	{
		Name:        "work",
		Description: "makes u work for money",
	},
	{
		Name:        "sex",
		Description: "sex someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "thing",
				Description: "the thing u wana sex",
				Required:    true,
			},
		},
	},
	{
		Name:        "eat",
		Description: "eat someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "thing",
				Description: "the thing u wana eat",
				Required:    true,
			},
		},
	},
	{
		Name:        "kill",
		Description: "kill someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "thing",
				Description: "the thing u wana kill",
				Required:    true,
			},
		},
	},
	{
		Name:        "goonto",
		Description: "goon to a thing",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "thing",
				Description: "the thing u wana goon to",
				Required:    true,
			},
		},
	},
	{
		Name:        "hug",
		Description: "hug someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "thing",
				Description: "the thing u wana hug",
				Required:    true,
			},
		},
	},
	{
		Name:        "rob",
		Description: "rob someone",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "person",
				Description: "the person u wana rob can be a user id or @user",
				Required:    true,
			},
		},
	},
}

func main() {

	util.GetConfig()
	util.GetData()

	rand.Seed(time.Now().UnixNano())

	bot, err := discordgo.New("Bot " + util.Config.Token) // grab the token and make a new session
	if err != nil {
		fmt.Println("ERROR: Couldnt create discord session: ", err)
		os.Exit(1)
	}

	bot.AddHandler(util.SlashCmd)     // slash commands
	bot.AddHandler(util.CtrlMessages) // grab messages into the message control

	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent // request perms

	err2 := bot.Open() // try to open the connection
	if err2 != nil {
		fmt.Println("ERROR: Coulnt open the connection: ", err2)
		os.Exit(1)
	}
	defer bot.Close() // close the connection when done

	for _, cmd := range cmds {
		_, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
		fmt.Println("Registering command: ", cmd)
		if err != nil {
			fmt.Println("ERROR: Couldnt register command ", cmd, ": ", err)
		}
	}

	fmt.Println("Bot should be running properly now, If you need to stop do ctrl + c or ctrl + d!")

	quit := make(chan os.Signal, 1)

	util.HourlyMessage(bot) // start the hourly message

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit // stop when indicated

	fmt.Println("Quiting...")

	os.Exit(0)
}
