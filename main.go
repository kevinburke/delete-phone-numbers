package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	twilio "github.com/saintpete/twilio-go"
)

func main() {
	flag.Parse()
	sid := os.Getenv("TWILIO_ACCOUNT_SID")
	token := os.Getenv("TWILIO_AUTH_TOKEN")
	client := twilio.NewClient(sid, token, nil)
	ctx := context.Background()
	iter := client.IncomingNumbers.GetPageIterator(nil)
	for {
		page, err := iter.Next(ctx)
		if err == twilio.NoMoreResults {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, number := range page.IncomingPhoneNumbers {
			fmt.Printf(`Sid: %s
Number: %s
Voice URL: %s
SMS URL: %s

Delete this number? `, number.Sid, number.PhoneNumber.Friendly(), number.VoiceURL, number.SMSURL)
			reader := bufio.NewReader(os.Stdin)
			text, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			text = strings.ToLower(strings.TrimSpace(text))
			if text == "y" || text == "yes" {
				innerCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
				err := client.IncomingNumbers.Release(innerCtx, number.Sid)
				if err != nil {
					log.Fatal(err)
				}
				cancel()
				fmt.Printf("Successfully deleted %s\n\n", number.PhoneNumber)
			}
		}
	}
}
