package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type ICDCode struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Result struct {
	Total      int       `json:"total"`
	Procedures []ICDCode `json:"procedures"`
}

func main() {
	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		term := strings.ToLower(r.URL.Query().Get("terms"))
		maxListStr := r.URL.Query().Get("maxList")

		maxList := 20 // default value
		if maxListStr != "" {
			var err error
			maxList, err = strconv.Atoi(maxListStr)
			if err != nil {
				http.Error(w, "maxList must be an integer", http.StatusBadRequest)
				return
			}
		}
		if maxList > 500 {
			maxList = 500
		}

		file, err := os.Open("icd10pcs_codes_2023.txt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		var allMatches []ICDCode
		for scanner.Scan() {
			line := scanner.Text()
			splitIndex := strings.Index(line, " ")
			code := ICDCode{
				ID:   line[:splitIndex],
				Name: fmt.Sprintf("%s %s", line[:splitIndex], line[splitIndex+1:]),
			}
			if strings.Contains(strings.ToLower(code.ID), term) || strings.Contains(strings.ToLower(code.Name), term) {
				allMatches = append(allMatches, code)
			}
		}

		if err := scanner.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		result := Result{
			Total:      len(allMatches),
			Procedures: allMatches[:min(maxList, len(allMatches))],
		}

		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
