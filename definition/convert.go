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
		if len(out) == 0 || out[len(out)-1].ModelInput != task.GetModelInput() {
			out = append(
				out,
				EnrichInput{
					ModelInput: task.GetModelInput(),
				},
			)
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

// LoadInputsFromProto converts a Load RPC request into LoadInputs. Sources on the wire are per
// task; ModelInputs is per model input, so a model input of several tasks gets its oldest task's
// collection time -- the conservative "how stale is this" answer.
func LoadInputsFromProto(
	req *ginosdkpluginv1.LoadRequest,
) (LoadInputs, error) {
	out := make(LoadInputs, len(req.GetInputs()))

	for idx, in := range req.GetInputs() {
		var endpoint ginosdkapisv1.WorkerAssignmentEndpoint

		err := json.Unmarshal(in.GetEndpoint(), &endpoint)
		if err != nil {
			return nil, fmt.Errorf("decoding endpoint of input %d: %w", idx, err)
		}

		modelInputs := make(map[string]time.Time, len(in.GetSources()))

		for _, source := range in.GetSources() {
			at := protoTime(source.GetCollectedAt())

			current, ok := modelInputs[source.GetModelInput()]
			if !ok || at.Before(current) {
				modelInputs[source.GetModelInput()] = at
			}
		}

		out[idx] = LoadInput{
			Endpoint:    endpoint,
			EnrichedAt:  protoTime(in.GetEnrichedAt()),
			ModelInputs: modelInputs,
			Payload:     in.GetPayload(),
		}
	}

	return out, nil
}

// protoTime is ts as a time, unset (nil) as the zero time rather than the unix epoch.
func protoTime(ts *timestamppb.Timestamp) time.Time {
	if ts == nil {
		return time.Time{}
	}

	return ts.AsTime()
}
