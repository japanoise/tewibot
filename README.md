# tewibot

A Discord bot for waifufriends; the spiritual successor to the lainbot family.

## features

- Fast and light on resources
- Setting and getting waifus and children
- Getting/giving `comfort`s
- Four genders (male, female, non-binary with singular they or it) - no misgendering or your money back!
- On-line help system - use the `help` command
- Set your own unique nickname with the bot, independant from your Discord username
- Remembers your nickname, gender, and family across servers
- Not descended from lainbot's or kkdwfb's code
- Saves to easy-to-read, easy-to-modify, and portable JavaScript Object Notation (json)
- Support for polyamory (Ree some more `:^)`)
- Add pictures of your family to the bot
- Get a random image from Danbooru, or request a picture of a family member from Danbooru
- Get the bot to introduce you to newcomers

## plans

- Quotes database
- Custom comforts
- Steal some features from lainbot forks
- The usual bot features (reminders, memos, weather etc.)

## setup

1. Install Go
2. Make sure $GOPATH is set and that this directory resides within $GOPATH/src/
3. `go get`
4. `go build`
5. Get your bot token from Discord
6. `./tewibot -t <bot token>`

## comforts

Comfort texts are stored in comforts.json. The format string uses pronoun
substitution ala FBMuck, but with a preceding 'w' for waifu's name & pronouns:

    %a (absolute)       = his, hers, theirs, its.
    %s (subjective)     = he, she, they, it.
    %o (objective)      = him, her, them, it.
    %p (possessive)     = his, her, their, its.
    %r (reflexive)      = himself, herself, themself, itself.
    %n (player's name)  = Name.

## Danbooru

Danbooru integration requires a Danbooru account and API key - admin must login
to Danbooru, generate an API key on their user page, then shut down the bot and
enter these into `waifus.json`.

Having done that, test out the functionality by using the `&danbooru tag`
command. If you have family members, you can add a tag to be associated with
them using `&tag tag name` and fetch pictures using the `&pic` command.

## Forks

- [EmiBot](https://github.com/EmiSu9999/EmiBot) is EmiSu9999's fork of Tewibot.
  It's broadly similar but more dialled towards usage in r/waifuism's Discord.
- I'm told that at some point a fork of Tewibot called Reisenbot was running in
  /mai/'s Discord.
