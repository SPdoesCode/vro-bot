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
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

// config structs

type channels struct {
	Server  string `toml:"server"`
	Channel string `toml:"channel"`
}

type CFG struct {
	Token    string     `toml:"token"`
	Prefix   string     `toml:"prefix"`
	ImageDir string     `toml:"image_dir"`
	Images   []string   `toml:"images"`
	Channel  []channels `toml:"channel"`
	Deaths   []string   `toml:"deaths"`
}

var config CFG

func getConfig() {

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		fmt.Println("ERROR: Couldnt load config.toml: ", err)
		os.Exit(1)
	}
	fmt.Println("Parsed toml")
	fmt.Println("Bot prefix is: ", config.Prefix)
}

func sendRand(self *discordgo.Session, Channel string) {

	if len(config.Images) == 0 {
		self.ChannelMessageSend(Channel, "No images configured.")
		fmt.Println("WARN: No images in config.Images")
		return
	}
	num := rand.Intn(len(config.Images))
	filename := config.ImageDir + config.Images[num] // get the random file
	fmt.Println("Selected image: ", filename)
	file, err := os.Open(filename)
	if err != nil {
		self.ChannelMessageSend(Channel, "Couldnt find image, sorry")
		fmt.Println("WARN: Failed to find image ", filename, ": ", err)
		return
	} // warn and alert to the error
	defer file.Close() // we need to close
	_, err = self.ChannelMessageSendComplex(Channel, &discordgo.MessageSend{
		Files: []*discordgo.File{
			{
				Name:   filename,
				Reader: file,
			},
		},
	})
	fmt.Println("Send image: ", file, filename)
	if err != nil {
		self.ChannelMessageSend(Channel, "Couldnt open image...")
		fmt.Println("WARN: Failed to find image ", filename, ": ", err)
	}

}

// for commands only
func ctrlMessages(self *discordgo.Session, message *discordgo.MessageCreate) {
	//fmt.Println("Message received: ", message.Content)

	if message.Author.ID == self.State.User.ID {
		return // return nothing cuz its our own message
	}

	if strings.HasPrefix(message.Content, config.Prefix) { // check for the prefix
		getConfig()
		fmt.Println("Reloaded config...")
		args := strings.Fields(message.Content)
		cmd := strings.TrimPrefix(args[0], config.Prefix) // remove it
		fmt.Println("Got command (", message.Content, ") and stripped prefix ", config.Prefix)
		switch cmd { // switch it to see what command it is
		case "getvro":
			fmt.Println("Command is getvro with args as ", args, " sent by ", message.Author)
			fmt.Println("Sending sendRand() to channel ", message.ChannelID)
			sendRand(self, message.ChannelID)

		case "sex": // this is a joke command
			fmt.Println("Command is sex with args as ", args, " sent by ", message.Author)
			if rand.Intn(2) == 0 {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and made them pregnant!")
				fmt.Println("Sent: \"<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and made them pregnant!\" to ", message.ChannelID)
			} else {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and failed to make them pregnant!")
				fmt.Println("Sent: \"<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and failed to make them pregnant!\" to ", message.ChannelID)
			}

		case "kill":
			fmt.Println("Command is kill with args as ", args, " sent by ", message.Author)
			if len(config.Deaths) == 0 {
				self.ChannelMessageSend(message.ChannelID, "No deaths configured.")
				fmt.Println("WARN: No images in config.Deaths")
				break
			}
			dnum := rand.Intn(len(config.Deaths))

			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+config.Deaths[dnum])
			fmt.Println("Sent: \"<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+config.Deaths[dnum]+"\" to ", message.ChannelID)

		case "goonto": // yet another obv joke command
			fmt.Println("Command is goonto with args as ", args, " sent by ", message.Author)
			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> gooned to "+strings.Join(args[1:], " ")+"!!!!11!")

		case "eat":
			fmt.Println("Command is hug with args as ", args, " sent by ", message.Author)
			self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> ate "+strings.Join(args[1:], " ")+" kindly with love <3")
			fmt.Println("Sent \"<@"+message.Author.ID+"> ate "+strings.Join(args[1:], " ")+" kindly with love <3\" to ", message.ChannelID)

		case "hug":
			fmt.Println("Command is eat with args as ", args, " sent by ", message.Author)
			if rand.Intn(56+56*2) == 8 {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... then started to eat them...")
				fmt.Println("Sent \"<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... then started to eat them...\" to ", message.ChannelID)
			} else {
				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... nothing else...")
				fmt.Println("Sent \"<@"+message.Author.ID+"> huged "+strings.Join(args[1:], " ")+" kindly... nothing else...\" to ", message.ChannelID)
			}

		case "help":
			fmt.Println("Command is help with args as ", args, " sent by ", message.Author)
			self.ChannelMessageSend(message.ChannelID, "Commands: help, kill, sex, getvro, eat, hug, goonto")
			fmt.Println("Sent help message to ", message.ChannelID)

		default:
			self.ChannelMessageSend(message.ChannelID, "Command "+cmd+" not found! Try the help command for more commands!")
			fmt.Println("No ", cmd, "command found!")
		}
	}
}

