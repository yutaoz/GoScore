package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func formatJSON(data []byte) string {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", " ")

	if err != nil {
		fmt.Println(err)
	}

	d := out.Bytes()
	return string(d)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	baseurl := "https://www.basketball-reference.com/leagues/NBA_1984_games-november.html"
	request, error := http.NewRequest("GET", baseurl, nil)

	if error != nil {
		fmt.Println(error)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}

	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}
	fmt.Println(string(responseBody))
	formattedData := formatJSON(responseBody)
	fmt.Println("Status: ", response.Status)
	fmt.Println("Response body: ", formattedData)

	f, err := os.Create("tmp/dat2.txt")
	check(err)
	defer f.Close()
	n, err := f.WriteString(string(responseBody))
	check(err)
	fmt.Println("wrote " + fmt.Sprint(n) + "bytes")

	// clean up memory after execution
	defer response.Body.Close()

}
