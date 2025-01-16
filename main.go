package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func getNoteInput() (string, string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter note title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	if title == "" {
		fmt.Println("Title cannot be empty. Please try again.")
		return getNoteInput()
	}

	fmt.Println("Enter note content (type 'END' on a new line to finish):")
	var content strings.Builder
	for {
		line, _ := reader.ReadString('\n')
		if strings.TrimSpace(line) == "END" {
			break
		}
		content.WriteString(line)
	}

	return title, content.String()
}

func saveNoteToFile(filename, content string) {
	ensureNotesDir()
	file, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Printf("Note successfully saved at %s\n", filename)
}

func sanitizeTitle(title string) string {
	return strings.ReplaceAll(strings.Map(func(r rune) rune {
		if strings.ContainsRune(" abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-", r) {
			return r
		}
		return -1
	}, title), " ", "_")
}

func ensureNotesDir() {
	err := os.MkdirAll("notes", os.ModePerm)
	if err != nil {
		fmt.Println("Error creating notes directory:", err)
	}
}

func getUniqueFilename(title string) string {
	base := sanitizeTitle(title)
	filename := fmt.Sprintf("notes/%s.md", base)

	counter := 1
	for {
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			break
		}
		filename = fmt.Sprintf("notes/%s_%d.md", base, counter)
		counter++
	}

	return filename
}

func main() {
	title, content := getNoteInput()
	filename := getUniqueFilename(title)
	saveNoteToFile(filename, content)
}
