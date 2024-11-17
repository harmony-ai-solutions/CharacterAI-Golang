// Package main
/*
Copyright © 2023-2024 Harmony AI Solutions & Contributors

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
	"strings"
)

func main() {
	// Retrieve the token and character ID from environment variables
	token := os.Getenv("CAI_TOKEN")
	webNextAuth := os.Getenv("CAI_WEBNEXTAUTH")
	proxyURL := os.Getenv("CAI_PROXY")
	characterID := os.Getenv("CAI_CHAR")

	if token == "" || characterID == "" {
		fmt.Println("Error: CAI_TOKEN or CAI_CHAR environment variable is not set.")
		os.Exit(1)
	}

	// Create a new client instance
	client := cai.NewClient(token, webNextAuth, proxyURL)
	err := client.Authenticate()
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(2)
	}

	// Fetch existing chats with the character
	chats, err := client.FetchChats(characterID, 0)
	if err != nil {
		fmt.Printf("Error fetching chats: %v\n", err)
		os.Exit(3)
	}

	var chat *cai.Chat

	if len(chats) > 0 {
		// Use the most recent chat with the character
		chat = chats[0]
		fmt.Printf("Using existing chat with ID: %s\n", chat.ChatID)
	} else {
		// Create a new chat with the character
		chat, _, err = client.CreateChat(characterID, true)
		if err != nil {
			fmt.Printf("Error creating chat: %v\n", err)
			os.Exit(4)
		}
		fmt.Printf("Created new chat with ID: %s\n", chat.ChatID)
	}

	// Print the previous messages in the chat (up to 5)
	messages, _, err := client.FetchMessages(chat.ChatID, false, "")
	if err != nil {
		fmt.Printf("Error fetching messages: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Previous messages (up to 5):")
	if len(messages) > 5 {
		messages = messages[len(messages)-5:]
	}
	for _, turn := range messages {
		var authorName string
		if turn.Author.IsHuman {
			authorName = "You"
		} else {
			authorName = turn.Author.Name
		}
		candidate := turn.Candidates[turn.PrimaryCandidateID]
		fmt.Printf("%s: %s\n", authorName, candidate.Text)
	}
	fmt.Println()

	// Start the interaction loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You: ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading user input: %v\n", err)
			os.Exit(1)
		}
		userInput = strings.TrimSpace(userInput)

		// Send the user's message to the character
		turn, err := client.SendMessage(characterID, chat.ChatID, userInput)
		if err != nil {
			fmt.Printf("Error sending message: %v\n", err)
			os.Exit(1)
		}

		// Retrieve the AI's response
		aiResponse := ""
		if turn != nil && len(turn.Candidates) > 0 {
			primaryCandidate := turn.Candidates[turn.PrimaryCandidateID]
			aiResponse = primaryCandidate.Text
		} else {
			fmt.Println("No response received from the AI.")
			continue
		}

		fmt.Printf("%s: %s\n", turn.Author.Name, aiResponse)
		fmt.Println()
	}
}
