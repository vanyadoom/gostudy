package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"gostudy/notifier"
)

var (
	ErrInsufficientFunds = errors.New("недостаточно средств на балансе данной валюты")
	ErrInvalidAmount     = errors.New("сумма операции должна быть больше нуля")
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
		Balances: make(map[string]float64),
		History:  []Transaction{},
	}
}

// 🔥 НОВЫЙ МЕТОД: Фоновый аудитор.
// Этот метод будет запущен как ОТДЕЛЬНАЯ ГОРУТИНА параллельно с меню!
func (w *Wallet) StartBackgroundAuditor() {
	// Запускаем бесконечный фоновый цикл робота [INDEX]
	for {
		// Заставляем робота поспать 15 секунд, чтобы он не грузил процессор непрерывно [INDEX]
		time.Sleep(2 * time.Minute)

		var totalUSDEquivalent float64
		for curr, bal := range w.Balances {
			// 🔥 ЗАМЕНИЛИ IF НА SWITCH: Код стал чистым, синяя линия исчезнет! [INDEX]
			switch curr {
			case "RUB":
				totalUSDEquivalent += bal / 100.0 // Переводим рубли в доллары [INDEX]
			case "EUR":
				totalUSDEquivalent += bal * 1.1 // Переводим евро в доллары [INDEX]
			default:
				totalUSDEquivalent += bal // Доллары считаем как есть
			}
		}

		// Робот делает фоновый анализ и выводит совет прямо в консоль [INDEX]
		fmt.Print("\n🤖 [ФОНОВЫЙ АУДИТОР]: Анализирую ваши активы... ")
		if totalUSDEquivalent > 200 {
			fmt.Printf("У вас солидный капитал (~$%.2f). Подумайте об инвестициях!\n", totalUSDEquivalent)
		} else {
			fmt.Printf("Ваш капитал скромный (~$%.2f). Рекомендуем пополнить счет.\n", totalUSDEquivalent)
		}
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

func (w *Wallet) Deposit(currency string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	w.Balances[currency] += amount

	tx := Transaction{
		Type:     "Пополнение",
		Currency: currency,
		Amount:   amount,
		Date:     time.Now(),
	}
	w.History = append(w.History, tx)

	if err := w.Save(); err != nil {
		return fmt.Errorf("ошибка сохранения базы данных: %w", err)
	}

	msg := fmt.Sprintf("Баланс пополнен на %.2f %s. Новый баланс: %.2f %s", amount, currency, w.Balances[currency], currency)
	w.notifier.Send(msg)

	return nil
}

func (w *Wallet) Withdraw(currency string, amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	currentBalance := w.Balances[currency]
	if amount > currentBalance {
		return ErrInsufficientFunds
	}

	w.Balances[currency] -= amount

	tx := Transaction{
		Type:     "Снятие",
		Currency: currency,
		Amount:   amount,
		Date:     time.Now(),
	}
	w.History = append(w.History, tx)

	if err := w.Save(); err != nil {
		return fmt.Errorf("ошибка сохранения базы данных: %w", err)
	}

	msg := fmt.Sprintf("Со счета снято %.2f %s. Остаток: %.2f %s", amount, currency, w.Balances[currency], currency)
	w.notifier.Send(msg)

	return nil
}
