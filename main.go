package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

var (
	currentOrganizationCode = 0
	timingStartTime         = time.Now()
)

func printTime() {
	fmt.Println("---------------------------------------")
	fmt.Printf("Время работы программы: %d мсек\n", time.Since(timingStartTime).Milliseconds())
	fmt.Println("---------------------------------------")
}

// main Точка входа
func main() {

	fmt.Println("---------------------------------------")
	defer printTime()

	// Загружаем если есть файл конфигурации или конфигурацию по умолчанию
	loadConfiguration()

	/*=====================================================
	 * Команда config _
	 *=====================================================
	 * (4 arg) priem.exe config _ VALUE
	 *=====================================================*/
	if len(os.Args) == 4 && os.Args[1] == "config" {
		// login
		arg := os.Args[2]
		confChanged := false
		if arg == "login" {
			Configuration.Login = arg
			confChanged = true
		}
		// password
		if arg == "password" {
			Configuration.Password = arg
			confChanged = true
		}
		// year
		if arg == "year" {
			Configuration.Year = os.Args[3]
			confChanged = true
		}
		if confChanged {
			saveConfiguration()
			fmt.Println("Конфигурация успешно изменена!")
			return
		}
	}

	/*=====================================================
	 * Команда  check
	 *=====================================================
	 * (4 arg) priem.exe check FILENAME.xlsx
	 *=====================================================*/
	if len(os.Args) == 3 && os.Args[1] == "check" {
		filename := os.Args[2]
		// Наличие файла
		if !fileExist(filename) {
			fmt.Println("Ошибка! Файл " + filename + " не существует!")
			return
		}
		currentOrganizationCode = getOrganizationCode(filename)
		if currentOrganizationCode == 0 {
			fmt.Println("Имя файла " + filename + " Должно содержать часть названия организации (НМК, Баз, Без, Бор)!")
			return
		}

		file := openXLSX(filename)

		if messages, errors := checkFile(file); errors {
			for _, msg := range messages {
				fmt.Println(msg)
			}
			fmt.Println("---------------------------------------")
			fmt.Println("Проверка файла прошла с ошибками!")
			fmt.Println("---------------------------------------")
			return
		}

		fmt.Println("---------------------------------------")
		fmt.Println("Проверка файла прошла успешно!")
		fmt.Println("---------------------------------------")
		parseFile(file, getOrganizationCode(filename))
		printStat()
		return
	}

	/*=====================================================
	 * Команда generate
	 *=====================================================
	 * (4 arg) priem.exe generate FILENAME_IN.xlsx FILENAME_OUT.xml
	 *=====================================================*/
	if len(os.Args) == 4 && os.Args[1] == "generate" {
		filenameIn := os.Args[2]
		filenameOut := os.Args[3]

		// Наличие файла
		if !fileExist(filenameIn) {
			fmt.Println("Ошибка! Файл " + filenameIn + " не существует!")
			return
		}

		file := openXLSX(filenameIn)
		parseFile(file, getOrganizationCode(filenameIn))
		printStat()

		orgCode := getOrganizationCode(filenameIn)
		generateIODs(orgCode)

		acceptedStudents := make([]*Application, 0)
		deniedStudents := make([]*Application, 0)
		for _, app := range domainDocumentsApplications {
			if app.Recommended {
				acceptedStudents = append(acceptedStudents, app)
			} else {
				deniedStudents = append(deniedStudents, app)
			}
		}
		acceptedFile := regexp.MustCompile(`.xml`).ReplaceAllString(filenameOut, "_ПРИНЯТО.xml")
		deniedFile := regexp.MustCompile(`.xml`).ReplaceAllString(filenameOut, "_НЕПРИНЯТО.xml")
		writeTree(acceptedFile, generateTree(orgCode, acceptedStudents))
		writeTree(deniedFile, generateTree(orgCode, deniedStudents))
		//writeAppUIDS(filenameIn)
		return
	}

	/*=====================================================
	 * Команда generate со стартовым коэффициентом
	 *=====================================================
	 * (4 arg) priem.exe generate FILENAME_IN.xlsx FILENAME_OUT.xml START_NUMBER
	 *=====================================================*/
	if len(os.Args) == 5 && os.Args[1] == "generate" {
		filenameIn := os.Args[2]
		filenameOut := os.Args[3]
		if number, err := strconv.ParseInt(os.Args[4], 10, 32); err != nil {
			fmt.Println(err)
			return
		} else {
			generationStartNumber = int(number)
		}

		// Наличие файла
		if !fileExist(filenameIn) {
			fmt.Println("Ошибка! Файл " + filenameIn + " не существует!")
			return
		}

		file := openXLSX(filenameIn)
		parseFile(file, getOrganizationCode(filenameIn))
		printStat()

		orgCode := getOrganizationCode(filenameIn)
		generateIODs(orgCode)

		acceptedStudents := make([]*Application, 0)
		deniedStudents := make([]*Application, 0)
		for _, app := range domainDocumentsApplications {
			if app.Recommended {
				acceptedStudents = append(acceptedStudents, app)
			} else {
				deniedStudents = append(deniedStudents, app)
			}
		}
		acceptedFile := regexp.MustCompile(`.xml`).ReplaceAllString(filenameOut, "_ПРИНЯТО.xml")
		deniedFile := regexp.MustCompile(`.xml`).ReplaceAllString(filenameOut, "_НЕПРИНЯТО.xml")
		writeTree(acceptedFile, generateTree(orgCode, acceptedStudents))
		writeTree(deniedFile, generateTree(orgCode, deniedStudents))
		//writeAppUIDS(filenameIn)
		return
	}

	/*=====================================================
	 * Команда send
	 *=====================================================
	 * (4 arg) priem.exe send FILENAME.xml
	 *=====================================================*/
	// Высылает на портал ФИС ГИА и Приема XML-запрос,
	// сохраняя в XML-файле номер пакета
	if len(os.Args) == 3 && os.Args[1] == "send" {
		filename := os.Args[2]

		// Наличие файла
		if !fileExist(filename) {
			fmt.Println("Ошибка! Файл " + filename + " не существует!")
			return
		}

		answer, err := webapiPutService(filename)
		if err != nil {
			fmt.Printf("Ошибка при выполнении запроса! %s \n", err.Error())
			return
		}
		fmt.Println("Ответ:" + answer)
		return
	}

	/*=====================================================
	 * Команда status
	 *=====================================================
	 * (4 arg) priem.exe status FILENAME.xml
	 *=====================================================*/
	// Высылает на портал ФИС ГИА и Приема XML-запрос о статусе пакета
	// и возвращает текст статуса
	if len(os.Args) == 3 && os.Args[1] == "status" {
		filename := os.Args[2]

		// Наличие файла
		if !fileExist(filename) {
			fmt.Println("Ошибка! Файл " + filename + " не существует!")
			return
		}

		answer, err := webapiCheckService(filename)
		if err != nil {
			fmt.Printf("Ошибка при выполнении запроса! %s \n", err.Error())
			return
		}
		fmt.Println(answer)
		return
	}

	/*=====================================================
	 * Команда delete
	 *=====================================================
	 * (4 arg) priem.exe delete FILENAME.xml
	 *=====================================================*/
	// Формирует из XML-файла и высылает на портал ФИС ГИА и Приема XML-запрос об удалении заявлений
	// сохраняя в XML-файле номер пакета
	if len(os.Args) == 3 && os.Args[1] == "delete" {
		filename := os.Args[2]

		// Наличие файла
		if !fileExist(filename) {
			fmt.Println("Ошибка! Файл " + filename + " не существует!")
			return
		}

		answer, err := webapiDeleteService(filename)
		if err != nil {
			fmt.Printf("Ошибка при выполнении запроса! %s \n", err.Error())
			return
		}
		fmt.Println(answer)
		return
	}

	/*=====================================================
	 * Команда deletestatus
	 *=====================================================
	 * (4 arg) priem.exe deletestatus FILENAME.xml
	 *=====================================================*/
	// Высылает на портал ФИС ГИА и Приема XML-запрос о статусе пакета удаления
	// и выводит текст статуса на экран
	if len(os.Args) == 3 && os.Args[1] == "deletestatus" {
		filename := os.Args[2]

		// Наличие файла
		if !fileExist(filename) {
			fmt.Println("Ошибка! Файл " + filename + " не существует!")
			return
		}

		answer, err := webapiCheckDeleteService(filename)
		if err != nil {
			fmt.Printf("Ошибка при выполнении запроса! %s \n", err.Error())
			return
		}
		fmt.Println(answer)
		return
	}


	fmt.Println("---------------------------------------")
	fmt.Println("Неизвестные параметры программы!")
	fmt.Println("---------------------------------------")
}
