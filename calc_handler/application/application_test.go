package application

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name           string
		expression     string
		expectedStatus int
		expectedBody   string
	}{
		{"ValidExpression", "3 + 5", http.StatusOK, "result: 8.000000"},
		{"EmptyExpression", "", http.StatusUnprocessableEntity, "error: empty expression or invalid request"},
		{"DivisionByZero", "3 / 0", http.StatusUnprocessableEntity, "error: division by zero"},
		{"InvalidExpression", "invalid", http.StatusUnprocessableEntity, "error: invalid expression"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(Request{Expression: tt.expression})
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			CalcHandler(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
