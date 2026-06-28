import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { ApiError, listJobs } from '../api/client'
import { JobStatusBadge } from '../components/JobStatus'
import type { JobSummary } from '../types'

export function JobList() {
  const [jobs, setJobs] = useState<JobSummary[]>([])
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    async function fetchJobs() {
      try {
        const data = await listJobs()
        const sorted = [...data.jobs].sort(
          (a, b) => new Date(b.createdAt).getTime() - new Date(a.createdAt).getTime(),
        )
        setJobs(sorted)
      } catch (err) {
        setError(
          err instanceof ApiError
            ? err.message
            : 'Failed to load jobs. Please try again.',
        )
      } finally {
        setIsLoading(false)
      }
    }

    void fetchJobs()
  }, [])

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center py-16 gap-3">
        <div className="w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin" />
        <p className="text-sm text-gray-500">Loading…</p>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-md p-4">
        <p className="text-sm text-red-700">{error}</p>
      </div>
    )
  }

  if (jobs.length === 0) {
    return (
      <div className="text-center py-16">
        <p className="text-gray-500 mb-4">No analyses yet.</p>
        <Link to="/" className="text-blue-600 hover:underline text-sm">
          Start your first analysis →
        </Link>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Job History</h2>
        <Link
          to="/"
          className="text-sm bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700"
        >
          New Analysis
        </Link>
      </div>

      <div className="space-y-2">
        {jobs.map((job) => (
          <Link
            key={job.id}
            to={`/jobs/${job.id}`}
            className="flex items-center justify-between bg-white border border-gray-200 rounded-lg px-4 py-3 hover:border-blue-300 hover:shadow-sm transition-all"
          >
            <div className="flex items-center gap-3">
              <JobStatusBadge status={job.status} />
              <span className="text-sm text-gray-700 capitalize">
                {job.buildingType.replace(/_/g, ' ')}
              </span>
            </div>
            <div className="flex items-center gap-4">
              <span className="text-xs text-gray-400 font-mono">{job.id}</span>
              <span className="text-xs text-gray-400">
                {new Date(job.createdAt).toLocaleString()}
              </span>
            </div>
          </Link>
        ))}
      </div>
    </div>
  )
}
