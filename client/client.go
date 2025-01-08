package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

//O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar.
//O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL
//e em seguida deverá retornar no formato JSON o resultado para o cliente.
//O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON).
// o client.go terá um timeout máximo de 300ms para receber o resultado do server.go.

func RequestDollarPriceBRL() {
	var bid string
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to make request: %v", err)
		return
	}
	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body) // Corrected to read from resp.Body
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading response: %v", err)
		return
	}

	err = json.Unmarshal(res, &bid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing response: %v", err)
		return
	}

	err = saveDollarValueToFile(bid)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create file: %v", err)
		return
	}
}

func saveDollarValueToFile(cotacao string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Printf("Error creating file")
		return err
	}
	_, err = f.Write([]byte(cotacao))

	f.Close()
	return err
}
