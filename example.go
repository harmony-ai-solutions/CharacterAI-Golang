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
	var aiParticipant *cai.ChatParticipant
	for _, participant := range chatData.Participants {
		if !participant.IsHuman {
			aiParticipant = participant
			break
		}
	}

	// Get History for chatroom
	if chatHistoryData, errHistory := caiClient.Chat.GetHistory(chatData.ExternalID); errHistory != nil {
		fmt.Println(fmt.Errorf("unable to parse user data"))
		os.Exit(4)
	} else {
		fmt.Println("Previous messages (up to 5):")
		messageHistory := chatHistoryData.Messages
		if len(chatHistoryData.Messages) > 5 {
			messageHistory = messageHistory[len(messageHistory)-5:]
		}
		for _, message := range messageHistory {
			fmt.Println(fmt.Sprintf("%v: %v", message.SourceName, message.Text))
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
		messageResult, errMessage := caiClient.Chat.SendMessage(chatData.ExternalID, aiParticipant.User.Username, scanner.Text(), nil)
		if errMessage != nil {
			fmt.Println(fmt.Errorf("unable to send message. Error: %v", errMessage))
			os.Exit(7)
		}
		// Handle result
		if len(messageResult.Replies) > 0 {
			firstReply := messageResult.Replies[0]
			fmt.Println(fmt.Sprintf("%v: %v", aiParticipant.Name, firstReply.Text))
			fmt.Println()
		}
	}

}
