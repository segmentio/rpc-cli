
  `rpc(1)` &sdot; a simple RPC CLI.

## Features

  - Interactive mode
  - Consumes input from stdin
  - Nice UX (`rpc :3000 sum [2,2]`)
  - HTTP / TCP (`rpc http://localhost:3000`)

## Installation

  ```go
  $ go get github.com/segmentio/rpc-cli/cmd/rpc
  ```

## Codecs

  - jsonrpc

## Usage

  ```bash
  $ rpc :3000
  :3000> Service.Echo foo=baz
  {
    foo: "baz"
  }
  :3000> Service.Sum [2,2]
  4
  ```

  ```bash
  $ echo '{"name":"{{ name }}"}' | phony | rpc :3000 Service.Echo
  {
    "name": "Dyan Patterson"
  }
  ....
  ```

  ```bash
  $ rpc :3000 Service.Echo foo=baz
  {
    "foo": "baz"
  }
  ```

  ```bash
  $ rpc :3000 Service.Sum [2,2]
  4
  ```
