package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"image/color"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/ebitengine/oto/v3"
	"github.com/hajimehoshi/go-mp3"
)

// Fixed username and password constants
const (
	adminUsername = "hasting"
	adminPassword = "holidays"
)

// MessageData contains the list of messages to be shown to the user.
type MessageData struct {
	Messages []Message
}

// Message represents a user message with associated color and cutoff time.
type Message struct {
	Text   string
	Colour color.RGBA
	CutOff time.Time
}

// Initial message data
var messages = MessageData{
	Messages: []Message{
		{Text: "Hello, World!", Colour: color.RGBA{R: 255, G: 165, B: 0}},
		{Text: "Welcome to the messenger.", Colour: color.RGBA{R: 130, G: 165, B: 0, A: 255}},
		{Text: "Note the date at the top.", Colour: color.RGBA{R: 100, G: 15, B: 0, A: 255}},
		{Text: "Have a great day!.", Colour: color.RGBA{R: 255, G: 165, B: 100, A: 255}},
	},
}

// main is the entry point for the program.
func main() {
	// Load messages from JSON file if it exists, otherwise save default messages.
	if fileExists("messages.json") {
		logF("Loading messages.")
		loadMessagesFromJSONfile()
	} else {
		logF("Loading default messages.")
		logFMessages()
		saveMessagesToJSONfile()
	}
	// Define HTTP routes
	http.HandleFunc("/", handler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/update-messages", updateMessagesHandler)

	// SSL/TLS paths can be provided via flags or environment variables
	certFlag := flag.String("cert", "", "path to TLS certificate")
	keyFlag := flag.String("key", "", "path to TLS key")
	flag.Parse()

	certPath := *certFlag
	if certPath == "" {
		certPath = os.Getenv("CERT_PATH")
	}
	if certPath == "" {
		logF("CERT_PATH not set. Using default cert.pem")
		certPath = "cert.pem"
	}

	keyPath := *keyFlag
	if keyPath == "" {
		keyPath = os.Getenv("KEY_PATH")
	}
	if keyPath == "" {
		logF("KEY_PATH not set. Using default key.pem")
		keyPath = "key.pem"
	}

	// Start the HTTP server
	logF("Starting Web server on port 443.")
	log.Fatal(http.ListenAndServeTLS(":443", certPath, keyPath, nil))
}

// HexColour returns the color of a message in hexadecimal format.
func (m Message) HexColour() string {
	return fmt.Sprintf("#%02X%02X%02X", m.Colour.R, m.Colour.G, m.Colour.B)
}

// checkAdminCredentials verifies if the provided BasicAuth credentials are valid.
func checkAdminCredentials(r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	logF("Authentication for '" + username + "' = " + fmt.Sprint(ok && username == adminUsername && password == adminPassword))
	return ok && username == adminUsername && password == adminPassword
}

// updateMessagesHandler handles updating the message content via POST requests.
func updateMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		logF("Invalid request method.")
		return
	}
	// Decode the incoming JSON data
	var data struct {
		Messages []struct {
			Text  string `json:"text"`
			Color string `json:"color"`
		} `json:"messages"`
	}
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusInternalServerError)
		logF("Failed to decode request body")
		return
	}
	// Update the global messages with the new texts and colors
	logF("New Messages are :-")
	for i, message := range data.Messages {
		if i < len(messages.Messages) {
			messages.Messages[i].Text = message.Text
			logF(message.Text)
			if color, err := parseHexColor(message.Color); err == nil {
				messages.Messages[i].Colour = color
			}
		}
	}
	// Save updated messages to JSON file
	saveMessagesToJSONfile()
	// Send a success response
	response := struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}
	json.NewEncoder(w).Encode(response)
}

