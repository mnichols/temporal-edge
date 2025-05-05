package main

import (
	"bytes"
	"connectrpc.com/connect"
	"context"
	"fmt"
	"github.com/google/uuid"
	tmprlv1 "github.com/mnichols/temporal-edge/tmprl/generated/tmprl/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mnichols/temporal-edge/tmprl/internal"
	"log"
)

var EncryptionKey = []byte("test-key-test-key-test-key-test!")

func main() {

	ctx := context.Background()
	c, err := internal.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	data := createData(5)

	var rawBytes [][]byte
	for _, d := range data {
		raw, err := proto.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		rawBytes = append(rawBytes, raw)
	}

	payloads := &tmprlv1.Payloads{
		Payloads: make([]*tmprlv1.Payload, len(rawBytes)),
	}
	for i, raw := range rawBytes {
		payloads.Payloads[i] = &tmprlv1.Payload{
			Metadata: map[string][]byte{"encoding": []byte("json/proto")},
			Data:     raw,
		}
	}

	workerRequest := &tmprlv1.StartWorkflowExecutionRequest{
		WorkflowId: timestamppb.Now().String(),
		Identity:   "launch",
		Input:      payloads,
	}

	allPayloadsBytes, err := proto.Marshal(workerRequest.Input)
	if err != nil {
		log.Fatal(err)
	}

	allPayloadsEncrypted, err := internal.Encrypt(allPayloadsBytes, EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}

	sentToTemporalService := &tmprlv1.Payloads{
		Payloads: make([]*tmprlv1.Payload, 1),
	}
	sentToTemporalService.Payloads[0] = &tmprlv1.Payload{
		Metadata: map[string][]byte{"encoding": []byte("envoy")},
		Data:     allPayloadsEncrypted,
	}
	proxiedRequest := &tmprlv1.StartWorkflowExecutionRequest{
		WorkflowId: timestamppb.Now().String(),
		Identity:   "launch",
		Input:      sentToTemporalService,
	}

	res, err := c.StartWorkflowExecution(ctx, connect.NewRequest(proxiedRequest))
	if err != nil {
		log.Fatal(err)
	}
	if res.Msg.Input.Payloads == nil {
		log.Fatal("payloads is empty")
	}
	if len(res.Msg.Input.Payloads) != 1 {
		log.Fatalf("expected 1 payload got %d", len(res.Msg.Input.Payloads))
	}
	responseEncoding := res.Msg.Input.Payloads[0].Metadata["encoding"]

	if string(responseEncoding) != "envoy" {
		log.Fatalf("unexpected encoding %V", string(responseEncoding))
	}

	decrypted, err := internal.Decrypt(res.Msg.Input.Payloads[0].Data, EncryptionKey)
	if err != nil {
		log.Fatal(err)
	}
	receivedPayloads := &tmprlv1.Payloads{}
	if err := proto.Unmarshal(decrypted, receivedPayloads); err != nil {
		log.Fatal(err)
	}
	responseProxySendsToWorker := &tmprlv1.StartWorkflowExecutionResponse{
		WorkflowId: res.Msg.WorkflowId,
		Identity:   res.Msg.Identity,
		Input: &tmprlv1.Payloads{
			Payloads: receivedPayloads.Payloads,
		},
	}

	if len(responseProxySendsToWorker.Input.Payloads) != len(rawBytes) {
		log.Fatalf("Expected back %d, got %d", len(rawBytes), len(responseProxySendsToWorker.Input.Payloads))
	}

	for i, payload := range responseProxySendsToWorker.Input.Payloads {
		if !bytes.Equal(payload.Data, rawBytes[i]) {
			log.Fatalf("decrypted unencryptedData does not match expected unencryptedData at %d", i)
		}
	}

	fmt.Printf("seemed to work! %#v\\n", res.Msg)
}
func createData(count int) []*tmprlv1.MyWorkflowData {
	result := make([]*tmprlv1.MyWorkflowData, 0, count)
	for i := 0; i < count; i++ {
		result = append(result, &tmprlv1.MyWorkflowData{Value: uuid.NewString()})
	}
	return result
}
