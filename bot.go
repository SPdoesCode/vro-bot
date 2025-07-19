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

type channels struct {
	Server 	string 		`toml:"server"`
	Channel string 		`toml:"channel"`
}

type CFG struct {
	Token     string   	`toml:"token"`
	Prefix    string   	`toml:"prefix"`
	ImageDir  string   	`toml:"image_dir"`
	Images    []string 	`toml:"images"`
	Channel   []channels   	`toml:"channel"`
	Deaths	  []string 	`toml:"deaths"`
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
			case "kill":
				fmt.Println("Command is kill with args as ", args , " sent by ", message.Author.ID)
				if len(config.Deaths) == 0 {
					self.ChannelMessageSend(message.ChannelID, "No deaths configured.")
					fmt.Println("WARN: No images in config.Deaths")
					break
				}
				dnum := rand.Intn(len(config.Deaths))

				self.ChannelMessageSend(message.ChannelID, "<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+config.Deaths[dnum])
				fmt.Println("Sent: \"<@"+message.Author.ID+"> killed "+strings.Join(args[1:], " ")+" with a "+config.Deaths[dnum]+"\" to ", message.ChannelID)

			default :
				fmt.Println("No ", cmd, "command found!")
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
			fmt.Println("Hourly Message Sending to channels ", config.Channel)
			if len(config.Channel) == 0 {
				fmt.Println("No channels configed... will not send!")
			} else {
				for _, ch := range config.Channel {
					server, err := self.Guild(ch.Server)
					if err != nil {
						fmt.Println("WARN: Count get guild info:", err)
					}

					fmt.Println("Sending hourly message to ", ch.Channel, "in the guild", ch.Server, "with name", server.Name)

					sendRand(self, ch.Channel)

					fmt.Println("Sent!")
				}
				fmt.Println("Sent to all configured servers!")
			}
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
	<- quit // stop when indicated

	fmt.Println("Quiting...")

	os.Exit(0)
}
