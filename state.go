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

func InitComforts() {
	f, err := os.Open("comforts.json")
	if err == nil {
		dec := json.NewDecoder(f)
		if err = dec.Decode(&Comforts); err != nil {
			fmt.Println(err.Error(), ", using minimal comforts db for now.")
			Comforts = []string{"_%wn hugs %n_"}
		}
	} else {
		fmt.Println(err.Error(), ", using minimal comforts db for now.")
		Comforts = []string{"_%wn hugs %n_"}
	}
}

func SaveGlobal() {
	f, err := os.Create("waifus.json")
	if err == nil {
		defer f.Close()
		data, err := json.MarshalIndent(&Global, "", "\t")
		if err != nil {
			fmt.Println(err.Error())
		} else {
			f.Write(data)
		}
	} else {
		fmt.Println(err.Error())
	}
}
