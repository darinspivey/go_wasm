package main

import (
	"encoding/json"
	"fmt"
	"regexp"
)

type Payload struct {
	Input     []OriginalEvent `json:"input"`
	Summaries []Summary       `json:"summaries"`
}

type Summary struct {
	LogType string   `json:"log_type"`
	Logs    []string `json:"logs"`
	Counts  Count    `json:"counts"`
}

type OriginalEvent struct {
	Line string `json:"line"`
}

type Count struct {
	Errors int `json:"errors"`
	Info   int `json:"info"`
}

var syslogRegexp = regexp.MustCompile(`^<(\d+)>`)

//export process_event
func process_event(text string) {
	if len(text) != 0 {
		// fmt.Println("============= input ================", text)
		payload := Payload{}
		err := json.Unmarshal([]byte(text), &payload)
		if err != nil {
			fmt.Printf("Error unmarshalling event: %+v\n", err)
			return
		}
		output := summarizeLines(&payload)
		output_str, err := json.Marshal(output)
		if err == nil {
			fmt.Println(string(output_str))
		}
	}
}

func summarizeLines(incoming *Payload) Payload {
	unmatchedLines := make([]OriginalEvent, 0)
	matchedLines := make([]string, 0)
	errorCount := 0

	for _, originalEvent := range incoming.Input {
		line := originalEvent.Line
		matches := syslogRegexp.FindStringSubmatch(line)
		if len(matches) == 0 {
			unmatchedLines = append(unmatchedLines, originalEvent)
			continue
		}
		// priority := matches[1]
		errorCount++
		matchedLines = append(matchedLines, line)
	}

	summary := Summary{
		LogType: "syslog",
		Logs:    matchedLines,
		Counts: Count{
			Errors: errorCount,
		},
	}

	output := Payload{
		Input:     unmatchedLines,
		Summaries: append(incoming.Summaries, summary),
	}

	return output
}

func main() {
}
