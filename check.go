package main

import (
	"fmt"
	"time"

	"github.com/tealeg/xlsx"
)
/*
 Для реализации 
 - сделать предупреждение об отчестве и omitempty в XML
 - проверка даты рождения - не младше 9 класса (по дате приема)


*/
const (
	// Количество столбцов в строке
	excelColumnesCount = 29
)

// Обобщенный тип функции, проверяющей ячейку
type cellCheckFunc func(*xlsx.Cell) []string

// Проверка временной ячейки на тип и принадлежность к определенному временному интервалу
// Дата в формате строки "дд.мм.ГГГГ"
func checkTimeCellBetween(before, after string) cellCheckFunc {
	// Переводим строковые метки во временные
	timeBefore, err := time.Parse("02.01.2006", before)
	if err != nil {
		panic("(checkTimeCellBetween): Ошибка генерации проверки функции  - timeBefore!")
	}
	timeAfter, err := time.Parse("02.01.2006", after)
	if err != nil {
		panic("(checkTimeCellBetween): Ошибка генерации проверки функции - timeAfter!")
	}
	// cellCheckFunc
	return func(cell *xlsx.Cell) []string {
		errors := make([]string, 0)

		// Проверяем тип
		if (cell.Type() != xlsx.CellTypeDate) && (cell.Type() != xlsx.CellTypeNumeric) {
			errors = append(errors, "Ошибка типа ячейки - необходим тип Дата(дд.мм.ГГГГ)!")
			return errors
		}
		// Получаем время
		cellTime, err := cell.GetTime(false)
		if err != nil {
			errors = append(errors, "Ошибка значения ячейки - ошибка чтения даты из ячейки!")
			return errors
		}
		// Проверка на принадлежность к интервалу
		if cellTime.Before(timeBefore) || cellTime.After(timeAfter) {
			errors = append(errors, fmt.Sprintf("Ошибка значения ячейки - дата (%s) из ячейки не входит во временной интервал!", cellTime.Format("02.01.2006")))
		}
		return errors
	}
}

// Проверка ячейки со строкой
func checkStringCell() cellCheckFunc {
	return func(cell *xlsx.Cell) []string {
		errors := make([]string, 0)
		// Проверяем тип
		if cell.Type() != xlsx.CellTypeString {
			errors = append(errors, "Ошибка типа ячейки - необходим тип Строка")
			return errors
		}
		// Получаем данные из ячейки
		data := cell.String()
		if data == "" {
			errors = append(errors, "Ошибка значения ячейки - Пустая строка!")
			return errors
		}
		return errors
	}
}

// Проверка ячейки со строкой
func checkNumericCell() cellCheckFunc {
	return func(cell *xlsx.Cell) []string {
		errors := make([]string, 0)
		// Проверяем тип
		if cell.Type() != xlsx.CellTypeNumeric {
			errors = append(errors, "Ошибка типа ячейки - необходим тип Число")
			return errors
		}
		return errors
	}
}

// Проверка значения по словарю
func checkDictionaryCell(dic []string) cellCheckFunc {
	return func(cell *xlsx.Cell) []string {
		errors := make([]string, 0)

		// Проверяем тип
		if cell.Type() != xlsx.CellTypeString {
			errors = append(errors, "Ошибка типа ячейки - необходим тип Строка")
			return errors
		}
		// Получаем данные из ячейки
		data := cell.String()
		if data == "" {
			errors = append(errors, "Ошибка значения ячейки - Пустая строка!")
			return errors
		}
		// Сопоставление со словарем
		found := false
		for _, key := range dic {
			if key == data {
				found = true
			}
		}
		if !found {
			errors = append(errors, fmt.Sprintf("Ошибка значения ячейки - Значение (%s) не было найдено в словаре!", data))
		}
		return errors
	}
}

// Возвращает ключи словаря
func values(dic interface{}) []string {
	vs := []string{}
	switch dic.(type) {
	case map[string]int:
		for key := range dic.(map[string]int) {
			vs = append(vs, key)
		}
	case map[string]string:
		for key := range dic.(map[string]string) {
			vs = append(vs, key)
		}
	}
	return vs
}

