import { useCallback, useEffect, useRef, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { ApiError, getJob } from '../api/client'
import { JobResultDisplay } from '../components/JobResult'
import { JobStatusBadge } from '../components/JobStatus'
import type { Job } from '../types'

const POLL_INTERVAL_MS = 2000
const TERMINAL_STATUSES = new Set<Job['status']>(['completed', 'failed'])

export function JobDetail() {
  const { id } = useParams<{ id: string }>()
  const [job, setJob] = useState<Job | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null)

  function stopPolling() {
    if (intervalRef.current !== null) {
      clearInterval(intervalRef.current)
      intervalRef.current = null
    }
  }

  const startPolling = useCallback(() => {
    if (!id) return

    async function fetchJob() {
      try {
        const data = await getJob(id!)
        setJob(data)
        setError(null)
        if (TERMINAL_STATUSES.has(data.status)) stopPolling()
      } catch (err) {
        setError(
          err instanceof ApiError
            ? err.message
            : 'Failed to load job. Please try again.',
        )
        stopPolling()
      } finally {
        setIsLoading(false)
      }
    }

    setIsLoading(true)
    setError(null)
    void fetchJob()
    intervalRef.current = setInterval(() => void fetchJob(), POLL_INTERVAL_MS)
  }, [id])

  useEffect(() => {
    startPolling()
    return stopPolling
  }, [startPolling])

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
        <p className="text-sm text-red-700 mb-3">{error}</p>
        <div className="flex gap-4">
          <button
            onClick={startPolling}
            className="text-sm text-red-600 underline hover:no-underline"
          >
            Retry
          </button>
          <Link to="/" className="text-sm text-red-600 underline hover:no-underline">
            ← New analysis
          </Link>
        </div>
      </div>
    )
  }

  if (!job) return null

  return (
    <div>
      <Link to="/" className="text-sm text-gray-500 hover:text-gray-700 mb-6 inline-block">
        ← New Analysis
      </Link>

      <div className="flex items-center gap-3 mb-2">
        <h2 className="text-2xl font-bold text-gray-900">Job Result</h2>
        <JobStatusBadge status={job.status} />
      </div>
      <p className="text-xs text-gray-400 font-mono mb-6">{job.id}</p>

      {(job.status === 'pending' || job.status === 'processing') && (
        <div className="flex items-center gap-3 bg-blue-50 border border-blue-200 rounded-md p-4 mb-6">
          <div className="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin shrink-0" />
          <p className="text-sm text-blue-700">
            {job.status === 'pending'
              ? 'Queued — waiting for the worker to pick up this job…'
              : 'Analysing building data…'}
          </p>
        </div>
      )}

      {job.status === 'failed' && (
        <div className="bg-red-50 border border-red-200 rounded-md p-4 mb-6">
          <p className="text-sm text-red-700">
            Analysis failed. Please submit a new job.
          </p>
        </div>
      )}

      {job.status === 'completed' && job.result && (
        <JobResultDisplay result={job.result} />
      )}

      {job.status === 'completed' && !job.result && (
        <div className="bg-yellow-50 border border-yellow-200 rounded-md p-4">
          <p className="text-sm text-yellow-700">
            Job completed but no result data was returned.
          </p>
        </div>
      )}
    </div>
  )
}
