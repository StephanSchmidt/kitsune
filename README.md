# kitsune

<img src="kitsune.webp" alt="Kitsune Logo" width="200">

A small Go utility library of shared helpers used across my projects: UUIDs (v7) with base62 encoding, error wrapping with structured details, JSON marshalling helpers for `database/sql`, and a URL producer type.

## Install

```sh
go get github.com/StephanSchmidt/kitsune
```

Requires Go 1.24+.

## Usage

### UUIDs

Generates UUIDv7 (time-ordered) using `github.com/gofrs/uuid`.

```go
id := kitsune.NewUuid()
s := kitsune.ToString(id)           // hex-encoded
b62 := kitsune.ToBase62(id)         // short, base62

parsed, err := kitsune.FromString(s)
back, err := kitsune.FromBase62(b62)
fromBytes, err := kitsune.FromByteArray(id.Bytes())
```

### Errors

Thin wrappers over `gitlab.com/tozd/go/errors` that attach structured key/value details and support `fmt`-style placeholders.

```go
if err := do(); err != nil {
    return kitsune.WrapWithDetails(err, "failed processing user %s", "user_id", userID)
}

err := kitsune.WithDetails("invalid state", "state", state, "step", 3)

details := kitsune.AllDetails(err) // map[string]interface{}
```

`WrapWithDetails` returns `nil` when the wrapped error is `nil`.

### JSON (database/sql)

Helpers for storing JSON in SQL columns.

```go
// driver.Valuer side
val, err := kitsune.Marshall(myStruct)

// Scan side
var out MyStruct
err = kitsune.Unmarshal(rawValue, &out) // accepts string or []byte
```

### URL producer

```go
type UrlProducer func(args ...interface{}) string
```

A function type for deferred/parameterised URL generation.

## Development

```sh
make test       # run tests
make alltest    # imports, lint, sec, nilcheck, test
make coverage   # coverage report
make build      # test + refresh coverage badge in README
```

## Code Coverage

<!-- COVERAGE:START -->
**Current Coverage: 100.0%**

Last updated: Mon Sep 29 20:26:02 CEST 2025
<!-- COVERAGE:END -->

## License

MIT — see [LICENSE](LICENSE).