// Проверка строки
func checkRow(row *xlsx.Row) ([]string, bool) {
	// Проверка на длину строки
	if len(row.Cells) < excelColumnesCount {
		return []string{"Строка не содержит достаточно количество столбцов!"}, true
	}

	errors := false
	messages := []string{}

	type column struct {
		Num       int
		Name      string
		Processor cellCheckFunc
	}
	columnes := []column{}
	// Модель проверки данных
	columnes = append(columnes, column{0, "Дата Регистрации", checkTimeCellBetween("01.06."+Configuration.Year, "01.09."+Configuration.Year)})
	//Номер заявления (Строка) - 1 --/--
	columnes = append(columnes, column{2, "Конкурсная группа", checkDictionaryCell(values(dictionarySpecialities))})
	columnes = append(columnes, column{3, "Фамилия", checkStringCell()})
	columnes = append(columnes, column{4, "Имя", checkStringCell()})
	columnes = append(columnes, column{5, "Отчество", checkStringCell()})
	columnes = append(columnes, column{6, "Пол", checkDictionaryCell(values(dictionaryGenders))})
	columnes = append(columnes, column{7, "Регион", checkDictionaryCell(values(dictionaryRegions))})
	columnes = append(columnes, column{8, "Тип населенного пункта", checkDictionaryCell(values(dictionaryCities))})
	columnes = append(columnes, column{9, "Адрес проживания", checkStringCell()})
	//Дата подачи оригиналов - 10 --/--
	columnes = append(columnes, column{11, "Вид документа (документ удостоверяющий личность)", checkDictionaryCell(values(dictionaryIdentityTypes))})
	columnes = append(columnes, column{12, "Гражданство (документ удостоверяющий личность)", checkDictionaryCell(values(dictionaryCountries))})
	//Серия - 13 --/--
	columnes = append(columnes, column{14, "Номер (документ удостоверяющий личность)", checkStringCell()})
	//Подразделение - 15 --/--
	columnes = append(columnes, column{16, "Кем выдан (документ удостоверяющий личность)", checkStringCell()})
	columnes = append(columnes, column{17, "Дата выдачи (документ удостоверяющий личность)", checkTimeCellBetween("01.01.1950", "01.09."+Configuration.Year)})
	columnes = append(columnes, column{18, "Дата рождения (документ удостоверяющий личность)", checkTimeCellBetween("01.01.1950", "01.09."+Configuration.Year)})
	columnes = append(columnes, column{19, "Место рождения (документ удостоверяющий личность)", checkStringCell()})
	columnes = append(columnes, column{20, "Вид документа (документ об образовании)", checkDictionaryCell(values(dictionaryEDocTypes))})
	//Дата предоставления оригиналов - 21 --/--
	//Серия - 22 --/--
	columnes = append(columnes, column{23, "Номер (документ об образовании)", checkStringCell()})
	columnes = append(columnes, column{24, "Дата выдачи (документ об образовании)", checkTimeCellBetween("01.01.1950", "01.09."+Configuration.Year)})
	//Регистрационный номер - 25 --/--
	columnes = append(columnes, column{26, "Год окончания (документ об образовании)", checkNumericCell()})
	columnes = append(columnes, column{27, "Средний балл (документ об образовании)", checkNumericCell()})
	columnes = append(columnes, column{28, "Организация (документ об образовании)", checkStringCell()})

	for _, col := range columnes {
		if errMessages := col.Processor(row.Cells[col.Num]); len(errMessages) > 0 {
			errors = true
			for _, err := range errMessages {
				messages = append(messages, fmt.Sprintf(" -- Столбец %d (%s): %s", col.Num, col.Name, err))
			}
		}
	}

	return messages, errors
}

// Проверка всех строк
func checkFile(file *xlsx.File) ([]string, bool) {
	rowNum := 0
	errors := false
	messages := []string{}

	for _, row := range file.Sheets[0].Rows {
		if rowNum > 1 {
			if errorMessages, isErrors := checkRow(row); isErrors {
				messages = append(messages, fmt.Sprintf("\nСрока %d", rowNum+1))
				for _, msg := range errorMessages {
					errors = true
					messages = append(messages, msg)
				}
			}
		}
		rowNum = rowNum + 1
	}
	if len(messages) > 0 {
		errors = true
	}
	return messages, errors
}
