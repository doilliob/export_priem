package main

import (
	"fmt"
	"regexp"

	"github.com/tealeg/xlsx"
)

/* addIdentityDocument Добавляет документ УЛ и возвращает указатель на него
 * если такой документ существует - возвращает ссылку на существующий документ */
func addIdentityDocument(idoc IdentityDocument) *IdentityDocument {
	// Ищем подходящий в массиве
	for _, doc := range domainDocumentsIdentity {
		if (doc.DocType == idoc.DocType) && (doc.Number == idoc.Number) && (doc.When == idoc.When) {
			return doc
		}
	}
	// Либо добавляем
	domainDocumentsIdentity = append(domainDocumentsIdentity, &idoc)
	return &idoc
}

// addEducationDocument Добавляет документ об образовании
func addEducationDocument(edoc EducationDocument) *EducationDocument {
	for _, doc := range domainDocumentsEducation {
		if (doc.Number == edoc.Number) && (doc.When == edoc.When) && (doc.IdentityDocument == edoc.IdentityDocument) {
			return doc
		}
	}
	domainDocumentsEducation = append(domainDocumentsEducation, &edoc)
	return &edoc
}

// addApplication Добавляет заявление в список, если на данного человека оно не зарегистрировано
func addApplication(iapp Application) {
	// Ищем подходящий
	for _, app := range domainDocumentsApplications {
		if (app.OrgCode == iapp.OrgCode) && (app.Speciality == iapp.Speciality) && (app.IdentityDocument == iapp.IdentityDocument) {
			return
		}
	}
	domainDocumentsApplications = append(domainDocumentsApplications, &iapp)
}

