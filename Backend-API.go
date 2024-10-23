package main

import (
	"encoding/json"
	"fmt"
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

// Handler für die Berechnung
func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	// Überprüfen, ob die Methode POST ist
	if r.Method != http.MethodPost {
		http.Error(w, "Nur POST-Anfragen sind erlaubt", http.StatusMethodNotAllowed)
		return
	}

	// Anfrage-Body in CalculationRequest dekodieren
	var calcReq CalculationRequest
	err := json.NewDecoder(r.Body).Decode(&calcReq)
	if err != nil {
		http.Error(w, "Fehlerhafte Anfrage", http.StatusBadRequest)
		return
	}

	// Ergebnis berechnen basierend auf der Funktionsweise
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
		// Fehler, wenn die Operation nicht erkannt wird
		calcRes.Fehlernachricht = "Ungültige Operation"
		http.Error(w, calcRes.Fehlernachricht, http.StatusBadRequest)
		return
	}

	// Antwort als JSON zurückgeben
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(calcRes)
}

func main() {
	// Der Handler für die Berechnung wird festgelegt
	http.HandleFunc("/taschenrechner", calculatorHandler)

	// Der Server wird auf Port 8080 gestartet
	fmt.Println("Server startet auf Port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
