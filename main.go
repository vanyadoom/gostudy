package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Wallet struct {
	Owner   string
	Balance float64
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

func main() {

	myWallet := Wallet{
		Owner:   "Иван",
		Balance: 100.0,
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\n🪙 КРИПТО-КОШЕЛЕК [%s] | Баланс: $%.2f\n", myWallet.Owner, myWallet.Balance)
		fmt.Println("1. Пополнить кошелёк")
		fmt.Println("2. Снять деньги")
		fmt.Println("3. Выход")
		fmt.Print("Выберите действие: ")

		scanner.Scan()

		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Введите сумму для пополнения: ")
			scanner.Scan()

			amount, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

			if amount > 0 {
				myWallet.Deposit(amount)
				fmt.Println("✅ Кошелек успешно пополнен!")
			} else {
				fmt.Println("❌ Ошибка: Сумма должна быть больше нуля")
			}

		case "2":
			fmt.Print("Введите сумму для снятия: ")
			scanner.Scan()

			amount, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

			if amount > 0 {

				if myWallet.Withdraw(amount) {
					fmt.Println("✅ Деньги успешно сняты!")
				} else {
					fmt.Println("❌ Ошибка: Недостаточно средств на балансе!")
				}
			} else {
				fmt.Println("❌ Ошибка: Сумма должна быть больше нуля")
			}

		case "3":
			fmt.Println("👋 До встречи в Wallet CLI!")
			return

		default:
			fmt.Println("❌ Неверный выбор! Попробуйте снова.")

		}
	}
}
