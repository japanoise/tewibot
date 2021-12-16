package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token   string
	BotID   string
	AdminID string
	LogLoc  string
	LogPref string
	Logging bool
	LN      int
)

const (
	GenderNeuter byte = iota
	GenderMale
	GenderFemale
	GenderItIts

	LogR int = 10000
)

var Spouse = [...]string{"spouse", "hazubando", "waifu", "spouse"}
var Child = [...]string{"child", "son", "daughter", "child"}
var Gender = [...]string{"enby", "male", "female", "non-binary"}
var pa = [...]string{"theirs", "his", "hers", "its"}
var ps = [...]string{"they", "he", "she", "it"}
var po = [...]string{"them", "him", "her", "it"}
var pp = [...]string{"their", "his", "her", "its"}
var pr = [...]string{"themself", "himself", "herself", "itself"}

var regexWaifuAffection *regexp.Regexp
var regexKismesis *regexp.Regexp
var regexSpouseKismesis *regexp.Regexp
var regexSpouseNB *regexp.Regexp
var regexSpouseMasc *regexp.Regexp
var regexSpouseFem *regexp.Regexp

type Quadrant byte

const (
	QuadFlushed Quadrant = iota
	QuadPitch
	QuadPale
)

type Human interface {
	GetName() string
	GetGender() byte
}

type BotWaifu struct {
	Name    string
	Gender  byte
	Picture string
	Tag     string
	Theme   string
	Anni    time.Time
	Bday    time.Time
	Quad    Quadrant
}

func (b *BotWaifu) GetName() string { return b.Name }
func (b *BotWaifu) GetGender() byte { return b.Gender }

type BotCmdPair struct {
	Cmd   string
	Reply string
}

type BotUser struct {
	Nickname string
	Gender   byte
	Waifus   []*BotWaifu
	Children []*BotWaifu
	Commands []*BotCmdPair
	Intro    string
	UseQ     bool
}

func (b *BotUser) GetName() string { return b.Nickname }
func (b *BotUser) GetGender() byte { return b.Gender }

type BotState struct {
	Users          map[string]*BotUser
	CommandPrefix  string
	DanbooruLogin  string
	DanbooruAPIKey string
}

type BotCmd func(*discordgo.Session, *discordgo.MessageCreate)

var Global BotState

var Commands map[string]BotCmd
var Usages map[string]string
var Comforts []string
var ChildComforts []string
var ChildReverseComforts []string
var PaleComforts []string
var PitchComforts []string

func reply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	fmt.Printf("me -> %s: %s\n", m.ChannelID, msg)
	_, _ = s.ChannelMessageSend(m.ChannelID, msg)
}

func adduserifne(m *discordgo.MessageCreate) {
	if Global.Users[m.Author.ID] == nil {
		ret := new(BotUser)
		ret.Nickname = m.Author.Username
		Global.Users[m.Author.ID] = ret
	}
}

func setIntro(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	u := Global.Users[m.Author.ID]
	if len(words) > 1 {
		u.Intro = strings.Join(words[1:], " ")
		reply(s, m, fmt.Sprintf("Setting %s's intro to %s", u.Nickname, u.Intro))
	} else {
		reply(s, m, fmt.Sprintf("%s", u.Intro))
	}
}

func fetchWaifu(u *BotUser) *BotWaifu {
	if u.Waifus == nil {
		return nil
	} else if len(u.Waifus) == 0 {
		return nil
	} else {
		return u.Waifus[0]
	}
}

func danbooruPic(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	if len(words) <= 1 {
		reply(s, m, "Please specify a tag to search on Danbooru")
	} else {
		reply(s, m, fetchImage(words[1]))
	}
}

func getWaifuPic(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	u := Global.Users[m.Author.ID]
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		wname := strings.Join(words[1:], " ")
		for _, wifu := range u.Waifus {
			if wifu.Name == wname {
				if wifu.Tag == "" {
					reply(s, m, fmt.Sprintf("Set a tag to use when looking for pictures of %s - &tag some_tag %s", wname, wname))
				} else {
					reply(s, m, fetchImage(wifu.Tag))
				}
			}
		}
		for _, wifu := range u.Children {
			if wifu.Name == wname {
				if wifu.Tag == "" {
					reply(s, m, fmt.Sprintf("Set a tag to use when looking for pictures of %s - &tag some_tag %s", wname, wname))
				} else {
					reply(s, m, fetchImage(wifu.Tag))
				}
			}
		}
	} else {
		tags := []string{}
		for _, wifu := range u.Waifus {
			if wifu.Tag != "" {
				tags = append(tags, wifu.Tag)
			}
		}
		for _, wifu := range u.Children {
			if wifu.Tag != "" {
				tags = append(tags, wifu.Tag)
			}
		}
		if len(tags) == 0 {
			reply(s, m, "Either you don't have any family members set, or none of your family members have danbooru tags.")
		} else {
			reply(s, m, fetchImage(randoms(tags)))
		}
	}
}

func waifuTagAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 2 {
		var tag string = words[1]
		if tag == "" {
			reply(s, m, "Please provide a tag")
			return
		}
		var wname string = strings.Join(words[2:], " ")
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname {
					reply(s, m, fmt.Sprintf("Setting %s's danbooru tag to %s",
						wname, tag))
					waifu.Tag = tag
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname {
					reply(s, m, fmt.Sprintf("Setting %s's danbooru tag to %s",
						wname, tag))
					c.Tag = tag
					return
				}
			}
		}
	}
}

func getGender(s *discordgo.Session, m *discordgo.MessageCreate) {
	var id string
	var u *BotUser
	if len(m.Mentions) > 0 {
		id = m.Mentions[0].ID
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

func getSpouseString(u *BotUser) string {
	wifu := fetchWaifu(u)
	ret := ""
	if wifu == nil {
		ret = fmt.Sprintf("Looks like %s doesn't have a waifu...\n", u.Nickname)
	} else if len(u.Waifus) == 1 {
		pic := ""
		if wifu.Picture != "" {
			pic = "\nPicture: " + wifu.Picture
		}
		if wifu.Theme != "" {
			pic += "\nTheme: " + wifu.Theme + ""
		}
		if !wifu.Anni.IsZero() {
			pic += "\nAnniversary: " + wifu.Anni.Format(shortForm)
		}
		if !wifu.Bday.IsZero() {
			pic += "\nBirthday: " + wifu.Bday.Format(shortForm)
		}
		spouse := "spouse"
		if !u.UseQ {
			spouse = Spouse[wifu.Gender]
		} else if wifu.Quad == QuadPale {
			spouse = "moirail (♢)"
		} else if wifu.Quad == QuadPitch {
			spouse = "kismesis (♤)"
		} else {
			spouse = Spouse[wifu.Gender] + " (♡)"
		}
		ret = fmt.Sprintf(
			"According to the databanks, %s's %s is %s.%s\n",
			u.Nickname, spouse, wifu.Name, pic)
	} else {
		ret = fmt.Sprintf("%s has %d spouses:\n", u.Nickname, len(u.Waifus))
		for i, waifu := range u.Waifus {
			pic := ""
			if waifu.Picture != "" {
				pic = "\nPicture: " + waifu.Picture
			}
			if waifu.Theme != "" {
				pic += "\nTheme: " + waifu.Theme
			}
			if !waifu.Anni.IsZero() {
				pic += "\nAnniversary: " + waifu.Anni.Format(shortForm)
			}
			if !waifu.Bday.IsZero() {
				pic += "\nBirthday: " + waifu.Bday.Format(shortForm)
			}
			spouse := "spouse"
			if !u.UseQ {
				spouse = Spouse[waifu.Gender]
			} else if waifu.Quad == QuadPale {
				spouse = "moirail (♢)"
			} else if waifu.Quad == QuadPitch {
				spouse = "kismesis (♤)"
			} else {
				spouse = Spouse[waifu.Gender] + " (♡)"
			}
			ret += fmt.Sprintf(
				"%d) %s %s, %s.%s\n", i+1,
				pp[u.Gender], spouse, waifu.Name, pic)
		}
	}
	return ret
}

func getWaifu(s *discordgo.Session, m *discordgo.MessageCreate) {
	var id string
	var u *BotUser
	if len(m.Mentions) > 0 {
		id = m.Mentions[0].ID
	} else {
		id = m.Author.ID
	}
	u = Global.Users[id]
	if u == nil {
		reply(s, m, "I've no idea who that is!")
	} else {
		reply(s, m, getSpouseString(u))
	}
}

func getChildString(u *BotUser, child *BotWaifu) string {
	pic := ""
	if child.Picture != "" {
		pic = "\nPicture: " + child.Picture
	}
	if child.Theme != "" {
		pic += "\nTheme " + child.Theme
	}
	if !child.Anni.IsZero() {
		pic += "\nAnniversary: " + child.Anni.Format(shortForm)
	}
	if !child.Bday.IsZero() {
		pic += "\nBirthday: " + child.Bday.Format(shortForm)
	}
	return fmt.Sprintf(
		"\n%s %s, %s.%s",
		pp[u.Gender], Child[child.Gender], child.Name, pic)
}

func getFamily(s *discordgo.Session, m *discordgo.MessageCreate) {
	var id string
	var u *BotUser
	if len(m.Mentions) > 0 {
		id = m.Mentions[0].ID
	} else {
		id = m.Author.ID
	}
	u = Global.Users[id]
	if u == nil {
		reply(s, m, "I've no idea who that is!")
	} else {
		ret := getSpouseString(u)
		if u.Children == nil {
			ret += fmt.Sprintf("Looks like %s doesn't have any children...", u.Nickname)
		} else if len(u.Children) == 0 {
			ret += fmt.Sprintf("Looks like %s doesn't have any children...", u.Nickname)
		} else {
			ret += fmt.Sprintf("%s's children are:", u.Nickname)
			for _, child := range u.Children {
				ret += getChildString(u, child)
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
		if strings.HasPrefix(strings.ToLower(words[1]), "i") {
			gen = GenderItIts
		}
		u.Gender = gen
		reply(s, m, fmt.Sprintf("Setting %s's gender to %s", u.Nickname, Gender[gen]))
	} else {
		reply(s, m, fmt.Sprintf("%s's gender is %s", u.Nickname, Gender[u.Gender]))
	}
}

func comfort(s *discordgo.Session, m *discordgo.MessageCreate) {
	comfortUser(s, m, false, fetchRandWaifu, Comforts)
}

func rcomfort(s *discordgo.Session, m *discordgo.MessageCreate) {
	comfortUser(s, m, true, fetchRandWaifu, Comforts)
}

func comfortUser(s *discordgo.Session, m *discordgo.MessageCreate, rev bool, f func(*BotUser) *BotWaifu, comforts []string) {
	var id string
	var name string
	var u *BotUser
	words := strings.Split(m.Content, " ")
	if len(m.Mentions) > 0 {
		id = m.Mentions[0].ID
		name = m.Mentions[0].Username
	} else {
		id = m.Author.ID
		name = m.Author.Username
	}
	u = Global.Users[id]
	if u == nil {
		reply(s, m, fmt.Sprintf("_cuddles %s close_", name))
	} else {
		name = u.Nickname
		wifu := f(u)
		if id == m.Author.ID && len(words) > 1 {
			var wname string = strings.Join(words[1:], " ")
			if Global.Users[m.Author.ID].Waifus != nil {
				for _, waifu := range u.Waifus {
					if waifu.Name == wname {
						wifu = waifu
					}
				}
			}
		}
		if wifu == nil {
			reply(s, m, fmt.Sprintf("_cuddles %s close_", name))
		} else {
			if wifu.Quad == QuadPitch {
				comforts = PitchComforts
			} else if wifu.Quad == QuadPale {
				comforts = PaleComforts
			}
			if rev {
				reply(s, m, pronouns(wifu, u, randoms(comforts)))
			} else {
				reply(s, m, pronouns(u, wifu, randoms(comforts)))
			}
		}
	}
}

func ccomfort(s *discordgo.Session, m *discordgo.MessageCreate) {
	comfortUser(s, m, false, fetchRandChild, ChildComforts)
}

func crcomfort(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Yes, we don't use the reverse flag. Counter-intuitive? A little.
	// This is one of the few places where legacy from the old bots creeps in:
	// dr/cr comforts are specified with Anon as the %n and the child as the %wn.
	comfortUser(s, m, false, fetchRandChild, ChildReverseComforts)
}

func waifuDel(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	gen := GenderFemale
	if strings.Contains(strings.ToLower(words[0]), "husbando") {
		gen = GenderMale
	}
	if strings.Contains(strings.ToLower(words[0]), "spouse") {
		gen = GenderNeuter
	}
	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		if Global.Users[m.Author.ID].Waifus == nil {
			reply(s, m, "But you don't have a waifu!")
		} else {
			u := Global.Users[m.Author.ID]
			for i, waifu := range u.Waifus {
				if waifu.Name == wname && waifu.Gender == gen {
					reply(s, m, fmt.Sprintf("Removing %s from %s's waifus",
						wname, m.Author.Username))
					copy(u.Waifus[i:], u.Waifus[i+1:])
					u.Waifus[len(u.Waifus)-1] = nil // or the zero value of T
					u.Waifus = u.Waifus[:len(u.Waifus)-1]
					return
				}
			}
			reply(s, m, "Couldn't find that waifu in your waifu list!")
		}
	}
}

func childDel(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		if Global.Users[m.Author.ID].Children == nil {
			reply(s, m, "But you don't have any children!")
		} else {
			u := Global.Users[m.Author.ID]
			for i, child := range u.Children {
				if child.Name == wname {
					reply(s, m, fmt.Sprintf("Removing %s from %s's children",
						wname, m.Author.Username))
					copy(u.Children[i:], u.Children[i+1:])
					u.Children[len(u.Children)-1] = nil // or the zero value of T
					u.Children = u.Children[:len(u.Children)-1]
					return
				}
			}
			reply(s, m, fmt.Sprintf("%s is not one of your children!", wname))
		}
	}
}

const shortForm = "2006-01-02"

// Hairy date code approaching. Patches welcome, but I won't fix it myself.
func prettyDate(date time.Time) string {
	now := time.Now()
	ret := date.Format(shortForm) + "."

	since := now.Sub(date)
	days := since.Hours() / 24
	ret += fmt.Sprintf("\nThat's %d days ago - roughly %d year(s) and %d day(s)!", int(days), int(days/365.25), int(days)%365)

	nextDate, _ := time.Parse(shortForm, fmt.Sprintf("%4d-%02d-%02d", now.Year(), date.Month(), date.Day()))
	if !nextDate.After(now) {
		nextDate, _ = time.Parse(shortForm, fmt.Sprintf("%4d-%02d-%02d", now.Year()+1, date.Month(), date.Day()))
	}
	until := nextDate.Sub(now)
	days = until.Hours() / 24
	ret += fmt.Sprintf("\nIt will next occur on %s - %d days away!", nextDate.Format(shortForm), int(days))

	return ret
}

func anni(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname && !waifu.Anni.IsZero() {
					reply(s, m, "Your anniversary with " + waifu.Name+" is "+prettyDate(waifu.Anni))
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname && !c.Anni.IsZero() {
					reply(s, m, "Your anniversary with " + c.Name+" is "+prettyDate(c.Anni))
					return
				}
			}
		}

		reply(s, m, "Not found, or date not set. Use waifuReg and anniReg")
	} else {
		reply(s, m, "Please tell me whose date you want to know!")
	}
}

