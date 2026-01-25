---
description: Generate unit tests for Go files/functions
---

# Generate Unit Tests

## Usage
```
/generate-tests [FILE_PATH or FUNCTION_NAME]
```

## Steps

1. Read the target file to understand the code structure
2. Identify all public functions/methods that need testing
3. Create test file in the same directory with `_test.go` suffix
4. Generate table-driven tests with multiple test cases

## Test Requirements

- Use `testify/assert` for assertions
- Use `testify/mock` for mocking dependencies
- Follow table-driven test pattern
- Test both happy path and error cases
- Target coverage > 80%

## Example Test Structure

```go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string
        input   InputType
        want    OutputType
        wantErr bool
    }{
        {
            name:    "success case",
            input:   validInput,
            want:    expectedOutput,
            wantErr: false,
        },
        {
            name:    "error case - invalid input",
            input:   invalidInput,
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Setup mocks
            mockRepo := new(MockRepository)
            mockRepo.On("Method", mock.Anything).Return(...)

            // Execute
            result, err := FunctionUnderTest(tt.input)

            // Assert
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, result)
            }
            mockRepo.AssertExpectations(t)
        })
    }
}
```

## Mock Generation
For interfaces, create mock files in `mocks/` subdirectory using mockery or manual implementation.
