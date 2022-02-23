package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Whois struct {
	Ip      string
	Region  string
	City    string
	Country string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + " [IP ADDRESS]")
		return
	}
	resp, err := http.Get("http://ipwhois.app/json/" + os.Args[1])
	if err != nil {
		fmt.Println("Error cannot get ipwhois.com.. are you connected to internet ?")
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	var whois Whois
	json.Unmarshal(bodyBytes, &whois)
	fmt.Println("IP: " + whois.Ip + " is located in " + whois.Country + " more precisely in " + whois.City)
}
