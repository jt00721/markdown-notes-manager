package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type NoteEntry struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type NoteFile struct {
	File []NoteEntry `json:"file"`
}

func getNoteInputForJSON(data *NoteFile) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter note title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)
	if title == "" {
		fmt.Println("Title cannot be empty. Please try again.")
		return
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

	entry := NoteEntry{
		Title:   title,
		Content: content.String(),
	}

	data.File = append(data.File, entry)
	fmt.Printf("Logged Note %s\n", title)
}

func loadData(filename string) (NoteFile, error) {
	var data NoteFile

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No existing data file found. Starting with default settings.")
			return NoteFile{}, nil
		}
		return data, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return data, fmt.Errorf("failed to decode data: %v", err)
	}

	fmt.Println("Data loaded successfully.")
	return data, nil
}

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

func saveNoteToJsonFile(data NoteFile, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode data: %v", err)
	}

	fmt.Println("Data saved successfully.")
	return nil
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
	ensureNotesDir()

	const filename = "notes.json"

	data, err := loadData(filename)
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		return
	}

	getNoteInputForJSON(&data)

	err = saveNoteToJsonFile(data, filename)
	if err != nil {
		fmt.Printf("Error saving data: %v\n", err)
	}

	// title, content := getNoteInput()

	// filename := getUniqueFilename(title)
	// saveNoteToFile(filename, content)
}
