import type {
  BuildingType,
  CreateJobResponse,
  Job,
  JobInput,
  JobStatus,
  ListJobsResponse,
  Suitability,
} from '../types'

interface ApiJobSummary {
  id: string
  status: JobStatus
  building_type: BuildingType
  created_at: string
}

interface ApiListJobsResponse {
  jobs: ApiJobSummary[]
}

const BASE_URL = '/api'

// Raw shapes returned by the Go API (snake_case). Internal to this module.
interface ApiOutput {
  estimated_heat_demand_kwh: number
  recommended_heat_pump_type: string
  suitability: Suitability
  confidence: number
  risk_flags: string[]
  next_steps: string[]
  summary: string
}

interface ApiJob {
  id: string
  status: JobStatus
  output: ApiOutput | null
  created_at: string
  updated_at: string
}

function mapJob(api: ApiJob): Job {
  return {
    id: api.id,
    status: api.status,
    createdAt: api.created_at,
    updatedAt: api.updated_at,
    result: api.output
      ? {
          estimatedHeatDemandKwh: api.output.estimated_heat_demand_kwh,
          recommendedHeatPumpType: api.output.recommended_heat_pump_type,
          suitability: api.output.suitability,
          confidence: api.output.confidence,
          riskFlags: api.output.risk_flags ?? [],
          nextSteps: api.output.next_steps ?? [],
          summary: api.output.summary,
        }
      : undefined,
  }
}

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: { 'Content-Type': 'application/json', ...init?.headers },
    ...init,
  })

  if (!res.ok) {
    let message = `HTTP ${res.status}`
    try {
      const body = (await res.json()) as { error?: string }
      if (typeof body.error === 'string') message = body.error
    } catch {
      // use default message
    }
    throw new ApiError(res.status, message)
  }

  return res.json() as Promise<T>
}

export function createJob(input: JobInput): Promise<CreateJobResponse> {
  return request<CreateJobResponse>('/jobs', {
    method: 'POST',
    body: JSON.stringify({
      building_type: input.buildingType,
      living_area_m2: input.livingAreaM2,
      construction_year: input.constructionYear,
      insulation_level: input.insulationLevel,
      current_heating_system: input.currentHeatingSystem,
      annual_energy_consumption_kwh: input.annualEnergyConsumptionKwh,
      installer_notes: input.installerNotes ?? '',
    }),
  })
}

export async function getJob(id: string): Promise<Job> {
  const api = await request<ApiJob>(`/jobs/${id}`)
  return mapJob(api)
}

export async function listJobs(): Promise<ListJobsResponse> {
  const api = await request<ApiListJobsResponse>('/jobs')
  return {
    jobs: api.jobs.map((j) => ({
      id: j.id,
      status: j.status,
      buildingType: j.building_type,
      createdAt: j.created_at,
    })),
  }
}
