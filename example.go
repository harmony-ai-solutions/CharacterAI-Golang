// Package main
/*
Copyright © 2023 Harmony AI Solutions & Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"bufio"
	"fmt"
	"github.com/harmony-ai-solutions/CharacterAI-Golang/cai"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Initial params - Usage of env vars recommended
	//token := ""
	//character := ""
	//isPlus := false
	token := os.Getenv("CAI_TOKEN")
	character := os.Getenv("CAI_CHAR")
	isPlus, errParse := strconv.ParseBool(os.Getenv("CAI_PLUS"))
	if errParse != nil {
		isPlus = false
	}

	// Create client
	caiClient, errClient := cai.NewGoCAI(token, isPlus)
	if errClient != nil {
		fmt.Println(fmt.Errorf("unable to create client, error: %q", errClient))
		os.Exit(1)
	}
	// Get chat data
	chatData, errChat := caiClient.Chat.GetChat(character)
	if errChat != nil {
		if strings.Contains(errChat.Error(), "404") {
			// Chat does not exist yet, create new
			chatData, errChat = caiClient.Chat.NewChat(character)
			if errChat != nil {
				fmt.Println(fmt.Errorf("unable to create chat, error: %q", errChat))
				os.Exit(3)
			}
		} else {
			fmt.Println(fmt.Errorf("unable to fetch chat data, error: %q", errChat))
			os.Exit(2)
		}
	}
	// Find AI paricipant in chat
	participantData, okParticipants := chatData["participants"]
	if !okParticipants {
		fmt.Println(fmt.Errorf("chat data doesn't hold participant data"))
		os.Exit(3)
	}
	participantList := participantData.([]interface{})
	var aiParticipantName string
	for _, participantRaw := range participantList {
		participant := participantRaw.(map[string]interface{})
		isHuman, okIsHuman := participant["is_human"]
		if okIsHuman {
			isHumanBool := isHuman.(bool)
			if !isHumanBool {
				var userData map[string]interface{}
				if userDataRaw, okUserData := participant["user"]; !okUserData {
					fmt.Println(fmt.Errorf("unable to parse user data"))
					os.Exit(4)
				} else {
					userData = userDataRaw.(map[string]interface{})
				}
				if userNameRaw, okUserName := userData["username"]; !okUserName {
					fmt.Println(fmt.Errorf("unable to parse user name"))
					os.Exit(5)
				} else {
					aiParticipantName = userNameRaw.(string)
					break
				}
			}
		}
	}

	for true {
		fmt.Println("Enter user message: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		errInput := scanner.Err()
		if errInput != nil {
			fmt.Println(fmt.Errorf("unable to scan user input. Error: %v", errInput))
			os.Exit(6)
		}
		// Send
		messageResult, errMessage := caiClient.Chat.SendMessage(chatData["external_id"].(string), aiParticipantName, scanner.Text(), nil)
		if errMessage != nil {
			fmt.Println(fmt.Errorf("unable to send message. Error: %v", errMessage))
			os.Exit(7)
		}
		// Handle result
		if replyDataRaw, okReplyData := messageResult["replies"]; okReplyData {
			replyData := replyDataRaw.([]interface{})
			firstReplyRaw := replyData[0]
			firstReply := firstReplyRaw.(map[string]interface{})
			fmt.Println(fmt.Sprintf("%v: %v", aiParticipantName, firstReply["text"]))
			fmt.Println()
		}
	}

}
