package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"gostudy/notifier"
)

type Transaction struct {
	Type     string
	Currency string
	Amount   float64
	Date     time.Time
}

type Wallet struct {
	filename string
	notifier notifier.Notifier
	Owner    string
	Balances map[string]float64
	History  []Transaction
}

func NewWallet(owner string, filename string, n notifier.Notifier) *Wallet {
	return &Wallet{
		filename: filename,
		notifier: n,
		Owner:    owner,
		Balances: make(map[string]float64), // Выделяем память под карту
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
		w.Balances["USD"] = 100.0
		w.Balances["EUR"] = 50.0
		w.Balances["RUB"] = 5000.0
		return w.Save()
	}

	data, err := os.ReadFile(w.filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, w)
}

func (w *Wallet) Deposit(currency string, amount float64) {
	w.Balances[currency] += amount

	tx := Transaction{
		Type:     "Пополнение",
		Currency: currency,
		Amount:   amount,
		Date:     time.Now(),
	}

	w.History = append(w.History, tx)

	_ = w.Save()

	msg := fmt.Sprintf("Баланс пополнен на %.2f %s. Новый баланс: %.2f %s", amount, currency, w.Balances[currency], currency)
	w.notifier.Send(msg)
}

func (w *Wallet) Withdraw(currency string, amount float64) bool {
	currentBalance := w.Balances[currency]
	if amount > currentBalance {
		return false
	}

	w.Balances[currency] -= amount

	tx := Transaction{
		Type:     "Снятие",
		Currency: currency,
		Amount:   amount,
		Date:     time.Now(),
	}
	w.History = append(w.History, tx)
	_ = w.Save()

	msg := fmt.Sprintf("Со счёта снято %.2f %s. Новый баланс: %.2f %s", amount, currency, w.Balances[currency], currency)
	w.notifier.Send(msg)

	return true
}
