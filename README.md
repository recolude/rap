# RAP

Recolude's recording file format.

## Testing Locally

You need to generate mocks before you can run parts of the test suite.

```
go generate ./...
```

## TODO

Stress Test Recordings

* 10000 subjects
   * Nested depth of 10
   * All using 3 or more streams
     * With atleast 1k captures
   * 1k unique metadata keys
