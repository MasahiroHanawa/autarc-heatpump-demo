import { useState } from 'react'
import type {
  BuildingType,
  CurrentHeatingSystem,
  InsulationLevel,
  JobInput,
} from '../types'

interface FormState {
  buildingType: BuildingType | ''
  livingAreaM2: string
  constructionYear: string
  insulationLevel: InsulationLevel | ''
  currentHeatingSystem: CurrentHeatingSystem | ''
  annualEnergyConsumptionKwh: string
  installerNotes: string
}

interface FormErrors {
  buildingType?: string
  livingAreaM2?: string
  constructionYear?: string
  insulationLevel?: string
  currentHeatingSystem?: string
  annualEnergyConsumptionKwh?: string
}

const INITIAL_STATE: FormState = {
  buildingType: '',
  livingAreaM2: '',
  constructionYear: '',
  insulationLevel: '',
  currentHeatingSystem: '',
  annualEnergyConsumptionKwh: '',
  installerNotes: '',
}

function validate(state: FormState): FormErrors {
  const errors: FormErrors = {}

  if (!state.buildingType) errors.buildingType = 'Required'
  if (!state.insulationLevel) errors.insulationLevel = 'Required'
  if (!state.currentHeatingSystem) errors.currentHeatingSystem = 'Required'

  const area = Number(state.livingAreaM2)
  if (!state.livingAreaM2 || isNaN(area) || area <= 0)
    errors.livingAreaM2 = 'Must be a positive number'

  const year = Number(state.constructionYear)
  const currentYear = new Date().getFullYear()
  if (!state.constructionYear || isNaN(year) || year < 1800 || year > currentYear)
    errors.constructionYear = `Must be between 1800 and ${currentYear}`

  const kwh = Number(state.annualEnergyConsumptionKwh)
  if (!state.annualEnergyConsumptionKwh || isNaN(kwh) || kwh <= 0)
    errors.annualEnergyConsumptionKwh = 'Must be a positive number'

  return errors
}

interface JobFormProps {
  onSubmit: (input: JobInput) => void
  isSubmitting: boolean
}

