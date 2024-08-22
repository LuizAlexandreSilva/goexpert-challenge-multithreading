package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	cep := "69058411"
	channelViaCEP := make(chan any)
	channelBrasilAPI := make(chan any)

	go requestCEPToAPI("viacep", channelViaCEP, cep)
	go requestCEPToAPI("brasilapi", channelBrasilAPI, cep)

	select {
	case msg := <-channelViaCEP:
		fmt.Printf("Received from ViaCEP: %s\n", msg)
	case msg := <-channelBrasilAPI:
		fmt.Printf("Received from BrasilAPI: %s\n", msg)
	case <-time.After(time.Second):
		println("Time exceeded.")
	}
}

func requestCEPToAPI(api string, c chan<- any, cep string) {
	var route string
	if api == "viacep" {
		route = fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	} else if api == "brasilapi" {
		route = fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s/", cep)
	}

	req, err := http.Get(route)
	if err != nil {
		print(err)
	}
	res, err := io.ReadAll(req.Body)
	c <- res
	req.Body.Close()
}
