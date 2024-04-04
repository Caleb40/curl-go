# About
This is a Golang implementation of cURL, developed without relying on external libraries such as net/http. Instead, it adheres to the HTTP specification from scratch. While still under development, the project currently offers essential cURL features. Additional functionalities will be progressively incorporated over time.

Here's an overview of key cURL commands and flags available in the curl-go CLI application:

## CLI Arguments

Here's the rearranged markdown based on the provided code:

## CLI Arguments

| Short | Long Form       | Default              | Type         | Description                                                                          |
|-------|-----------------|----------------------|--------------|--------------------------------------------------------------------------------------|
| `-V`  | `--version`     | `false`              | Boolean      | Returns version and exits                                                            |
| `-v`  | `--verbose`     | `false`              | Boolean      | Logs all headers and body to output                                                  |
| `-X`  | `--method`      | `GET`                | String       | Specifies the HTTP method to use (usually `GET` unless modified by other parameters) |
| `-o`  | `--output`      | `stdout`             | String       | Specifies where to output results                                                    |
| `-u`  | `--user`        | (none)               | String       | Specifies user:password for HTTP authentication                                      |
| `-d`  | `--data`        | (none)               | String Slice | HTML form data, sets mime type to `application/x-www-form-urlencoded`                |
| `-F`  | `--form`        | (none)               | String Slice | HTML form data, sets mime type to `multipart/form-data`                              |
|       | `--stderr`      | `stderr`             | String       | Logs errors to this replacement for stderr                                           |
| `-D`  | `--dump-header` | (none)               | String       | Specifies where to output headers (not enabled by default)                           |
| `-A`  | `--user-agent`  | `go-curling/##DEV##` | String       | Specifies the user-agent to use                                                      |
| `-e`  | `--referer`     | (none)               | String       | Specifies the referer URL to use with HTTP request                                   |
|       | `--url`         | (none)               | String       | Specifies the requesting URL                                                         |
| `-f`  | `--fail`        | `false`              | Boolean      | If fail, does not emit contents and returns fail exit code (-6)                      |
| `-k`  | `--insecure`    | `false`              | Boolean      | Ignores invalid SSL certificates                                                     |
| `-s`  | `--silent`      | `false`              | Boolean      | Silences all program console output                                                  |
| `-S`  | `--show-error`  | `false`              | Boolean      | Shows error info even if silent mode is on                                           |
| `-I`  | `--head`        | `false`              | Boolean      | Only returns headers, ignoring body content                                          |
| `-i`  | `--include`     | `false`              | Boolean      | Includes headers (prepended to body content)                                         |
| `-b`  | `--cookie`      | (none)               | String Slice | HTTP cookie, raw HTTP cookie only (use `-c` for cookie jar files)                    |
| `-c`  | `--cookie-jar`  | (none)               | String       | File for storing (read and write) cookies                                            |
| `-T`  | `--upload-file` | (none)               | String       | Raw file to PUT (default) to the given URL, not encoded                              |