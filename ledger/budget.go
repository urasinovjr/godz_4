package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type Budget struct {
	Category string  `json:"category"`
	Limit    float64 `json:"limit"`
	Period   string  `json:"period"`
}

func (b Budget) Validate() error {
	if b.Limit <= 0 {
		return errors.New("лимит бюджета должен быть больше нуля")
	}
	if b.Category == "" {
		return errors.New("у бюджета не указана категория")
	}
	return nil
}

var budgets = map[string]Budget{}

func SetBudget(b Budget) error {
	err := b.Validate()
	if err != nil {
		return err
	}
	budgets[b.Category] = b
	return nil
}

func LoadBudgets(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("не удалось прочитать бюджеты: %w", err)
	}

	var list []Budget
	err = json.Unmarshal(data, &list)
	if err != nil {
		return fmt.Errorf("не удалось разобрать JSON с бюджетами: %w", err)
	}

	for _, b := range list {
		err = SetBudget(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func PrintBudgets() {
	fmt.Println("Бюджеты:")
	for _, b := range budgets {
		fmt.Printf("  %s: %.2f руб. на %s\n", b.Category, b.Limit, b.Period)
	}
}
