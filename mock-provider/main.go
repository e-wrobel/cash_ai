package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"` // "credit" or "debit"
	Timestamp int64   `json:"timestamp"`
}

// generateMockTransactions returns a list of mock transactions with some inconsistencies
//
//nolint:staticcheck
func generateMockTransactions() []Transaction {
	staticTransactions := []Transaction{
		{ID: "txn-12345", UserID: "user-1", Amount: 100.50, Type: "credit", Timestamp: time.Now().Add(-10 * time.Minute).Unix()},
		{ID: "txn-67890", UserID: "user-2", Amount: 75.25, Type: "debit", Timestamp: time.Now().Add(-5 * time.Minute).Unix()},
	}
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(5) + 5 // Generate between 5-10 transactions
	transactions := make([]Transaction, n)

	for i := 0; i < n; i++ {
		transactions[i] = Transaction{
			ID:        randomID(),
			UserID:    randomID(),
			Amount:    float64(rand.Intn(10000)) / 100.0,
			Type:      randomType(),
			Timestamp: time.Now().Add(time.Duration(-rand.Intn(1000)) * time.Second).Unix(),
		}
	}

	// Introduce occasional duplicates
	if rand.Intn(10) > 7 {
		transactions = append(transactions, transactions[rand.Intn(len(transactions))])
	}

	// Introduce missing fields
	if len(transactions) > 0 && rand.Intn(10) > 7 {
		transactions[0].UserID = ""
	}

	return append(transactions, staticTransactions...)
}

func randomID() string {
	return time.Now().Format("20060102150405") + fmt.Sprint(rune(rand.Intn(100)))
}

func randomType() string {
	types := []string{"credit", "debit"}
	return types[rand.Intn(len(types))]
}

//nolint:errcheck
func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	transactions := generateMockTransactions()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

// TransactionsClient provides methods to fetch transactions without calling HTTP manually
type TransactionsClient struct {
	Endpoint string
}

// NewTransactionsClient creates a new client instance
func NewTransactionsClient(endpoint string) *TransactionsClient {
	return &TransactionsClient{Endpoint: endpoint}
}

// GetTransactions fetches the list of transactions from the mock server
func (c *TransactionsClient) GetTransactions() ([]Transaction, error) {
	resp, err := http.Get(c.Endpoint + "/transactions")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch transactions")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var transactions []Transaction
	if err := json.Unmarshal(body, &transactions); err != nil {
		return nil, err
	}

	return transactions, nil
}

func main() {
	http.HandleFunc("/transactions", transactionsHandler)
	log.Println("Mock transaction server running on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
