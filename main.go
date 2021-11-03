package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Error(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Breach struct {
	Success bool   `json:"success"`
	Found   string `json:"found"`
	Result  []struct {
		Password_Founded bool     `json:"has_password"`
		Password         string   `json:"password"`
		SHA1             string   `json:"sha1"`
		Sources          []string `json:"sources"`
	}
}

func main() {
	var (
		email   string
		API_Key string
		result  Breach
	)
	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"

	flag.StringVar(&email, "e", "", "The email address")
	flag.StringVar(&API_Key, "k", "", "The API key")
	flag.Parse()
	if email == "" || API_Key == "" {
		fmt.Println("If you don't have an API key, please visit https://rapidapi.com/ to create one.")
		fmt.Println("Please provide the email and API key")
		flag.PrintDefaults()
		return
	}
	url := "https://breachdirectory.p.rapidapi.com/?func=auto&term=" + email
	req, err := http.NewRequest("GET", url, nil)
	Error(err)
	req.Header.Add("x-rapidapi-host", "breachdirectory.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", string(API_Key))
	res, err := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	Error(err)
	json_response := string(body)
	json.Unmarshal([]byte(json_response), &result)
	for _, v := range result.Result {
		if v.Password_Founded == true {
			fmt.Println(string(colorGreen), "Password found: "+v.Password)
			for i := range v.Sources {
				fmt.Println(string(colorGreen), "From this leak: "+v.Sources[i])
			}
		} else if v.Password_Founded == false || len(v.Sources) != 0 {
			fmt.Println(string(colorRed), "no password in clear text but it is part of this leak", string(colorReset))
			for i := range v.Sources {
				fmt.Println(string(colorGreen), "Leak Source: "+v.Sources[i])
			}
		} else {
			fmt.Println(string(colorRed), "no password found", string(colorReset))
		}
	}
}
