package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/keybase/go-keybase-chat-bot/kbchat"
)

func fail(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(3)
}

func main() {
	var kbLoc string
	var kbc *kbchat.API
	var err error

	flag.StringVar(&kbLoc, "keybase", "keybase", "the location of the Keybase app")
	flag.Parse()

	if kbc, err = kbchat.Start(kbLoc); err != nil {
		fail("Error creating API: %s", err.Error())
	}

	team_name := "<team_name>"
	start_string := "test pls ignore"
	fmt.Printf("saying hello on conversation: %s\n", team_name)

	if err = kbc.SendMessageByTeamName(team_name, start_string, nil); err != nil {
		fail("Error sending message; %s", err.Error())
	}

	messages, err := kbc.ReadMessagesByTeamName(team_name, nil)

	if err != nil {
		fail("Error reading messages; %s", err.Error())
	}

	messageId := 4337

	for i := range messages {
		m := messages[i].M
		fmt.Printf("%d - %d: %s\n", i, m.ID, m.Content.Text.Body)
		if m.Content.Text.Body == start_string {
			messageId = m.ID
		}
	}

	fmt.Printf("Message ID is %d\n", messageId)

	for {
		max := 30
		for i := 0; i < max; i++ {
			newmsg := strings.Repeat("-", i) + "#" + strings.Repeat("-", max-i)
			if err = kbc.EditMessageByTeamName(team_name, messageId, newmsg, nil); err != nil {
				fail("Error editing message: %s", err.Error())
			}

			time.Sleep(10 * time.Millisecond)
		}
	}
}
