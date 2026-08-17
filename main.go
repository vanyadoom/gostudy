package main

import "fmt"

type Wallet struct {
	Owner   string
	Balance float64
}

func (w *Wallet) Deposit(amount float64) {
	w.Balance += amount
}

func main() {
	myWallet := Wallet{
		Owner:   "Иван",
		Balance: 100.0,
	}

	fmt.Printf("Стартовый кошелек %s: $%.2f\n", myWallet.Owner, myWallet.Balance)

	myWallet.Deposit(50.5)

	fmt.Printf("Кошелек %s после пополнения: $%.2f\n", myWallet.Owner, myWallet.Balance)

}
