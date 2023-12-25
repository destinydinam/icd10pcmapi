package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type ICDCode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("terms")

		file, err := os.Open("icd10pcs_codes_2023.txt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var matches []ICDCode
		for scanner.Scan() {
			line := scanner.Text()
			splitIndex := strings.Index(line, " ")
			code := ICDCode{
				ID:   line[:splitIndex],
				Name: line[splitIndex+1:],
			}
			if strings.Contains(strings.ToLower(code.ID), strings.ToLower(term)) || strings.Contains(strings.ToLower(code.Name), strings.ToLower(term)) {
				matches = append(matches, code)
			}
		}

		if err := scanner.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(matches)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
