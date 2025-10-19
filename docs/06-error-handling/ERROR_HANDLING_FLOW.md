# Error Handling Flow Diagram

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         Client Request                           │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Gin Handler                                 │
│  • Bind JSON (c.ShouldBindJSON)                                 │
│  • Validate struct (validator.Struct)                           │
│  • Business logic                                               │
└─────────────────────────┬───────────────────────────────────────┘
                          │
                ┌─────────┴─────────┐
                │                   │
                ▼                   ▼
        ┌───────────┐       ┌─────────────┐
        │   Error   │       │   Success   │
        └─────┬─────┘       └──────┬──────┘
              │                    │
              ▼                    ▼
    ┌──────────────────┐   ┌──────────────────┐
    │  Error Handler   │   │ Success Handler  │
    │  (utils package) │   │  (utils package) │
    └─────┬────────────┘   └─────┬────────────┘
          │                      │
          ▼                      ▼
    ┌─────────────────┐   ┌──────────────────┐
    │  Format Error   │   │  Format Success  │
    │  • Validation   │   │  • Add data      │
    │  • AppError     │   │  • Add message   │
    │  • Add details  │   │  • Set success=T │
    └─────┬───────────┘   └─────┬────────────┘
          │                      │
          │                      │
          └──────────┬───────────┘
                     │
                     ▼
          ┌─────────────────────┐
          │   JSON Response     │
          │   • success: bool   │
          │   • data/error      │
          │   • details         │
          └─────────────────────┘
                     │
                     ▼
          ┌─────────────────────┐
          │   Client Receives   │
          └─────────────────────┘
```

## Error Flow Details

### 1. Validation Error Flow

```
Input: Invalid JSON/Struct
         │
         ▼
c.ShouldBindJSON / validator.Struct
         │
         ▼
utils.RespondValidationError(c, err, code)
         │
         ▼
errors.FormatValidationError(err, code)
         │
         ├─> Parse validator.ValidationErrors
         ├─> Convert field names (PascalCase → readable)
         ├─> Generate readable messages per tag
         └─> Create AppError with details
         │
         ▼
utils.RespondError(c, 400, appError)
         │
         ▼
JSON Response with error details
```

### 2. Business Logic Error Flow

```
Service/Repository Error
         │
         ▼
Check Error Type
         │
    ┌────┴────┐
    │         │
    ▼         ▼
Not Found  Other Error
    │         │
    │         ├─> Create AppError
    │         └─> Add details (.WithDetails)
    │         │
    └─────┬───┘
          │
          ▼
utils.RespondWithAppError(c, appError)
          │
          ├─> GetStatusCodeFromErrorCode(code)
          └─> Map code to HTTP status
          │
          ▼
JSON Response with error
```

### 3. Success Flow

```
Successful Operation
         │
         ▼
utils.RespondSuccess(c, status, data, message?)
         │
         ├─> Create SuccessResponse struct
         ├─> Set success = true
         ├─> Add data
         └─> Add optional message
         │
         ▼
c.JSON(status, response)
         │
         ▼
JSON Response with data
```

## Component Interaction

```
┌─────────────────────┐
│   user_handler.go   │
│  (Your Handler)     │
└──────────┬──────────┘
           │
           │ uses
           │
           ▼
┌─────────────────────┐       ┌─────────────────────┐
│   response.go       │───────│    errors.go        │
│  (Response Utils)   │  uses │  (Error Definitions)│
│                     │       │                     │
│ • RespondSuccess    │       │ • AppError struct   │
│ • RespondError      │       │ • Error codes       │
│ • RespondValidation │       │ • Predefined errors │
│ • RespondWithApp    │       │ • Format functions  │
└─────────────────────┘       └─────────────────────┘
```

## Error Code Mapping

```
Error Code                    HTTP Status
────────────────────────────────────────────
USER001 (Not Found)      →   404 Not Found
USER002 (Already Exists) →   409 Conflict
USER003-005 (DB Ops)     →   500 Internal Server Error
USER006 (Invalid Creds)  →   401 Unauthorized
USER007 (Validation)     →   400 Bad Request

AUTH001-004 (Token/Auth) →   401 Unauthorized
AUTH005 (Forbidden)      →   403 Forbidden
AUTH006-007 (Validation) →   400 Bad Request

NOVEL001 (Not Found)     →   404 Not Found
NOVEL002-004 (DB Ops)    →   500 Internal Server Error
NOVEL005 (Validation)    →   400 Bad Request

[Similar pattern for other domains]
```

## Validation Error Transformation

```
Input:
validator.ValidationErrors {
  Field: "Username"
  Tag: "required"
  Value: ""
}

         │
         ▼

Process:
1. Extract field name: "Username"
2. Convert to readable: "username" (PascalCase → lowercase with spaces)
3. Get tag: "required"
4. Generate message: getValidationMessage(fieldError)
   → "username is required"

         │
         ▼

Output:
ValidationError {
  Field: "username"
  Message: "username is required"
  Tag: "required"
  Value: ""
}

         │
         ▼

AppError {
  Code: "USER007"
  Message: "Validation failed"
  Details: {
    fields: [ValidationError, ...]
    summary: "username: username is required; ..."
  }
}
```

## Usage Pattern

```
┌──────────────────────────────────────────────────┐
│ Step 1: Bind JSON                                │
│ if err := c.ShouldBindJSON(&dto); err != nil {  │
│     utils.RespondValidationError(...)           │
│     return                                       │
│ }                                               │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│ Step 2: Validate Struct                          │
│ validate := validator.New()                      │
│ if err := validate.Struct(dto); err != nil {    │
│     utils.RespondValidationError(...)           │
│     return                                       │
│ }                                               │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│ Step 3: Business Validation                      │
│ if someCondition {                               │
│     utils.RespondWithAppError(...)              │
│     return                                       │
│ }                                               │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│ Step 4: Service Call                             │
│ result, err := h.Service.DoSomething()          │
│ if err != nil {                                 │
│     appError := errors.Err...WithDetails(...)   │
│     utils.RespondWithAppError(...)              │
│     return                                       │
│ }                                               │
└─────────────────┬────────────────────────────────┘
                  │
                  ▼
┌──────────────────────────────────────────────────┐
│ Step 5: Success Response                         │
│ utils.RespondSuccess(c, status, result, msg)    │
└──────────────────────────────────────────────────┘
```

## Key Benefits Visualization

```
Before                              After
──────                              ─────

Raw validation error        →       Structured error with code
Inconsistent formats        →       Standardized responses
Hard to debug              →       Unique error codes
Generic messages           →       Human-readable messages
No error details           →       Detailed field errors
Manual status codes        →       Automatic mapping
```

## File Structure

```
internal/common/
├── errors/
│   └── errors.go
│       ├── AppError struct
│       ├── Error codes (constants)
│       ├── Predefined errors (variables)
│       ├── FormatValidationError()
│       ├── getValidationMessage()
│       └── toReadableField()
│
└── utils/
    ├── response.go
    │   ├── SuccessResponse struct
    │   ├── ErrorResponse struct
    │   ├── RespondSuccess()
    │   ├── RespondError()
    │   ├── RespondValidationError()
    │   ├── RespondWithAppError()
    │   └── GetStatusCodeFromErrorCode()
    │
    └── helpers.go (existing)
```

## Integration Points

```
Your Handlers
     │
     ├─> internal/common/utils (Response helpers)
     │        │
     │        └─> internal/common/errors (Error definitions)
     │
     └─> go-playground/validator (Validation)
              │
              └─> Errors formatted by our system
```
