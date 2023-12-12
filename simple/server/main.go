package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

const difficulty = 4  // difficult for PoW
const word = "naruto" // base word

func main() {

	//create a Listener
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println("listening error:", err.Error())
		return
	}
	defer listener.Close()

	fmt.Println("server is working...")

	//endless loop for handling connections
	for {

		//try to create one connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return
		}

		//if its ok go for hadling the connection in own thread
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("New client added:", conn.RemoteAddr())

	//create objects for reading and writing
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	//write welcome message with w - object
	writer.WriteString("You are connected!!!!: \n")
	writer.WriteString("Enter a nonce to check you access. For requesting challenge task just press Enter")
	writer.Flush()

	//endless loop for getting and sending messages with a client
	for {

		writer.WriteString(" \n")
		writer.Flush()

		//read answer from client
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error data reading:", err.Error())

			return
		}
		input = strings.TrimSpace(input)
		fmt.Println("Message from client ", conn.RemoteAddr(), " - ", input)

		if input == "" {
			writer.WriteString(fmt.Sprintf("challenge:" + word + " " + strconv.Itoa(difficulty) + "\n"))
			writer.Flush()

			continue
		} else {
			//figuring out if he had done work
			isRight := nonceIsRight(word, difficulty, input)

			//if hashes are not the same
			if !isRight {
				writer.WriteString(fmt.Sprintf("access denied: try again. For requesting challenge task just press Enter "))
				writer.Flush()

				continue
			}

			//show the truth
			writer.WriteString(fmt.Sprintf("access granted: "))
			writer.WriteString(getQuote() + "\n")
			writer.Flush()
		}
	}
}

// for checking the PoW
func nonceIsRight(word string, difficulty int, nonce string) bool {

	hash := sha256.Sum256([]byte(fmt.Sprintf("%s%s", word, nonce)))
	hashStr := hex.EncodeToString(hash[:])

	if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
		return true
	}

	return false
}

// getting quotes from quotes api service
func getQuote() string {
	url := "https://zenquotes.io/?api=quotes"
	response, err := http.Get(url)
	if err != nil {
		fmt.Println("error qoutes reading:", err.Error())

		return ""
	}
	defer response.Body.Close()
	data, _ := io.ReadAll(response.Body)
	var quotes []interface{}
	err = json.Unmarshal(data, &quotes)
	if err != nil {
		fmt.Println("error quotes parsing:", err.Error())

		return ""
	}
	firstQuote := quotes[0].(map[string]interface{})

	return firstQuote["q"].(string)
}
