// Package jobs defines the core domain model for heat pump analysis jobs.
package jobs

import "time"

// Status represents the lifecycle state of an analysis job.
type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusCompleted  Status = "completed"
	StatusFailed     Status = "failed"
)

// BuildingType classifies the type of building being assessed.
type BuildingType string

const (
	BuildingTypeDetachedHouse BuildingType = "detached_house"
	BuildingTypeSemiDetached  BuildingType = "semi_detached"
	BuildingTypeApartment     BuildingType = "apartment"
	BuildingTypeCommercial    BuildingType = "commercial"
)

// InsulationLevel describes the thermal insulation quality of the building.
type InsulationLevel string

const (
	InsulationLevelPoor      InsulationLevel = "poor"
	InsulationLevelAverage   InsulationLevel = "average"
	InsulationLevelGood      InsulationLevel = "good"
	InsulationLevelExcellent InsulationLevel = "excellent"
)

// CurrentHeatingSystem describes the existing heating system being replaced.
type CurrentHeatingSystem string

const (
	HeatingSystemGasBoiler      CurrentHeatingSystem = "gas_boiler"
	HeatingSystemOilBoiler      CurrentHeatingSystem = "oil_boiler"
	HeatingSystemElectric       CurrentHeatingSystem = "electric"
	HeatingSystemDistrictHeat   CurrentHeatingSystem = "district_heating"
)

// HeatPumpType is the recommended heat pump category from the analysis.
type HeatPumpType string

const (
	HeatPumpTypeAirSource    HeatPumpType = "air_source"
	HeatPumpTypeGroundSource HeatPumpType = "ground_source"
	HeatPumpTypeWaterSource  HeatPumpType = "water_source"
)

// Suitability indicates how suitable the building is for a heat pump installation.
type Suitability string

const (
	SuitabilityHigh       Suitability = "high"
	SuitabilityMedium     Suitability = "medium"
	SuitabilityLow        Suitability = "low"
	SuitabilityNotSuitable Suitability = "not_suitable"
)

// Input holds the building and site-survey data provided by the installer.
type Input struct {
	BuildingType               BuildingType
	LivingAreaM2               float64
	ConstructionYear           int
	InsulationLevel            InsulationLevel
	CurrentHeatingSystem       CurrentHeatingSystem
	AnnualEnergyConsumptionKwh float64
	InstallerNotes             string
}

// Output holds the structured analysis result produced by the analyzer.
// It is nil until the job reaches StatusCompleted.
type Output struct {
	EstimatedHeatDemandKwh  float64
	RecommendedHeatPumpType HeatPumpType
	Suitability             Suitability
	Confidence              float64 // 0.0–1.0
	RiskFlags               []string
	NextSteps               []string
	Summary                 string
}

// Job is the central aggregate that tracks an analysis request from creation
// through background processing to a final result.
type Job struct {
	ID        string
	Input     Input
	Output    *Output // nil until StatusCompleted
	Status    Status
	CreatedAt time.Time
	UpdatedAt time.Time
}
