package main // Наш главный пакет приложения

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gostudy/wallet" // Импортируем наш пакет кошелька
)

func main() {
	// Создаем кошелек Ивана (стартовые балансы теперь задаются внутри Load, если файла нет)
	myWallet := wallet.NewWallet("Иван", "wallet.json") // 🔥 ВОТ ЭТУ СТРОЧКУ НУЖНО ВЕРНУТЬ!

	err := myWallet.Load() // Теперь всё отлично, myWallet существует!
	if err != nil {

		fmt.Println("❌ Критическая ошибка при загрузке кошелька:", err)
		return
	}

	scanner := bufio.NewScanner(os.Stdin) // Сканер для чтения консоли

	for {
		fmt.Printf("\n🪙 КРИПТО-КОШЕЛЕК [%s]\n", myWallet.Owner)
		fmt.Println("💰 Текущие балансы:")

		// 🔥 МАГИЯ MAP: Перебираем карту. На каждом круге получаем имя валюты (curr) и баланс (bal) [INDEX]
		for curr, bal := range myWallet.Balances {
			fmt.Printf("   • %s: $%.2f\n", curr, bal)
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

			// 🔥 ЗАЩИТА: Проверяем, входит ли введенная строка в список разрешенных
			if currency != "USD" && currency != "EUR" && currency != "RUB" {
				fmt.Println("❌ Ошибка: Поддерживаются только валюты USD, EUR и RUB!")
				continue // Прерываем этот шаг и возвращаем пользователя в начало меню
			}

			fmt.Print("Введите сумму для пополнения: ")
			scanner.Scan()
			amount, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

			if amount > 0 {
				myWallet.Deposit(currency, amount)
				fmt.Println("✅ Кошелек успешно пополнен!")
			} else {
				fmt.Println("❌ Ошибка: Сумма должна быть больше нуля")
			}

		case "2":
			fmt.Print("Введите валюту (USD, EUR, RUB): ")
			scanner.Scan()
			currency := strings.ToUpper(strings.TrimSpace(scanner.Text()))

			// 🔥 ЗАЩИТА: Такая же проверка для снятия денег
			if currency != "USD" && currency != "EUR" && currency != "RUB" {
				fmt.Println("❌ Ошибка: Поддерживаются только валюты USD, EUR и RUB!")
				continue
			}

			fmt.Print("Введите сумму для снятия: ")
			scanner.Scan()
			amount, _ := strconv.ParseFloat(strings.TrimSpace(scanner.Text()), 64)

			if amount > 0 {
				if myWallet.Withdraw(currency, amount) {
					fmt.Println("✅ Деньги успешно сняты!")
				} else {
					fmt.Println("❌ Ошибка: Недостаточно средств на балансе данной валюты!")
				}
			} else {
				fmt.Println("❌ Ошибка: Сумма должна быть больше нуля")
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
				// В выводе истории теперь тоже красиво подсвечиваем валюту операции
				fmt.Printf("[%s] %s: %.2f %s\n", dateStr, tx.Type, tx.Amount, tx.Currency)
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
