package nvo

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestSurgicalNoteModelUsesDatabaseColumnNamesInJSON(t *testing.T) {
	var note SurgicalNoteModel
	if err := json.Unmarshal([]byte(`{"patient_name":"Ada","note_type":"1"}`), &note); err != nil {
		t.Fatal(err)
	}
	if note.PatientName == nil || *note.PatientName != "Ada" {
		t.Fatalf("PatientName = %v, want Ada", note.PatientName)
	}
	if note.NoteType == nil || *note.NoteType != "1" {
		t.Fatalf("NoteType = %v, want 1", note.NoteType)
	}
}

func TestSurgicalNoteRoutesRejectInvalidIDsAndPagination(t *testing.T) {
	previousMode := gin.Mode()
	gin.SetMode(gin.TestMode)
	defer gin.SetMode(previousMode)
	router := gin.New()
	SurgicalNotesRegister(router.Group("/api/surgical-notes"))

	tests := []struct {
		path string
	}{
		{path: "/api/surgical-notes/0"},
		{path: "/api/surgical-notes/not-an-id"},
		{path: "/api/surgical-notes?limit=0"},
		{path: "/api/surgical-notes?offset=-1"},
	}
	for _, test := range tests {
		t.Run(test.path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, test.path, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d; body: %s", response.Code, http.StatusBadRequest, response.Body.String())
			}
		})
	}
}
