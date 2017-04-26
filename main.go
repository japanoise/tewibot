package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
	BotID string
)

const (
	GenderNeuter byte = iota
	GenderMale
	GenderFemale
)

var Spouse = [...]string{"spouse", "hazubando", "waifu"}

type BotWaifu struct {
	Name   string
	Gender byte
}

type BotUser struct {
	Nickname string
	Gender   byte
	Waifus   []BotWaifu
}

type BotState struct {
	Users         map[string]*BotUser
	CommandPrefix string
}

type BotCmd func(*discordgo.Session, *discordgo.MessageCreate)

var Global BotState

var Commands map[string]BotCmd

func reply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	fmt.Printf("me -> %s: %s\n", m.ChannelID, msg)
	_, _ = s.ChannelMessageSend(m.ChannelID, msg)
}

func hl2id(hl string) string {
	return strings.TrimPrefix(strings.TrimSuffix(hl, ">"), "<@!")
}

func adduserifne(m *discordgo.MessageCreate) {
	if Global.Users[m.Author.ID] == nil {
		ret := new(BotUser)
		ret.Nickname = m.Author.Username
		Global.Users[m.Author.ID] = ret
	}
}

func fetchWaifu(u *BotUser) *BotWaifu {
	if u.Waifus == nil {
		return nil
	} else if len(u.Waifus) == 0 {
		return nil
	} else {
		return &u.Waifus[0]
	}
}

func getWaifu(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	var id string
	var u *BotUser
	if len(words) > 1 {
		id = hl2id(words[1])
	} else {
		id = m.Author.ID
	}
	u = Global.Users[id]
	if u == nil {
		reply(s, m, "I've no idea who that is!")
	} else {
		wifu := fetchWaifu(u)
		if wifu == nil {
			reply(s, m, fmt.Sprintf("Looks like %s doesn't have a waifu...", u.Nickname))
		} else {
			reply(s, m, fmt.Sprintf(
				"According to the databanks, %s's %s is %s",
				u.Nickname, Spouse[wifu.Gender], wifu.Name))
		}
	}
}

func comfort(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	var id string
	var name string
	var u *BotUser
	if len(words) > 1 {
		id = hl2id(words[1])
		name = strings.Join(words[1:], " ")
	} else {
		id = m.Author.ID
		name = m.Author.Username
	}
	u = Global.Users[id]
	if u == nil {
		reply(s, m, fmt.Sprintf("_cuddles %s close_", name))
	} else {
		name = u.Nickname
		wifu := fetchWaifu(u)
		if wifu == nil {
			reply(s, m, fmt.Sprintf("_cuddles %s close_", name))
		} else {
			reply(s, m, fmt.Sprintf("_%s cuddles %s close_", wifu.Name, name))
		}
	}
}

func waifuReg(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	gen := GenderFemale
	if strings.Contains(strings.ToLower(words[0]), "husbando") {
		gen = GenderMale
	}
	if strings.Contains(strings.ToLower(words[0]), "spouse") {
		gen = GenderNeuter
	}
	spouse := Spouse[gen]
	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		Global.Users[m.Author.ID].Waifus = []BotWaifu{
			BotWaifu{wname, gen},
		}
		reply(s, m, fmt.Sprintf("Setting %s's %s to %s",
			m.Author.Username, spouse, wname))
		fmt.Println(m.Author.ID, spouse, wname)
	}
}

func addCommand(c BotCmd, aliases ...string) {
	for _, alias := range aliases {
		Commands[alias] = c
	}
}

func init() {
	Commands = make(map[string]BotCmd)
	addCommand(waifuReg, "waifureg", "husbandoreg", "setwaifu", "sethusbando", "spousereg", "setspouse")
	addCommand(getWaifu, "waifu", "husbando", "spouse")
	addCommand(comfort, "comfort", "hug")
	InitGlobal()

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Get the account information.
	u, err := dg.User("@me")
	if err != nil {
		fmt.Println("error obtaining account details,", err)
	}

	// Store the account ID for later use.
	BotID = u.ID

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sigchan := make(chan os.Signal, 10)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	fmt.Println("Recieved interrupt, exiting gracefully")
	SaveGlobal()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	fmt.Printf("%s (%s) -> %s: %s\n", m.Author.Username, m.Author.ID, m.ChannelID, m.Content)

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		reply(s, m, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		reply(s, m, "Ping!")
	}

	if len(m.Content) > len(Global.CommandPrefix) {
		if strings.HasPrefix(m.Content, Global.CommandPrefix) {
			run := Commands[strings.ToLower(strings.Split(strings.TrimPrefix(
				m.Content, Global.CommandPrefix), " ")[0])]
			if run != nil {
				run(s, m)
			}
		}
	}
}
