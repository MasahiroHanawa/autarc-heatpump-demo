import type { JobResult, Suitability } from '../types'

const SUITABILITY_COLOR: Record<Suitability, string> = {
  excellent: 'text-green-700',
  good: 'text-green-600',
  fair: 'text-yellow-600',
  poor: 'text-orange-600',
  unsuitable: 'text-red-600',
}

interface JobResultDisplayProps {
  result: JobResult
}

export function JobResultDisplay({ result }: JobResultDisplayProps) {
  const suitabilityColor = SUITABILITY_COLOR[result.suitability]
  const suitabilityLabel =
    result.suitability.charAt(0).toUpperCase() + result.suitability.slice(1)

  return (
    <div className="space-y-6">
      <div className="bg-white border border-gray-200 rounded-lg p-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-3">Summary</h3>
        <p className="text-gray-700 leading-relaxed">{result.summary}</p>
      </div>

      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div className="bg-white border border-gray-200 rounded-lg p-4">
          <p className="text-xs text-gray-500 mb-1">Heat Demand</p>
          <p className="text-lg font-bold text-gray-900">
            {result.estimatedHeatDemandKwh.toLocaleString()} kWh
          </p>
        </div>
        <div className="bg-white border border-gray-200 rounded-lg p-4">
          <p className="text-xs text-gray-500 mb-1">Pump Type</p>
          <p className="text-sm font-semibold text-gray-900 capitalize">
            {result.recommendedHeatPumpType.replace(/_/g, ' ')}
          </p>
        </div>
        <div className="bg-white border border-gray-200 rounded-lg p-4">
          <p className="text-xs text-gray-500 mb-1">Suitability</p>
          <p className={`text-sm font-semibold ${suitabilityColor}`}>
            {suitabilityLabel}
          </p>
        </div>
        <div className="bg-white border border-gray-200 rounded-lg p-4">
          <p className="text-xs text-gray-500 mb-1">Confidence</p>
          <p className="text-lg font-bold text-gray-900">
            {Math.round(result.confidence * 100)}%
          </p>
        </div>
      </div>

      {result.riskFlags.length > 0 && (
        <div className="bg-orange-50 border border-orange-200 rounded-lg p-4">
          <h4 className="text-sm font-semibold text-orange-800 mb-2">Risk Flags</h4>
          <ul className="space-y-1">
            {result.riskFlags.map((flag) => (
              <li key={flag} className="text-sm text-orange-700 flex gap-2">
                <span aria-hidden>•</span>
                <span className="capitalize">{flag.replace(/_/g, ' ')}</span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {result.nextSteps.length > 0 && (
        <div className="bg-white border border-gray-200 rounded-lg p-4">
          <h4 className="text-sm font-semibold text-gray-900 mb-2">Next Steps</h4>
          <ol className="space-y-1">
            {result.nextSteps.map((step, i) => (
              <li key={i} className="text-sm text-gray-700 flex gap-2">
                <span className="font-medium text-gray-400 shrink-0">{i + 1}.</span>
                <span>{step}</span>
              </li>
            ))}
          </ol>
        </div>
      )}
    </div>
  )
}
