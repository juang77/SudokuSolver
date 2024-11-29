package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Struct para manejar el JSON con la matriz de Sudoku
type SudokuRequest struct {
	Sudoku [][]int `json:"sudoku"`
}

// Valida si un Sudoku está correctamente resuelto
func isSudokuSolved(sudoku [][]int) bool {
	// Función para validar una lista de 9 elementos
	isValidSet := func(nums []int) bool {
		seen := make(map[int]bool)
		for _, num := range nums {
			if num < 1 || num > 9 || seen[num] {
				return false
			}
			seen[num] = true
		}
		return true
	}

	// Validar filas
	for _, row := range sudoku {
		if !isValidSet(row) {
			return false
		}
	}

	// Validar columnas
	for col := 0; col < 9; col++ {
		column := make([]int, 9)
		for row := 0; row < 9; row++ {
			column[row] = sudoku[row][col]
		}
		if !isValidSet(column) {
			return false
		}
	}

	// Validar subcuadrículas de 3x3
	for startRow := 0; startRow < 9; startRow += 3 {
		for startCol := 0; startCol < 9; startCol += 3 {
			subgrid := make([]int, 0, 9)
			for row := startRow; row < startRow+3; row++ {
				for col := startCol; col < startCol+3; col++ {
					subgrid = append(subgrid, sudoku[row][col])
				}
			}
			if !isValidSet(subgrid) {
				return false
			}
		}
	}

	return true
}

// Resuelve un Sudoku utilizando backtracking
func solveSudoku(sudoku [][]int) bool {
	// Encuentra una celda vacía
	findEmpty := func() (int, int, bool) {
		for row := 0; row < 9; row++ {
			for col := 0; col < 9; col++ {
				if sudoku[row][col] == 0 {
					return row, col, true
				}
			}
		}
		return -1, -1, false
	}

	isValid := func(row, col, num int) bool {
		// Validar la fila
		for x := 0; x < 9; x++ {
			if sudoku[row][x] == num {
				return false
			}
		}
		// Validar la columna
		for x := 0; x < 9; x++ {
			if sudoku[x][col] == num {
				return false
			}
		}
		// Validar la subcuadrícula de 3x3
		startRow, startCol := (row/3)*3, (col/3)*3
		for x := startRow; x < startRow+3; x++ {
			for y := startCol; y < startCol+3; y++ {
				if sudoku[x][y] == num {
					return false
				}
			}
		}
		return true
	}

	row, col, found := findEmpty()
	if !found {
		return true // No hay más celdas vacías, Sudoku resuelto
	}

	for num := 1; num <= 9; num++ {
		if isValid(row, col, num) {
			sudoku[row][col] = num
			if solveSudoku(sudoku) {
				return true
			}
			sudoku[row][col] = 0
		}
	}

	return false
}

// Handler para el endpoint GET
func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Hello World"))
}

// Handler para el endpoint POST
func sudokuHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var request SudokuRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&request); err != nil {
		http.Error(w, "Error al procesar el JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validación básica del tamaño de la matriz
	if len(request.Sudoku) != 9 {
		http.Error(w, "La matriz debe tener exactamente 9 filas", http.StatusBadRequest)
		return
	}
	for _, row := range request.Sudoku {
		if len(row) != 9 {
			http.Error(w, "Cada fila de la matriz debe tener exactamente 9 columnas", http.StatusBadRequest)
			return
		}
	}

	// Validar si el Sudoku está resuelto
	isSolved := isSudokuSolved(request.Sudoku)
	var solution [][]int
	var solutionFound bool

	if !isSolved {
		// Hacer una copia de la matriz para no modificar la original
		solution = make([][]int, 9)
		for i := range request.Sudoku {
			solution[i] = make([]int, 9)
			copy(solution[i], request.Sudoku[i])
		}

		// Intentar resolver el Sudoku
		solutionFound = solveSudoku(solution)
	}

	response := map[string]interface{}{
		"message": "Validación de Sudoku completada",
		"solved":  isSolved,
		"status":  "Sudoku solucionado",
	}

	if !isSolved {
		if solutionFound {
			response["status"] = "Sudoku no solucionado, solución generada"
			response["solution"] = solution
		} else {
			response["status"] = "Sudoku no solucionado, sin solución válida"
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/hello", helloWorldHandler)
	http.HandleFunc("/sudoku", sudokuHandler)

	port := "8080"
	fmt.Printf("Servidor escuchando en http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error iniciando el servidor: %s", err)
	}
}
