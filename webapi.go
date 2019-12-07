package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const (
	serviceURI               = "http://10.0.3.1:8080/import/importservice.svc"
	serviceMethodPut         = "/import"
	serviceMethodDelete      = "/delete"
	serviceMethodCheck       = "/import/result"
	serviceMethodCheckDelete = "/delete/result"
)

//
//webapiPutService отправляет запрос
func webapiPutService(XMLFilename string) (string, error) {
	var answer []byte
	file, _ := os.Open(XMLFilename)
	reader := io.Reader(file)
	resp, err := http.Post(serviceURI+serviceMethodPut, "text/xml", reader)
	if err != nil {
		return "", err
	}
	answer, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	text := string(answer)
	if matched, err := regexp.MatchString("PackageID", text); matched && (err == nil) {
		packageCode := regexp.MustCompile(`.*\<PackageID\>(\d+)\<\/PackageID\>.*`).ReplaceAllString(text, "$1")
		text = "Номер пакета: " + packageCode
		// Вставляем в файл XML код пакета в виде комментария
		buff, err := ioutil.ReadFile(XMLFilename)
		if err == nil {
			content := string(buff)
			content = regexp.MustCompile("<!--.*-->").ReplaceAllString(content, "<!-- "+packageCode+" -->")
			ioutil.WriteFile(XMLFilename, []byte(content), 0644)
		}

	}
	return text, nil
}

// webapiCheckService делает запрос и проверяет статус пакета
func webapiCheckService(XMLFilename string) (string, error) {
	// Проверяем наличие номера пакета в виде комментария
	buff, err := ioutil.ReadFile(XMLFilename)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл")
	}
	content := string(buff)
	if matched, err := regexp.MatchString(`<!--\s+(\d+)\s+-->`, content); !matched || (err != nil) {
		return "", fmt.Errorf("не удалось найти номер пакета, возможно - пакет не отправлен")
	}
	// Выделяем номер пакета и авторизационные данные
	packageID := regexp.MustCompile(`<!--\s+(\d+)\s+-->`).FindString(content)
	packageID = regexp.MustCompile(`<!--\s+(\d+)\s+-->`).ReplaceAllString(packageID, "<PackageID>$1</PackageID>")
	login := regexp.MustCompile(`<Login>.*</Login>`).FindString(content)
	passworg := regexp.MustCompile(`<Pass>.*</Pass>`).FindString(content)
	institutionID := regexp.MustCompile(`<InstitutionID>.*</InstitutionID>`).FindString(content)

	data := `<?xml version="1.0" encoding="utf-8"?>
	<Root>
		<AuthData>
			<Login></Login>
			<Pass></Pass>
			<InstitutionID></InstitutionID>
		</AuthData>
		<GetResultImportApplication>
			<PackageID></PackageID>
		</GetResultImportApplication>
	</Root>`
	data = regexp.MustCompile("<Login></Login>").ReplaceAllString(data, login)
	data = regexp.MustCompile("<Pass></Pass>").ReplaceAllString(data, passworg)
	data = regexp.MustCompile("<InstitutionID></InstitutionID>").ReplaceAllString(data, institutionID)
	data = regexp.MustCompile("<PackageID></PackageID>").ReplaceAllString(data, packageID)
	resp, err := http.Post(serviceURI+serviceMethodCheck, "text/xml", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	var answer []byte
	answer, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(answer), nil
}

