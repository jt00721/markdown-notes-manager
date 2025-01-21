package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
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

func displayMenu() {
	fmt.Println("\nMarkdown Note Manager")
	fmt.Println("=====================")
	fmt.Println("1. Create a new note")
	fmt.Println("2. View a note")
	fmt.Println("3. Edit a note")
	fmt.Println("4. List all notes")
	fmt.Println("5. Search notes")
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
			fmt.Println("Edit a note...")
		case "4":
			fmt.Println("List all notes...")
		case "5":
			fmt.Println("Searching for note...")
			return
		case "6":
			fmt.Println("Exiting")
			return
		default:
			fmt.Println("Invalid choice. Please select a valid option (1-6).")
		}
	}

	// openInBrowser(readFile("notes/My_Note.md"))
}
