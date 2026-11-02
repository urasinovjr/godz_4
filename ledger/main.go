package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Transaction struct {
	ID          int
	Amount      float64
	Category    string
	Description string
	Date        string
}

func (t Transaction) Validate() error {
	if t.Amount == 0 {
		return errors.New("сумма транзакции не может быть нулевой")
	}
	if t.Amount < 0 {
		return errors.New("сумма транзакции не может быть отрицательной")
	}
	if t.Category == "" {
		return errors.New("у транзакции не указана категория")
	}
	return nil
}

var transactions = []Transaction{}

func AddTransaction(tx Transaction) error {
	err := tx.Validate()
	if err != nil {
		return err
	}

	budget, ok := budgets[tx.Category]
	if ok {
		spent := 0.0
		for _, t := range transactions {
			if t.Category == tx.Category {
				spent = spent + t.Amount
			}
		}

		if spent+tx.Amount > budget.Limit {
			return fmt.Errorf("превышен бюджет по категории %s", tx.Category)
		}
	}

	tx.ID = len(transactions) + 1
	transactions = append(transactions, tx)
	return nil
}

func ListTransactions() []Transaction {
	result := make([]Transaction, len(transactions))
	copy(result, transactions)
	return result
}

func tryAdd(tx Transaction) {
	err := AddTransaction(tx)
	if err != nil {
		fmt.Println("Ошибка! Транзакция \""+tx.Description+"\" не добавлена:", err)
		return
	}
	fmt.Printf("Транзакция \"%s\" добавлена: %s, %.2f руб.\n", tx.Description, tx.Category, tx.Amount)
}

func trySetBudget(b Budget) {
	err := SetBudget(b)
	if err != nil {
		fmt.Println("Ошибка! Бюджет \""+b.Category+"\" не установлен:", err)
		return
	}
	fmt.Printf("Бюджет \"%s\" установлен: %.2f руб.\n", b.Category, b.Limit)
}

func main() {
	fmt.Println("Ledger service started")

	fmt.Println("\nПроверка через CheckValid:")
	items := []Validatable{
		Transaction{Amount: 500, Category: "еда", Date: "2026-11-01"},
		Transaction{Amount: 0, Category: "еда", Date: "2026-11-01"},
		Budget{Category: "еда", Limit: -100, Period: "месяц"},
	}
	for i, item := range items {
		err := CheckValid(item)
		if err != nil {
			fmt.Println("Проверка", i+1, "- ошибка:", err)
		} else {
			fmt.Println("Проверка", i+1, "- все в порядке")
		}
	}

	fmt.Println("\nУстанавливаем бюджеты:")
	trySetBudget(Budget{Category: "еда", Limit: 5000, Period: "месяц"})
	trySetBudget(Budget{Category: "транспорт", Limit: 2000, Period: "месяц"})
	trySetBudget(Budget{Category: "одежда", Limit: 0, Period: "месяц"})
	PrintBudgets()

	fmt.Println("\nТраты на еду (лимит 5000):")
	tryAdd(Transaction{Amount: 3000, Category: "еда", Description: "продукты на неделю", Date: "2026-11-01"})
	tryAdd(Transaction{Amount: 1500, Category: "еда", Description: "кафе", Date: "2026-11-02"})
	tryAdd(Transaction{Amount: 1000, Category: "еда", Description: "доставка пиццы", Date: "2026-11-03"})

	fmt.Println("\nКатегория без бюджета:")
	tryAdd(Transaction{Amount: 10000, Category: "подарки", Description: "подарок маме", Date: "2026-11-03"})

	fmt.Println("\nНеправильные транзакции:")
	tryAdd(Transaction{Amount: -200, Category: "еда", Description: "что-то странное", Date: "2026-11-03"})
	tryAdd(Transaction{Amount: 700, Category: "", Description: "непонятно что", Date: "2026-11-03"})

	fmt.Println("\nПоднимаем лимит на еду до 6000 и пробуем пиццу еще раз:")
	trySetBudget(Budget{Category: "еда", Limit: 6000, Period: "месяц"})
	tryAdd(Transaction{Amount: 1000, Category: "еда", Description: "доставка пиццы", Date: "2026-11-04"})

	fmt.Println("\nЗагружаем бюджеты из budgets.json")
	file, err := os.Open("budgets.json")
	if err != nil {
		fmt.Println("Ошибка! Не получилось открыть файл:", err)
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		err = LoadBudgets(reader)
		if err != nil {
			fmt.Println("Ошибка!", err)
		} else {
			fmt.Println("Бюджеты загружены")
		}
	}
	PrintBudgets()

	fmt.Println("\nТраты после загрузки из файла:")
	tryAdd(Transaction{Amount: 2200, Category: "транспорт", Description: "проездной", Date: "2026-11-04"})
	tryAdd(Transaction{Amount: 2000, Category: "развлечения", Description: "концерт", Date: "2026-11-05"})
	tryAdd(Transaction{Amount: 1500, Category: "развлечения", Description: "кино и боулинг", Date: "2026-11-05"})

	fmt.Println("\nПробуем загрузить битый JSON (лимит строкой):")
	err = LoadBudgets(strings.NewReader(`[{"category": "еда", "limit": "много"}]`))
	if err != nil {
		fmt.Println("Ошибка!", err)
	}

	fmt.Println("\nСписок транзакций:")
	list := ListTransactions()
	for _, t := range list {
		fmt.Printf("  %d. %s, %.2f руб., %s (%s)\n", t.ID, t.Category, t.Amount, t.Description, t.Date)
	}
	fmt.Println("Всего транзакций:", len(list))
}
