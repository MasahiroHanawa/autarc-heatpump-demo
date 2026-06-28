package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/server"
	"log/slog"
	"os"
)

// newTestServer builds a real router backed by an in-memory store.
func newTestServer(t *testing.T) http.Handler {
	t.Helper()
	store := jobs.NewInMemoryStore()
	svc := jobs.NewService(store)
	h := server.NewJobsHandler(svc)
	return server.NewRouter(slog.New(slog.NewTextHandler(os.Stdout, nil)), h)
}

func TestCreateJob_Valid(t *testing.T) {
	router := newTestServer(t)

	body := `{
		"building_type": "detached_house",
		"living_area_m2": 120,
		"construction_year": 1990,
		"insulation_level": "average",
		"current_heating_system": "gas_boiler"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)

	if resp["id"] == "" {
		t.Error("expected non-empty id")
	}
	if resp["status"] != "pending" {
		t.Errorf("expected status pending, got %s", resp["status"])
	}
	if resp["output"] != nil {
		t.Error("expected output to be null for a new job")
	}

	t.Logf("created job id=%s status=%s", resp["id"], resp["status"])
}

func TestCreateJob_InvalidBody(t *testing.T) {
	router := newTestServer(t)

	req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBufferString("not json"))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestCreateJob_ValidationErrors(t *testing.T) {
	router := newTestServer(t)

	// living_area_m2 = 0 and invalid enum values should both fail
	body := `{
		"building_type": "castle",
		"living_area_m2": 0,
		"construction_year": 1990,
		"insulation_level": "average",
		"current_heating_system": "gas_boiler"
	}`

	req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var resp map[string]any
	json.NewDecoder(rec.Body).Decode(&resp)

	errs, ok := resp["errors"].([]any)
	if !ok || len(errs) == 0 {
		t.Error("expected non-empty errors array")
	}

	t.Logf("got %d validation errors", len(errs))
}

func TestGetJob_NotFound(t *testing.T) {
	router := newTestServer(t)

	req := httptest.NewRequest(http.MethodGet, "/api/jobs/non-existent", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetJob_Found(t *testing.T) {
	router := newTestServer(t)

	// First create a job
	body := `{
		"building_type": "apartment",
		"living_area_m2": 60,
		"construction_year": 2010,
		"insulation_level": "good",
		"current_heating_system": "electric"
	}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBufferString(body))
	createReq.Header.Set("Content-Type", "application/json")
	createRec := httptest.NewRecorder()
	router.ServeHTTP(createRec, createReq)

	var created map[string]any
	json.NewDecoder(createRec.Body).Decode(&created)
	id := created["id"].(string)

	// Then fetch it
	getReq := httptest.NewRequest(http.MethodGet, "/api/jobs/"+id, nil)
	getRec := httptest.NewRecorder()
	router.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	var fetched map[string]any
	json.NewDecoder(getRec.Body).Decode(&fetched)

	if fetched["id"] != id {
		t.Errorf("expected id %s, got %s", id, fetched["id"])
	}

	t.Logf("fetched job id=%s status=%s", fetched["id"], fetched["status"])
}

func TestListJobs(t *testing.T) {
	router := newTestServer(t)

	body := `{
		"building_type": "detached_house",
		"living_area_m2": 100,
		"construction_year": 2000,
		"insulation_level": "poor",
		"current_heating_system": "oil_boiler"
	}`

	for range 3 {
		req := httptest.NewRequest(http.MethodPost, "/api/jobs", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(httptest.NewRecorder(), req)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/jobs", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var list []any
	json.NewDecoder(rec.Body).Decode(&list)

	if len(list) != 3 {
		t.Errorf("expected 3 jobs, got %d", len(list))
	}

	t.Logf("listed %d jobs", len(list))
}
