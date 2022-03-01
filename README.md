# kong_api

### Task List

- [x] Return a list of services
    - [x] filtering by keyword (name and description)
    - [x] sorting
    - [x] pagination
- [x] Fetch single service
    - [x] include method for retrieving its versions
- [x] BONUS
    - [x] CRUD
        - [x] Create Service
        - [x] Update Service
        - [x] Delete Service
        - [x] Create version associated with service (no validation on whether service exists)
    - [ ] Unit test - Will omit this
    - [ ] Authentication/Authorization - Will omit this

# NOTE:
- Going to be quiet honest here...I've never used gorm and its auto migrations...so this was just kind of hacked together. Typically the DB tables would NOT be plural but the gorm creates it with s. Probably some set up that I need to configure but it doesn't break anything. I wanted to try out gorm so I said why the hell not. Lets learn something new. Also same with mux. My current role doesn't use that so again...lets learn something new.

# To run 
`docker-compose up --build`
`http://localhost:8080`

# API Endpoints
## Services
`GET` `/v1/services` - retrieves all services paginated with number of versions
#### Supported Query Params:
`page` - the current page
`page_size` - max services per page. Default is 10
`keyword` - a keyword search based on name or description
`sort_by` - property to sort by (id, name or description supported). Default is id
`sort_order` - ASC or DESC. Default is ASC
```
{
	"status": 200,
	"data": {
		"total": 1,
		"services": [
			{
				"id": 1,
				"name": "Service Name",
				"description": "This is a service test",
				"versions": 2,
				"created_at": "2022-03-01T00:22:08Z",
				"updated_at": "2022-03-01T00:22:08Z",
				"deleted_at": null
			}
		]
	}
}
```
`GET` `/v1/services/:id` - retrieves a single service with list of versions
```
{
	"status": 200,
	"data": {
		"id": 1,
		"name": "Service Name",
		"description": "This is a service test",
		"versions": [
			"1.0",
			"1.1"
		],
		"created_at": "2022-03-01T00:22:08Z",
		"updated_at": "2022-03-01T00:22:08Z",
		"deleted_at": null
	}
}
```
`POST` `/v1/services` - creates a new service with an initial version number. Returns the new service created
Payload:
```
{
	"name": "Service Name",
	"description": "This is a service test",
	"initial_version": "1.0"
}
```
`PUT` `/v1/services/:id` - updates a single service's name or description. Returns the updated service
Payload:
```
{
	"name": "Renaming Service",
	"description": "What is this?"
}
```
`DELETE` `/v1/services/:id` - marks a service as deleted. Does NOT delete associated versions. Returns OK
`POST` `/v1/versions` - creates a new version with an associated service
Payload:
```
{
	"version": "1.x",
	"service_id": 1
}
```

# Stack
- Backend - Go
- DB - MySQL
    - Chose MySQL since this is what I'm most familiar with. Chose a relational DB since the DB model doesn't seem to need to change frequently so no need for a NoSQL.

# Assumptions
- Assuming users will input the correct payload. There is some validation for inputs but it isn't very robust.
- When associating a new version to a service, there is no validation for service_id, so assumption is that the user will always input an existing service_id...not realistic but ya know.

# Design
- I tried to separate models. DB has its own model which I did not want to re-use and dictate the endpoint model. Which is why there are different models (endpoint model and database model). Idea is that the endpoint models are what we will send as the response. The database model will only be used in the DB layer.
- Tried my best to separate the backend as a controller layer (routes), service layer (business logic) and the database layer (database queries).
- For the individual GET endpoint for services, I could have just set the version IDs in the array, but for easier access and visualization, I just put the actual version string in an array. If I used just the ID, we would need another request to retrieve the version.
- For the index endpoint, I did the count on the backend query because it is much simplier to allow the DB to do what it does best. The other option could have been to just have the actual versions in an array and then on the front-end, do a .length to get the total number of versions. I think that would be a nice to have but gives away too much information. If there was a requirement to display the versions via a tooltip, then we may want to have the raw version numbers in the array instead.

# Future
- Add authentication. This would use JWT which would be generated when a user logs in.
    - This would require implementing a middleware so every request would verify the user making the request. UserID can be set in the header which the middleware will retrieve and verify the token.
- Add authorization. Some services may be permission locked for specific users. In the controller layer, we should verify whether the user has permission to view specific services. If they do not, on the single service endpoint, they will receive a 403. On the index endpoint, the service should not appear and only services which the user has permission will appear in the list.
- Add Unit tests and api tests
    - Unit tests will test out each service layer function with expected input and output
    - api tests will essentially be an e2e test by calling the endpoints with fake DB data.
- Add a shit ton of validations
- Add transactions to make sure if anything fails during creation of a service with an initial version that we roll things back
- Delete associated versions when deleting services

