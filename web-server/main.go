package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	flag.Parse()

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/events", eventsHandler)

	http.ListenAndServe(":5050", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")

	token := []string{
		"I", "am", "a", "Web3", "Developer", "at", "CodeCrunch", "Techlabs", "Pvt.", "Ltd.",
		"focusing", "on", "building", "decentralized", "apps", "(DApps)", "like", "blockchain-based",
		"certifications,", "NFT", "marketplaces,", "the", "metaverse,", "and", "games.",
		"Over", "the", "past", "three", "years,", "I’ve", "worked", "with", "Ethereum", "on",
		"different", "projects,", "gaining", "hands-on", "experience", "in", "blockchain", "technology.",

		"I’ve", "mentored", "developers", "at", "HackFS", "2023", "and", "2024", "(ETHGlobal),",
		"helping", "them", "solve", "problems", "and", "improve", "their", "projects.",
		"I’ve", "also", "hosted", "events,", "such", "as", "a", "Filecoin", "workshop", "and",
		"an", "Ethereum", "session,", "where", "I", "explained", "blockchain", "concepts",
		"in", "an", "easy-to-understand", "way.", "Through", "CodeCrunch", "Tech", "School,",
		"I", "create", "videos", "in", "Hindi", "to", "make", "technology", "simple",
		"for", "regional", "learners.",

		"As", "part", "of", "the", "FVM", "Early", "Builders", "Program,", "I", "helped",
		"build", "Trustified,", "a", "platform", "for", "issuing", "digital", "credentials.",
		"My", "team", "and", "I", "have", "worked", "on", "many", "Web3", "projects",
		"and", "created", "educational", "content", "to", "teach", "Web3.",
		"This", "will", "be", "my", "first", "time", "contributing", "to", "an",
		"open-source", "project", "in", "the", "Protocol", "Labs", "ecosystem.",
	}

	for _, token := range token {
		content := fmt.Sprintf("%s ", string(token))
		w.Write([]byte(content))
		w.(http.Flusher).Flush()
		time.Sleep(time.Millisecond * 120)

	}
}