// webapiDeleteService делает запрос и проверяет статус пакета
func webapiDeleteService(XMLFilename string) (string, error) {
	// Проверяем наличие номера пакета в виде комментария
	buff, err := ioutil.ReadFile(XMLFilename)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл")
	}
	content := string(buff)
	// Выделяем номер пакета и авторизационные данные
	login := regexp.MustCompile(`<Login>.*</Login>`).FindString(content)
	passworg := regexp.MustCompile(`<Pass>.*</Pass>`).FindString(content)
	institutionID := regexp.MustCompile(`<InstitutionID>.*</InstitutionID>`).FindString(content)
	// Находим все номера заявлений
	appNumberArray := regexp.MustCompile(`<ApplicationNumber>.*</ApplicationNumber>`).FindAllString(content, -1)
	regNumberArray := regexp.MustCompile(`<RegistrationDate>.*</RegistrationDate>`).FindAllString(content, -1)
	// Проверяем количество
	if len(appNumberArray) != len(regNumberArray) {
		return "", fmt.Errorf("количество регистрационных номеров и дат заявлений не совпадают")
	}
	// Формируем структуру найденных номеров и дат приложений
	appsString := ""
	for i, appNumber := range appNumberArray {
		str := `
			<Application>
				<ApplicationNumber></ApplicationNumber>
				<RegistrationDate></RegistrationDate>
			</Application>
		`
		str = regexp.MustCompile("<ApplicationNumber></ApplicationNumber>").ReplaceAllString(str, appNumber)
		str = regexp.MustCompile("<RegistrationDate></RegistrationDate>").ReplaceAllString(str, regNumberArray[i])
		appsString = appsString + str
	}
	// Формируем итоговый пакет
	data := `<?xml version="1.0" encoding="utf-8"?>
	<Root>
		<AuthData>
			<Login></Login>
			<Pass></Pass>
			<InstitutionID></InstitutionID>
		</AuthData>
		  <DataForDelete>
			<Applications>
				APPSTRING
			</Applications>
		 </DataForDelete>
	</Root>`
	data = regexp.MustCompile("<Login></Login>").ReplaceAllString(data, login)
	data = regexp.MustCompile("<Pass></Pass>").ReplaceAllString(data, passworg)
	data = regexp.MustCompile("<InstitutionID></InstitutionID>").ReplaceAllString(data, institutionID)
	data = regexp.MustCompile("APPSTRING").ReplaceAllString(data, appsString)
	// Отправляем на сервис
	resp, err := http.Post(serviceURI+serviceMethodDelete, "text/xml", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	var answer []byte
	answer, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	text := string(answer)
	if matched, err := regexp.MatchString("PackageID", text); matched && (err == nil) {
		packageCode := regexp.MustCompile(`.*\<PackageID\>(\d+)\<\/PackageID\>.*`).ReplaceAllString(text, "$1")
		text = "Номер пакета: " + packageCode
		// Вставляем в файл XML код пакета в виде комментария
		buff, err := ioutil.ReadFile(XMLFilename)
		if err == nil {
			content := string(buff)
			content = regexp.MustCompile("<!--.*-->").ReplaceAllString(content, "<!-- "+packageCode+" -->")
			ioutil.WriteFile(XMLFilename, []byte(content), 0644)
		}

	}
	return text, nil
}

// webapiCheckDeleteService делает запрос и проверяет статус пакета
func webapiCheckDeleteService(XMLFilename string) (string, error) {
	// Проверяем наличие номера пакета в виде комментария
	buff, err := ioutil.ReadFile(XMLFilename)
	if err != nil {
		return "", fmt.Errorf("не удалось прочитать файл")
	}
	content := string(buff)
	if matched, err := regexp.MatchString(`<!--\s+(\d+)\s+-->`, content); !matched || (err != nil) {
		return "", fmt.Errorf("не удалось найти номер пакета, возможно - пакет не отправлен на удаление")
	}
	// Выделяем номер пакета и авторизационные данные
	packageID := regexp.MustCompile(`<!--\s+(\d+)\s+-->`).FindString(content)
	packageID = regexp.MustCompile(`<!--\s+(\d+)\s+-->`).ReplaceAllString(packageID, "<PackageID>$1</PackageID>")
	login := regexp.MustCompile(`<Login>.*</Login>`).FindString(content)
	passworg := regexp.MustCompile(`<Pass>.*</Pass>`).FindString(content)
	institutionID := regexp.MustCompile(`<InstitutionID>.*</InstitutionID>`).FindString(content)

	data := `<?xml version="1.0" encoding="utf-8"?>
	<Root>
		<AuthData>
			<Login></Login>
			<Pass></Pass>
			<InstitutionID></InstitutionID>
		</AuthData>
		<GetResultDeleteApplication>
			<PackageID></PackageID>
		</GetResultDeleteApplication>
	</Root>`
	data = regexp.MustCompile("<Login></Login>").ReplaceAllString(data, login)
	data = regexp.MustCompile("<Pass></Pass>").ReplaceAllString(data, passworg)
	data = regexp.MustCompile("<InstitutionID></InstitutionID>").ReplaceAllString(data, institutionID)
	data = regexp.MustCompile("<PackageID></PackageID>").ReplaceAllString(data, packageID)
	resp, err := http.Post(serviceURI+serviceMethodCheckDelete, "text/xml", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	var answer []byte
	answer, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(answer), nil
}
