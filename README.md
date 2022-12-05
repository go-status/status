# status: Yet Another Traceable `error` Interface

The `status` package provides the `status.Status` interface, which extends the
built-in `error` interface in the Go standard library.  `status.Status`
includes additional information, such as an error code and stack trace, which
can be useful for debugging and troubleshooting.

## Usage

To use the status package, import it in your code:

```go
import "github.com/go-status/status"
```

### Creating a `status.Status` object

### Converting a `status.Status` object to an `error`

## Differences Between `status.Status` and `error`

The `status.Status` interface extends the built-in `error` interface by
including additional information, such as an error code and stack trace.
This makes it easier to include and access detailed information about an error
in your code.

In contrast, the `error` interface provided by the Go standard library only
includes an error message.  This means that developers must manually encode
additional information, such as an error code, into the error message.  This
can make it more difficult to access and use this information in their code.
