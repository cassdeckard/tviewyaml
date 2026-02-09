package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/cassdeckard/tviewyaml"
	"github.com/cassdeckard/tviewyaml/template"
)

var (
	// Mock data for search filtering
	fruits = []string{
		"Apple", "Apricot", "Banana", "Blackberry", "Blueberry",
		"Cherry", "Cranberry", "Grape", "Grapefruit", "Kiwi",
		"Lemon", "Lime", "Mango", "Orange", "Peach",
		"Pear", "Pineapple", "Plum", "Raspberry", "Strawberry",
	}
	
	// Email validation regex
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
)

// RegisterInputFieldLive adds live input field demo functions as custom template functions
func RegisterInputFieldLive(b *tviewyaml.AppBuilder) *tviewyaml.AppBuilder {
	maxZero := 0
	b.WithTemplateFunction("updateCharCount", 0, &maxZero, nil, updateCharCount)
	b.WithTemplateFunction("validateEmail", 0, &maxZero, nil, validateEmail)
	b.WithTemplateFunction("searchFilter", 0, &maxZero, nil, searchFilter)
	return b
}

func updateCharCount(ctx *template.Context) {
	// Get the current input text from state
	text, _ := ctx.GetState("__inputText")
	textStr, ok := text.(string)
	if !ok {
		textStr = ""
	}
	
	// Calculate character count
	count := len(textStr)
	maxCount := 50
	
	// Create status message with color coding
	var status string
	if count == 0 {
		status = "[gray]0 / 50 characters[::-]"
	} else if count < 40 {
		status = fmt.Sprintf("[green]%d / %d characters[::-]", count, maxCount)
	} else if count < maxCount {
		status = fmt.Sprintf("[yellow]%d / %d characters (approaching limit)[::-]", count, maxCount)
	} else {
		status = fmt.Sprintf("[red]%d / %d characters (at limit!)[::-]", count, maxCount)
	}
	
	// Update state
	ctx.SetStateDirect("charCount", status)
}

func validateEmail(ctx *template.Context) {
	// Get the current input text from state
	text, _ := ctx.GetState("__inputText")
	textStr, ok := text.(string)
	if !ok {
		textStr = ""
	}
	
	// Validate email format
	var status string
	if textStr == "" {
		status = "[gray]Enter an email address[::-]"
	} else if emailRegex.MatchString(textStr) {
		status = "[green]✓ Valid email format[::-]"
	} else {
		status = "[red]✗ Invalid email format[::-]"
	}
	
	// Update state
	ctx.SetStateDirect("emailStatus", status)
}

func searchFilter(ctx *template.Context) {
	// Get the current input text from state
	text, _ := ctx.GetState("__inputText")
	searchQuery, ok := text.(string)
	if !ok {
		searchQuery = ""
	}
	
	// Filter fruits based on search query
	var results []string
	if searchQuery == "" {
		// Show all fruits when empty
		results = fruits
	} else {
		// Case-insensitive search
		query := strings.ToLower(searchQuery)
		for _, fruit := range fruits {
			if strings.Contains(strings.ToLower(fruit), query) {
				results = append(results, fruit)
			}
		}
	}
	
	// Format results
	var status string
	if len(results) == 0 {
		status = "[red]No matches found[::-]"
	} else if searchQuery == "" {
		status = fmt.Sprintf("[cyan]%d fruits available[::-]", len(results))
	} else {
		// Show first 5 results
		displayCount := 5
		if len(results) < displayCount {
			displayCount = len(results)
		}
		
		resultList := make([]string, displayCount)
		for i := 0; i < displayCount; i++ {
			resultList[i] = results[i]
		}
		
		moreText := ""
		if len(results) > displayCount {
			moreText = fmt.Sprintf(" (+%d more)", len(results)-displayCount)
		}
		
		status = fmt.Sprintf("[green]%d matches:%s[::-]\n%s", 
			len(results), moreText, strings.Join(resultList, ", "))
	}
	
	// Update state
	ctx.SetStateDirect("searchResults", status)
}
