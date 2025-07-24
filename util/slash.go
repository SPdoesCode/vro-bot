package util

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func RespondTo(self *discordgo.Session, cmd *discordgo.InteractionCreate, content string) {
	err := self.InteractionRespond(cmd.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
	if err != nil {
		fmt.Println("ERROR: Failed to respond to slash command: ", err)
	}
}

func SendRepRand(self *discordgo.Session, cmd *discordgo.InteractionCreate) {

	RespondTo(self, cmd, "Sending image of vro...")

	if len(Images) == 0 {
		RespondTo(self, cmd, "No images configured.")
		fmt.Println("WARN: No images in Config.Images")
		return
	}
	num := rand.Intn(len(Images))
	filename := Config.ImageDir + Images[num] // get the random file
	fmt.Println("Selected image: ", filename)
	file, err := os.Open(filename)
	if err != nil {
		RespondTo(self, cmd, "Couldnt find image, sorry")
		fmt.Println("WARN: Failed to find image ", filename, ": ", err)
		return
	} // warn and alert to the error
	defer file.Close() // we need to close
	_, err = self.FollowupMessageCreate(cmd.Interaction, true, &discordgo.WebhookParams{
		Files: []*discordgo.File{
			{
				Name:   filepath.Base(filename),
				Reader: file,
			},
		},
	})
	fmt.Println("Send image: ", file, filename)
	if err != nil {
		RespondTo(self, cmd, "Couldnt open image...")
		fmt.Println("WARN: Failed to find image ", filename, ": ", err)
	}

}

// slash command handler
func SlashCmd(self *discordgo.Session, cmd *discordgo.InteractionCreate) {
	var args string
	var username, userID string
	if cmd.Member != nil {
		username = cmd.Member.User.Username
		userID = cmd.Member.User.ID
	} else {
		username = cmd.User.Username
		userID = cmd.User.ID
	}
	options := cmd.ApplicationCommandData().Options
	for _, opt := range options {
		args = opt.StringValue()
	}
	switch cmd.ApplicationCommandData().Name {
	case "getvro":
		fmt.Println("Command is getvro with args as ", args, " sent by ", username)
		fmt.Println("Sending sendRand() to channel ", cmd.ChannelID)
		//RespondTo(self, cmd, "Sending image of vro...")
		SendRepRand(self, cmd)
	case "kill":
		fmt.Println("Command is kill with args as ", args, " sent by ", username)
		if len(Config.Deaths) == 0 {
			RespondTo(self, cmd, "No deaths configured.")
			fmt.Println("WARN: No images in Config.Deaths")
			break
		}
		dnum := rand.Intn(len(Config.Deaths))

		RespondTo(self, cmd, "<@"+userID+"> killed "+args+" with a "+Config.Deaths[dnum])
		fmt.Println("Sent: \"<@"+userID+"> killed "+args+" with a "+Config.Deaths[dnum]+"\" to ", cmd.ChannelID)
	case "sex": // this is a joke command
		fmt.Println("Command is sex with args as ", args, " sent by ", username)
		if rand.Intn(2) == 0 {
			RespondTo(self, cmd, "<@"+userID+"> had sex with "+args+" and made them pregnant!")
			fmt.Println("Sent: \"<@"+userID+"> had sex with "+args+" and made them pregnant!\" to ", cmd.ChannelID)
		} else {
			RespondTo(self, cmd, "<@"+userID+"> had sex with "+args+" and failed to make them pregnant!")
			fmt.Println("Sent: \"<@"+userID+"> had sex with "+args+" and failed to make them pregnant!\" to ", cmd.ChannelID)
		}
	case "hug":
		fmt.Println("Command is eat with args as ", args, " sent by ", username)
		if rand.Intn(56+56*2) == 8 {
			RespondTo(self, cmd, "<@"+userID+"> huged "+args+" kindly... then started to eat them...")
			fmt.Println("Sent \"<@"+userID+"> huged "+args+" kindly... then started to eat them...\" to ", cmd.ChannelID)
		} else {
			RespondTo(self, cmd, "<@"+userID+"> huged "+args+" kindly... nothing else...")
			fmt.Println("Sent \"<@"+userID+"> huged "+args+" kindly... nothing else...\" to ", cmd.ChannelID)
		}
	case "work":
		fmt.Println("Command is work with args as ", args, " sent by ", username)
		num := rand.Intn(201)
		RespondTo(self, cmd, "<@"+userID+"> wroked and earned "+strconv.Itoa(num)+" money!")
		addMoney(userID, num, self)
		fmt.Println("Sent \"<@"+userID+"> wroked and earned "+strconv.Itoa(num)+" money!\" to ", cmd.ChannelID)
	case "rob":
		fmt.Println("Command is rob with args as ", args, " sent by ", username)
		victim := args
		num := rand.Intn(101)
		victim = strings.TrimPrefix(victim, "<@")
		victim = strings.TrimPrefix(victim, "!")
		victim = strings.TrimSuffix(victim, ">")
		e := rmMoney(victim, num, self)
		if !e {
			RespondTo(self, cmd, "User not found!")
			break
		}
		RespondTo(self, cmd, "<@"+userID+"> robbed <@"+victim+"> and stole "+strconv.Itoa(num)+" money!")
		addMoney(userID, num, self)
		fmt.Println("Sent \"<@"+userID+"> robbed <@"+victim+"> and stole "+strconv.Itoa(num)+" money!\" to ", cmd.ChannelID)
	case "goonto": // yet another obv joke command
		fmt.Println("Command is goonto with args as ", args, " sent by ", username)
		RespondTo(self, cmd, "<@"+userID+"> gooned to "+args+"!!!!11!")

	case "eat":
		fmt.Println("Command is hug with args as ", args, " sent by ", username)
		RespondTo(self, cmd, "<@"+userID+"> ate "+args+" kindly with love <3")
		fmt.Println("Sent \"<@"+userID+"> huged "+args+" kindly with love <3\" to ", cmd.ChannelID)

	case "help":
		fmt.Println("Command is help with args as ", args, " sent by ", userID)
		RespondTo(self, cmd, "Commands: help, kill, sex, getvro, eat, hug, goonto")
		fmt.Println("Sent help message to ", cmd.ChannelID)
	}
}
