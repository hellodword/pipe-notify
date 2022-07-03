package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	VERBOSITY_NONE = iota
	VERBOSITY_OTHER
	VERBOSITY_ALL
)

var notified = map[string]string{}

func main() {
	verbosity := flag.Int("verbosity", VERBOSITY_NONE, "")
	pattern := flag.String("pattern", "", "")
	index := flag.Int("index", 0, "index")

	flag.Parse()

	_ = *index
	var r *regexp.Regexp
	if *pattern != "" {
		r = regexp.MustCompile(*pattern)
	}

	if *index <= 0 {
		*index = 0
	}

	if r.NumSubexp() < *index {
		panic("NumSubexp")
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		handleLine(*verbosity, r, *index, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func handleLine(verbosity int, r *regexp.Regexp, index int, line string) {
	var match string
	if r != nil {
		matches := r.FindStringSubmatch(line)
		if len(matches) > 0 {
			match = matches[index]
		}
	}

	if match != "" {
		notifierTelegramBot(match)
	}

	if verbosity <= VERBOSITY_NONE {
		return
	}

	if verbosity >= VERBOSITY_ALL {
		fmt.Println(line)
	}

	if r != nil {
		fmt.Println(r.ReplaceAllString(line, ""))
	}

}

func notifierTelegramBot(s string) {
	_, ok := notified[s]
	notified[s] = s
	if !ok && strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")) != "" && strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID")) != "" {
		http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN"))),
			"application/json",
			strings.NewReader(fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "disable_notification": false}`, strings.TrimSpace(os.Getenv("TELEGRAM_CHAT_ID")), s)))
	}
	return
}
