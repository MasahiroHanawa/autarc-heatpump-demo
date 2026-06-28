package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
)

// JobsHandler holds the HTTP handlers for the jobs API.
type JobsHandler struct {
	svc *jobs.Service
}

func NewJobsHandler(svc *jobs.Service) *JobsHandler {
	return &JobsHandler{svc: svc}
}

// — request / response types —————————————————————————————————————————————————

type createJobRequest struct {
	BuildingType               jobs.BuildingType         `json:"building_type"`
	LivingAreaM2               float64                   `json:"living_area_m2"`
	ConstructionYear           int                       `json:"construction_year"`
	InsulationLevel            jobs.InsulationLevel      `json:"insulation_level"`
	CurrentHeatingSystem       jobs.CurrentHeatingSystem `json:"current_heating_system"`
	AnnualEnergyConsumptionKwh float64                   `json:"annual_energy_consumption_kwh"`
	InstallerNotes             string                    `json:"installer_notes"`
}

type jobResponse struct {
	ID        string          `json:"id"`
	Status    jobs.Status     `json:"status"`
	Input     inputResponse   `json:"input"`
	Output    *outputResponse `json:"output"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type inputResponse struct {
	BuildingType               jobs.BuildingType         `json:"building_type"`
	LivingAreaM2               float64                   `json:"living_area_m2"`
	ConstructionYear           int                       `json:"construction_year"`
	InsulationLevel            jobs.InsulationLevel      `json:"insulation_level"`
	CurrentHeatingSystem       jobs.CurrentHeatingSystem `json:"current_heating_system"`
	AnnualEnergyConsumptionKwh float64                   `json:"annual_energy_consumption_kwh"`
	InstallerNotes             string                    `json:"installer_notes"`
}

type outputResponse struct {
	EstimatedHeatDemandKwh  float64          `json:"estimated_heat_demand_kwh"`
	RecommendedHeatPumpType jobs.HeatPumpType `json:"recommended_heat_pump_type"`
	Suitability             jobs.Suitability  `json:"suitability"`
	Confidence              float64           `json:"confidence"`
	RiskFlags               []string          `json:"risk_flags"`
	NextSteps               []string          `json:"next_steps"`
	Summary                 string            `json:"summary"`
}

// — handlers ——————————————————————————————————————————————————————————————————

// Create handles POST /api/jobs
func (h *JobsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createJobRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errs := validateCreateJobRequest(req); len(errs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"errors": errs})
		return
	}

	job, err := h.svc.CreateJob(r.Context(), jobs.Input{
		BuildingType:               req.BuildingType,
		LivingAreaM2:               req.LivingAreaM2,
		ConstructionYear:           req.ConstructionYear,
		InsulationLevel:            req.InsulationLevel,
		CurrentHeatingSystem:       req.CurrentHeatingSystem,
		AnnualEnergyConsumptionKwh: req.AnnualEnergyConsumptionKwh,
		InstallerNotes:             req.InstallerNotes,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create job")
		return
	}

	writeJSON(w, http.StatusCreated, toJobResponse(job))
}

// Get handles GET /api/jobs/{id}
func (h *JobsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	job, err := h.svc.GetJob(r.Context(), id)
	if err != nil {
		if errors.Is(err, jobs.ErrNotFound) {
			writeError(w, http.StatusNotFound, "job not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to get job")
		return
	}

	writeJSON(w, http.StatusOK, toJobResponse(job))
}

// List handles GET /api/jobs
func (h *JobsHandler) List(w http.ResponseWriter, r *http.Request) {
	all, err := h.svc.ListJobs(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list jobs")
		return
	}

	resp := make([]jobResponse, len(all))
	for i, j := range all {
		resp[i] = toJobResponse(j)
	}

	writeJSON(w, http.StatusOK, resp)
}

// — mapping ———————————————————————————————————————————————————————————————————

func toJobResponse(j jobs.Job) jobResponse {
	resp := jobResponse{
		ID:     j.ID,
		Status: j.Status,
		Input: inputResponse{
			BuildingType:               j.Input.BuildingType,
			LivingAreaM2:               j.Input.LivingAreaM2,
			ConstructionYear:           j.Input.ConstructionYear,
			InsulationLevel:            j.Input.InsulationLevel,
			CurrentHeatingSystem:       j.Input.CurrentHeatingSystem,
			AnnualEnergyConsumptionKwh: j.Input.AnnualEnergyConsumptionKwh,
			InstallerNotes:             j.Input.InstallerNotes,
		},
		CreatedAt: j.CreatedAt,
		UpdatedAt: j.UpdatedAt,
	}

	if j.Output != nil {
		o := outputResponse{
			EstimatedHeatDemandKwh:  j.Output.EstimatedHeatDemandKwh,
			RecommendedHeatPumpType: j.Output.RecommendedHeatPumpType,
			Suitability:             j.Output.Suitability,
			Confidence:              j.Output.Confidence,
			RiskFlags:               j.Output.RiskFlags,
			NextSteps:               j.Output.NextSteps,
			Summary:                 j.Output.Summary,
		}
		resp.Output = &o
	}

	return resp
}
