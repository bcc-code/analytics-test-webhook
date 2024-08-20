
# Dumb webhook storage

This is a simple Go service designed to receive arbitrary JSON data via a webhook and store it in memory.
It's retrievable via an API_KEY protected endpoint.

## Deployment

```sh
gcloud run deploy # And just follow the steps. Allow anonymous access when asked.
```

## Features

- **Webhook Endpoint**: Accepts POST requests with JSON payloads at an endpoint identified by a unique `id`.
- **In-Memory Data Storage**: Stores the received JSON data in memory, associated with the given `id`.
- **Automatic Data Expiry**: Clears the stored data for an `id` if no data is received within the last 10 minutes before processing a new request.
- **Data Retrieval**: Provides a GET endpoint to retrieve all stored data for a specific `id`, protected by an API key.

## Use Case

This service is particularly useful for scenarios like end-to-end (E2E) testing where you need to verify that certain events or data points are being sent to a specific endpoint. For example, you could integrate this service with tools like RudderStack or SegmentIO during your CI/CD pipeline to validate that expected analytics events are received during test runs.

## Endpoints

- `POST /webhook/:id` - Receives JSON data on the specified `id`.
- `GET /get_data/:id?api_key=your_api_key` - Retrieves all data received for the specified `id`, protected by an API key.

## Environment Variables

- `API_KEY` - The API key required to access the `GET /get_data/:id` endpoint.

## Example Usage

Below is an example shell script (`example_usage.sh`) demonstrating how to use the service:

```bash
#!/bin/bash

# Set the API key for accessing the data
API_KEY="your_api_key"

# Define the ID for the webhook
WEBHOOK_ID="test-webhook"

# Send a POST request to the webhook with some JSON data
curl -X POST -H "Content-Type: application/json"      -d '{"event": "test_event", "value": 42}'      http://localhost:8080/webhook/$WEBHOOK_ID

# Retrieve the stored data for the webhook
curl -G -H "Content-Type: application/json"      --data-urlencode "api_key=$API_KEY"      http://localhost:8080/get_data/$WEBHOOK_ID
```

## Running the Service

1. Set the `API_KEY` environment variable before running the service:
   ```bash
   export API_KEY=your_api_key
   ```
2. Run the service:
   ```bash
   go run main.go
   ```
3. Use the example shell script to send data to the webhook and retrieve it.

## Notes

- Ensure that the `API_KEY` environment variable is set and matches the key used in your requests to protect access to the stored data.
- The service clears data for an `id` if no data has been received in the last 10 minutes when a new request is processed.
