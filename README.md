# temporal-edge

### Prerequisites

Generate the protobuf goodies by using `buf`.
https://buf.build/docs/reference/cli/buf/generate/

```shell

# first generate the proto defs
buf generate

# now generate the descriptor to be used by envoy
buf build --as-file-descriptor-set \
  --config buf.yaml \
  --output tmprl/generated/tmprl/set.json
```

### Verify 1

```shell

# this starts our simple fake Temporal service
# TODO do it with grpc instead of connect but really it should be the same
go run tmprl/cmd/server/main.go

# now run our little fake Temporal SDK request
# in this case we send `StartWorkflowExecutionRequest`
go run tmprl/cmd/client/main.go
```

**Implementation**: 
Overwrite the `Payloads` of outbound `[]*Payload` to have exactly ONE item: the _encrypted_ original `*Payloads` input value

**Verify**:
Original bytes sent from the "Worker" should exactly match, in order of appearance,
the response we get from the service after being _decrypted_.

### Verify 2

**Implementation**:
Use envoy field extraction to pull specific field(s) to perform the [first example](#verify-1).
Reference: https://www.envoyproxy.io/docs/envoy/latest/api-v3/extensions/filters/http/grpc_field_extraction/v3/config.proto#example

In order to do this we need to:
1. Verify simple overwrite of the request for a simple field. Let's use WorkflowID and verify it in our fake `service`.
2. Next, let's practice using a TinyGo impl to run some custom code and do the same thing.
3. Finally, put it all together and send a complex message with `Payloads` that must do encryption properly.
