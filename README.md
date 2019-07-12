<h1 align="center">
  <a href="libp2p.io"><img width="250" src="https://github.com/libp2p/libp2p/blob/master/logo/black-bg-2.png?raw=true" alt="libp2p hex logo" hspace="150" /> <img width="200" src="https://storage.googleapis.com/dghubble/gopher-on-bird.png" alt="go-twitter logo"/>
</a>
</h1>

# dohrnii
This project is a Go implementation of another project I made. This one uses the libp2p Networking Stack to store social media information on a permissionless distributed ledger.

It can be used as a framework for storing any data into a distributed ledger. I chose tweets in this case but you can literally fork it and use any API to fetch data of any service provider you want. Only a single element in the block structure needs to be updated. The proof-of-work system and peer-to-peer implementation stays the same.

*Turritopsis dohrnii is an immortal jellyfish.* 

## Requirements

- [go-twitter](https://github.com/dghubble/go-twitter)
- [go-libp2p](https://github.com/libp2p/go-libp2p)

## Build
Run `go build` in `cmd/dohrnii/`.

## Launch
Run `./dohrnii -l [PORT]` with `[PORT]` being any port you want. 

## License
Dohrnii is available under the MIT license. See the [LICENSE](LICENSE) file for more info.