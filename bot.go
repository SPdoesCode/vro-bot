package main

/*
 *
 * 	Hourly vro bot by SP649/SPdoesCode
 *
 */

import (
	"fmt"
	"time"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)

type CFG struct {
	Token     string   `toml:"token"`
	Prefix    string   `toml:"prefix"`
	ImageDir  string   `toml:"image_dir"`
	Images    []string `toml:"images"`
	Channel   string   `toml:"channel"`
	Deaths	  []string `toml:"deaths"`
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
	filename := config.ImageDir+config.Images[num] // get the random file
	fmt.Println("Selected image: ", filename)
	file, err := os.Open(filename)
	if err != nil { self.ChannelMessageSend(Channel, "Couldnt find image, sorry"); fmt.Println("WARN: Failed to find image ", filename,": ", err); return } // warn and alert to the error
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
		fmt.Println("WARN: Failed to find image ", filename,": ", err)
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
				fmt.Println("Command is getvro, sending sendRand() to channel ", message.ChannelID)
				sendRand(self, message.ChannelID)
			case "sex": // this is a joke command
				fmt.Println("Command is sex with args as ", args , " sent by ", message.Author.ID)
				if rand.Intn(2) == 0 {
					self.ChannelMessageSend(Channel, "<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and made them pregnant!")
					fmt.Println("Sent: \"<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and made them pregnant!\" to ", message.ChannelID)
				} else {
					self.ChannelMessageSend(Channel, "<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and failed to make them pregnant!")
					fmt.Println("Sent: \"<@"+message.Author.ID+"> had sex with "+strings.Join(args[1:], " ")+" and failed to make them pregnant!\" to ", message.ChannelID)
				}
			case "kill":
				fmt.Println("Command is kill with args as ", args , " sent by ", message.Author.ID)
				if len(config.Deaths) == 0 {
					self.ChannelMessageSend(Channel, "No deaths configured.")
					fmt.Println("WARN: No images in config.Deaths")
					break
				}
				dnum := rand.Intn(len(config.Deaths))

				self.ChannelMessageSend(Channel, "<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+Deaths[dnum])
				fmt.Println("Sent: \"<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+Deaths[dnum]+"\" to ", message.ChannelID)

			default :
				fmt.Println("No ", endcmd, "command found!")
		}
	}
}

// every hour message hanndler
func hourlyMessage(self *discordgo.Session) {
	go func(){
		ticker := time.NewTicker(1 * time.Hour) // set up the hour control
		defer ticker.Stop()

		for {
			<- ticker.C
			sendRand(self, config.Channel)

		}
	}()
}

func main() {

	getConfig()

	rand.Seed(time.Now().UnixNano())

	bot, err := discordgo.New("Bot "+config.Token) // grab the token and make a new session
	if err != nil {
		fmt.Println("ERROR: Couldnt create discord session: ", err)
		os.Exit(1)
	}

	bot.AddHandler(ctrlMessages) // grab messages into the message control

	err2 := bot.Open() // try to open the connection
	if err2 != nil {
		fmt.Println("ERROR: Coulnt open the connection: ", err2)
		os.Exit(1)
	}

	bot.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	defer bot.Close() // close the connection when done

	fmt.Println("Bot should be running properly now, If you need to stop do ctrl + c or ctrl + d!")

	quit := make(chan os.Signal, 1)

	hourlyMessage(bot) // start the hourly message

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<- quit // stop when indicated

	fmt.Println("Quiting...")

	os.Exit(0)
}