func anniReg(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 2 {
		var dateS string = words[1]
		if dateS == "" {
			reply(s, m, "Please provide a date (anni YYYY-MM-DD waifu name)")
			return
		}

		date, err := time.Parse(shortForm, dateS)
		if err != nil {
			reply(s, m, "Please provide the date in YYYY-MM-DD form: "+err.Error())
			return
		}

		var wname string = strings.Join(words[2:], " ")
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname {
					reply(s, m, "Adding the date...")
					waifu.Anni = date
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname {
					reply(s, m, "Adding the date...")
					c.Anni = date
					return
				}
			}
		}
	}
}

func bday(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname && !waifu.Bday.IsZero() {
					reply(s, m, waifu.Name+"'s birthday is "+prettyDate(waifu.Bday))
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname && !c.Bday.IsZero() {
					reply(s, m, c.Name+"'s birthday is "+prettyDate(c.Bday))
					return
				}
			}
		}

		reply(s, m, "Not found, or date not set. Use waifuReg and bdayReg")
	} else {
		reply(s, m, "Please tell me whose date you want to know!")
	}
}

func bdayReg(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 2 {
		var dateS string = words[1]
		if dateS == "" {
			reply(s, m, "Please provide a date (bdayreg YYYY-MM-DD waifu name)")
			return
		}

		date, err := time.Parse(shortForm, dateS)
		if err != nil {
			reply(s, m, "Please provide the date in YYYY-MM-DD form: "+err.Error())
			return
		}

		var wname string = strings.Join(words[2:], " ")
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname {
					reply(s, m, "Adding the date...")
					waifu.Bday = date
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname {
					reply(s, m, "Adding the date...")
					c.Bday = date
					return
				}
			}
		}
	}
}


