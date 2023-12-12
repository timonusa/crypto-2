package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println("Error of creating connection:", err.Error())
		return
	}
	defer conn.Close()

	go readServerResponses(conn)

	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(conn)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading the input:", err.Error())
			return
		}

		writer.WriteString(input)
		writer.Flush()
	}
}

func readServerResponses(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("error getting the response:", err.Error())
			return
		}

		//if get challange from server
		if strings.Contains(response, "challenge:") {
			response = strings.Replace(response, "challenge:", "", -1)
			parts := strings.Split(response, " ")
			word := parts[0]
			difficulty, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			nonce := calculatePoW(word, difficulty)
			writer.WriteString(strconv.Itoa(nonce))
			writer.Flush()
			fmt.Println("Got challenge to solve. Press enter")
		} else {
			fmt.Println(response)
		}

	}
}

func calculatePoW(data string, difficulty int) int {
	nonce := 0
	for {
		hash := sha256.Sum256([]byte(fmt.Sprintf("%s%d", data, nonce)))
		hashStr := hex.EncodeToString(hash[:])

		if strings.HasPrefix(hashStr, strings.Repeat("0", difficulty)) {
			return nonce
		}

		nonce++
	}
}