// parseHexColor converts a hex color string to a color.RGBA object.
func parseHexColor(s string) (color.RGBA, error) {
	c, err := strconv.ParseUint(strings.TrimPrefix(s, "#"), 16, 32)
	if err != nil {
		return color.RGBA{}, err
	}
	return color.RGBA{
		R: uint8(c >> 16),
		G: uint8(c >> 8 & 0xFF),
		B: uint8(c & 0xFF),
		A: 255,
	}, nil
}

// fileExists checks if a file with the given name exists.
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// handler serves the main page and logs the request.
func handler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("normalLayout.html")
	if err != nil {
		logF(err)
	}
	logFMessages()
	tmpl.Execute(w, messages)
	logRequestDetails(r)
}

// adminHandler serves the admin page and logs the request.
func adminHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAdminCredentials(r) {
		w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
		http.Error(w, "Unauthorized.", http.StatusUnauthorized)
		return
	}
	tmpl, err := template.ParseFiles("adminLayout.html")
	if err != nil {
		logF(err)
	}
	loadMessagesFromJSONfile()
	tmpl.Execute(w, messages)
	logRequestDetails(r)
}

// loadMessagesFromJSONfile loads the messages from a JSON file.
func loadMessagesFromJSONfile() {
	fileBytes, err := os.ReadFile("messages.json")
	if err != nil {
		logF("No JSON file found.")
		return
	}
	err = json.Unmarshal(fileBytes, &messages)
	if err != nil {
		logF("Error unmarshalling:", err)
	}
	logF("JSON file loaded.")
}

// saveMessagesToJSONfile saves the messages to a JSON file.
func saveMessagesToJSONfile() {
	if err := PlayMP3("alarm05.mp3"); err != nil {
		logF(err)
	}
	jsonBytes, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		logF("Error marshalling JSON file:", err)
		return
	}
	err = os.WriteFile("messages.json", jsonBytes, 0644)
	if err != nil {
		logF("Error writing JSON file:", err)
		return
	}
	logF("Messages saved to JSON file:-")
	logFMessages()
}

// logF logs a message with the calling function's name.
func logF(v ...interface{}) {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		funcName := runtime.FuncForPC(pc).Name()
		log.Println(fmt.Sprintf("[%s]", funcName), v)
	} else {
		log.Println(v)
	}
}

// logRequestDetails logs the details of the HTTP request.
func logRequestDetails(r *http.Request) {
	logF("--------------")
	logF(r.URL.Path)
	ip := getIP(r)
	logF("IP address:", ip)
}

// getIP extracts the IP address from the request.
func getIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr // Use the whole address if splitting fails
	}
	return host
}

var otoCtx *oto.Context

func getOtoContext() (*oto.Context, error) {
	if otoCtx != nil {
		return otoCtx, nil
	}

	op := &oto.NewContextOptions{
		SampleRate:   44100,                   // Usually 44100 or 48000.
		ChannelCount: 2,                       // 1 is mono sound, and 2 is stereo.
		Format:       oto.FormatSignedInt16LE, // Format of the source. go-mp3's format is signed 16bit integers.
	}

	var err error
	otoCtx, _, err = oto.NewContext(op)
	return otoCtx, err
}

func logFMessages() {
	for i, _ := range messages.Messages {
		m := messages.Messages[i].Text
		logF(m)
	}
}

// PlayMP3 plays the given MP3 file using streaming
func PlayMP3(filePath string) error {
	// Get the oto context
	ctx, err := getOtoContext()
	if err != nil {
		return err
	}

	// Open the file for reading. Do NOT close before you finish playing!
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode file. This process is done as the file plays, so it won't
	// load the whole thing into memory.
	decodedMp3, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}

	// Create a new 'player' using the context
	player := ctx.NewPlayer(decodedMp3)

	// Play starts playing the sound and returns without waiting for it (Play() is async).
	player.Play()

	// We can wait for the sound to finish playing using something like this
	for player.IsPlaying() {
		time.Sleep(time.Millisecond)
	}

	// Close the player after you're done with it
	if err = player.Close(); err != nil {
		return err
	}

	return nil
}