func waifuPicAdd(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 2 {
		var pic string = words[1]
		if pic == "" {
			reply(s, m, "Please provide a picture")
			return
		}
		var wname string = strings.Join(words[2:], " ")
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname {
					reply(s, m, fmt.Sprintf("Adding a picture of %s - %s",
						wname, pic))
					waifu.Picture = pic
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname {
					reply(s, m, fmt.Sprintf("Adding a picture of %s - %s",
						wname, pic))
					c.Picture = pic
					return
				}
			}
		}
	}
}

func waifuReg(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	gen := GenderFemale
	quadrant := QuadFlushed
	if strings.Contains(strings.ToLower(words[0]), "husbando") {
		gen = GenderMale
	}
	if strings.Contains(strings.ToLower(words[0]), "spouse") {
		gen = GenderNeuter
	}
	if strings.Contains(strings.ToLower(words[0]), "pale") {
		Global.Users[m.Author.ID].UseQ = true
		quadrant = QuadPale
	}
	if strings.Contains(strings.ToLower(words[0]), "pitch") {
		Global.Users[m.Author.ID].UseQ = true
		quadrant = QuadPitch
	}

	spouse := Spouse[gen]
	if quadrant == QuadPale {
		spouse = "moirail"
	} else if quadrant == QuadPitch {
		spouse = "kismesis"
	}

	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		waifu := newWaifu(wname, gen)
		waifu.Quad = quadrant
		if Global.Users[m.Author.ID].Waifus == nil {
			Global.Users[m.Author.ID].Waifus = []*BotWaifu{
				waifu,
			}
		} else {
			for idx, wifu := range Global.Users[m.Author.ID].Waifus {
				if wname == wifu.Name {
					Global.Users[m.Author.ID].Waifus[idx].Quad = quadrant
					Global.Users[m.Author.ID].Waifus[idx].Gender = gen
					reply(s, m, fmt.Sprintf("Updating %s's %s to %s",
						m.Author.Username, spouse, wname))
					return
				}
			}
			Global.Users[m.Author.ID].Waifus = append(Global.Users[m.Author.ID].Waifus,
				waifu)
		}
		reply(s, m, fmt.Sprintf("Setting %s's %s to %s",
			m.Author.Username, spouse, wname))
		fmt.Println(m.Author.ID, spouse, wname)
	}
}

