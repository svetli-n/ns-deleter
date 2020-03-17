package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
)

func send(rcpt, text string) {

	log.Println(text)

	/*
		fmt.Println(text)
		fmt.Println("Email sent")
		return
	*/

	// Connect to the remote SMTP server.
	c, err := smtp.Dial("smtp.hm.com:25")
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient first
	if err := c.Mail("noreply@hm.com"); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(rcpt); err != nil {
		log.Fatal(err)
	}
	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}

	from := fmt.Sprintf("From: <%s>\r\n", "noreply@hm.com")
	to := fmt.Sprintf("To: <%s>\r\n", rcpt)
	subject := "Subject: " + "Namespace deletion -- " + getRandLine() + "\r\n"
	body := text + "\r\n"
	msg := from + to + subject + "\r\n" + body

	_, err = fmt.Fprintf(wc, msg)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

func getRandLine() string {
	quotes := []string{
		"It is always important to know when something has reached its end",
		"There is no real ending",
		"A man is like a novel: until the very last page you don't know how it will end",
		"Ends are not bad things, they just mean that something else is about to begin",
		"There's a trick to the 'graceful exit",
		"Build your life on your dreams; because dreams never have bad endings",
		"We lose the precious sense that an end is only a beginning in disguise",
		"Every ending is a beginning",
		"An end is only a beginning in disguise",
		"A sunset is nothing more and nothing less than the backside of a sunrise",
		"It all begins and ends in your mind",
		"I try to avoid stories that end with",
		"There should have been a better farewell",
		"There is always more after the ending",
		"A moment's beginning ends in a moment",
		"The great miraculous bell of translucent ice is suspended in mid-air",
		"Endings to be useful must be inconclusive",
		"It's a lot easier to say when something ended rather than when it began",
		"The feeling is less like an ending than just another starting point",
		"I may regret the way we ended, but I will never regret what we had",
		"It's much easier to not know things sometimes",
		"Everything has to come to an end, sometime",
	}
	return quotes[rand.Intn(len(quotes))]

}
