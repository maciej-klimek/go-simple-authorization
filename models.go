package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type LoginData struct {
	PasswordHash string
	SessionToken string
	CSRFToken    string
}

var Users = map[string]LoginData{}
var Filename string = "userData.json"

func loadUserData() error {
	_, err := os.Stat(Filename)
	if os.IsNotExist(err) {
		fmt.Println("> File", Filename, " does not exist, creating a new one.")
		return nil
	}

	fileContent, err := os.ReadFile(Filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(fileContent, &Users)
	if err != nil {
		return err
	}

	fmt.Println("> User data loaded from", Filename)
	return nil
}

func saveUserData() error {
	file, err := os.Create(Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(Users)
	if err != nil {
		return err
	}

	fmt.Println("> User data saved to", Filename)
	return nil
}
