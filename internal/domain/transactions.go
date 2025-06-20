package domain

// Transaction is the canonical representation of a user cash movement.
type Transaction struct {
	ID        string  `json:"id"`
	UserID    string  `json:"user_id"`
	Amount    float64 `json:"amount"`
	Type      string  `json:"type"` // "credit" or "debit"
	Timestamp int64   `json:"timestamp"`
}

// Equal compares the business‑relevant fields.
func (t Transaction) Equal(other Transaction) bool {
	return t.ID == other.ID &&
		t.UserID == other.UserID &&
		t.Amount == other.Amount &&
		t.Type == other.Type &&
		t.Timestamp == (other.Timestamp)
}
