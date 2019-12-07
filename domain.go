package main

import (
	"time"
)

var (
	domainDocumentsIdentity     []*IdentityDocument  /* Документы УЛ */
	domainDocumentsEducation    []*EducationDocument /* Документы об образовании */
	domainDocumentsApplications []*Application       /* Заявления */
)

func init() {
	domainDocumentsIdentity = make([]*IdentityDocument, 0)
	domainDocumentsEducation = make([]*EducationDocument, 0)
	domainDocumentsApplications = make([]*Application, 0)
}

// IdentityDocument документ удостоверяющий личность
type IdentityDocument struct {
	UID         string // UID
	Firstname   string // ФИО
	Lastname    string
	Middlename  string
	BirthDate   time.Time // Дата рождения
	BirthPlace  string    // Место рождения
	Gender      int       // Пол
	Region      int       // Адрес проживания
	CityType    int
	Address     string
	DocType     int       // Тип документа удостоверяющего личность
	Serial      string    // Серия
	Number      string    // Номер
	Citizenship int       // Гражданство
	When        time.Time // Дата выдачи
	Who         string    // Кем выдано
	DepCode     string    // Код подразделения
}

// EducationDocument Документ удостоверяющий личность
// Делаем его отдельным
type EducationDocument struct {
	UID              string            // UID
	DocType          string            // Тип
	Serial           string            // Серия
	Number           string            // Номер
	RegisterNumber   string            // Регистрационный номер
	When             time.Time         // Дата выдачи
	EndYear          int               // Год окончания
	Who              string            // Кем выдано
	GPA              float64           // Средний балл
	IdentityDocument *IdentityDocument // Ссылка на персональные данные
}

// Application Заявление на поступление
type Application struct {
	UID                   string             //UID
	OrgCode               int                // Код организации
	Date                  time.Time          // Дата подачи заявления
	Number                string             // Номер заявления (строка)
	Speciality            string             // Специальность
	IdentityDocument      *IdentityDocument  // Документ, удостоверяющий личность
	IdentityOriginalDate  time.Time          // Дата предоставления оригинала
	EducationDocument     *EducationDocument // Документ об образовании
	EducationOriginalDate time.Time          // Дата предоставления оригинала
	Recommended           bool               // Рекомендован/зачислен
}
