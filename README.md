# ObjectHash-Proto

[![Build
Status](https://travis-ci.org/deepmind/objecthash-proto.svg?branch=master)](https://travis-ci.org/deepmind/objecthash-proto)
[![Go Report
Card](https://goreportcard.com/badge/github.com/deepmind/objecthash-proto)](https://goreportcard.com/report/github.com/deepmind/objecthash-proto)
[![GoDoc](https://godoc.org/github.com/deepmind/objecthash-proto?status.svg)](https://godoc.org/github.com/deepmind/objecthash-proto)

The library is an implementation of [Ben Laurie's
ObjectHash](https://github.com/benlaurie/objecthash) for protocol buffers.

This implementation is still experimental and until it is stable, protobuf
messages are not guaranteed to result in the same value.

## Usage

Get a new `ProtoHasher` instance using the `NewHasher` method, then call
`HashProto` with a protobuf message to get its ObjectHash:

```golang
hasher := protohash.NewHasher()
hash, err := hasher.HashProto(message)
```

## Options

In order to simplify compatibility with other ObjectHash applications, this
library exposes the following options that control how the hashing is done:

1.  `EnumsAsStrings()`: Makes enum values get hashed as strings instead of being
    hashed as their integer values.

1.  `FieldNamesAsKeys()`: Makes protobuf message fields use their names as keys
    instead of using the field tag numbers as keys.

1.  `MessageIdentifier(i)`: Instead of hashing protobuf messages as maps, this
    makes it possible to distinguish them by using `i` as the type-identifier
    that gets used in calculating the ObjectHash of a message.

Those options can be specified in any order as arguments to the `NewHasher`
function. Example:

```golang
hasher := protohash.NewHasher(EnumsAsStrings(), MessageIdentifier(`m`), FieldNamesAsKeys())
```
