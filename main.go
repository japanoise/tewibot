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
var Child = [...]string{"child", "son", "daughter"}
var Gender = [...]string{"enby", "male", "female"}
var pa = [...]string{"theirs", "his", "hers"}
var ps = [...]string{"they", "he", "she"}
var po = [...]string{"them", "him", "her"}
var pp = [...]string{"their", "his", "her"}
var pr = [...]string{"themself", "himself", "herself"}

type Human interface {
	GetName() string
	GetGender() byte
}

type BotWaifu struct {
	Name   string
	Gender byte
}

func (b *BotWaifu) GetName() string { return b.Name }
func (b *BotWaifu) GetGender() byte { return b.Gender }

type BotUser struct {
	Nickname string
	Gender   byte
	Waifus   []BotWaifu
	Children []BotWaifu
}

func (b *BotUser) GetName() string { return b.Nickname }
func (b *BotUser) GetGender() byte { return b.Gender }

type BotState struct {
	Users         map[string]*BotUser
	CommandPrefix string
}

type BotCmd func(*discordgo.Session, *discordgo.MessageCreate)

var Global BotState

var Commands map[string]BotCmd
var Usages map[string]string
var Comforts []string

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

func getGender(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		gen := u.Gender
		reply(s, m, fmt.Sprintf("%s's gender is %s (%s, %s)", u.Nickname, Gender[gen],
			ps[gen], po[gen]))
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

func getFamily(s *discordgo.Session, m *discordgo.MessageCreate) {
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
		ret := ""
		if wifu == nil {
			ret = fmt.Sprintf("Looks like %s doesn't have a waifu...\n", u.Nickname)
		} else {
			ret = fmt.Sprintf(
				"According to the databanks, %s's %s is %s\n",
				u.Nickname, Spouse[wifu.Gender], wifu.Name)
		}
		if u.Children == nil {
			ret += fmt.Sprintf("Looks like %s doesn't have any children...", u.Nickname)
		} else if len(u.Children) == 0 {
			ret += fmt.Sprintf("Looks like %s doesn't have any children...", u.Nickname)
		} else {
			ret += fmt.Sprintf("%s's children are:", u.Nickname)
			for _, child := range u.Children {
				ret += fmt.Sprintf(
					"\n%s %s, %s",
					pp[u.Gender], Child[child.Gender], child.Name)
			}
		}
		reply(s, m, ret)
	}
}

func pronouns(user Human, waifu Human, str string) string {
	ug := user.GetGender()
	wg := waifu.GetGender()
	ret := str
	ret = strings.Replace(ret, "%a", pa[ug], -1)
	ret = strings.Replace(ret, "%wa", pa[wg], -1)
	ret = strings.Replace(ret, "%s", ps[ug], -1)
	ret = strings.Replace(ret, "%ws", ps[wg], -1)
	ret = strings.Replace(ret, "%o", po[ug], -1)
	ret = strings.Replace(ret, "%wo", po[wg], -1)
	ret = strings.Replace(ret, "%p", pp[ug], -1)
	ret = strings.Replace(ret, "%wp", pp[wg], -1)
	ret = strings.Replace(ret, "%r", pr[ug], -1)
	ret = strings.Replace(ret, "%wr", pr[wg], -1)
	ret = strings.Replace(ret, "%n", user.GetName(), -1)
	ret = strings.Replace(ret, "%wn", waifu.GetName(), -1)
	return ret
}

func nickname(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	u := Global.Users[m.Author.ID]
	if len(words) > 1 {
		newnick := strings.Join(words[1:], " ")
		reply(s, m, fmt.Sprintf("Setting %s's nickname to %s", u.Nickname, newnick))
		u.Nickname = newnick
	} else {
		reply(s, m, fmt.Sprintf("Your nickname is %s", u.Nickname))
	}
}

func setGender(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	u := Global.Users[m.Author.ID]
	if len(words) > 1 {
		gen := GenderNeuter
		if strings.HasPrefix(strings.ToLower(words[1]), "f") {
			gen = GenderFemale
		}
		if strings.HasPrefix(strings.ToLower(words[1]), "m") {
			gen = GenderMale
		}
		u.Gender = gen
		reply(s, m, fmt.Sprintf("Setting %s's gender to %s", u.Nickname, Gender[gen]))
	} else {
		reply(s, m, fmt.Sprintf("%s's gender is %s", u.Nickname, Gender[u.Gender]))
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
			reply(s, m, pronouns(u, wifu, randoms(Comforts)))
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

func addChild(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	gen := GenderFemale
	if strings.Contains(strings.ToLower(words[0]), "son") {
		gen = GenderMale
	}
	if strings.Contains(strings.ToLower(words[0]), "child") {
		gen = GenderNeuter
	}
	child := Child[gen]
	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		if Global.Users[m.Author.ID].Children == nil {
			Global.Users[m.Author.ID].Children = []BotWaifu{
				BotWaifu{wname, gen},
			}
		} else {
			Global.Users[m.Author.ID].Children = append(
				Global.Users[m.Author.ID].Children, BotWaifu{wname, gen})
		}
		reply(s, m, fmt.Sprintf("Setting %s's %s to %s",
			m.Author.Username, child, wname))
		fmt.Println(m.Author.ID, child, wname)
	}
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		cmd := strings.TrimPrefix(strings.Join(words[1:], " "), Global.CommandPrefix)
		if Usages[cmd] == "" {
			reply(s, m, fmt.Sprintf("The help system doesn't know about %s%s",
				Global.CommandPrefix, cmd))
		} else {
			reply(s, m, fmt.Sprintf("%s%s - %s", Global.CommandPrefix,
				cmd, Usages[cmd]))
		}
	} else {
		rep := "tewibot - a spiritual successor to the lainbot family of irc bots.\nSupported commands (type &help _command_ for usage text):\n"
		for key, _ := range Commands {
			rep += Global.CommandPrefix + key + ", "
		}
		reply(s, m, rep)
	}
}

func addCommand(c BotCmd, usage string, aliases ...string) {
	for _, alias := range aliases {
		Commands[alias] = c
		Usages[alias] = usage
	}
}

func init() {
	Commands = make(map[string]BotCmd)
	Usages = make(map[string]string)
	addCommand(waifuReg, "Register your waifu with the bot", "waifureg", "husbandoreg", "setwaifu", "sethusbando", "spousereg", "setspouse")
	addCommand(getGender, "Print your (or someone else's) gender", "gender", "getgender")
	addCommand(getWaifu, "Print your (or someone else's) waifu", "waifu", "husbando", "spouse")
	addCommand(comfort, "Dispense hugs and other niceness", "comfort", "hug")
	addCommand(setGender, "Set your gender - m, f, x\nThis affects which pronouns the bot will use for you (he, she, they)", "setgender", "genderreg")
	addCommand(addChild, "Register one of your children with the bot", "setchild", "childreg", "setdaughteru", "daughterureg", "setsonfu", "sonfureg")
	addCommand(getFamily, "Print your (or someone else's) family", "family", "getfamily")
	addCommand(nickname, "If given a nickname, set your nickname to that. Otherwise, print your nickname.", "nick", "nickname", "setnick", "setnickname")
	addCommand(help, "Access the on-line help system", "help", "usage", "sos")
	InitGlobal()
	InitComforts()

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