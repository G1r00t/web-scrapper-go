package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.zenrows.com/v1/?apikey=ed801bac7b7a137c57b99151979ab7e58a0427a6&url=https%3A%2F%2Fwww.g2.com%2Fproducts%2Fasana%2Freviews&js_render=true&premium_proxy=true", nil)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Create or open the file to save the response body
	file, err := os.Create("response_body.txt")
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

		log.Println("Response body saved to response_body.txt")
	}
