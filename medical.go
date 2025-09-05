package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/phpdave11/gofpdf"
)

// Use your existing templates
var tmpl = template.Must(template.ParseFiles("templates/form.tmpl"))
var resultTmpl = template.Must(template.ParseFiles("templates/result.tmpl"))

// PatientRecord struct - enhanced with storage fields
type PatientRecord struct {
	ID                  string    `json:"id"`
	FirstName           string    `json:"first_name"`
	MiddleName          string    `json:"middle_name"`
	Surname             string    `json:"surname"`
	DOB                 string    `json:"dob"`
	Gender              string    `json:"gender"`
	Phone               string    `json:"phone"`
	Email               string    `json:"email"`
	Address             string    `json:"address"`
	NOKName             string    `json:"nok_name"`
	NOKRelationship     string    `json:"nok_relationship"`
	NOKContact          string    `json:"nok_contact"`
	MaritalStatus       string    `json:"marital_status"`
	InsuranceProvider   string    `json:"insurance_provider"`
	InsuranceMemberNo   string    `json:"insurance_member_number"`
	ClinicalHistory     string    `json:"clinical_history"`
	Allergies           string    `json:"allergies"`
	Assessments         string    `json:"assessments"`
	TreatmentPlan       string    `json:"treatment_plan"`
	Medication          string    `json:"medication"`
	Referrals           string    `json:"referrals"`
	TestResults         string    `json:"test_results"`
	ConsultationDate    string    `json:"consultation_datetime"`
	ConsultationPlace   string    `json:"consultation_place"`
	Reaction            string    `json:"patient_reaction"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

var lastRecord PatientRecord

// Simple storage functions
func savePatientRecord(record *PatientRecord) error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data/patients", 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	// Generate ID if not provided
	if record.ID == "" {
		record.ID = fmt.Sprintf("PAT_%d", time.Now().UnixNano())
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

func getAllPatients() ([]*PatientRecord, error) {
	files, err := filepath.Glob(filepath.Join("data/patients", "*.json"))
	if err != nil {
		return nil, fmt.Errorf("error reading patient files: %v", err)
	}

	var records []*PatientRecord
	for _, file := range files {
		id := strings.TrimSuffix(filepath.Base(file), ".json")
		record, err := loadPatientRecord(id)
		if err != nil {
			log.Printf("Error loading patient %s: %v", id, err)
			continue
		}
		records = append(records, record)
	}

	return records, nil
}

// Your existing handlers - just enhanced to save data
func formHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func finishHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Parse error", http.StatusBadRequest)
		return
	}

	// Create patient record exactly as before
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
	// Check if we want a specific patient's PDF
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
	fields := []struct{
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

// Simple patient list handler
func listPatientsHandler(w http.ResponseWriter, r *http.Request) {
	patients, err := getAllPatients()
	if err != nil {
		http.Error(w, "Error loading patients", http.StatusInternalServerError)
		return
	}

	// Simple HTML page to list patients
	html := `<!DOCTYPE html>
<html>
<head>
    <title>Patient List</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        .patient { border: 1px solid #ddd; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .name { font-size: 1.2em; font-weight: bold; }
        .info { color: #666; margin: 5px 0; }
        .actions { margin-top: 10px; }
        .actions a { margin-right: 10px; padding: 5px 10px; background: #007bff; color: white; text-decoration: none; border-radius: 3px; }
        .nav { margin-bottom: 20px; }
        .nav a { margin-right: 15px; text-decoration: none; }
    </style>
</head>
<body>
    <div class="nav">
        <a href="/">← New Patient</a>
        <a href="/patients">All Patients</a>
    </div>
    <h1>Patient Records (` + fmt.Sprintf("%d", len(patients)) + ` total)</h1>`

	for _, patient := range patients {
		html += fmt.Sprintf(`
    <div class="patient">
        <div class="name">%s %s %s</div>
        <div class="info">ID: %s | DOB: %s | Phone: %s</div>
        <div class="info">Created: %s</div>
        <div class="actions">
            <a href="/pdf?id=%s" target="_blank">View PDF</a>
        </div>
    </div>`,
			patient.FirstName, patient.MiddleName, patient.Surname,
			patient.ID, patient.DOB, patient.Phone,
			patient.CreatedAt.Format("2006-01-02 15:04"),
			patient.ID)
	}

	html += `</body></html>`
	fmt.Fprint(w, html)
}

func main() {
	// Your existing routes
	http.HandleFunc("/", formHandler)
	http.HandleFunc("/finish", finishHandler)
	http.HandleFunc("/pdf", pdfHandler)

	// New simple route to list patients
	http.HandleFunc("/patients", listPatientsHandler)

	// Your existing static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("Server running on http://localhost:8080")
	fmt.Println("New features:")
	fmt.Println("  /patients - View all saved patient records")
	fmt.Println("  /pdf?id=<patient_id> - Generate PDF for specific patient")
	fmt.Println("Patient records will be saved in data/patients/ folder")
	
	log.Fatal(http.ListenAndServe(":8080", nil))
}