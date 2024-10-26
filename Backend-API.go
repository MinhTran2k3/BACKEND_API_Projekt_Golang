package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Struct für die Anfrage
type CalculationRequest struct {
	Funktionsweise string  `json:"Funktionsweise"` // z.B. "add", "subtract", etc.
	Zahl1          float64 `json:"Zahl1"`          // Erste Zahl
	Zahl2          float64 `json:"Zahl2"`          // Zweite Zahl
}

// Struct für die Antwort
type CalculationResponse struct {
	Ergebnis        float64 `json:"Ergebnis"`                  // Ergebnis der Berechnung
	Fehlernachricht string  `json:"Fehlernachricht,omitempty"` // Optional: Fehlernachricht
}

func saveCalculation(calcReq CalculationRequest, calcRes CalculationResponse) {
	// Lade vorhandene Daten
	var calculations []map[string]interface{}
	data, err := ioutil.ReadFile("berechnung.json")
	if err != nil {
		fmt.Println("Datei konnte nicht gelesen werden")
		return
	}
	json.Unmarshal(data, calculations)

	// Neues Ergebnis zur Liste hinzufügen
	NeueEingabe := map[string]interface{}{
		"Funktionsweise":  calcReq.Funktionsweise,
		"Zahl1":           calcReq.Zahl1,
		"Zahl2":           calcReq.Zahl2,
		"Ergebnis":        calcRes.Ergebnis,
		"Fehlernachricht": calcRes.Fehlernachricht,
	}
	calculations = append(calculations, NeueEingabe)
	updatedData, _ := json.Marshal(calculations)
	ioutil.WriteFile("berechnung.json", updatedData)

	fmt.Print("Daten wurden gespeichert")
}

// Handler für die Berechnung
func calculatorHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	var calcReq CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&calcReq)
	if err != nil {
		http.Error(w, "Fehlerhafte Anfrage", http.StatusBadRequest)
		return
	}

	var calcRes CalculationResponse

	switch calcReq.Funktionsweise {
	case "plus":
		calcRes.Ergebnis = calcReq.Zahl1 + calcReq.Zahl2 // Addition
	case "minus":
		calcRes.Ergebnis = calcReq.Zahl1 - calcReq.Zahl2 // Subtraktion
	case "mal":
		calcRes.Ergebnis = calcReq.Zahl1 * calcReq.Zahl2 // Multiplikation
	case "geteilt":
		if calcReq.Zahl2 == 0 {
			// Fehler, wenn durch 0 geteilt wird
			calcRes.Fehlernachricht = "Teilen durch 0 ist nicht erlaubt"
			http.Error(w, calcRes.Fehlernachricht, http.StatusBadRequest)
			return
		}
		calcRes.Ergebnis = calcReq.Zahl1 / calcReq.Zahl2 // Division
	default:

		calcRes.Fehlernachricht = "Ungültige Operation"
		http.Error(w, calcRes.Fehlernachricht, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calcRes)
}

func main() {

	http.HandleFunc("/taschenrechner", calculatorHandler)

	// Der Server wird auf Port 8080 gestartet
	fmt.Println("Server startet auf Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
