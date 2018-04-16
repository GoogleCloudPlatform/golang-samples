# Dialogflow: Go samples

Dialogflow samples using the Go client.

## Table of contents

+ [Before you begin](#before-you-begin)

+ [Samples](#samples)
  + [Detect intent (text, audio, streaming audio)](#detect-intent)
  + [Intent management](#intent-management)
  + [Entity type management](#entity-type-management)
  + [Entity management](#entity-management)
  + [Session entity type management](#session-entity-type-management)

## Before you begin

1. If your project does not already have a Dialogflow agent, create one
   following [these
   instructions](https://dialogflow.com/docs/getting-started/building-your-first-agent#create_an_agent).

   (If you want to create an enterprise agent, follow [these
   instructions](https://cloud.google.com/dialogflow-enterprise/docs/quickstart)
   instead.)

2. This sample comes with a [sample
   agent](./resources/RoomReservation.zip)
   which you can use to run the samples. Follow [these
   instructions](https://dialogflow.com/docs/best-practices/import-export-for-versions)
   to import the agent from the [Dialogflow
   console](https://console.dialogflow.com/api-client/).

   (**Warning:** Importing the sample agent will add intents and entities to
   your Dialogflow agent. You may want to use a separate Google Cloud Platform
   project, or export your existing Dialogflow agent before importing the sample
   agent.)

## Samples

### Detect intent

(Text, audio, streaming audio)

[Source code](detect_intent/detect_intent.go)

**Usage:** `go run detect_intent/detect_intent.go -help`

### Intent management

[Source code](intent_management/intent_management.go)

**Usage:** `go run intent_management/intent_management.go -help`

### Entity type management

[Source code](entity_type_management/entity_type_management.go)

**Usage:** `go run entity_type_management/entity_type_management.go -help`

### Entity management

[Source code](entity_management/entity_management.go)

**Usage:** `go run entity_management/entity_management.go -help`

### Session entity type management

[Source code](session_entity_type_management/session_entity_type_management.go)

**Usage:** `go run session_entity_type_management/session_entity_type_management.go -help`