func vacillate(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")

	if len(words) > 1 {
		var wname string = strings.Join(words[1:], " ")
		u := Global.Users[m.Author.ID]
		u.UseQ = true
		for i, waifu := range u.Waifus {
			if waifu.Name == wname {
				if waifu.Quad == QuadFlushed {
					reply(s, m, fmt.Sprintf(
						"Aww, looks like someone caught pitch feelings for %s~ ♤", waifu.Name))
					u.Waifus[i].Quad = QuadPitch
				} else if waifu.Quad == QuadPitch {
					reply(s, m, fmt.Sprintf(
						"Aww, looks like someone caught flushed feelings for %s~ ♡", waifu.Name))
					u.Waifus[i].Quad = QuadFlushed
				} else {
					reply(s, m, "I don't think that quadrant vacillates...")
				}
				return
			}
		}
		reply(s, m, "Couldn't find that waifu in your waifu list!")
	} else {
		reply(s, m, "Who is vacillating?")
	}
}

func newWaifu(name string, gen byte) *BotWaifu {
	ret := &BotWaifu{}
	ret.Name = name
	ret.Gender = gen
	return ret
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
			Global.Users[m.Author.ID].Children = []*BotWaifu{
				newWaifu(wname, gen),
			}
		} else {
			Global.Users[m.Author.ID].Children = append(
				Global.Users[m.Author.ID].Children, newWaifu(wname, gen))
		}
		reply(s, m, fmt.Sprintf("Setting %s's %s to %s",
			m.Author.Username, child, wname))
		fmt.Println(m.Author.ID, child, wname)
	}
}

func themeAddOrGet(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		var pic string = words[1]
		if pic == "" {
			reply(s, m, "Please provide a picture or family member name")
			return
		}
		dopic := strings.HasPrefix(pic, "http://") ||
			strings.HasPrefix(pic, "https://")
		var wname string
		if dopic {
			if len(words) < 3 {
				reply(s, m, "You need to specify who gets the theme song!")
				return
			}
			wname = strings.Join(words[2:], " ")
		} else {
			wname = strings.Join(words[1:], " ")
		}
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			for _, waifu := range u.Waifus {
				if waifu.Name == wname {
					if dopic {
						reply(s, m, fmt.Sprintf("Adding a theme for %s - %s",
							wname, pic))
						waifu.Theme = pic
					} else if waifu.Theme == "" {
						reply(s, m, "No theme found, you can add one with the theme command")
					} else {
						reply(s, m, fmt.Sprintf("%s's theme is %s",
							wname, waifu.Theme))
					}
					return
				}
			}
		}

		if Global.Users[m.Author.ID].Children != nil {
			u := Global.Users[m.Author.ID]
			for _, c := range u.Children {
				if c.Name == wname {
					if dopic {
						reply(s, m, fmt.Sprintf("Adding a theme for %s - %s",
							wname, pic))
						c.Theme = pic
					} else if c.Theme == "" {
						reply(s, m, "No theme found, you can add one with the theme command")
					} else {
						reply(s, m, fmt.Sprintf("%s's theme is %s",
							wname, c.Theme))
					}
					return
				}
			}
		}
	} else {

		reply(s, m, "Not enough arguments. Format: &theme LINK WAIFU or &theme WAIFU")
	}
}

func addBotCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	u := Global.Users[m.Author.ID]
	if len(u.Commands) >= 5 {
		reply(s, m, "You already have 5 custom commands.")
		return
	}
	words := strings.Split(m.Content, " ")
	if len(words) > 2 {
		rply := strings.Join(words[2:], " ")
		for i, cmd := range u.Commands {
			if cmd.Cmd == words[1] {
				u.Commands[i].Reply = rply
				reply(s, m, fmt.Sprintf("Updated command %s", cmd.Cmd))
				return
			}
		}

		reply(s, m, fmt.Sprintf("Added command %s", words[1]))
		u.Commands = append(u.Commands, &BotCmdPair{words[1], rply})
	} else {
		reply(s, m, "Not enough arguments. Format: &addcmd COMMAND REPLY...")
	}
}

func delBotCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		u := Global.Users[m.Author.ID]
		for _, delcmd := range words[1:] {
			done := false
			for i, cmd := range u.Commands {
				if delcmd == cmd.Cmd {
					copy(u.Commands[i:], u.Commands[i+1:])
					u.Commands[len(u.Commands)-1] = nil
					u.Commands = u.Commands[:len(u.Commands)-1]
					done = true
					reply(s, m, fmt.Sprintf("Custom command %s was deleted", delcmd))					
					break
				}
			}
			if !done {
				reply(s, m, fmt.Sprintf("Warning: %s was not deleted (this command is case sensitive)", delcmd))
			}
		}
	} else {
		reply(s, m, "Please provide one or more commands to delete")
	}
}

func lsBotCmd(s *discordgo.Session, m *discordgo.MessageCreate) {
	adduserifne(m)
	u := Global.Users[m.Author.ID]
	if len(u.Commands) <= 0 {
		reply(s, m, "You have no custom commands.")
	} else {
		reps := []string{"Your custom commands are:"}
		for _, cmd := range u.Commands {
			reps = append(reps, fmt.Sprintf("%s: \"%s\"", cmd.Cmd, cmd.Reply))
		}
		reply(s, m, strings.Join(reps, "\n"))
	}
}

func help(s *discordgo.Session, m *discordgo.MessageCreate) {
	words := strings.Split(m.Content, " ")
	if len(words) > 1 {
		cmd := strings.TrimPrefix(strings.Join(words[1:], " "), Global.CommandPrefix)
		if Usages[cmd] == "" {
			reply(s, m, fmt.Sprintf("The help system doesn't know about %s%s",
				Global.CommandPrefix, cmd))
		} else if Commands[cmd] == nil {
			reply(s, m, Usages[cmd])
		} else {
			reply(s, m, fmt.Sprintf("%s%s - %s", Global.CommandPrefix,
				cmd, Usages[cmd]))
		}
	} else {
		reply(s, m, HelpMenu(Global.CommandPrefix))
	}
}

func adminInfo(s *discordgo.Session, m *discordgo.MessageCreate) {
	rep := ""
	if AdminID == "" {
		rep = "There is no admin set for the bot."
	} else if isSenderAdmin(m) {
		rep = "You are the admin."
	} else {
		admin, err := s.User(AdminID)
		if err == nil {
			rep = fmt.Sprintf("%s is the admin.", admin.Mention())
		} else {
			rep = fmt.Sprintf("%s is the admin, but I can't mention them...", AdminID)
		}
	}
	reply(s, m, rep)
}

