package server

import (
	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
)

// fieldError describes a single validation failure.
type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// validateCreateJobRequest checks all required fields and enum values.
// Returns a non-empty slice when the request is invalid.
func validateCreateJobRequest(req createJobRequest) []fieldError {
	var errs []fieldError

	if req.LivingAreaM2 <= 0 {
		errs = append(errs, fieldError{"living_area_m2", "must be greater than 0"})
	}

	if req.ConstructionYear < 1800 || req.ConstructionYear > 2100 {
		errs = append(errs, fieldError{"construction_year", "must be between 1800 and 2100"})
	}

	if !validBuildingType(req.BuildingType) {
		errs = append(errs, fieldError{"building_type", "must be one of: detached_house, semi_detached, apartment, commercial"})
	}

	if !validInsulationLevel(req.InsulationLevel) {
		errs = append(errs, fieldError{"insulation_level", "must be one of: poor, average, good, excellent"})
	}

	if !validHeatingSystem(req.CurrentHeatingSystem) {
		errs = append(errs, fieldError{"current_heating_system", "must be one of: gas_boiler, oil_boiler, electric, district_heating"})
	}

	return errs
}

func validBuildingType(v jobs.BuildingType) bool {
	switch v {
	case jobs.BuildingTypeDetachedHouse, jobs.BuildingTypeSemiDetached,
		jobs.BuildingTypeApartment, jobs.BuildingTypeCommercial:
		return true
	}
	return false
}

func validInsulationLevel(v jobs.InsulationLevel) bool {
	switch v {
	case jobs.InsulationLevelPoor, jobs.InsulationLevelAverage,
		jobs.InsulationLevelGood, jobs.InsulationLevelExcellent:
		return true
	}
	return false
}

func validHeatingSystem(v jobs.CurrentHeatingSystem) bool {
	switch v {
	case jobs.HeatingSystemGasBoiler, jobs.HeatingSystemOilBoiler,
		jobs.HeatingSystemElectric, jobs.HeatingSystemDistrictHeat:
		return true
	}
	return false
}
