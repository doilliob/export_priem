package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/tealeg/xlsx"
)

var (
	generationStartNumber int = 0
)

type xmlIdentityDocument struct {
	XMLName                xml.Name `xml:"IdentityDocument"`
	UID                    string   `xml:"UID"`
	IdentityDocumentTypeID int      `xml:"IdentityDocumentTypeID"`
	OriginalReceivedDate   string   `xml:"OriginalReceivedDate"`
	LastName               string   `xml:"LastName"`
	FirstName              string   `xml:"FirstName"`
	MiddleName             string   `xml:"MiddleName"`
	GenderID               int      `xml:"GenderID"`
	DocumentSeries         string   `xml:"DocumentSeries,omitempty"`
	DocumentNumber         string   `xml:"DocumentNumber"`
	DocumentDate           string   `xml:"DocumentDate"`
	DocumentOrganization   string   `xml:"DocumentOrganization"`
	NationalityTypeID      int      `xml:"NationalityTypeID"`
	BirthDate              string   `xml:"BirthDate"`
}

type xmlEducationDocument struct {
	XMLName              xml.Name `xml:""`
	UID                  string   `xml:"UID"`
	OriginalReceivedDate string   `xml:"OriginalReceivedDate"`
	DocumentNumber       string   `xml:"DocumentNumber"`
	DocumentDate         string   `xml:"DocumentDate"`
	Organization         string   `xml:"DocumentOrganization"`
	GPA                  float64  `xml:"GPA"`
}

type xmlEduDocument struct {
	XMLName  xml.Name `xml:"EduDocument"`
	Document xmlEducationDocument
}

type xmlEduDocuments struct {
	XMLName      xml.Name `xml:"EduDocuments"`
	EduDocuments []xmlEduDocument
}
type xmlApplicationDocuments struct {
	XMLName          xml.Name `xml:"ApplicationDocuments"`
	IdentityDocument xmlIdentityDocument
	EduDocuments     xmlEduDocuments
}

type xmlMailAddress struct {
	XMLName    xml.Name `xml:"MailAddress"`
	RegionID   int      `xml:"RegionID"`
	TownTypeID int      `xml:"TownTypeID"`
	Address    string   `xml:"Address"`
}

type xmlEmailOrMailAddress struct {
	XMLName     xml.Name `xml:"EmailOrMailAddress"`
	MailAddress xmlMailAddress
}

type xmlEntrant struct {
	XMLName            xml.Name `xml:"Entrant"`
	UID                string   `xml:"UID"`
	LastName           string   `xml:"LastName"`
	FirstName          string   `xml:"FirstName"`
	MiddleName         string   `xml:"MiddleName"`
	GenderID           int      `xml:"GenderID"`
	EmailOrMailAddress xmlEmailOrMailAddress
}

type xmlFinSourceEduForm struct {
	XMLName             xml.Name `xml:"FinSourceEduForm"`
	CompetitiveGroupUID string   `xml:"CompetitiveGroupUID"`
}

type xmlFinSourceAndEduForms struct {
	XMLName          xml.Name `xml:"FinSourceAndEduForms"`
	FinSourceEduForm xmlFinSourceEduForm
}

type xmlApplication struct {
	XMLName              xml.Name `xml:"Application"`
	UID                  string   `xml:"UID"`
	ApplicationNumber    string   `xml:"ApplicationNumber"`
	Entrant              xmlEntrant
	RegistrationDate     string `xml:"RegistrationDate"`
	StatusID             int    `xml:"StatusID"`
	NeedHostel           bool   `xml:"NeedHostel"`
	FinSourceAndEduForms xmlFinSourceAndEduForms
	ApplicationDocuments xmlApplicationDocuments
}

type xmlPackageApplications struct {
	XMLName      xml.Name `xml:"Applications"`
	Applications []xmlApplication
}
type xmlPackageData struct {
	XMLName      xml.Name `xml:"PackageData"`
	Applications xmlPackageApplications
}

type xmlAuthData struct {
	XMLName       xml.Name `xml:"AuthData"`
	Login         string   `xml:"Login"`
	Pass          string   `xml:"Pass"`
	InstitutionID int      `xml:"InstitutionID"`
}

type xmlRoot struct {
	XMLName     xml.Name `xml:"Root"`
	AuthData    xmlAuthData
	PackageData xmlPackageData
}

