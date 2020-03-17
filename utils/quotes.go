package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	quotesFile      string = "quotes-full.csv"
	quotesFileShort string = "quotes-short.csv"
	url             string = "https://www.inc.com/kevin-daum/35-quotes-on-ending-things-that-will-make-2016-more-exciting.html"
)

func download() {
	quotes := make([]string, 30)
	c := http.Client{}
	r, err := c.Get(url)
	if err != nil || r.StatusCode != http.StatusOK {
		panic(err.Error())
	}
	doc, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		panic(err.Error())
	}
	reStart := regexp.MustCompile(`^[0-9]+\. \"`)
	reEnd := regexp.MustCompile(`\"$`)
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		if len(s.Text()) == 0 {
			return
		}
		q1 := reStart.ReplaceAllString(s.Text(), "")
		q2 := reEnd.ReplaceAllString(q1, "")
		quotes = append(quotes, q2)
	})

	f, err := os.Create(quotesFile)
	defer f.Close()
	if err != nil {
		panic(err.Error())
	}
	f.WriteString("quote" + "\t" + "auhtor" + "\n")
	for i := 33; i < len(quotes); i += 2 {
		f.WriteString(quotes[i-1] + "\t" + quotes[i] + "\n")
	}
}

func takeFirst() {
	qfShort, err := os.Create(quotesFileShort)
	if err != nil {
		panic(err.Error())
	}
	qf, err := os.Open(quotesFile)
	if err != nil {
		panic(err.Error())
	}
	defer qf.Close()
	scanner := bufio.NewScanner(qf)
	scanner.Scan()
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ".")
		if len(s) < 1 {
			panic(errors.New("Strange data"))
		}
		qfShort.WriteString(s[0] + "\n")
	}

}

func main() {
	for i := 0; i < 10; i++ {
		fmt.Println(getRandLine())
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