// parseRow Обрабатывает строку файла, формируя заявление
func parseRow(row *xlsx.Row, orgCode int) {
	//=================================
	// ДОКУМЕНТ УДОСТОВЕРЯЮЩИЙ ЛИЧНОСТЬ
	//=================================
	idoc := IdentityDocument{}
	// Фамилия - 3
	idoc.Lastname = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[3].String(), "")
	// Имя - 4
	idoc.Firstname = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[4].String(), "")
	// Отчество - 5
	idoc.Middlename = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[5].String(), "")
	// Дата рождения - 18
	if bdate, err := row.Cells[18].GetTime(false); err != nil {
		panic(" -- Ошибка преобразования поля (Дата рождения) для студента " + idoc.Lastname)
	} else {
		idoc.BirthDate = bdate
	}
	// Место рождения - 19
	idoc.BirthPlace = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[19].String(), "")
	// Пол - 6
	if _, ok := dictionaryGenders[row.Cells[6].String()]; !ok {
		panic(" -- Ошибка преобразования поля (Пол) для студента " + idoc.Lastname)
	}
	idoc.Gender = dictionaryGenders[row.Cells[6].String()]
	// Регион - 7
	if _, ok := dictionaryRegions[row.Cells[7].String()]; !ok {
		panic(" -- Ошибка преобразования поля (Регион) для студента " + idoc.Lastname)
	}
	idoc.Region = dictionaryRegions[row.Cells[7].String()]
	// Тип населенного пункта - 8
	if _, ok := dictionaryCities[row.Cells[8].String()]; !ok {
		panic(" -- Ошибка преобразования поля (Тип населенного пункта) для студента " + idoc.Lastname)
	}
	idoc.CityType = dictionaryCities[row.Cells[8].String()]
	// Адрес проживания (строка) - 9
	idoc.Address = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[9].String(), "")
	// Тип документа удостоверяющего личность - 11
	if _, ok := dictionaryIdentityTypes[row.Cells[11].String()]; !ok {
		panic(" -- Ошибка преобразования поля (Тип документа удостоверяющего личность) для студента " + idoc.Lastname)
	}
	idoc.DocType = dictionaryIdentityTypes[row.Cells[11].String()]
	// Серия - 13
	if row.Cells[13].Type() == xlsx.CellTypeString {
		idoc.Serial = regexp.MustCompile(`\s+`).ReplaceAllString(row.Cells[13].String(), "")
	}
	// Номер - 14
	idoc.Number = regexp.MustCompile(`\s+`).ReplaceAllString(row.Cells[14].String(), "")
	// Гражданство - 12
	if _, ok := dictionaryCountries[row.Cells[12].String()]; !ok {
		panic(" -- Ошибка преобразования поля (Гражданство) для студента " + idoc.Lastname)
	}
	idoc.Citizenship = dictionaryCountries[row.Cells[12].String()]
	// Кем выдан - 16
	idoc.Who = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[16].String(), "")
	// Когда - 17
	if iwhen, err := row.Cells[17].GetTime(false); err != nil {
		panic(" -- Ошибка преобразования поля (Дата выдачи УЛ) для студента " + idoc.Lastname)
	} else {
		idoc.When = iwhen
	}
	// Подразделение - 15
	if row.Cells[15].Type() == xlsx.CellTypeString {
		idoc.DepCode = regexp.MustCompile(`\s+`).ReplaceAllString(row.Cells[15].String(), "")
	}
	//============================
	// ДОКУМЕНТ ОБ ОБРАЗОВАНИИ
	//============================
	edoc := EducationDocument{}
	//Вид документа - 20
	edoc.DocType = dictionaryEDocTypes[row.Cells[20].String()]
	//Серия - 22 --/--
	if row.Cells[22].Type() == xlsx.CellTypeString {
		edoc.Serial = regexp.MustCompile(`\s+`).ReplaceAllString(row.Cells[22].String(), "")
	}
	//Номер - 23
	edoc.Number = regexp.MustCompile(`\s+`).ReplaceAllString(row.Cells[23].String(), "")
	//Дата выдачи - 24
	if ewhen, err := row.Cells[24].GetTime(false); err != nil {
		panic(" -- Ошибка преобразования поля (Дата выдачи обр) для студента " + idoc.Lastname)
	} else {
		edoc.When = ewhen
	}
	//Регистрационный номер - 25 --/--
	if row.Cells[25].Type() == xlsx.CellTypeString {
		edoc.RegisterNumber = regexp.MustCompile(`\s+`).ReplaceAllString(row.Cells[25].String(), "")
	}
	//Год окончания - 26
	if year, err := row.Cells[26].Int(); err != nil {
		panic(" -- Ошибка преобразования поля (Год окончания) для студента " + idoc.Lastname)
	} else {
		edoc.EndYear = year
	}
	//Средний балл - 27
	if gpa, err := row.Cells[27].Float(); err != nil {
		panic(" -- Ошибка преобразования поля (Средний балл) для студента " + idoc.Lastname)
	} else {
		edoc.GPA = gpa
	}
	//Наименование организации, выдавшей документ - 28
	edoc.Who = regexp.MustCompile(`^\s+|\s+$`).ReplaceAllString(row.Cells[28].String(), "")

	//==============================
	// ЗАЯВЛЕНИЕ НА ПОСТУПЛЕНИЕ
	//==============================
	app := Application{}
	//Дата Регистрации - 0
	if dt, err := row.Cells[0].GetTime(false); err != nil {
		panic(" -- Ошибка преобразования поля (Дата Регистрации Заявления) для студента " + idoc.Lastname)
	} else {
		app.Date = dt
	}
	//Номер заявления (Строка) - 1 --/--
	//Конкурсная группа - 2
	app.Speciality = dictionarySpecialities[row.Cells[2].String()]
	// Зачислен?
	if (row.Cells[29].Type() == xlsx.CellTypeString) && ((row.Cells[29].String() == "Да") || (row.Cells[29].String() == "да")) {
		app.Recommended = true
	}
	// Добавляем данные по умолчанию и сохраняем в массив
	edoc.IdentityDocument = addIdentityDocument(idoc)
	//
	app.EducationDocument = addEducationDocument(edoc)
	app.EducationOriginalDate = app.Date
	//
	app.IdentityDocument = edoc.IdentityDocument
	app.IdentityOriginalDate = app.Date
	//
	app.OrgCode = orgCode
	addApplication(app)
}

// parseFile обрабатывает файл и формирует структуры для выгрузки
func parseFile(xlsxFile *xlsx.File, orgCode int) {
	// Проверка на правильность заполнения ключевых полей
	if messages, errors := checkFile(xlsxFile); errors {
		fmt.Println("----------------------------------------------")
		fmt.Println("При проверке значений таблицы возникли ошибки:")
		fmt.Println("----------------------------------------------")
		for _, msg := range messages {
			fmt.Println(msg)
		}
		return
	}
	// Загрузка заявок
	rownum := 0
	for _, row := range xlsxFile.Sheets[0].Rows {
		if rownum > 1 {
			parseRow(row, orgCode)
		}
		rownum = rownum + 1
	}
}

func printStat() {
	fmt.Println("----------------------------------------------")
	fmt.Printf("\nСтатистика:\n")
	fmt.Println("----------------------------------------------")
	fmt.Printf(" - Сформировано %d документов УЛ\n", len(domainDocumentsIdentity))
	fmt.Printf(" - Сформировано %d документов об образовании\n", len(domainDocumentsEducation))
	fmt.Printf(" - Сформировано %d заявлений на зачисление\n", len(domainDocumentsApplications))
	// Рекомендовано к зачислению
	recommended := 0
	for _, app := range domainDocumentsApplications {
		if app.Recommended {
			recommended = recommended + 1
		}
	}
	fmt.Printf(" - Рекомендованных к зачислению: %d \n", recommended)
	fmt.Printf("\n")
}
