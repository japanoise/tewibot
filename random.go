package main

import (
	"math/rand"
)

func randoms(fr []string) string {
	return fr[rand.Intn(len(fr))]
}

func fetchRandWaifu(u *BotUser) *BotWaifu {
	if u.Waifus == nil {
		return nil
	} else if len(u.Waifus) == 0 {
		return nil
	} else if len(u.Waifus) == 1 {
		return u.Waifus[0]
	} else {
		return u.Waifus[rand.Intn(len(u.Waifus))]
	}
}

func fetchRandChild(u *BotUser) *BotWaifu {
	if u.Children == nil {
		return nil
	} else if len(u.Children) == 0 {
		return nil
	} else if len(u.Children) == 1 {
		return u.Children[0]
	} else {
		return u.Children[rand.Intn(len(u.Children))]
	}
}
