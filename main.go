package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CEP struct {
	Cep          string
	Street       string
	City         string
	State        string
	Neighborhood string
}

type ViaCepResponse struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type BrasilAPIResponse struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

func doRequest(url string) []byte {
	res, _ := http.Get(url)

	body := res.Body
	defer res.Body.Close()

	data, _ := io.ReadAll(body)

	return data
}

func main() {
	cep := os.Args[1]

	brasilAPI := make(chan CEP)
	viacep := make(chan CEP)

	go func() {
		url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
		data := doRequest(url)

		var response ViaCepResponse
		_ = json.Unmarshal(data, &response)

		cep := CEP{
			Cep:          response.Cep,
			Street:       response.Logradouro,
			City:         response.Localidade,
			State:        response.Uf,
			Neighborhood: response.Bairro,
		}

		viacep <- cep
	}()

	go func() {
		url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
		data := doRequest(url)

		var response BrasilAPIResponse
		_ = json.Unmarshal(data, &response)

		cep := CEP{
			Cep:          response.Cep,
			Street:       response.Street,
			City:         response.City,
			State:        response.State,
			Neighborhood: response.Neighborhood,
		}

		brasilAPI <- cep
	}()

	select {
	case res := <-brasilAPI:
		fmt.Println("BrasilAPI responded to request first!")
		fmt.Println(res)
	case res := <-viacep:
		fmt.Println("ViaCEP responded to request first!")
		fmt.Println(res)
	case <-time.After(time.Second):
		fmt.Println("Timeout!")
	}
}
