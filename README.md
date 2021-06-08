# Wheel API
The Wheel API exposes endpoints to manage:
- Teams
- People
- Turns

## Teams
GET `/api/teams`
```json
// Response:
{
  "data": [
    {
      "id": 1, 
      "name": "Trading"
    }, 
     ...
  ]
}
```

GET `/api/teams/{team}`
```json
// Response:
{
  "data": [
    {
      "id": 1, 
      "name": "Trading", 
      "people": [
        {
          "id": 1,
          "first_name": "Natasha",
          "last_name": "Romanoff",
          "email": "b.widow88@vendhq.com",
          "team_id": 1
        },
        ...
      ]
    } 
  ]
}
```

POST `/api/teams`
```json
// Request:
{ 
  "name":  "NewTeam"
}

// Response
{
  "data": {
    "id": 1, 
    "name": "NewTeam"
  }
}
```

PUT `/api/teams/{team}`
```json
// Request:
{ 
  "name":  "NewName"
}

// Response
{
  "data": {
    "id": 1, 
    "name": "NewName"
  }
}
```

DELETE `/api/teams/{team}`
```json
// Response
{
  "data": ""
}
```

## People
GET `/api/teams/{team}/people`
```json
// Response:
{
  "data": [
    {
      "id": 1,
      "first_name": "Bartholomew Henry",
      "last_name": "Allen",
      "email": "speed@vendhq.com",
      "team_id": 1
    },
     ...
  ]
}
```

GET `/api/teams/{team}/people/{person}`
```json
// Response:
{
  "data": {
    "id": 1,
    "first_name": "Anthony Edward",
    "last_name": "Stark",
    "email": "iron@vendhq.com",
    "team_id": 1
  }
}
```

POST `/api/teams/{team}/people`
```json
// Request:
{
  "first_name": "Bruce",
  "last_name": "Wayne",
  "email": "not.batman@vendhq.com"
}

// Response
{
  "data": {
    "id": 3,
    "first_name": "Bruce",
    "last_name": "Wayne",
    "email": "not.batman@vendhq.com",
    "team_id": 1
  }
}
```

PUT `/api/teams/{team}/people/{person}`
```json
// Request:
{
  "first_name": "Other",
  "last_name": "Name",
  "email": "other.email@vendhq.com",
  "team_id": 2
}

// Response
{
  "data": {
    "id": 3,
    "first_name": "Other",
    "last_name": "Name",
    "email": "other.email@vendhq.com",
    "team_id": 2
  }
}
```

DELETE `/api/teams/{team}/people/{person}`
```json
// Response
{
  "data": ""
}
```

## Turns
GET `/api/teams/{team}/turns`

Optional Query params:
- limit (Defaults to `10`)
- offset (Defaults to `0`)
- date_from (`YYYY-MM-DD`)
- date_to (`YYYY-MM-DD`)
```json
// Response:
{
  "data": [
    {
      "id": 1,
      "person_id": 2,
      "team_id": 1,
      "date": "2021-05-18T00:00:00+12:00",
      "created_at": "2021-05-17T04:11:32+12:00"
    },
    ...
  ]
}
```
POST `/api/teams/{team}/turns`
```json
// Request:
{
  "person_id": 1
}

// Response:
{
  "data": {
    "id": 1,
    "person_id": 1,
    "team_id": 1,
    "date": "2021-05-18T00:00:00+12:00",
    "created_at": "2021-05-17T04:11:32+12:00"
  }
}
```