package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/tealeg/xlsx"
)

//
func decode(JSONFileIn, XLSFileOut string) {
	var data []map[string]interface{}

	fileData, err := ioutil.ReadFile(JSONFileIn)
	if err != nil {
		fmt.Println("Ошибка открытия файла!")
		return
	}
	// Обработка JSON
	if err := json.Unmarshal(bytes.TrimPrefix(fileData, []byte("\xef\xbb\xbf")), &data); err != nil {
		fmt.Println("Ошибка преобразования JSON файла!")
		fmt.Println(err.Error())
		return
	}
	// Открытие файла
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Заявления")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var row *xlsx.Row
	var cell *xlsx.Cell
	// Пропускаем две строки
	row = sheet.AddRow()
	row = sheet.AddRow()
	// Обрабатываем массив JSON
	for _, element := range data {
		// Для каждой специальности
		for _, specArr := range element["specs"].([]interface{}) {
			mp := specArr.(map[string]interface{})
			if (mp["id"].(string) == "") && (mp["name"].(string) == "") {
				continue
			}
			specCode := mp["id"].(string)
			specName := ""
			for name, code := range dictionarySpecialities {
				if code == specCode {
					specName = name
				}
			}

			row = sheet.AddRow()
			//Дата Регистрации - 0
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["app_data"].(string)); err == nil {
				cell.SetDate(tm)
			}
			//Номер заявления (Строка) - 1 --/--
			cell = row.AddCell()
			cell.SetString(element["app_number"].(string))
			//Конкурсная группа - 2
			cell = row.AddCell()
			cell.SetString(specName)
			//Фамилия - 3
			cell = row.AddCell()
			cell.SetString(element["lastname"].(string))
			//Имя - 4
			cell = row.AddCell()
			cell.SetString(element["firstname"].(string))
			//Отчество - 5
			cell = row.AddCell()
			cell.SetString(element["middlename"].(string))
			//Пол - 6
			cell = row.AddCell()
			for key, val := range dictionaryGenders {
				if val == int(element["gender"].(float64)) {
					cell.SetString(key)
				}
			}
			//Регион - 7
			cell = row.AddCell()
			for key, val := range dictionaryRegions {
				if val == int(element["region_id"].(float64)) {
					cell.SetString(key)
				}
			}
			//Тип населенного пункта - 8
			cell = row.AddCell()
			for key, val := range dictionaryCities {
				if val == int(element["town_type_id"].(float64)) {
					cell.SetString(key)
				}
			}
			//Адрес проживания (строка) - 9
			cell = row.AddCell()
			cell.SetString(element["address"].(string))
			//Дата подачи оригиналов - 10
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["app_data"].(string)); err == nil {
				cell.SetDate(tm)
			}
			//Вид документа - 11
			cell = row.AddCell()
			for key, val := range dictionaryIdentityTypes {
				if val == int(element["identity_document_type_id"].(float64)) {
					cell.SetString(key)
				}
			}
			//Гражданство - 12
			cell = row.AddCell()
			for key, val := range dictionaryCountries {
				if val == int(element["identity_document_nationality"].(float64)) {
					cell.SetString(key)
				}
			}
			//Серия - 13
			cell = row.AddCell()
			cell.SetString(element["identity_document_series"].(string))
			//Номер - 14
			cell = row.AddCell()
			cell.SetString(element["identity_document_number"].(string))
			//Подразделение - 15
			cell = row.AddCell()
			cell.SetString(element["identity_document_dep_code"].(string))
			//Кем выдан - 16
			cell = row.AddCell()
			cell.SetString(element["identity_document_organization"].(string))
			//Когда - 17
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["identity_document_date"].(string)); err == nil {
				cell.SetDate(tm)
			}
			//Дата рождения - 18
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["birth_date"].(string)); err == nil {
				cell.SetDate(tm)
			}
			//Место рождения - 19
			cell = row.AddCell()
			cell.SetString(element["identity_document_birth_place"].(string))
			//Вид документа - 20
			cell = row.AddCell()
			//Дата предоставления оригиналов - 21 --/--
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["app_data"].(string)); err == nil {
				cell.SetDate(tm)
			}
			//Серия - 22 --/--
			cell = row.AddCell()
			cell.SetString("")
			//Номер - 23
			cell = row.AddCell()
			cell.SetString(element["edoc_number"].(string))
			//Дата выдачи - 24
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["edoc_date"].(string)); err == nil {
				cell.SetDate(tm)
			}
			//Регистрационный номер - 25 --/--
			cell = row.AddCell()
			cell.SetString("")
			//Год окончания - 26
			cell = row.AddCell()
			if tm, err := time.Parse("02.01.2006", element["edoc_date"].(string)); err == nil {
				cell.SetInt(tm.Year())
			}
			//Средний балл - 27
			cell = row.AddCell()
			cell.SetFloat(element["edoc_gpa"].(float64))
			//Наименование организации, выдавшей документ - 28
			cell = row.AddCell()
			if element["edoc_organization"].(string) == "" {
				cell.SetString(dictionarySchools[rand.Intn(len(dictionarySchools))])
			} else {
				cell.SetString(element["edoc_organization"].(string))
			}
			//Принят - 29
			cell = row.AddCell()

		}

	}
	err = file.Save(XLSFileOut)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

}
