package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/phpdave11/gofpdf"
)

var tmpl = template.Must(template.ParseFiles("templates/form.tmpl"))
var resultTmpl = template.Must(template.ParseFiles("templates/result.tmpl"))

// PatientRecord struct
type PatientRecord struct {
	ID                  string
	FirstName           string
	MiddleName          string
	Surname             string
	DOB                 string
	Gender              string
	Phone               string
	Email               string
	Address             string
	NOKName             string
	NOKRelationship     string
	NOKContact          string
	MaritalStatus       string
	InsuranceProvider   string
	InsuranceMemberNo   string
	ClinicalHistory     string
	Allergies           string
	Assessments         string
	TreatmentPlan       string
	Medication          string
	Referrals           string
	TestResults         string
	ConsultationDate    string
	ConsultationPlace   string
	Reaction            string
}

var lastRecord PatientRecord

func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	lastRecord = PatientRecord{
		ID:                r.FormValue("id"),
		FirstName:         r.FormValue("first_name"),
		MiddleName:        r.FormValue("middle_name"),
		Surname:           r.FormValue("surname"),
		DOB:               r.FormValue("dob"),
		Gender:            r.FormValue("gender"),
		Phone:             r.FormValue("phone"),
		Email:             r.FormValue("email"),
		Address:           r.FormValue("address"),
		NOKName:           r.FormValue("nok_name"),
		NOKRelationship:   r.FormValue("nok_relationship"),
		NOKContact:        r.FormValue("nok_contact"),
		MaritalStatus:     r.FormValue("marital_status"),
		InsuranceProvider: r.FormValue("insurance_provider"),
		InsuranceMemberNo: r.FormValue("insurance_member_number"),
		ClinicalHistory:   r.FormValue("clinical_history"),
		Allergies:         r.FormValue("allergies"),
		Assessments:       r.FormValue("assessments"),
		TreatmentPlan:     r.FormValue("treatment_plan"),
		Medication:        r.FormValue("medication"),
		Referrals:         r.FormValue("referrals"),
		TestResults:       r.FormValue("test_results"),
		ConsultationDate:  r.FormValue("consultation_datetime"),
		ConsultationPlace: r.FormValue("consultation_place"),
		Reaction:          r.FormValue("reaction"),
	}

	resultTmpl.Execute(w, nil)
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Patient Record")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	fields := []struct{
		Label string
		Value string
	}{
		{"ID", lastRecord.ID},
		{"Full Name", lastRecord.FirstName + " " + lastRecord.MiddleName + " " + lastRecord.Surname},
		{"Date of Birth", lastRecord.DOB},
		{"Gender", lastRecord.Gender},
		{"Phone", lastRecord.Phone},
		{"Email", lastRecord.Email},
		{"Address", lastRecord.Address},
		{"Next of Kin", lastRecord.NOKName},
		{"Relationship", lastRecord.NOKRelationship},
		{"NOK Contact", lastRecord.NOKContact},
		{"Marital Status", lastRecord.MaritalStatus},
		{"Insurance Provider", lastRecord.InsuranceProvider},
		{"Member Number", lastRecord.InsuranceMemberNo},
		{"Clinical History", lastRecord.ClinicalHistory},
		{"Allergies", lastRecord.Allergies},
		{"Assessments", lastRecord.Assessments},
		{"Treatment Plan", lastRecord.TreatmentPlan},
		{"Medication", lastRecord.Medication},
		{"Referrals", lastRecord.Referrals},
		{"Test Results", lastRecord.TestResults},
		{"Consultation Date & Time", lastRecord.ConsultationDate},
		{"Consultation Place", lastRecord.ConsultationPlace},
		{"Reaction to Treatment", lastRecord.Reaction},
	}

	for _, f := range fields {
		pdf.CellFormat(50, 8, f.Label+":", "0", 0, "L", false, 0, "")
		pdf.MultiCell(0, 8, f.Value, "0", "L", false)
		pdf.Ln(2)
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "inline; filename=patient_record.pdf")
	err := pdf.Output(w)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/finish", finishHandler)
	http.HandleFunc("/pdf", pdfHandler)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