export function JobForm({ onSubmit, isSubmitting }: JobFormProps) {
  const [form, setForm] = useState<FormState>(INITIAL_STATE)
  const [errors, setErrors] = useState<FormErrors>({})
  const [touched, setTouched] = useState<Partial<Record<keyof FormState, boolean>>>({})

  function handleChange(field: keyof FormState, value: string) {
    setForm((prev) => ({ ...prev, [field]: value }))
    if (touched[field]) {
      const next = validate({ ...form, [field]: value })
      setErrors((prev) => ({ ...prev, [field]: next[field as keyof FormErrors] }))
    }
  }

  function handleBlur(field: keyof FormState) {
    setTouched((prev) => ({ ...prev, [field]: true }))
    const next = validate(form)
    setErrors((prev) => ({ ...prev, [field]: next[field as keyof FormErrors] }))
  }

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const errs = validate(form)
    if (Object.keys(errs).length > 0) {
      setErrors(errs)
      setTouched(
        Object.fromEntries(Object.keys(errs).map((k) => [k, true])) as Partial<
          Record<keyof FormState, boolean>
        >,
      )
      return
    }
    onSubmit({
      buildingType: form.buildingType as BuildingType,
      livingAreaM2: Number(form.livingAreaM2),
      constructionYear: Number(form.constructionYear),
      insulationLevel: form.insulationLevel as InsulationLevel,
      currentHeatingSystem: form.currentHeatingSystem as CurrentHeatingSystem,
      annualEnergyConsumptionKwh: Number(form.annualEnergyConsumptionKwh),
      ...(form.installerNotes ? { installerNotes: form.installerNotes } : {}),
    })
  }

  const fieldClass = (error?: string) =>
    `block w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 ${
      error ? 'border-red-300' : 'border-gray-300'
    }`

  return (
    <form onSubmit={handleSubmit} noValidate className="space-y-6">
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Building Type
          </label>
          <select
            value={form.buildingType}
            onChange={(e) => handleChange('buildingType', e.target.value)}
            onBlur={() => handleBlur('buildingType')}
            className={fieldClass(errors.buildingType)}
          >
            <option value="">Select…</option>
            <option value="detached_house">Detached House</option>
            <option value="semi_detached">Semi-Detached</option>
            <option value="terraced">Terraced</option>
            <option value="apartment">Apartment</option>
            <option value="bungalow">Bungalow</option>
            <option value="commercial">Commercial</option>
          </select>
          {errors.buildingType && (
            <p className="mt-1 text-xs text-red-600">{errors.buildingType}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Living Area (m²)
          </label>
          <input
            type="number"
            min={1}
            value={form.livingAreaM2}
            onChange={(e) => handleChange('livingAreaM2', e.target.value)}
            onBlur={() => handleBlur('livingAreaM2')}
            placeholder="e.g. 150"
            className={fieldClass(errors.livingAreaM2)}
          />
          {errors.livingAreaM2 && (
            <p className="mt-1 text-xs text-red-600">{errors.livingAreaM2}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Construction Year
          </label>
          <input
            type="number"
            min={1800}
            max={new Date().getFullYear()}
            value={form.constructionYear}
            onChange={(e) => handleChange('constructionYear', e.target.value)}
            onBlur={() => handleBlur('constructionYear')}
            placeholder="e.g. 1985"
            className={fieldClass(errors.constructionYear)}
          />
          {errors.constructionYear && (
            <p className="mt-1 text-xs text-red-600">{errors.constructionYear}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Insulation Level
          </label>
          <select
            value={form.insulationLevel}
            onChange={(e) => handleChange('insulationLevel', e.target.value)}
            onBlur={() => handleBlur('insulationLevel')}
            className={fieldClass(errors.insulationLevel)}
          >
            <option value="">Select…</option>
            <option value="poor">Poor</option>
            <option value="moderate">Moderate</option>
            <option value="good">Good</option>
            <option value="excellent">Excellent</option>
          </select>
          {errors.insulationLevel && (
            <p className="mt-1 text-xs text-red-600">{errors.insulationLevel}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Current Heating System
          </label>
          <select
            value={form.currentHeatingSystem}
            onChange={(e) => handleChange('currentHeatingSystem', e.target.value)}
            onBlur={() => handleBlur('currentHeatingSystem')}
            className={fieldClass(errors.currentHeatingSystem)}
          >
            <option value="">Select…</option>
            <option value="gas_boiler">Gas Boiler</option>
            <option value="oil_boiler">Oil Boiler</option>
            <option value="electric">Electric</option>
            <option value="district_heating">District Heating</option>
            <option value="none">None</option>
          </select>
          {errors.currentHeatingSystem && (
            <p className="mt-1 text-xs text-red-600">{errors.currentHeatingSystem}</p>
          )}
        </div>

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Annual Energy Consumption (kWh)
          </label>
          <input
            type="number"
            min={1}
            value={form.annualEnergyConsumptionKwh}
            onChange={(e) => handleChange('annualEnergyConsumptionKwh', e.target.value)}
            onBlur={() => handleBlur('annualEnergyConsumptionKwh')}
            placeholder="e.g. 18000"
            className={fieldClass(errors.annualEnergyConsumptionKwh)}
          />
          {errors.annualEnergyConsumptionKwh && (
            <p className="mt-1 text-xs text-red-600">
              {errors.annualEnergyConsumptionKwh}
            </p>
          )}
        </div>
      </div>

      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Installer Notes{' '}
          <span className="font-normal text-gray-400">(optional)</span>
        </label>
        <textarea
          value={form.installerNotes}
          onChange={(e) => handleChange('installerNotes', e.target.value)}
          rows={3}
          placeholder="e.g. South-facing roof, good access for outdoor unit."
          className="block w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <button
        type="submit"
        disabled={isSubmitting}
        className="px-6 py-2.5 bg-blue-600 text-white text-sm font-medium rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
      >
        {isSubmitting ? 'Submitting…' : 'Run Analysis'}
      </button>
    </form>
  )
}
