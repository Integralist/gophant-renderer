# Gophant Renderer

> Ingest SQS messages and render their content into templates

## Local setup

- `spurious-server start`
- `spurious init`
- `spurious start`
- `export AWS_ACCESS_KEY_ID=access; export AWS_SECRET_ACCESS_KEY=secret; go run renderer.go`

> Note: DynamoDB is the only faked AWS service that requires exporting env vars
