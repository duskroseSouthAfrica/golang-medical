package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/phpdave11/gofpdf"
)

// Use your existing templates and add the search template
var tmpl = template.Must(template.ParseFiles("templates/form.tmpl"))
var resultTmpl = template.Must(template.ParseFiles("templates/result.tmpl"))

// PatientRecord struct - enhanced with storage fields
type PatientRecord struct {
	ID                string    `json:"id"`
	FirstName         string    `json:"first_name"`
	MiddleName        string    `json:"middle_name"`
	Surname           string    `json:"surname"`
	DOB               string    `json:"dob"`
	Gender            string    `json:"gender"`
	Phone             string    `json:"phone"`
	Email             string    `json:"email"`
	Address           string    `json:"address"`
	NOKName           string    `json:"nok_name"`
	NOKRelationship   string    `json:"nok_relationship"`
	NOKContact        string    `json:"nok_contact"`
	MaritalStatus     string    `json:"marital_status"`
	InsuranceProvider string    `json:"insurance_provider"`
	InsuranceMemberNo string    `json:"insurance_member_number"`
	ClinicalHistory   string    `json:"clinical_history"`
	Allergies         string    `json:"allergies"`
	Assessments       string    `json:"assessments"`
	TreatmentPlan     string    `json:"treatment_plan"`
	Medication        string    `json:"medication"`
	Referrals         string    `json:"referrals"`
	TestResults       string    `json:"test_results"`
	ConsultationDate  string    `json:"consultation_datetime"`
	ConsultationPlace string    `json:"consultation_place"`
	Reaction          string    `json:"patient_reaction"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

var lastRecord PatientRecord

// A struct to hold the data for the search template
type SearchPageData struct {
	Query   string
	Results []*PatientRecord
}

// Simple storage functions (no changes needed)
func savePatientRecord(record *PatientRecord) error {
	// ... your existing save function
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data/patients", 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Generate ID if not provided
	if record.ID == "" {
		record.ID = fmt.Sprintf("PATIENT_%d", time.Now().UnixNano())
	}

	record.UpdatedAt = time.Now()
	if record.CreatedAt.IsZero() {
		record.CreatedAt = record.UpdatedAt
	}

	// Save as JSON file
	filename := filepath.Join("data/patients", record.ID+".json")
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(record)
}

func loadPatientRecord(id string) (*PatientRecord, error) {
	// ... your existing load function
	filename := filepath.Join("data/patients", id+".json")
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("patient not found: %v", err)
	}
	defer file.Close()

	var record PatientRecord
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&record)
	if err != nil {
		return nil, fmt.Errorf("error reading patient data: %v", err)
	}

	return &record, nil
}

// Your existing handlers - no changes needed
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	// ... your existing finish handler logic to create and save the record

	record := PatientRecord{
		FirstName:         r.FormValue("first_name"),
		MiddleName:        r.FormValue("middle_name"),
		Surname:           r.FormValue("surname"),
		DOB:               r.FormValue("dob"),
		Gender:            r.FormValue("gender"),
		Phone:             r.FormValue("phone"),
		Email:             r.FormValue("email"),
		Address:           r.FormValue("address"),
		NOKName:           r.FormValue("next_of_kin"),
		NOKRelationship:   r.FormValue("relationship"),
		NOKContact:        r.FormValue("kin_contact"),
		MaritalStatus:     r.FormValue("marital_status"),
		InsuranceProvider: r.FormValue("insurance_provider"),
		InsuranceMemberNo: r.FormValue("member_number"),
		ClinicalHistory:   r.FormValue("clinical_history"),
		Allergies:         r.FormValue("allergies"),
		Assessments:       r.FormValue("assessments"),
		TreatmentPlan:     r.FormValue("treatment_plan"),
		Medication:        r.FormValue("medication"),
		Referrals:         r.FormValue("referrals"),
		TestResults:       r.FormValue("test_results"),
		ConsultationDate:  r.FormValue("consultation_datetime"),
		ConsultationPlace: r.FormValue("consultation_place"),
		Reaction:          r.FormValue("patient_reaction"),
	}

	// Save to file
	if err := savePatientRecord(&record); err != nil {
		log.Printf("Error saving patient record: %v", err)
		http.Error(w, "Error saving patient record", http.StatusInternalServerError)
		return
	}

	// Keep for PDF generation (as before)
	lastRecord = record
	log.Printf("Saved patient record with ID: %s", record.ID)

	resultTmpl.Execute(w, nil)
}

func pdfHandler(w http.ResponseWriter, r *http.Request) {
	// ... your existing PDF handler logic
	patientID := r.URL.Query().Get("id")
	var record PatientRecord

	if patientID != "" {
		// Load specific patient
		loadedRecord, err := loadPatientRecord(patientID)
		if err != nil {
			http.Error(w, "Patient not found", http.StatusNotFound)
			return
		}
		record = *loadedRecord
	} else {
		// Use last record (as before)
		record = lastRecord
	}

	// Your existing PDF generation code
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Patient Record")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	fields := []struct {
		Label string
		Value string
	}{
		{"ID", record.ID},
		{"Full Name", record.FirstName + " " + record.MiddleName + " " + record.Surname},
		{"Date of Birth", record.DOB},
		{"Gender", record.Gender},
		{"Phone", record.Phone},
		{"Email", record.Email},
		{"Address", record.Address},
		{"Next of Kin", record.NOKName},
		{"Relationship", record.NOKRelationship},
		{"NOK Contact", record.NOKContact},
		{"Marital Status", record.MaritalStatus},
		{"Insurance Provider", record.InsuranceProvider},
		{"Member Number", record.InsuranceMemberNo},
		{"Clinical History", record.ClinicalHistory},
		{"Allergies", record.Allergies},
		{"Assessments", record.Assessments},
		{"Treatment Plan", record.TreatmentPlan},
		{"Medication", record.Medication},
		{"Referrals", record.Referrals},
		{"Test Results", record.TestResults},
		{"Consultation Date & Time", record.ConsultationDate},
		{"Consultation Place", record.ConsultationPlace},
		{"Reaction to Treatment", record.Reaction},
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

	// Update this line to use the new searchHandler

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", nil))
}