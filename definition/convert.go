package definition

import (
	"encoding/json"
	"fmt"
	"time"

	ginosdkapisv1 "github.com/gino-works/gino-sdk/apis/v1"
	ginosdkpluginv1 "github.com/gino-works/gino-sdk/plugin/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TransformRequestFromProto converts a Transform RPC request into the definition's
// TransformRequest. The fn chain (req.GetFn()) is separate: it picks which registered transform
// runs, see Registry.Transform.
func TransformRequestFromProto(
	req *ginosdkpluginv1.TransformRequest,
) (TransformRequest, error) {
	var endpoint ginosdkapisv1.WorkerAssignmentEndpoint

	err := json.Unmarshal(req.GetEndpoint(), &endpoint)
	if err != nil {
		return TransformRequest{}, fmt.Errorf("decoding endpoint: %w", err)
	}

	return TransformRequest{
		Endpoint:    endpoint,
		ModelInput:  req.GetModelInput(),
		Task:        req.GetTask(),
		CollectedAt: protoTime(req.GetCollectedAt()),
		Payload:     req.GetPayload(),
	}, nil
}

// EnrichInputsFromProto converts an Enrich RPC request's flat, ordered task list into
// EnrichInputs: consecutive tasks of the same input are grouped, order is kept.
func EnrichInputsFromProto(
	req *ginosdkpluginv1.EnrichRequest,
) EnrichInputs {
	var out EnrichInputs

	for _, task := range req.GetTasks() {
		if len(out) == 0 || out[len(out)-1].Name != task.GetModelInput() {
			out = append(out, EnrichInput{Name: task.GetModelInput()})
		}

		last := &out[len(out)-1]

		last.Tasks = append(
			last.Tasks,
			EnrichTask{
				Name:        task.GetTask(),
				CollectedAt: protoTime(task.GetCollectedAt()),
				Payload:     task.GetPayload(),
			},
		)
	}

	return out
}

// protoTime is ts as a time, unset (nil) as the zero time rather than the unix epoch.
func protoTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}

	return ts.AsTime()
}
