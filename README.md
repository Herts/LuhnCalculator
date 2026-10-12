# LuhnCalculator

A dependency-free Go library for validating Luhn checksums and generating check
digits. String inputs preserve leading zeros and work beyond integer limits.

## Install

Requires Go 1.26 or later.

```sh
go get github.com/Herts/LuhnCalculator
```

```go
import luhn "github.com/Herts/LuhnCalculator"

number, err := luhn.AppendCheckDigit("0007992739871")
// number: "00079927398713"; err: nil
valid := luhn.Validate(number) // true
digit, err := luhn.CheckDigit("7992739871") // 3, nil
```

## API

| Function | Result |
| --- | --- |
| `Validate(number string) bool` | Whether the complete number's checksum is valid. |
| `CheckDigit(payload string) (int, error)` | The single digit, from 0 to 9, to append. |
| `AppendCheckDigit(payload string) (string, error)` | The original payload plus its check digit. |

Inputs must be nonempty ASCII decimal strings. Spaces, signs, punctuation, and
Unicode digits are rejected. Generation returns `ErrInvalidDigits`; validation
returns `false`. A checksum does not verify an identifier's issuer or permitted
length. For example, `"0"` has a valid checksum.

### Compatibility

The package name remains `LuhnCalculator`; the import alias `luhn` is optional.
The original integer APIs remain available but are deprecated:

- `GenerateLastDigit(12)` returns **125**, the complete number, as before. Despite
  its name, it does not return only the last digit. Negative inputs and integer
  overflow return `0`; use the string APIs for explicit errors.
- `IsValidLuhn(125)` returns `true`; negative inputs return `false`.

## Development

```sh
go test -race -cover ./...
go vet ./...
go test -fuzz=FuzzAppendCheckDigit -fuzztime=10s
```

Tests cover known checksums, invalid characters, leading zeros, long inputs,
legacy behavior, overflow, and executable examples. Fuzzing checks preservation
of the payload and uniqueness of its check digit.
