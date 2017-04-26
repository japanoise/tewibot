package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func InitGlobal() {
	f, err := os.Open("waifus.json")
	if err == nil {
		dec := json.NewDecoder(f)
		if err = dec.Decode(&Global); err != nil {
			fmt.Println(err.Error(), ", using a blank db for now.")
			Global = BotState{make(map[string]*BotUser), "&"}
		}
	} else {
		fmt.Println(err.Error(), ", using a blank db for now.")
		Global = BotState{make(map[string]*BotUser), "&"}
	}
}

func SaveGlobal() {
	f, err := os.Create("waifus.json")
	if err == nil {
		dec := json.NewEncoder(f)
		if err = dec.Encode(&Global); err != nil {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}
