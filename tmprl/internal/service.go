package internal

import (
	"connectrpc.com/connect"
	"context"
	v1 "github.com/mnichols/temporal-edge/tmprl/generated/tmprl/v1"
	"github.com/mnichols/temporal-edge/tmprl/generated/tmprl/v1/tmprlv1connect"
	"log"
)

type TmprlService struct {
	tmprlv1connect.UnimplementedTmprlServiceHandler
}

func (s *TmprlService) StartWorkflowExecution(ctx context.Context, c *connect.Request[v1.StartWorkflowExecutionRequest]) (*connect.Response[v1.StartWorkflowExecutionResponse], error) {
	log.Println("Registered", c.Msg.GetWorkflowId(), "with identity", c.Msg.GetIdentity())

	/* So the first Payload has:
	{
		"metadata": {
			"envoy": "enc",
		},
		"data": <encrypted(msg.Payloads.Payloads array) >
	}
	*/
	if c.Msg.Input != nil {
		log.Println("Received this many items in Payloads", len(c.Msg.Input.Payloads))
		log.Println("The first payload value metadata is ", string(c.Msg.Input.Payloads[0].GetMetadata()["envoy"]))
		log.Println("The first payload value data is ", string(c.Msg.Input.Payloads[0].Data))
	}

	return connect.NewResponse(&v1.StartWorkflowExecutionResponse{
		WorkflowId: c.Msg.GetWorkflowId(),
		Input:      c.Msg.GetInput(),
		Identity:   c.Msg.GetIdentity(),
	}), nil
}
