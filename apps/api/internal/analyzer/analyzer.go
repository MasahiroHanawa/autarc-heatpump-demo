// Package analyzer defines the Analyzer interface and a deterministic mock.
// A real OpenAI or Anthropic implementation can swap in later without
// changing any other package.
package analyzer

import (
	"context"

	"github.com/MasahiroHanawa/autarc-heatpump-demo/apps/api/internal/jobs"
)

// Analyzer produces a structured heat pump recommendation from job input.
type Analyzer interface {
	Analyze(ctx context.Context, input jobs.Input) (jobs.Output, error)
}

// Mock is a deterministic Analyzer that produces realistic-looking output
// without calling any external service. Safe to use in tests and local dev.
type Mock struct{}

func NewMock() *Mock { return &Mock{} }

func (m *Mock) Analyze(_ context.Context, input jobs.Input) (jobs.Output, error) {
	return jobs.Output{
		EstimatedHeatDemandKwh:  input.LivingAreaM2 * heatDemandFactor(input.InsulationLevel),
		RecommendedHeatPumpType: recommendType(input),
		Suitability:             suitability(input),
		Confidence:              confidence(input),
		RiskFlags:               riskFlags(input),
		NextSteps:               nextSteps(input),
		Summary:                 "Mock analysis complete. Building assessed based on provided survey data.",
	}, nil
}

func heatDemandFactor(level jobs.InsulationLevel) float64 {
	switch level {
	case jobs.InsulationLevelExcellent:
		return 50
	case jobs.InsulationLevelGood:
		return 70
	case jobs.InsulationLevelAverage:
		return 90
	default:
		return 120
	}
}

func recommendType(input jobs.Input) jobs.HeatPumpType {
	if input.LivingAreaM2 > 200 || input.BuildingType == jobs.BuildingTypeCommercial {
		return jobs.HeatPumpTypeGroundSource
	}
	return jobs.HeatPumpTypeAirSource
}

func suitability(input jobs.Input) jobs.Suitability {
	if input.InsulationLevel == jobs.InsulationLevelPoor && input.ConstructionYear < 1980 {
		return jobs.SuitabilityLow
	}
	if input.InsulationLevel == jobs.InsulationLevelGood || input.InsulationLevel == jobs.InsulationLevelExcellent {
		return jobs.SuitabilityHigh
	}
	return jobs.SuitabilityMedium
}

func confidence(input jobs.Input) float64 {
	if input.InstallerNotes != "" {
		return 0.90
	}
	return 0.75
}

func riskFlags(input jobs.Input) []string {
	var flags []string
	if input.InsulationLevel == jobs.InsulationLevelPoor {
		flags = append(flags, "Poor insulation may reduce heat pump efficiency")
	}
	if input.ConstructionYear < 1980 {
		flags = append(flags, "Older building may require radiator upgrades")
	}
	if input.CurrentHeatingSystem == jobs.HeatingSystemOilBoiler {
		flags = append(flags, "Oil tank removal required before installation")
	}
	return flags
}

func nextSteps(input jobs.Input) []string {
	steps := []string{"Schedule on-site survey", "Request last 12 months of energy bills"}
	if input.InsulationLevel == jobs.InsulationLevelPoor {
		steps = append(steps, "Obtain insulation upgrade quotes before proceeding")
	}
	steps = append(steps, "Check local planning permissions for outdoor unit")
	return steps
}
