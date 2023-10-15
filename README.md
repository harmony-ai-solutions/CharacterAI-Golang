# 💬 CharacterAI - Golang Port by Project Harmony.AI
![Tag](https://img.shields.io/github/license/harmony-ai-solutions/CharacterAI-Golang)

An unofficial API Client for CharacterAI, written in Golang, ported over from Python.

Original Python source code in this repo: https://github.com/kramcat/CharacterAI

Original Readme kept for reference: [Original Readme](README.old.md)

---

## ⚠️ ATTENTION: Pre-Release!
This repo is currently very barebone and not optimized for productive usage in golang applications yet.
Any support with testing, verifying functionality and adding proper golang struct handling is heavily appreciated.

#### You have questions, need help, or just want to show your support? Reach us here:

[Discord Server & Patreon page](#how-to-reach-out-to-us).

## 💻 Installation
```bash
go get github.com/harmony-ai-solutions/CharacterAI-Golang
```

## 📚 Documentation
Detailed documentation and Apidocs coming soon

## 📙 Example
Example code for a simple Chat app can be found in [example.go](example.go)
```Golang
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
		//... 
	}
	
	// ... Very Explicit Handling neeeded currently, see example.go for details 
	
	for true {
		fmt.Println("Enter user message: ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		errInput := scanner.Err()
		if errInput != nil {
			//...
		}
		// Send
		messageResult, errMessage := caiClient.Chat.SendMessage(chatData["external_id"].(string), aiParticipantName, scanner.Text(), nil)
		if errMessage != nil {
			//...
		}
		// Handle result
		if replyDataRaw, okReplyData := messageResult["replies"]; okReplyData {
			//...
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