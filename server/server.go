package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Currency struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

func Handler() {
	http.ListenAndServe(":8080", nil)
	http.HandleFunc("/cotacao", GetDollarBidHandler)
}

func GetDollarBidHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	currency, err := getCurrency("USD", "BRL")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bid, err := getCurrencyBid(currency)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(bid)
}

func getCurrency(startingCurrency string, convCurrency string) (*Currency, error) {
	var currency Currency
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/"+startingCurrency+"-"+convCurrency, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed getting currency: %v", err)
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed reading response: %v", err)
		return nil, err
	}

	err = json.Unmarshal(res, &currency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed parsing response: %v", err)
		return nil, err
	}

	err = saveCurrencyToDB(&currency)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving currency to db: %v", err)
		return nil, err
	}

	return &currency, err
}

func getCurrencyBid(currency *Currency) (string, error) {
	if currency.Bid == "" {
		return "", errors.New("Error: currency's bid is empty")
	}

	return currency.Bid, nil
}

func saveCurrencyToDB(currency *Currency) error {
	db, err := initDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize database: %v", err)
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	insertQuery := `
    INSERT INTO currency (code, codein, name, high, low, varBid, pctChange, bid, ask, timestamp, create_date)
    VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	_, err = db.ExecContext(ctx, insertQuery, currency.Code, currency.Codein, currency.Name, currency.High, currency.Low, currency.VarBid, currency.PctChange, currency.Bid, currency.Ask, currency.Timestamp, currency.CreateDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to insert currency data: %v", err)
		return err
	}

	return nil
}

func initDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./currency.db")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
    CREATE TABLE IF NOT EXISTS currency (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        code TEXT,
        codein TEXT,
        name TEXT,
        high TEXT,
        low TEXT,
        varBid TEXT,
        pctChange TEXT,
        bid TEXT,
        ask TEXT,
        timestamp TEXT,
        create_date TEXT
    );`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return db, nil
}
