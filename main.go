package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func main() {
	showHelp := flag.Bool("help", false, "Menampilkan bantuan")
	flag.Parse()
	if *showHelp {
		showUsage()
		return
	}

	args := os.Args

	if len(args) < 4 {
		fmt.Println("Mohon berikan argumen pilih fungsi dan input file")
		fmt.Println("--help untuk bantuan")
		os.Exit(1)
	}

	selectFunc := args[1]
	inputFile := args[2]
	outputFileName := args[3]

	if selectFunc == "" || inputFile == "" || outputFileName == "" {
		fmt.Println("Mohon berikan argumen selectFunc dan inputFile")
		os.Exit(1)
	}

	switch selectFunc {
	case "tojson":
		toJson(inputFile, outputFileName)
	case "tocsv":
		toCsv(inputFile, outputFileName)
	}

}

func toJson(inputFile, outputFileName string) {
		// Buka file CSV
		csvFile, err := os.Open(inputFile)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer csvFile.Close()
	
		// Membaca data CSV
		reader := csv.NewReader(csvFile)
		csvData, err := reader.ReadAll()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	
		// Mengonversi data CSV menjadi JSON
		jsonData := make([]map[string]string, 0)
		headers := csvData[0]
		for _, row := range csvData[1:] {
			record := make(map[string]string)
			for i, value := range row {
				record[headers[i]] = value
			}
			jsonData = append(jsonData, record)
		}
	
		// Membuka file JSON
		jsonFile, err := os.Create(outputFileName)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer jsonFile.Close()
	
		// Menulis data JSON ke file
		jsonEncoder := json.NewEncoder(jsonFile)
		jsonEncoder.SetIndent("", "  ")
		err = jsonEncoder.Encode(jsonData)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	
		fmt.Println("Konversi CSV ke JSON selesai.")
}

func toCsv(inputFile, outputFileName string) {
	// Buka file JSON
	jsonFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer jsonFile.Close()

	// Membaca data JSON
	jsonData := make([]map[string]interface{}, 0)
	jsonDecoder := json.NewDecoder(jsonFile)
	err = jsonDecoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Mendapatkan kolom header dari data JSON
	headers := make([]string, 0)
	if len(jsonData) > 0 {
		for key := range jsonData[0] {
			headers = append(headers, key)
		}
	}

	// Membuka file CSV
	csvFile, err := os.Create(outputFileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvFile.Close()

	// Menulis header ke file CSV
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write(headers)

	// Menulis data JSON ke file CSV
	for _, record := range jsonData {
		row := make([]string, 0)
		for _, key := range headers {
			value, ok := record[key]
			if ok {
				row = append(row, fmt.Sprintf("%v", value))
			} else {
				row = append(row, "")
			}
		}
		csvWriter.Write(row)
	}

	csvWriter.Flush()

	fmt.Println("Konversi JSON ke CSV selesai.")
}

// Menampilkan bantuan
func showUsage() {
	fmt.Fprintln(os.Stderr, "\nOpsi:")
	flag.PrintDefaults()
	fmt.Fprintln(os.Stderr, "\nDeskripsi:")
	fmt.Fprintln(os.Stderr, "Script untuk merubah data json ke csv dan sebaliknya.")
	fmt.Fprintln(os.Stderr, "")
	fmt.Fprintln(os.Stderr, "Contoh penggunaan:")
	fmt.Fprintln(os.Stderr, "pilih-fungsi: ['tojson', 'tocsv']")
	fmt.Fprintln(os.Stderr, "input-file: nama file")
	fmt.Fprintf(os.Stderr, "go run <path> pilih-fungsi input-file")
}