// generateOIDs  генерирует OIDs для документов
// Application UID = 2017-3920-СД-9-ОЧ-Б-0
// ApplicationNumber = 2017-3920-СД-9-ОЧ-Б-0- (1)
// Entrant UID = 2017-3920-E-0
// IdentityDocument UID = 2017-3920-I-0
// SchoolCertificateBasicDocument UID = 2017-3920-E-0
func generateIODs(orgCode int) {
	currentYear := Configuration.Year
	edocNumber := generationStartNumber
	for _, edoc := range domainDocumentsEducation {
		edoc.UID = currentYear + "-" + strconv.Itoa(orgCode) + "-EDOC-" + strconv.Itoa(edocNumber)
		edocNumber = edocNumber + 1
	}
	idocNumber := generationStartNumber
	for _, idoc := range domainDocumentsIdentity {
		idoc.UID = currentYear + "-" + strconv.Itoa(orgCode) + "-IDOC-" + strconv.Itoa(idocNumber)
		idocNumber = idocNumber + 1
	}
	appNumber := generationStartNumber
	for _, app := range domainDocumentsApplications {
		app.UID = currentYear + "-" + strconv.Itoa(orgCode) + "-APP-" + strconv.Itoa(appNumber)
		app.Number = currentYear + "-" + strconv.Itoa(orgCode) + "-" + app.Speciality + "-" + strconv.Itoa(appNumber)
		appNumber = appNumber + 1
	}
}

// generateTree генерирует XML Для всех структур
func generateTree(orgCode int, ApplicationsList []*Application) *xmlRoot {

	login := Configuration.Login
	password := Configuration.Password

	root := xmlRoot{}
	root.AuthData = xmlAuthData{Login: login, Pass: password, InstitutionID: orgCode}
	root.PackageData = xmlPackageData{}
	root.PackageData.Applications = xmlPackageApplications{}

	for _, app := range ApplicationsList {
		idoc := app.IdentityDocument
		edoc := app.EducationDocument
		// Application
		XMLApplication := xmlApplication{}
		XMLApplication.UID = app.UID
		XMLApplication.ApplicationNumber = app.Number
		XMLApplication.NeedHostel = false
		XMLApplication.RegistrationDate = app.Date.Format("2006-01-02") + "T10:00:01+00:00" //2017-06-20T10:00:01+00:00
		// StatusID
		XMLApplication.StatusID = documentStatusIntroduced // Принято
		// Entrant
		XMLApplication.Entrant = xmlEntrant{}
		XMLApplication.Entrant.UID = idoc.UID + "-ENTR"
		XMLApplication.Entrant.LastName = idoc.Lastname
		XMLApplication.Entrant.FirstName = idoc.Firstname
		XMLApplication.Entrant.MiddleName = idoc.Middlename
		XMLApplication.Entrant.GenderID = idoc.Gender
		// EmailOrMailAddress
		XMLApplication.Entrant.EmailOrMailAddress = xmlEmailOrMailAddress{}
		XMLApplication.Entrant.EmailOrMailAddress.MailAddress = xmlMailAddress{}
		XMLApplication.Entrant.EmailOrMailAddress.MailAddress.Address = idoc.Address
		XMLApplication.Entrant.EmailOrMailAddress.MailAddress.RegionID = idoc.Region
		XMLApplication.Entrant.EmailOrMailAddress.MailAddress.TownTypeID = idoc.CityType

		// FinSourceAndEduForm
		finSource := xmlFinSourceEduForm{}
		finSource.CompetitiveGroupUID = Configuration.Year + "-" + strconv.Itoa(orgCode) + "-" + app.Speciality
		// FinSourceAndEduForms
		XMLApplication.FinSourceAndEduForms = xmlFinSourceAndEduForms{}
		XMLApplication.FinSourceAndEduForms.FinSourceEduForm = finSource
		// AppDoc
		XMLApplication.ApplicationDocuments = xmlApplicationDocuments{}
		// Identity
		XMLApplication.ApplicationDocuments.IdentityDocument = xmlIdentityDocument{}
		XMLApplication.ApplicationDocuments.IdentityDocument.UID = idoc.UID
		XMLApplication.ApplicationDocuments.IdentityDocument.IdentityDocumentTypeID = idoc.DocType
		XMLApplication.ApplicationDocuments.IdentityDocument.OriginalReceivedDate = app.IdentityOriginalDate.Format("2006-01-02")
		XMLApplication.ApplicationDocuments.IdentityDocument.LastName = idoc.Lastname
		XMLApplication.ApplicationDocuments.IdentityDocument.FirstName = idoc.Firstname
		XMLApplication.ApplicationDocuments.IdentityDocument.MiddleName = idoc.Middlename
		XMLApplication.ApplicationDocuments.IdentityDocument.GenderID = idoc.Gender
		XMLApplication.ApplicationDocuments.IdentityDocument.DocumentSeries = idoc.Serial
		XMLApplication.ApplicationDocuments.IdentityDocument.DocumentNumber = idoc.Number
		XMLApplication.ApplicationDocuments.IdentityDocument.DocumentDate = idoc.When.Format("2006-01-02")
		XMLApplication.ApplicationDocuments.IdentityDocument.DocumentOrganization = idoc.Who
		XMLApplication.ApplicationDocuments.IdentityDocument.NationalityTypeID = idoc.Citizenship
		XMLApplication.ApplicationDocuments.IdentityDocument.BirthDate = idoc.BirthDate.Format("2006-01-02")
		// EduDoc
		XMLEDoc := xmlEducationDocument{}
		XMLEDoc.UID = edoc.UID
		XMLEDoc.OriginalReceivedDate = app.EducationOriginalDate.Format("2006-01-02")
		XMLEDoc.DocumentDate = edoc.When.Format("2006-01-02")
		XMLEDoc.DocumentNumber = edoc.Number
		XMLEDoc.Organization = edoc.Who
		XMLEDoc.GPA = edoc.GPA
		XMLEDoc.XMLName.Local = edoc.DocType
		// EDoccontaner
		XMLEDocContainer := xmlEduDocument{}
		XMLEDocContainer.Document = XMLEDoc
		// EDocsContainer
		XMLApplication.ApplicationDocuments.EduDocuments = xmlEduDocuments{}
		XMLApplication.ApplicationDocuments.EduDocuments.EduDocuments = append(XMLApplication.ApplicationDocuments.EduDocuments.EduDocuments, XMLEDocContainer)
		// AppendToRoot
		root.PackageData.Applications.Applications = append(root.PackageData.Applications.Applications, XMLApplication)
	}
	return &root
}

