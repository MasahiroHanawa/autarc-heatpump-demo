import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { createJob, ApiError } from '../api/client'
import { JobForm } from '../components/JobForm'
import type { JobInput } from '../types'

export function NewJob() {
  const navigate = useNavigate()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  async function handleSubmit(input: JobInput) {
    setIsSubmitting(true)
    setError(null)
    try {
      const job = await createJob(input)
      navigate(`/jobs/${job.id}`)
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : 'An unexpected error occurred. Please try again.',
      )
      setIsSubmitting(false)
    }
  }

  return (
    <div>
      <h2 className="text-2xl font-bold text-gray-900 mb-6">New Analysis</h2>

      {error && (
        <div className="mb-6 bg-red-50 border border-red-200 rounded-md p-4">
          <p className="text-sm text-red-700">{error}</p>
        </div>
      )}

      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <JobForm onSubmit={handleSubmit} isSubmitting={isSubmitting} />
      </div>
    </div>
  )
}
