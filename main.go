package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
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

func readFile(filename string) string {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("The file '%s' does not exist.\n", filename)
		return ""
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return ""
	}

	return string(data)
}

func displayNoteContent(title, content string) {
	fmt.Printf("---- %s ----\n\n%s\n\n-----------------\n", title, content)
}

func viewNote(title string) {
	filename := fmt.Sprintf("notes/%s.md", sanitizeTitle(title))
	content := readFile(filename)
	if content == "" {
		fmt.Println("No content to display")
		return
	}

	displayNoteContent(title, content)
}

func renderMarkdown(content string) string {
	var buf bytes.Buffer
	md := goldmark.New()
	if err := md.Convert([]byte(content), &buf); err != nil {
		fmt.Println("Error rendering Markdown:", err)
		return content
	}

	return buf.String()
}

func openInBrowser(content string) {
	file, err := os.Create("preview.html")
	if err != nil {
		fmt.Println("Error creating HTML file:", err)
		return
	}
	defer file.Close()

	markdownContent := renderMarkdown(content)

	_, err = file.WriteString(markdownContent)
	if err != nil {
		fmt.Println("Error writing HTML content:", err)
		return
	}

	exec.Command("open", "preview.html").Start()
}

func editNoteInline(title string) {
	filename := fmt.Sprintf("notes/%s.md", sanitizeTitle(title))

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("The note '%s' does not exist.\n", title)
		return
	}

	// Read and display the current content
	currentContent, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading the note:", err)
		return
	}
	fmt.Println("Current Content:")
	fmt.Println(string(currentContent))

	// Collect updated content
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter new content for the note (type 'END' on a new line to finish):")
	var updatedContent strings.Builder
	for {
		line, _ := reader.ReadString('\n')
		if strings.TrimSpace(line) == "END" {
			break
		}
		updatedContent.WriteString(line)
	}

	if strings.TrimSpace(updatedContent.String()) == "" {
		fmt.Println("Cannot save an empty note. Changes discarded.")
		return
	}

	backupNote(filename)

	// Save the updated content
	err = os.WriteFile(filename, []byte(updatedContent.String()), 0644)
	if err != nil {
		fmt.Println("Error saving the note:", err)
		return
	}
	fmt.Printf("The note '%s' has been updated successfully.\n", title)
}

func editNoteWithEditor(title string) {
	filename := fmt.Sprintf("notes/%s.md", sanitizeTitle(title))

	// Check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("The note '%s' does not exist.\n", title)
		return
	}

	backupNote(filename)

	// Launch the default editor
	cmd := exec.Command("nano", filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error opening the editor:", err)
		return
	}
	fmt.Printf("The note '%s' has been updated.\n", title)
}

func backupNote(filename string) {
	backupFilename := fmt.Sprintf("%s.bak", filename)
	input, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error creating backup:", err)
		return
	}
	os.WriteFile(backupFilename, input, 0644)
}

func listNotes() {
	files, err := os.ReadDir("notes")
	if err != nil {
		fmt.Println("Error reading notes directory:", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No notes found.")
		return
	}

	fmt.Println("\nYour Notes:")
	for _, file := range files {
		title := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		fmt.Println("- " + title)
	}
}

func searchByTitle(query string) {
	files, err := os.ReadDir("notes")
	if err != nil {
		fmt.Println("Error reading notes directory:", err)
		return
	}

	fmt.Printf("Searching for notes with titles containing '%s':\n", query)
	found := false
	for _, file := range files {
		title := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
		if strings.Contains(strings.ToLower(title), strings.ToLower(query)) {
			fmt.Println("- " + title)
			found = true
		}
	}

	if !found {
		fmt.Println("No notes found with the specified title.")
	}
}

func searchByContent(query string) {
	files, err := os.ReadDir("notes")
	if err != nil {
		fmt.Println("Error reading notes directory:", err)
		return
	}

	fmt.Printf("Searching for notes with content containing '%s':\n", query)
	found := false
	for _, file := range files {
		content, err := os.ReadFile("notes/" + file.Name())
		if err != nil {
			fmt.Println("Error reading file:", file.Name())
			continue
		}

		if strings.Contains(strings.ToLower(string(content)), strings.ToLower(query)) {
			title := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			fmt.Println("- " + title)
			found = true
		}
	}

	if !found {
		fmt.Println("No notes found with the specified content.")
	}
}

func displayMenu() {
	fmt.Println("\nMarkdown Note Manager")
	fmt.Println("=====================")
	fmt.Println("1. Create a new note")
	fmt.Println("2. View a note")
	fmt.Println("3. Edit a note")
	fmt.Println("4. List all notes")
	fmt.Println("5. Search notes")
	fmt.Println("6. Exit")
	fmt.Print("\nSelect an option (1-6): ")
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		displayMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			title, content := getNoteInput()
			filename := getUniqueFilename(title)
			saveNoteToFile(filename, content)
		case "2":
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter note title to read: ")
			title, _ := reader.ReadString('\n')
			title = strings.TrimSpace(title)

			viewNote(title)
		case "3":
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("1. Edit note inline")
			fmt.Println("2. Edit note with editor")
			fmt.Print("\nSelect an option (1-2): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter note title to edit inline: ")
				title, _ := reader.ReadString('\n')
				title = strings.TrimSpace(title)

				editNoteInline(title)
			case "2":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter note title to edit with editor: ")
				title, _ := reader.ReadString('\n')
				title = strings.TrimSpace(title)

				editNoteWithEditor(title)
			default:
				fmt.Println("Invalid choice. Please select a valid option (1-2).")
			}
		case "4":
			listNotes()
		case "5":
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("1. Search for note title")
			fmt.Println("2. Search for note content")
			fmt.Print("\nSelect an option (1-2): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "1":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter note title to search: ")
				title, _ := reader.ReadString('\n')
				title = strings.TrimSpace(title)

				searchByTitle(title)
			case "2":
				reader := bufio.NewReader(os.Stdin)

				fmt.Print("Enter note content to search: ")
				query, _ := reader.ReadString('\n')
				query = strings.TrimSpace(query)

				searchByContent(query)
			default:
				fmt.Println("Invalid choice. Please select a valid option (1-2).")
			}
		case "6":
			fmt.Println("Exiting")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-6).")
		}
	}

	// openInBrowser(readFile("notes/My_Note.md"))
}
