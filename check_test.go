package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/tealeg/xlsx"
)

// Тестирование фспомогательных функций

// Правильное вхождение
func TestCheckTimeCellBetweenNormal(t *testing.T) {

	cell := xlsx.Cell{}
	tm, err := time.Parse("02.01.2006", "27.11.2019")
	if err != nil {
		t.Errorf("Ошибка: неправильно выполнился time.Parse в TestCheckTimeCellBetweenNormal! %s", err)
	}
	cell.SetDate(tm)

	if msgs := checkTimeCellBetween("26.11.2019", "28.11.2019"); len(msgs(&cell)) > 0 {
		for _, msg := range msgs(&cell) {
			fmt.Println(msg)
		}
		t.Errorf("Ошибка: функция неправильно определяет вхождение %s в интервал 26.11.2019-28.11.2019!", tm.Format("02.01.2006"))
	}
}

// Неправильное вхождение
func TestCheckTimeCellBetweenFail(t *testing.T) {

	cell := xlsx.Cell{}
	tm, err := time.Parse("02.01.2006", "20.11.2019")
	if err != nil {
		t.Errorf("Ошибка: неправильно выполнился time.Parse в TestCheckTimeCellBetweenNormal! %s", err)
	}
	cell.SetDate(tm)

	if msgs := checkTimeCellBetween("26.11.2019", "28.11.2019"); len(msgs(&cell)) == 0 {
		for _, msg := range msgs(&cell) {
			fmt.Println(msg)
		}
		t.Errorf("Ошибка: функция неправильно определяет вхождение %s в интервал 26.11.2019-28.11.2019!", tm.Format("02.01.2006"))
	}
}

// Проверка на панику
func TestCheckTimeCellBetweenPanic1(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Ошибка! Не возникло исключение при неправильном формате даты")
		}
	}()
	checkTimeCellBetween("2019.11.26", "28.11.2019")
}

// Проверка на панику
func TestCheckTimeCellBetweenPanic2(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Ошибка! Не возникло исключение при неправильном формате даты")
		}
	}()
	checkTimeCellBetween("28.11.2019", "2019.11.26")
}