// every hour message hanndler
func hourlyMessage(self *discordgo.Session) {
	go func() {
		ticker := time.NewTicker(1 * time.Hour) // set up the hour control
		defer ticker.Stop()

		for {
			<-ticker.C
			fmt.Println("Hourly Message Sending to channels ", config.Channel)

			for _, ch := range config.Channel {
				server, err := self.Guild(ch.Server)
				if err != nil {
					fmt.Println("WARN: Couldnt get guild info:", err)
				}
				fmt.Println("Sending hourly message to ", ch.Channel, "in the guild", ch.Server, "with name", server.Name)

				sendRand(self, ch.Channel)
				fmt.Println("Sent!")

			}
			fmt.Println("Sent to all configured servers!")
		}
	}()
}

// slash command handler
func slashCmd(self *discordgo.Session, cmd *discordgo.InteractionCreate) {
	var args string
	options := cmd.ApplicationCommandData().Options
	for _, opt := range options {
		args = opt.StringValue()
	}
	switch cmd.ApplicationCommandData().Name {
	case "getvro":
		fmt.Println("Command is getvro with args as ", args, " sent by ", cmd.Member.User.Username)
		fmt.Println("Sending sendRand() to channel ", cmd.ChannelID)
		sendRand(self, cmd.ChannelID)
	case "kill":
		fmt.Println("Command is kill with args as ", args, " sent by ", cmd.Member.User.Username)
		if len(config.Deaths) == 0 {
			self.ChannelMessageSend(cmd.ChannelID, "No deaths configured.")
			fmt.Println("WARN: No images in config.Deaths")
			break
		}
		dnum := rand.Intn(len(config.Deaths))

		self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> killed "+args+" with a "+config.Deaths[dnum])
		fmt.Println("Sent: \"<@"+cmd.Member.User.ID+"> killed "+args+" with a "+config.Deaths[dnum]+"\" to ", cmd.ChannelID)
	case "sex": // this is a joke command
		fmt.Println("Command is sex with args as ", args, " sent by ", cmd.Member.User.Username)
		if rand.Intn(2) == 0 {
			self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> had sex with "+args+" and made them pregnant!")
			fmt.Println("Sent: \"<@"+cmd.Member.User.ID+"> had sex with "+args+" and made them pregnant!\" to ", cmd.ChannelID)
		} else {
			self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> had sex with "+args+" and failed to make them pregnant!")
			fmt.Println("Sent: \"<@"+cmd.Member.User.ID+"> had sex with "+args+" and failed to make them pregnant!\" to ", cmd.ChannelID)
		}
	case "hug":
		fmt.Println("Command is eat with args as ", args, " sent by ", cmd.Member.User.Username)
		if rand.Intn(56+56*2) == 8 {
			self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> huged "+args+" kindly... then started to eat them...")
			fmt.Println("Sent \"<@"+cmd.Member.User.ID+"> huged "+args+" kindly... then started to eat them...\" to ", cmd.ChannelID)
		} else {
			self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> huged "+args+" kindly... nothing else...")
			fmt.Println("Sent \"<@"+cmd.Member.User.ID+"> huged "+args+" kindly... nothing else...\" to ", cmd.ChannelID)
		}

	case "goonto": // yet another obv joke command
		fmt.Println("Command is goonto with args as ", args, " sent by ", cmd.Member.User.Username)
		self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> gooned to "+args+"!!!!11!")

	case "eat":
		fmt.Println("Command is hug with args as ", args, " sent by ", cmd.Member.User.Username)
		self.ChannelMessageSend(cmd.ChannelID, "<@"+cmd.Member.User.ID+"> ate "+args+" kindly with love <3")
		fmt.Println("Sent \"<@"+cmd.Member.User.ID+"> huged "+args+" kindly with love <3\" to ", cmd.ChannelID)

	case "help":
		fmt.Println("Command is help with args as ", args, " sent by ", cmd.Member.User.ID)
		self.ChannelMessageSend(cmd.ChannelID, "Commands: help, kill, sex, getvro, eat, hug, goonto")
		fmt.Println("Sent help message to ", cmd.ChannelID)
	}
}

func main() {

	getConfig()

	rand.Seed(time.Now().UnixNano())

	bot, err := discordgo.New("Bot " + config.Token) // grab the token and make a new session
	if err != nil {
		fmt.Println("ERROR: Couldnt create discord session: ", err)
		os.Exit(1)
	}

	cmds := []*discordgo.ApplicationCommand{ // slash commands init
		{
			Name:        "getvro",
			Description: "send a randome vro image",
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
	}

	for _, cmd := range cmds {
		_, err := bot.ApplicationCommandCreate(bot.State.User.ID, "", cmd)
		if err != nil {
			fmt.Println("ERROR: Couldnt register command ", cmd, ": ", err)
		}
	}

	bot.AddHandler(slashCmd)     // slash commands
	bot.AddHandler(ctrlMessages) // grab messages into the message control

	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent // request perms

	err2 := bot.Open() // try to open the connection
	if err2 != nil {
		fmt.Println("ERROR: Coulnt open the connection: ", err2)
		os.Exit(1)
	}

	defer bot.Close() // close the connection when done

	fmt.Println("Bot should be running properly now, If you need to stop do ctrl + c or ctrl + d!")

	quit := make(chan os.Signal, 1)

	hourlyMessage(bot) // start the hourly message

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit // stop when indicated

	fmt.Println("Quiting...")

	os.Exit(0)
}
