package main

import (
	"os"

	"github.com/tealeg/xlsx"
)

// Проверяет существование файла
func fileExist(filename string) bool {
	// Проверка на существование файла
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// Открывает файл XLSX и возвращает его
func openXLSX(filename string) *xlsx.File {
	xlsxFile, err := xlsx.OpenFile(filename)
	if err != nil {
		panic("Файл не существует!: " + filename)
	}
	return xlsxFile
}

// ReadListOfStrings Читает и возвращает из XLSX файла
func ReadListOfStrings(filename string) []string {
	xlsxFile := openXLSX(filename)
	var list []string
	for _, sheet := range xlsxFile.Sheets {
		for _, rows := range sheet.Rows {
			if rows.Cells[0].Type() == xlsx.CellTypeString {
				list = append(list, rows.Cells[0].String())
			}
		}
	}
	return list
}

// ReadsMapOfStrings Читает и возвращает карту [string] => string из XLSX файла
func ReadsMapOfStrings(filename string) map[string]string {
	xlsxFile := openXLSX(filename)
	var list map[string]string
	list = make(map[string]string)
	for _, sheet := range xlsxFile.Sheets {
		for _, rows := range sheet.Rows {
			if (rows.Cells[0].Type() == xlsx.CellTypeString) && (rows.Cells[1].Type() == xlsx.CellTypeString) {
				list[rows.Cells[0].String()] = rows.Cells[1].String()
			}
		}
	}
	return list
}

// ReadsMapOfStringInt Читает и возвращает карту [string] => int из XLSX файла
func ReadsMapOfStringInt(filename string) map[string]int {
	xlsxFile := openXLSX(filename)
	var list map[string]int
	list = make(map[string]int)
	for _, sheet := range xlsxFile.Sheets {
		for _, rows := range sheet.Rows {
			if (rows.Cells[0].Type() == xlsx.CellTypeString) && (rows.Cells[1].Type() == xlsx.CellTypeNumeric) {
				i, err := rows.Cells[1].Int()
				if err == nil {
					list[rows.Cells[0].String()] = i
				}
			}
		}
	}
	return list
}

// ReadsMapOfIntString Читает и возвращает карту [int] => string из XLSX файла
func ReadsMapOfIntString(filename string) map[int]string {
	xlsxFile := openXLSX(filename)
	var list map[int]string
	list = make(map[int]string)
	for _, sheet := range xlsxFile.Sheets {
		for _, rows := range sheet.Rows {
			if (rows.Cells[0].Type() == xlsx.CellTypeNumeric) && (rows.Cells[1].Type() == xlsx.CellTypeString) {
				i, err := rows.Cells[0].Int()
				if err == nil {
					list[i] = rows.Cells[1].String()
				}
			}
		}
	}
	return list
}
