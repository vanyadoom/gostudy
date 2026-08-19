package wallet

import (
	"encoding/json"
	"errors"
	"os"
	"time"
)

type Transaction struct {
	Type   string
	Amount float64
	Date   time.Time
}

type Wallet struct {
	filename string
	Owner    string
	Balance  float64
	History  []Transaction
}

func NewWallet(owner string, startBalance float64, filename string) *Wallet {
	return &Wallet{
		filename: filename,
		Owner:    owner,
		Balance:  startBalance,
		History:  []Transaction{},
	}
}

func (w *Wallet) Save() error {
	data, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(w.filename, data, 0644)
}

func (w *Wallet) Load() error {
	if _, err := os.Stat(w.filename); errors.Is(err, os.ErrNotExist) {
		return w.Save()
	}

	data, err := os.ReadFile(w.filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, w)
}

func (w *Wallet) Deposit(amount float64) {
	w.Balance += amount

	tx := Transaction{
		Type:   "Пополнение",
		Amount: amount,
		Date:   time.Now(),
	}

	w.History = append(w.History, tx)

	_ = w.Save()
}

func (w *Wallet) Withdraw(amount float64) bool {
	if amount > w.Balance {
		return false
	}
	w.Balance -= amount

	tx := Transaction{
		Type:   "Снятие",
		Amount: amount,
		Date:   time.Now(),
	}
	w.History = append(w.History, tx)

	_ = w.Save()
	return true
}
