peerpoxy intercepts and parses the raft RPCs that etcd peers use to communicate between
themselves. It's designed for development and pedagogical purposes, so should not be used
in production.

## Usage

Though you're welcome to use the peerpoxy cmd or pkg in other ways, it comes with a
pre-configured Procfile that spins up a 3 node development etcd cluster, with peerproxy
pre-configured to proxy and print all messages between peers.

You'll need to install goreman to run the Procfile:

```
go install github.com/mattn/goreman@latest
```

And then you can start the cluster:

```
goreman start
```

If you want to cleanup the development data files, they're in the config directory:

```
rm -r config/infra*
```

Here's an example of the logs peerproxy outputs:

```
01:10:47 peerproxy | &{Upstream:localhost:32380 Path:/raft/stream/message/e46cab30ccf0234f SrcID:bdbd480b20e09794 DstID:e46cab30ccf0234f Message:type:MsgHeartbeat to:16459718964217389903 from:13672163256400123796 term:16 commit:42 snapshot:<> }
```

## Future Work

peerproxy already serves its purpose, but in the future it may support:

- Latency injection (e.g., for simulating global etcd clusters)
- Improved log messages with more metadata and less noise
- Tooling to view/filter/analyze log messages
