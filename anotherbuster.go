package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func message(message string) {
	if string(message[0]) == "*" {
		fmt.Println("\033[30m[\033[36m*\033[30m]\033[0m", message[1:], "\033[0m")
	} else if string(message[0]) == "x" {
		fmt.Println("\033[30m[\033[31mx\033[30m]\033[0m", message[1:], "\033[0m")
		os.Exit(1)
	} else {
		fmt.Println("\033[30m[\033[32m"+string(message[0])+"\033[30m]\033[0m", message[1:], "\033[0m")
	}
}

func flagString(flag string) (result string) {
	for i, a := range os.Args[1:] {
		if a == flag {
			return os.Args[i+2]
		}
	}
	return ""
}

func flagBool(flag string) (result bool) {
	for _, a := range os.Args[1:] {
		if a == flag {
			return true
		}
	}
	return false
}

func flagInt(flag string) (result int) {
	for i, a := range os.Args[1:] {
		if a == flag {
			temp, _ := strconv.Atoi(os.Args[i+2])
			return temp
		}
	}
	return -1
}

func format_url(url string) (result string) {
	result = url
	if !(strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://")) {
		result = "https://" + result
	}
	if !strings.HasSuffix(url, "/") {
		result = result + "/"
	}
	return
}

func request(url string) *http.Response {
	if url == "" {
		message("xNo url given, try with -u {target_url} or --url {target_url}")
	}
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Connection", "Keep-Alive")
	request.Header.Set("Accept-Language", "en-US")
	request.Header.Set("User-Agent", "Mozilla/5.0")
	response, err := client.Do(request)
	if err != nil {
		message("xError while reaching " + url + " \033[30m(" + err.Error() + ")")
	}
	defer response.Body.Close()
	return response
}

func get_content(wordlist string) []string {
	if wordlist == "" {
		message("xNo wordlist given !")
	}
	file, err := os.Open(wordlist)
	if err != nil {
		message("xFile \"" + wordlist + "\" not found !")
	}
	scanner := bufio.NewScanner(file)
	result := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		result = append(result, line)
	}
	return result
}

func main() {
	message("*\033[36mAnother\033[0m Buster")

	url := flagString("-u")
	wordlist := flagString("-w")
	suffix := flagString("-s")
	target_length := flagInt("-l")
	target_response := flagInt("-r")
	verbose := flagBool("-v")

	// Format the url and try to reach it
	url = format_url(url)
	_ = request(url)

	// Get wordlist and the length
	content := get_content(wordlist)
	size := len(content)

	// Set default target
	if target_length == -1 && target_response == -1 {
		target_response = 200
	}

	message("?Target : " + url)
	message("?Wordlist : " + wordlist)
	message("?\033[0mStart scanning")

	// Make a request for each word in wordlist
	for i := 0; i < size; i++ {
		fmt.Printf("\r\033[30m[\033[32m?\033[30m]\033[0m Try\033[30m (%d/%d)\033[0m", i, size)
		response := request(url + content[i] + "/")
		if verbose {
			fmt.Println("\n\033[1A\033[K\r\033[30m-\033[0m "+content[i]+suffix, "\033[30m", response.StatusCode, http.StatusText(response.StatusCode), " Length:", response.ContentLength)
		} else if target_response == response.StatusCode {
			fmt.Println("\n\033[1A\033[K\r\033[30m-\033[0m "+content[i]+suffix, "\033[30m", response.StatusCode, http.StatusText(response.StatusCode))
		} else if target_length != -1 && int64(target_length) == response.ContentLength {
			fmt.Println("\n\033[1A\033[K\r\033[30m-\033[0m "+content[i]+suffix, "\033[30m", response.StatusCode, response.ContentLength)
		}
	}
}
