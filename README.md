# Gophant Renderer

A static publishing package

> Ingest SQS messages, render content and store in S3 for easy lookup

## Architecture

Gophant is built on the premise of producers (renderers) and consumers (brokers):

![Gophant Architecture](https://cloud.githubusercontent.com/assets/180050/13029403/4cfb12f4-d283-11e5-9958-d996b6cb97ab.png)

## Local setup

### Creating Resources

First thing we need to do is to create some fake services (SQS, DynamoDB, S3) for our renderer to utilise. In the real world this is handled by AWS, but for the purposes of running things locally we'll have to handle this ourselves manually.

- `spurious-server start`
- `spurious init`
- `spurious start`

> Note: [Spurious](https://github.com/spurious-io/spurious) is required for the faking of AWS Services

- `export AWS_ACCESS_KEY_ID=access; export AWS_SECRET_ACCESS_KEY=secret; go run local.go`

> DynamoDB is the only faked AWS service that requires exporting env vars.  
> It uses them to segregate the tables you create.  
> For your code (or Spurious Browser) to access your DynamoDB data,  
> you'll need to export/login using the same access/secret key values

#### Output

You should see something like the following output when running `local.go` for the first time:

```
SQS Create Queue:
{
  QueueUrl: "http://sqs.spurious.localhost:32768/producer"
}

SQS Sent Message:
{
  MD5OfMessageBody: "1356c67d7ad1638d816bfb822dd2c25d",
  MessageId: "fe422b17-4c3c-4542-9330-7a253a034025"
}

Sequence Table:
{
  TableDescription: {
    AttributeDefinitions: [{
        AttributeName: "key",
        AttributeType: "S"
      }],
    CreationDateTime: 2016-02-13 17:25:54 +0000 UTC,
    ItemCount: 0,
    KeySchema: [{
        AttributeName: "key",
        KeyType: "HASH"
      }],
    ProvisionedThroughput: {
      NumberOfDecreasesToday: 0,
      ReadCapacityUnits: 10,
      WriteCapacityUnits: 10
    },
    TableName: "sequencer",
    TableSizeBytes: 0,
    TableStatus: "ACTIVE"
  }
}

Lookup Table:
{
  TableDescription: {
    AttributeDefinitions: [{
        AttributeName: "component_key",
        AttributeType: "S"
      },{
        AttributeName: "batch_version",
        AttributeType: "N"
      }],
    CreationDateTime: 2016-02-13 17:25:55 +0000 UTC,
    ItemCount: 0,
    KeySchema: [{
        AttributeName: "component_key",
        KeyType: "HASH"
      },{
        AttributeName: "batch_version",
        KeyType: "RANGE"
      }],
    ProvisionedThroughput: {
      NumberOfDecreasesToday: 0,
      ReadCapacityUnits: 10,
      WriteCapacityUnits: 10
    },
    TableName: "lookup",
    TableSizeBytes: 0,
    TableStatus: "ACTIVE"
  }
}

Finished creating local resources
```

### Optional flags

There are flags you can use to configure the local setup (if you prefer not to use Spurious):

- `-queue`: defaults to `producer`
- `-region`: defaults to `eu-west-1`
- `-s3-bucket`: defaults to `gophant`
- `-s3-endpoint`: defaults to `spurious ports --json`
- `-sqs-endpoint`: defaults to `spurious ports --json`
- `-dynamo-endpoint`: defaults to `spurious ports --json`

> Note: with DynamoDB we default table names to `sequencer` and `lookup` 

### Running the Renderer

- `export AWS_ACCESS_KEY_ID=access; export AWS_SECRET_ACCESS_KEY=secret; go run renderer.go`