func isSenderAdmin(m *discordgo.MessageCreate) bool {
	return AdminID != "" && m.Author.ID == AdminID
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
	addCommand(bday, "Show a birthday", "bday")
	addCommand(bdayReg, "Register a birthday (YYYY-MM-DD)", "bdayreg")
	addCommand(anni, "Show your anniversary", "anni")
	addCommand(anniReg, "Register your anniversary (YYYY-MM-DD)", "annireg")
	addCommand(waifuReg, "Register your waifu with the bot", "waifureg",
		"husbandoreg", "setwaifu", "sethusbando", "spousereg", "setspouse",
		"palewaifureg", "palehusbandoreg", "palespousereg",
		"pitchwaifureg", "pitchhusbandoreg", "pitchspousereg")
	addCommand(waifuDel, "Delete a previously registered waifu", "waifudel", "husbandodel", "spousedel")
	addCommand(childDel, "Delete a previously registered child", "daughterdel", "sondel", "childdel")
	addCommand(getGender, "Print your (or someone else's) gender", "gender", "getgender")
	addCommand(getWaifu, "Print your (or someone else's) waifu", "waifu", "husbando", "spouse")
	addCommand(comfort, "Dispense hugs and other niceness from your waifu", "comfort", "hug")
	addCommand(rcomfort, "Dispense hugs and other niceness to your waifu", "rcomfort", "rhug")
	addCommand(ccomfort, "Dispense hugs and other niceness from your child", "ccomfort", "dcomfort", "chug", "dhug")
	addCommand(crcomfort, "Dispense hugs and other niceness to your child", "crcomfort", "drcomfort", "crhug", "drhug")
	addCommand(setGender, "Set your gender - m, f, x, i\nThis affects which pronouns the bot will use for you (he, she, they, it)", "setgender", "genderreg")
	addCommand(addChild, "Register one of your children with the bot", "setchild", "childreg", "setdaughteru", "daughterureg", "setsonfu", "sonfureg")
	addCommand(getFamily, "Print your (or someone else's) family", "family", "getfamily")
	addCommand(nickname, "If given a nickname, set your nickname to that. Otherwise, print your nickname.", "nick", "nickname", "setnick", "setnickname")
	addCommand(help, "Access the on-line help system", "help", "usage", "sos")
	addCommand(adminInfo, "Print information about the admin", "admin")
	addCommand(waifuPicAdd, "Add a picture to your waifu; e.g. &picadd http://i.imgur.com/Gqf1rGi.jpg Miku", "picadd")
	addCommand(danbooruPic, "Fetch an image with the given tag from danbooru", "danbooru")
	addCommand(waifuTagAdd, "Set your child or waifu's tag to use when searching danbooru", "tag")
	addCommand(getWaifuPic, "Get an image of your waifu or child from danbooru", "pic")
	addCommand(setIntro, "Set or display your introduction", "intro")
	addCommand(themeAddOrGet, "Set or get your waifu or child's theme, e.g. &theme https://www.youtube.com/watch?v=U_CfriU4Cng Miku", "theme")
	addCommand(addBotCmd, "Add a custom command, e.g. &addcmd !tsun Miku: It's not like I like you or anything, BAKA!", "addcmd")
	addCommand(delBotCmd, "Delete a custom command, works with multiple commands e.g. &delcmd !yan !tsun !kuu", "delcmd")
	addCommand(vacillate, "Vacillate a waifu between flushed and pitch (if you don't know what this means, you probably don't want to do it)", "vax", "vacillate")
	addCommand(lsBotCmd, "Lists your custom commands", "lscmd")
	InitGlobal()
	InitComforts()
	InitHelp()

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&AdminID, "a", "", "Admin's Discord ID")
	flag.StringVar(&LogLoc, "l", "", "Place to put logs in")
	flag.StringVar(&LogPref, "p", "", "Prefix for lines to include in logs")
	flag.Parse()
	Logging = LogLoc != ""
	if Logging {
		log.SetFlags(log.Ldate | log.Ltime | log.LUTC)
		logRotate()
	}

	// Let's see you ROTC leeches coming up with regexes even half this good. -- Kona
	regexWaifuAffection = regexp.MustCompile("^[Ii] ((re+a+l+y+ )*lo+ve|ne+d|wa+nt|a+do+r+e+) (my )?([^.,!?]*)")
	regexKismesis = regexp.MustCompile("^[Ii] ((re+a+l+y+ )*ha+te) (my )?([^.,!?]*)")
	regexSpouseKismesis = regexp.MustCompile("^ki+sme+si+s+$")
	regexSpouseNB = regexp.MustCompile("^(spo+u+s+e+|da+te+ma+te+)$")
	regexSpouseMasc = regexp.MustCompile("^(h[au]+[sz]u*bando*|bo+yfri+e+nd+)$")
	regexSpouseFem = regexp.MustCompile("^(wa*i+f[ue]+|gi+r+l+fri+e+nd+)$")
}

var logfile *os.File = nil

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
	dg.AddHandler(messageDelete)

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
	if logfile != nil {
		logfile.Close()
	}
}

