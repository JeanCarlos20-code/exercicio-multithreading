package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ViaCEP struct {
	Cep        string `json:"cep"`
	Estado     string `json:"estado"`
	Cidade     string `json:"localidade"`
	Bairro     string `json:"bairro"`
	Logradouro string `json:"logradouro"`
}

type BrasilCEP struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	channelViaCepApi, channelCEPBrasilApi := make(chan map[string]string), make(chan map[string]string)
	for _, cep := range os.Args[1:] {
		go func() {
			channelViaCepApi <- ViaCEPApiHandler(ctx, cep)
		}()

		go func() {
			channelCEPBrasilApi <- CEPBrasilApiHandler(ctx, cep)
		}()
	}

	select {
	case messageChannelViaCepApi := <-channelViaCepApi:
		fmt.Println(messageChannelViaCepApi)

	case messageChannelCEPBrasilApi := <-channelCEPBrasilApi:
		fmt.Println(messageChannelCEPBrasilApi)

	case <-ctx.Done():
		fmt.Println("Timeout reached")
	}
}

func TransformCEP(ctx context.Context, api string) *http.Response {
	req, err := http.NewRequestWithContext(ctx, "GET", api, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	return res
}

func CEPBrasilApiHandler(ctx context.Context, cep string) map[string]string {
	res := TransformCEP(ctx, "https://brasilapi.com.br/api/cep/v1/"+cep)
	defer res.Body.Close()

	var data BrasilCEP

	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Println("CEPBrasilApiHandler")

	return map[string]string{
		"Cep":        data.Cep,
		"Estado":     data.State,
		"Cidade":     data.City,
		"Bairro":     data.Neighborhood,
		"Logradouro": data.Street,
	}
}

func ViaCEPApiHandler(ctx context.Context, cep string) map[string]string {
	res := TransformCEP(ctx, "https://viacep.com.br/ws/"+cep+"/json/")
	defer res.Body.Close()

	var data ViaCEP

	err := json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	fmt.Println("ViaCEPApiHandler")

	return map[string]string{
		"Cep":        data.Cep,
		"Estado":     data.Estado,
		"Cidade":     data.Cidade,
		"Bairro":     data.Bairro,
		"Logradouro": data.Logradouro,
	}
}
