
# Things I'm learning about protocol developments

- implement a simple htlc pass on from client.go
- implement the load balancing service from fiatjaf
- 

## OTher stuff

- define inputs of plugin (hooks, subscrptions)
- clear marshalling and unmarshalling of message type
-  

## Codec (the way mesasges are encoded)

- TLV (Type-Lenght-Value) encoding:
- every field of variable length has a preceding length field
  - `[u16:len]`
  - `[len*byte:refund_scriptpubkey]`
  - `[u16:len]`
  - `[len*byte:secret]`
- need encoding and decoding methods for this
- `chain_hash` is the hash of the blockchains genesis block

## Questions

- what is the difference between protocol buffers and manual codec
- 

## TODO

- look at eclair-hc-plugin cli commands
- [ ] create cli commands and work into the details from there
  - [ ] invoke-hc
  - [ ] send-htlc
  - [ ] sf
