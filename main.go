package main

import (
	"fmt"
	"strconv"
)

type Station struct {
	ID        string
	Message   string
	WalshCode []int
}

type Receiver struct {
	Stations []Station
}

func main() {
	codes := generateWalshCode(8)
	stations := []Station{
		{ID: "A", Message: "GOD", WalshCode: codes[0]},
		{ID: "B", Message: "CAT", WalshCode: codes[1]},
		{ID: "C", Message: "HAM", WalshCode: codes[2]},
		{ID: "D", Message: "SUN", WalshCode: codes[3]},
	}
	receiver := Receiver{Stations: stations}

	encodedMessages := make([][]int, len(stations))

	for i, station := range stations {
		fmt.Printf("Station %s sends message: %s \r\n", station.ID, station.Message)
		encodedMessages[i] = station.Broadcast()
		fmt.Printf("Station %s broadcasts message: %s \r\n", station.ID, printSlice(encodedMessages[i]))

	}

	broadcastedMessage := broadcast(encodedMessages)

	fmt.Printf("Full message broadcasted: %s \r\n", printSlice(broadcastedMessage))

	receiver.receive(broadcastedMessage)
}

func broadcast(encodedMessages [][]int) []int {
	broadcastedMessage := make([]int, len(encodedMessages[0]))
	for i := 0; i < len(encodedMessages); i++ {
		for j := 0; j < len(encodedMessages[i]); j++ {
			broadcastedMessage[j] += encodedMessages[i][j]
		}
	}

	return broadcastedMessage
}

func generateWalshCode(n int) [][]int {
	result := make([][]int, n)
	for i := 0; i < n; i++ {
		result[i] = make([]int, n)
	}

	result[0][0] = 1
	for k := 1; k < n; k += k {
		for i := 0; i < k; i++ {
			for j := 0; j < k; j++ {
				result[i+k][j] = result[i][j]
				result[i][j+k] = result[i][j]
				result[i+k][j+k] = -result[i][j]
			}
		}
	}

	return result
}

func convertBinaryToAscii(binaryMessage []int) string {
	asciiString := ""
	for i := 0; i < len(binaryMessage); i += 8 {
		var byteValue int
		for j := 0; j < 8; j++ {
			byteValue = (byteValue << 1) | binaryMessage[i+j]
		}
		asciiString += string(byte(byteValue))
	}
	return asciiString
}

func (s Station) Broadcast() []int {
	binaryMessage := s.broadcastBinaryMessage()
	return s.encodeBinary(binaryMessage)

}

func (s Station) binaryMessageLength() int {
	return len(s.broadcastBinaryMessage())
}

func (s Station) broadcastBinaryMessage() []int {
	res := ""
	var result []int
	for _, c := range s.Message {
		res = fmt.Sprintf("%s%.8b", res, c)
	}

	for _, code := range res {
		if int(code-'0') > 0 {
			result = append(result, 1)
		} else {
			result = append(result, -1)
		}
	}
	return result
}

func (s Station) encodeBinary(message []int) []int {
	var encoded []int
	for i := 0; i < len(message); i++ {
		for j := 0; j < len(s.WalshCode); j++ {
			encodedByte := s.WalshCode[j] * message[i]

			encoded = append(encoded, encodedByte)
		}
	}
	return encoded
}

func (r Receiver) receive(fullMessage []int) {
	for _, station := range r.Stations {
		decodedMessage := decodeForStation(station, fullMessage)
		ascii := convertBinaryToAscii(decodedMessage)
		fmt.Printf("Received message %s on station %s", ascii, station.ID)
		fmt.Println()
	}
}

func decodeForStation(station Station, fullMessage []int) []int {
	message := make([]int, station.binaryMessageLength())
	for i := 0; i < station.binaryMessageLength(); i++ {
		sum := 0
		for j := 0; j < len(station.WalshCode); j++ {
			sum += fullMessage[i*len(station.WalshCode)+j] * station.WalshCode[j]
		}

		if sum/len(station.WalshCode) > 0 {
			message[i] = 1
		} else {
			message[i] = 0
		}
	}
	return message
}

func printSlice(slice []int) string {
	result := ""
	for _, val := range slice {
		result += " " + strconv.Itoa(val)
	}
	return result
}
