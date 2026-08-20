package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gostudy/notifier"
	"gostudy/wallet"
)

func main() {
	tgBot := notifier.TelegramNotifier{Username: "vanya_crypto_bot"}
	myWallet := wallet.NewWallet("Иван", "wallet.json", tgBot)

	err := myWallet.Load()
	if err != nil {
		fmt.Println("❌ Критическая ошибка при загрузке кошелька:", err)
		return
	}

	// 🔥 МАГИЯ МНОГОПОТОЧНОСТИ GO (ГОРУТИНА):
	// Мы пишем слово `go` перед вызовом метода [INDEX].
	// Теперь этот робот улетает работать параллельно в фоновом режиме [INDEX]!
	// Он будет засыпать и просыпаться каждые 15 секунд, никак не мешая циклу меню ниже [INDEX].
	go myWallet.StartBackgroundAuditor()

	currencySymbols := map[string]string{
		"USD": "$",
		"EUR": "€",
		"RUB": "₽",
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf("\n🪙 КРИПТО-КОШЕЛЕК [%s]\n", myWallet.Owner)
		fmt.Println("💰 Текущие балансы:")

		for curr, bal := range myWallet.Balances {
			symbol := currencySymbols[curr]
			fmt.Printf("   • %s: %s%.2f\n", curr, symbol, bal)
		}

		fmt.Println("\nМеню:")
		fmt.Println("1. Пополнить кошелек")
		fmt.Println("2. Снять деньги")
		fmt.Println("3. Показать историю транзакций")
		fmt.Println("4. Выход")
		fmt.Print("Выберите действие: ")

		scanner.Scan()
		choice := strings.TrimSpace(scanner.Text())

		switch choice {
		case "1":
			fmt.Print("Введите валюту (USD, EUR, RUB): ")
			scanner.Scan()
			currency := strings.ToUpper(strings.TrimSpace(scanner.Text()))

			if currency != "USD" && currency != "EUR" && currency != "RUB" {
				fmt.Println("❌ Ошибка: Поддерживаются только валюты USD, EUR и RUB!")
				continue
			}

			fmt.Print("Введите сумму для пополнения: ")
			scanner.Scan()
			amount, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

			err := myWallet.Deposit(currency, amount)
			if err != nil {
				fmt.Printf("❌ Ошибка выполнения: %v\n", err)
			} else {
				fmt.Println("✅ Кошелек успешно пополнен!")
			}

		case "2":
			fmt.Print("Введите валюту (USD, EUR, RUB): ")
			scanner.Scan()
			currency := strings.ToUpper(strings.TrimSpace(scanner.Text()))

			if currency != "USD" && currency != "EUR" && currency != "RUB" {
				fmt.Println("❌ Ошибка: Поддерживаются только валюты USD, EUR и RUB!")
				continue
			}

			fmt.Print("Введите сумму для снятия: ")
			scanner.Scan()
			amount, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

			err := myWallet.Withdraw(currency, amount)
			if err != nil {
				fmt.Printf("❌ Ошибка выполнения: %v\n", err)
			} else {
				fmt.Println("✅ Деньги успешно сняты!")
			}

		case "3":
			fmt.Println("\n📜 ИСТОРИЯ ТРАНЗАКЦИЙ:")
			if len(myWallet.History) == 0 {
				fmt.Println("📭 История пуста. Вы еще не совершали операций.")
				continue
			}

			fmt.Println(strings.Repeat("-", 60))
			for _, tx := range myWallet.History {
				dateStr := tx.Date.Format("2006-01-02 15:04:05")
				symbol := currencySymbols[tx.Currency]
				fmt.Printf("[%s] %s: %s%.2f\n", dateStr, tx.Type, symbol, tx.Amount)
			}
			fmt.Println(strings.Repeat("-", 60))

		case "4":
			fmt.Println("👋 До встречи в Wallet CLI!")
			return

		default:
			fmt.Println("❌ Неверный выбор! Попробуйте снова.")
		}
	}
}