func genLogName() string {
	return LogLoc + string(os.PathSeparator) + "log-" + strconv.FormatInt(time.Now().Unix(), 10) + ".txt"
}

func logRotate() {
	if logfile != nil {
		logfile.Close()
	}
	logfile, err := os.OpenFile(genLogName(), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logfile = nil
		fmt.Println(err.Error())
		log.SetOutput(os.Stdout)
		return
	}
	log.SetOutput(logfile)
}

func logMsg(fomt string, args ...interface{}) {
	msg := fmt.Sprintf(fomt, args...)
	fmt.Print(msg)
	if Logging && strings.HasPrefix(msg, LogPref) {
		log.Print(msg)
		logfile.Sync()
		LN++
		if LN > LogR {
			logRotate()
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == BotID {
		return
	}

	ch, err := s.Channel(m.ChannelID)
	if err != nil {
		logMsg("[%s] %s (%s) -> %s: %s\n", "unknown guild", m.Author.Username, m.Author.ID, m.ChannelID, m.Content)
	} else {
		logMsg("[%s] %s (%s) -> %s (%s): %s\n", ch.GuildID, m.Author.Username, m.Author.ID, ch.Name, m.ChannelID, m.Content)
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		reply(s, m, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		reply(s, m, "Ping!")
	}

	// Informal standard
	if m.Content == ".bots" {
		reply(s, m, "tewibot reporting in! [Golang] https://github.com/japanoise/tewibot")
	}

	if mat := regexWaifuAffection.FindStringSubmatch(m.Content); mat != nil {
		affectionVerb := mat[1]
		spouseOrName := mat[4]
		if len(spouseOrName) == 0 {
			return
		}
		if regexSpouseNB.MatchString(spouseOrName) {
			reply(s, m, fmt.Sprintf("I'm sure they %s you too!", affectionVerb))
		} else if regexSpouseMasc.MatchString(spouseOrName) {
			reply(s, m, fmt.Sprintf("I'm sure he %ss you too!", affectionVerb))
		} else if regexSpouseFem.MatchString(spouseOrName) {
			reply(s, m, fmt.Sprintf("I'm sure she %ss you too!", affectionVerb))
		} else {
			adduserifne(m)
			if Global.Users[m.Author.ID].Waifus != nil {
				u := Global.Users[m.Author.ID]
				for _, waifu := range u.Waifus {
					if waifu.Name == spouseOrName {
						reply(s, m, fmt.Sprintf("%s %ss you too!", waifu.Name, affectionVerb))
						return
					}
				}
			}
		}
	} else if mat := regexKismesis.FindStringSubmatch(m.Content); mat != nil {
		adduserifne(m)
		kmgender := GenderNeuter
		if Global.Users[m.Author.ID].Waifus != nil {
			u := Global.Users[m.Author.ID]
			km := false
			for _, waifu := range u.Waifus {
				km = km || waifu.Quad == QuadPitch
				if waifu.Quad == QuadPitch {
					kmgender = waifu.Gender
				}
			}
			if !km {
				return
			}
		}
		affectionVerb := mat[1]
		spouseOrName := mat[4]
		if len(spouseOrName) == 0 {
			return
		}
		if regexSpouseKismesis.MatchString(spouseOrName) {
			if kmgender == GenderNeuter {
				reply(s, m, fmt.Sprintf("I'm sure %s %s you too!", ps[kmgender], affectionVerb))
			} else {
				reply(s, m, fmt.Sprintf("I'm sure %s %ss you too!", ps[kmgender], affectionVerb))
			}
		} else {
			if Global.Users[m.Author.ID].Waifus != nil {
				u := Global.Users[m.Author.ID]
				for _, waifu := range u.Waifus {
					if waifu.Name == spouseOrName && waifu.Quad == QuadPitch {
						reply(s, m, fmt.Sprintf("%s %ss you too!", waifu.Name, affectionVerb))
						return
					}
				}
			}
		}
	}

	if len(m.Content) > len(Global.CommandPrefix) {
		if strings.HasPrefix(m.Content, Global.CommandPrefix) {
			run := Commands[strings.ToLower(strings.Split(strings.TrimPrefix(
				m.Content, Global.CommandPrefix), " ")[0])]
			if run != nil {
				run(s, m)
			}
			return
		}
	}

	if len(m.Content) > 0 {
		adduserifne(m)
		for _, cmd := range Global.Users[m.Author.ID].Commands {
			if m.Content == cmd.Cmd {
				reply(s, m, cmd.Reply)
			}
		}
	}
}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	fmt.Printf("%+v\n", m.Message)
}
