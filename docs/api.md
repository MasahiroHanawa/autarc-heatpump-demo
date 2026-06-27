# REST API Specification

Base URL: `http://localhost:8080/api`

## Endpoints

### Create Job

```
POST /api/jobs
Content-Type: application/json
```

**Request Body:**

```json
{
  "buildingType": "detached_house",
  "livingAreaM2": 150,
  "constructionYear": 1985,
  "insulationLevel": "moderate",
  "currentHeatingSystem": "gas_boiler",
  "annualEnergyConsumptionKwh": 18000,
  "installerNotes": "South-facing roof, good access for outdoor unit."
}
```

**Response: `202 Accepted`**

```json
{
  "id": "job_abc123",
  "status": "pending",
  "createdAt": "2026-06-27T10:00:00Z"
}
```

### Get Job

```
GET /api/jobs/:id
```

**Response: `200 OK`** (pending/processing)

```json
{
  "id": "job_abc123",
  "status": "processing",
  "input": { ... },
  "createdAt": "2026-06-27T10:00:00Z",
  "updatedAt": "2026-06-27T10:00:02Z"
}
```

**Response: `200 OK`** (completed)

```json
{
  "id": "job_abc123",
  "status": "completed",
  "input": { ... },
  "result": {
    "estimatedHeatDemandKwh": 12000,
    "recommendedHeatPumpType": "air_to_water",
    "suitability": "good",
    "confidence": 0.82,
    "riskFlags": ["insulation_below_optimal"],
    "nextSteps": [
      "Verify wall insulation thickness",
      "Check electrical supply capacity"
    ],
    "summary": "An air-to-water heat pump is suitable for this property. The 1985 construction with moderate insulation suggests a heat demand of approximately 12,000 kWh/year. Upgrading insulation would improve efficiency."
  },
  "createdAt": "2026-06-27T10:00:00Z",
  "updatedAt": "2026-06-27T10:00:05Z"
}
```

**Response: `404 Not Found`**

```json
{
  "error": "job not found"
}
```

### List Jobs

```
GET /api/jobs
```

**Response: `200 OK`**

```json
{
  "jobs": [
    {
      "id": "job_abc123",
      "status": "completed",
      "buildingType": "detached_house",
      "createdAt": "2026-06-27T10:00:00Z"
    }
  ]
}
```

## Domain Types

### JobStatus

| Value | Description |
|-------|-------------|
| `pending` | Job created, waiting for worker |
| `processing` | Worker picked up the job |
| `completed` | Analysis finished successfully |
| `failed` | Analysis failed |

### BuildingType

`detached_house`, `semi_detached`, `terraced`, `apartment`, `bungalow`, `commercial`

### InsulationLevel

`poor`, `moderate`, `good`, `excellent`

### CurrentHeatingSystem

`gas_boiler`, `oil_boiler`, `electric`, `district_heating`, `none`

### Suitability

`excellent`, `good`, `fair`, `poor`, `unsuitable`

## Error Responses

All errors follow this format:

```json
{
  "error": "human-readable error message"
}
```

| Status | Meaning |
|--------|---------|
| `400` | Invalid request body or missing required fields |
| `404` | Job not found |
| `500` | Internal server error |
