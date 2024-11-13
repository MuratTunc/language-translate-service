package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bregydoc/gtranslate"
	"github.com/pemistahl/lingua-go"
	"github.com/rs/cors"
)

// Language represents a language with a two-letter code and name.
type Language struct {
	Code string // e.g., "en" for English
	Name string // e.g., "English"
}

// Languages is a slice of available languages
var Languages = []Language{
	{Code: "af", Name: "Afrikaans"},
	{Code: "sq", Name: "Albanian"},
	{Code: "ar", Name: "Arabic"},
	{Code: "hy", Name: "Armenian"},
	{Code: "bn", Name: "Bengali"},
	{Code: "bs", Name: "Bosnian"},
	{Code: "ca", Name: "Catalan"},
	{Code: "hr", Name: "Croatian"},
	{Code: "cs", Name: "Czech"},
	{Code: "da", Name: "Danish"},
	{Code: "nl", Name: "Dutch"},
	{Code: "en", Name: "English"},
	{Code: "eo", Name: "Esperanto"},
	{Code: "et", Name: "Estonian"},
	{Code: "tl", Name: "Filipino"},
	{Code: "fi", Name: "Finnish"},
	{Code: "fr", Name: "French"},
	{Code: "de", Name: "German"},
	{Code: "el", Name: "Greek"},
	{Code: "gu", Name: "Gujarati"},
	{Code: "hi", Name: "Hindi"},
	{Code: "hu", Name: "Hungarian"},
	{Code: "is", Name: "Icelandic"},
	{Code: "id", Name: "Indonesian"},
	{Code: "it", Name: "Italian"},
	{Code: "ja", Name: "Japanese"},
	{Code: "jw", Name: "Javanese"},
	{Code: "km", Name: "Khmer"},
	{Code: "ko", Name: "Korean"},
	{Code: "la", Name: "Latin"},
	{Code: "lv", Name: "Latvian"},
	{Code: "lt", Name: "Lithuanian"},
	{Code: "mk", Name: "Macedonian"},
	{Code: "ml", Name: "Malayalam"},
	{Code: "mr", Name: "Marathi"},
	{Code: "my", Name: "Myanmar (Burmese)"},
	{Code: "ne", Name: "Nepali"},
	{Code: "no", Name: "Norwegian"},
	{Code: "pl", Name: "Polish"},
	{Code: "pt", Name: "Portuguese"},
	{Code: "pa", Name: "Punjabi"},
	{Code: "ro", Name: "Romanian"},
	{Code: "ru", Name: "Russian"},
	{Code: "sr", Name: "Serbian"},
	{Code: "si", Name: "Sinhala"},
	{Code: "sk", Name: "Slovak"},
	{Code: "sl", Name: "Slovenian"},
	{Code: "es", Name: "Spanish"},
	{Code: "su", Name: "Sundanese"},
	{Code: "sw", Name: "Swahili"},
	{Code: "sv", Name: "Swedish"},
	{Code: "ta", Name: "Tamil"},
	{Code: "te", Name: "Telugu"},
	{Code: "th", Name: "Thai"},
	{Code: "tr", Name: "Turkish"},
	{Code: "uk", Name: "Ukrainian"},
	{Code: "ur", Name: "Urdu"},
	{Code: "vi", Name: "Vietnamese"},
	{Code: "cy", Name: "Welsh"},
	{Code: "xh", Name: "Xhosa"},
	{Code: "yi", Name: "Yiddish"},
	{Code: "zu", Name: "Zulu"},
}

// GetLanguageCode returns the two-letter code for a given language name
func GetLanguageCode(languageName string) (string, bool) {
	for _, lang := range Languages {
		if lang.Name == languageName {
			return lang.Code, true
		}
	}
	return "", false
}

type TranslateRequest struct {
	Text string `json:"text"`
	To   string `json:"to"`
}

type TranslateResponse struct {
	TranslatedText string `json:"translatedText,omitempty"`
	Status         bool   `json:"status"`
	Message        string `json:"message"`
}

func TranslateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request TranslateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	translated, err := gtranslate.TranslateWithParams(request.Text, gtranslate.TranslationParams{
		From: "auto",
		To:   request.To,
	})
	if err != nil {
		sendErrorResponse(w, "Translation failed", http.StatusInternalServerError)
		return
	}

	response := TranslateResponse{
		TranslatedText: translated,
		Status:         true,
		Message:        "",
	}

	sendJSONResponse(w, response, http.StatusOK)
}

func GetLanguageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request TranslateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		sendErrorResponse(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var code string
	detector := lingua.NewLanguageDetectorBuilder().
		FromAllLanguages().
		WithLowAccuracyMode().
		Build()

	if detectedLang, exists := detector.DetectLanguageOf(request.Text); exists {
		fmt.Println(detectedLang)
		langName := detectedLang.String()
		if code, found := GetLanguageCode(langName); found {
			fmt.Println("Detected language code:", code)
		} else {
			fmt.Println("Language not found in code map.")
			sendErrorResponse(w, "Language Detection is failed", http.StatusInternalServerError)
		}
	}

	response := TranslateResponse{
		TranslatedText: "Detectlanguage",
		Status:         true,
		Message:        code,
	}

	sendJSONResponse(w, response, http.StatusOK)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	response := TranslateResponse{
		Status:  false,
		Message: message,
	}
	sendJSONResponse(w, response, statusCode)
}

func sendJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("Failed to encode JSON response: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/translate", TranslateHandler)
	mux.HandleFunc("/getLanguageCode", GetLanguageHandler)

	c := cors.Default().Handler(mux)

	fmt.Println("Starting server on http://serverip:80")
	log.Fatal(http.ListenAndServe(":80", c))
}
