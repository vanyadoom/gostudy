package wallet

type Wallet struct {
	Owner   string
	Balance float64
}

func NewWallet(owner string, startBalance float64) *Wallet {
	return &Wallet{
		Owner:   owner,
		Balance: startBalance,
	}
}

func (w *Wallet) Deposit(amount float64) {
	w.Balance += amount
}

func (w *Wallet) Withdraw(amount float64) bool {
	if amount > w.Balance {
		return false
	}
	w.Balance -= amount
	return true
}
