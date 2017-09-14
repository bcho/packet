# Packet

[![Build Status](https://travis-ci.org/bcho/packet.svg)](https://travis-ci.org/bcho/packet)
[![](https://godoc.org/github.com/bcho/packet?status.svg)](http://godoc.org/github.com/bcho/packet)

Read & write binary packet with length.


## History

### 2.0

Adds:

- `ReadPacket` / `WritePacket` functions for low-level usage.
- `NewWriter` for building io.Writer adapter.

### 1.0

This package is based on [Leaf][]'s `tcp_msg` module.

[Leaf]: https://github.com/name5566/leaf/blob/47f1a7cc53fb761dd9d3b125f40c682fbcf8f158/network/tcp_msg.go

## License

MIT
