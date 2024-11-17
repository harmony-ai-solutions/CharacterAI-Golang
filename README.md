# 💬 CharacterAI - Golang Port by Project Harmony.AI

![Tag](https://img.shields.io/github/license/harmony-ai-solutions/CharacterAI-Golang)

An unofficial API Client for [CharacterAI](https://character.ai/), written in Golang, ported over from Python.

Original Python source code by [Xtr4F](https://github.com/Xtr4F) and supporters in this repo: https://github.com/Xtr4F/PyCharacterAI

---

⚠️ ATTENTION - Unofficial community repository! ⚠️

This is an unofficial library which has no relation to the CharacterAI development team. 
 
CharacterAI has no official api and all breakpoints were found manually using reverse engineering.
The authors are not responsible for possible consequences of using this library.

Documentation may be incomplete or missing. This repo is not optimized for productive usage in golang applications yet.
Use at your own risk.

You have questions, need help, or just want to show your support? Reach us
here: [Discord Server & Patreon page](#how-to-reach-out-to-us).

### TODO's:

- [x] Port over API from source repo
  - [x] Confirm basic functionality
- [x] Golang QOL improvements
  - [x] Create Wrapper Structs for API Endpoints + Parse them within the API methods
  - [x] Add proper WebSocket client for V2 / Websocket API
- [ ] Documentation & Testing
  - [x] Tests for main chat functions
  - [ ] Write tests for all API Methods => Not all methods have tests yet, but most.
  - [ ] Documentation for Endpoints & Data Types

## 💻 Installation

```bash
go get github.com/harmony-ai-solutions/CharacterAI-Golang
```

## 📚 Documentation

Detailed documentation and API-Docs TBD

## 📙 Example

Example code for a simple, functional Chat app. The code can also be found in [example.go](example.go)

```Golang
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
```

---

## About Project Harmony.AI

![Project Harmony.AI](docs/images/Harmony-Main-Banner-200px.png)

### Our goal: Elevating Human <-to-> AI Interaction beyond known boundaries.
Project Harmony.AI emerged from the idea to allow for a seamless living together between AI-driven characters and humans.
Since it became obvious that a lot of technologies required for achieving this goal are not existing or still very experimental,
the long term vision of Project Harmony is to establish the full set of technologies which help minimizing biological and
technological barriers in Human <-to-> AI Interaction.

### Our principles: Fair use and accessibility

We want to counter today's tendencies of AI development centralization at the hands of big
corporations. We're pushing towards maximum transparency in our own development efforts, and aim for our software to be
accessible and usable in the most democratic ways possible.

Therefore, for all our current and future software offerings, we'll perform a constant and well-educated evaluation whether
we can safely open source them in parts or even completely, as long as this appears to be non-harmful towards achieving
the project's main goal.

Also, we're constantly striving to keep our software offerings as accessible as possible when it comes to services which
cannot be run or managed by everyone - For example our Harmony Speech TTS Engine. As long as this project exists,
we'll be trying out utmost to provide free tiers for personal and public research use of our software and APIs.

However, at the same time we'll also ensure everyone who supports us or actively joins forces with us on our journey, gets
something proper back in turn. Therefore we're also maintaining a Patreon Page with different supporter tiers, as we are
open towards collaboration with other businesses.

### How to reach out to us

#### If you want to collaborate or support this Project financially:

Feel free to join our Discord Server and / or subscribe to our Patreon - Even $1 helps us drive this project forward.

![Harmony.AI Discord Server](docs/images/discord32.png) [Harmony.AI Discord Server](https://discord.gg/f6RQyhNPX8)

![Harmony.AI Discord Server](docs/images/patreon32.png) [Harmony.AI Patreon](https://patreon.com/harmony_ai)

#### If you want to use our software commercially or discuss a business or development partnership:

Contact us directly via: [contact@project-harmony.ai](mailto:contact@project-harmony.ai)

---
&copy; 2023 Harmony AI Solutions & Contributors

Licensed under the Apache 2.0 License