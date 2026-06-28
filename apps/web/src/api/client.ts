import type {
  CreateJobResponse,
  Job,
  JobInput,
  ListJobsResponse,
} from '../types'

const BASE_URL = '/api'

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
    body: JSON.stringify(input),
  })
}

export function getJob(id: string): Promise<Job> {
  return request<Job>(`/jobs/${id}`)
}

export function listJobs(): Promise<ListJobsResponse> {
  return request<ListJobsResponse>('/jobs')
}
