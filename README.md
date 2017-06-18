# tewibot

A Discord bot for waifufriends; the spiritual successor to the lainbot family.

## features

- Fast and light on resources
- Setting and getting waifus and children
- Getting/giving `comfort`s
- Three genders (male, female, non-binary with singular they) - no misgendering or your money back!
- On-line help system - use the `help` command
- Set your own unique nickname with the bot, independant from your Discord username
- Remembers your nickname, gender, and family across servers
- Not descended from lainbot's or kkdwfb's code
- Saves to easy-to-read, easy-to-modify, and portable JavaScript Object Notation (json)
- Support for polyamory (Ree some more `:^)`)
- Add pictures of your family to the bot

## plans

- Quotes database
- Custom comforts
- `rcomfort`
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

    %a (absolute)       = his, hers, theirs.
    %s (subjective)     = he, she, they.
    %o (objective)      = him, her, them.
    %p (possessive)     = his, her, their.
    %r (reflexive)      = himself, herself, themself.
    %n (player's name)  = Name.
