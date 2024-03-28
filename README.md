# Gossip Glomers

My attempt at the [Fly.io Gossip Glomers distributed systems challenges](https://fly.io/dist-sys/).

## Challenges

- `#1` Echo ✅
- `#2` Unique IDs ✅
- `#3` Broadcast
  - `a` Single-Node Broadcast ✅
  - `b` Multi-Node Broadcast ✅
  - `c` Fault Tolerant Broadcast ✅
    - A separate implementation `propagateWithResponseAndRetry` is used and should pass this challenge with decent but <100% success rate.
  - `d-e` Efficient Broadcast. Current benchmark:
    - Messages-per-operation `77.756516`
    - Median latency `303`
    - Maximum latency `526`
