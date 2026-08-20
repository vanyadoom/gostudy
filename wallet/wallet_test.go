package wallet

import (
	"os"
	"testing"
)

type MockNotifier struct{}

func (m MockNotifier) Send(message string) {}

func TestWalletOperations(t *testing.T) {
	testFile := "test_wallet.json"

	_ = os.Remove(testFile)

	defer os.Remove(testFile)

	mockBot := MockNotifier{}
	w := NewWallet("Тестер", testFile, mockBot)
	_ = w.Load()

	err := w.Deposit("USD", 50.0)
	if err != nil {
		t.Errorf("Deposit вернул непредвиденную ошибку: %v", err)
	}

	if w.Balances["USD"] != 150.0 {

		t.Errorf("Ошибка Deposit: ожидали баланс $150.00, но получили $%.2f", w.Balances["USD"])
	}

	err = w.Withdraw("USD", 40.0)
	if err != nil {
		t.Errorf("Withdraw вернул ошибку при валидном снятии: %v", err)
	}

	if w.Balances["USD"] != 110.0 {
		t.Errorf("Ошибка Withdraw: ожидали остаток $110.00, но получили $%.2f", w.Balances["USD"])
	}

	err = w.Withdraw("USD", 500.0)
	if err != nil && err != ErrInsufficientFunds {

		t.Errorf("Ожидали ошибку '%v', но получили '%v'", ErrInsufficientFunds, err)
	}

	if err == nil {

		t.Errorf("Ошибка защиты: кошелек позволил снять больше, чем есть на балансе, и не вернул ошибку")
	}
}
