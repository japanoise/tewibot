package main

import (
	"fmt"
	"strings"
)

func InitHelp() {
	Usages["basic"] = strings.Join([]string{
		"tewibot keeps some basic information about users.",
		"To set your gender, use genderreg followed by m, f, x, or i (he, she, they, it).",
		"To set or update your nickname, use setnick",
		"",
		"tewibot's primary feature is keeping a database of people's partners and kiddos.",
		"To register a partner, use waifureg, husbandoreg, spousereg for a female, male, non-binary spouse respectively.",
		"To delete a partner, use waifudel, husbandodel, spousedel.",
		"To add a kiddo, use daughterureg, sonfureg, childreg.",
		"To delete a kiddo, use daughterdel, sondel, childdel.",
		"The 'family' command displays your family, or if you ping someone in the command, theirs.",
		"",
		"You can add additional information to a partner. The usual format is e.g. &picadd <http://i.imgur.com/Gqf1rGi.jpg> Miku",
		"To add a picture, use picadd.",
		"To set a family member's theme, use theme.",
		"To set a family member's birthday, use bdayreg. The date format (VERY fussy) is YYYY-MM-DD. Date first, then waifu.",
		"Setting a family member's anniversary with you is the same, using the command annireg.",
		"",
		"Once your family is setup, tewibot can comfort you. The basic 'comfort' command will pick a random family member; but if you specify, you can choose a particular family member to comfort you.",
		"rcomfort comforts you in reverse (i.e. 'Waifufriend hugs Miku' rather than 'Miku hugs Waifufriend').",
		"ccomfort comforts only with kiddos.",
		"crcomfort reverse comforts only with kiddos.",
	}, "\n")
	Usages["booru"] = strings.Join([]string{
		"1. First, make sure you have registered your Significant Other! If you haven't, you can register them with \"&waifu\" or \"&husbando\", and then putting your S/O's name afterwards.",
		"For example, '&waifu Monika'",
		"This also allows you to use commands such as \"&comfort\" or \"&hug\"",
		"You will be using the name you registered them as in the next steps.",
		"",
		"2. To use &pic, you must specify a tag on Danbooru.",
		"Type \"&tag\", and then put the Danbooru tag after '&tag', and then after that, put their name!",
		"For example, \"&tag monika_(doki_doki_literature_club) Monika\"",
		"Make sure you are using their correct tag on Danbooru, and the name that you registered them as earlier.",
		"",
		"Remember, the bot automatically adds rating: safe to every tag, but some people don't properly tag their images. You may encounter a picture that's unpleasant. Some bot channels allow you to delete these pictures though, so be aware!",
	}, "\n")
	Usages["customcmd"] = strings.Join([]string{
		"tewibot can register custom commands unique to you. These do not need a prefix, i.e. a command called 'tsun' can be activated by just typing 'tsun'.",
		"Add a new custom command with addcmd. The format is addcmd [cmd] [response] (commands must be one word).",
		"Delete a custom command with delcmd.",
		"List your custom commands with lscmd.",
	}, "\n")
	Usages["quadrants"] = strings.Join([]string{
		"tewibot supports adding a kismesis or a moirail. If you don't know what that means, you probably don't want it.",
		"To register a pale partner, use palewaifureg, palehusbandoreg, palespousereg for a female, male, non-binary moirail respectively.",
		"To register a pitch partner, use pitchwaifureg, pitchhusbandoreg, pitchspousereg for a female, male, non-binary kismesis respectively.",
		"",
		"Pale and pitch partners have special comfort texts. If you want to vaccilate your kismesis to a flushed quadrant, use the 'vax' command.",
	}, "\n")
}

func HelpMenu(prefix string) string {
	return fmt.Sprintf(strings.Join([]string{
		"tewibot - a spiritual successor to the lainbot family\n",
		"prefix: %s\n",
		"type %shelp [topic] to get help on these topics:",
		"- basic",
		"- booru",
		"- customcmd",
		"- quadrants",
	}, "\n"), prefix, prefix)
}