// writeTree выгружает дерево в файл
func writeTree(filename string, tree *xmlRoot) bool {
	file, _ := os.Create(filename)
	xmlWriter := io.Writer(file)
	xmlWriter.Write([]byte(xml.Header))
	xmlWriter.Write([]byte("<!--  -->\n")) // Для номера пакета
	enc := xml.NewEncoder(xmlWriter)
	enc.Indent("  ", "    ")
	if err := enc.Encode(tree); err != nil {
		fmt.Printf("Ошибка выгрузки XML файла: %v\n", err)
		return false
	}
	return true
}

// writeAppUIDS записывает в файл XLSX заявления, названный FILE.xlsx -> FILE_apps.xlsx
func writeAppUIDS(filename string) bool {
	filenameOut := regexp.MustCompile(`\.xlsx`).ReplaceAllString(filename, "_apps.xlsx")
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("Заявления")
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	var row *xlsx.Row
	var cell *xlsx.Cell
	for _, app := range domainDocumentsApplications {
		row = sheet.AddRow()
		cell = row.AddCell() // ApplicationNumber
		cell.Value = app.Number
		cell = row.AddCell() // RegistrationDate
		cell.Value = app.Date.Format("2006-01-02") + "T10:00:01+00:00"
	}
	err = file.Save(filenameOut)
	if err != nil {
		fmt.Printf(err.Error())
		return false
	}
	return true
}

// readAppUIDS читает данные и возвращает массив с данными
func readAppUIDS(filename string) (map[string]string, error) {
	filenameIn := regexp.MustCompile(`\.xlsx`).ReplaceAllString(filename, "_apps.xlsx")
	if !fileExist(filenameIn) {
		return nil, fmt.Errorf("Ошибка! Файла " + filenameIn + " с данными заявлений не существует!")
	}
	data := ReadsMapOfStrings(filenameIn)
	return data, nil
}
