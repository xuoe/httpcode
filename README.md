### Installation

```
$ go install github.com/xuoe/httpcode@latest
```

### Usage

- Search by status code:
    ```
    $ httpcode 300 400
    300  Multiple Choices
    400  Bad Request
    ```

- Search by status text (case-insensitive, partial matches):
    ```
    $ httpcode accept Bad
    202  Accepted
    400  Bad Request
    406  Not Acceptable
    502  Bad Gateway
    ```

- Search by numeric glob pattern (`?` may be replaced by `_` to bypass shell quoting):
    ```
    $ httpcode __9
    409  Conflict
    429  Too Many Requests
    ```
    These are equivalent:
    ```
    # httpcode '1*'
    # httpcode '1??'
    $ httpcode 1__
    100  Continue
    101  Switching Protocols
    102  Processing
    103  Early Hints
    ```

- Include MDN links:
    ```
    $ httpcode -m 1__
    100  Continue             https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/100
    101  Switching Protocols  https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/101
    102  Processing           https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/102
    103  Early Hints          https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/103
    ```

Note that since multiple arguments are supported, any combination of the above may be used.
