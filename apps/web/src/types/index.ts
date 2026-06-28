export type JobStatus = 'pending' | 'processing' | 'completed' | 'failed'

export type BuildingType =
  | 'detached_house'
  | 'semi_detached'
  | 'terraced'
  | 'apartment'
  | 'bungalow'
  | 'commercial'

export type InsulationLevel = 'poor' | 'moderate' | 'good' | 'excellent'

export type CurrentHeatingSystem =
  | 'gas_boiler'
  | 'oil_boiler'
  | 'electric'
  | 'district_heating'
  | 'none'

export type Suitability = 'excellent' | 'good' | 'fair' | 'poor' | 'unsuitable'

export interface JobInput {
  buildingType: BuildingType
  livingAreaM2: number
  constructionYear: number
  insulationLevel: InsulationLevel
  currentHeatingSystem: CurrentHeatingSystem
  annualEnergyConsumptionKwh: number
  installerNotes?: string
}

export interface JobResult {
  estimatedHeatDemandKwh: number
  recommendedHeatPumpType: string
  suitability: Suitability
  confidence: number
  riskFlags: string[]
  nextSteps: string[]
  summary: string
}

export interface Job {
  id: string
  status: JobStatus
  input?: JobInput
  result?: JobResult
  createdAt: string
  updatedAt?: string
}

export interface JobSummary {
  id: string
  status: JobStatus
  buildingType: BuildingType
  createdAt: string
}

export interface CreateJobResponse {
  id: string
  status: JobStatus
  createdAt: string
}

export interface ListJobsResponse {
  jobs: JobSummary[]
}
