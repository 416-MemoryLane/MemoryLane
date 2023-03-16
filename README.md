# MemoryLane

## Quick Start

From the project root directory:

- Build p2p node class `go build -o app`
- Spin up an p2p node `./app`

## Proof of Concept

From the project root directory:

- Build the p2p node class `go build -o app`
- Spin up a p2p host in one terminal `./app`
- Record what the `ip-multiaddr` and `<peer-id>`are for the host node
- Spin up another p2p node in another terminal connecting it to the node we just created `./app --peer-address <ip-multiaddr>/p2p/<peer-id>`
- You should see that both the nodes (via terminal output) are sending and receiving a number to each other and incrementing it.
