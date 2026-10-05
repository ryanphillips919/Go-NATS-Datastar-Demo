package main

import (
	"html/template"
	"log"
	"net/http"
)

type Expense struct {
	ID          int
	Amount      float64
	Description string
}

// This expense lives in memory and resets when the server restarts.
var expense = Expense{
	ID:          1,
	Amount:      482.16,
	Description: "Materials",
}

var expensePage = template.Must(template.New("expense").Parse(`<!doctype html>
<html>
<body>
  <h1>Expense</h1>
  <p>ID: {{.ID}}</p>
  <p>Amount: {{printf "%.2f" .Amount}}</p>
  <form method="post" action="/expense">
    <label>Description: <input name="description" value="{{.Description}}"></label>
    <button type="submit">Save</button>
  </form>
</body>
</html>`))

func expenseHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
	case http.MethodPost:
		// Read the submitted form value and update the in-memory struct.
		expense.Description = r.FormValue("description")
	default:
		w.Header().Set("Allow", "GET, POST")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := expensePage.Execute(w, expense); err != nil {
		http.Error(w, "could not render expense", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/expense", expenseHandler)
	log.Println("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
