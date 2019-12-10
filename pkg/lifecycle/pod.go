package lifecycle

import (
	"context"
	"encoding/json"
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func GetFirstTerminationMessage(pod *corev1.Pod) string {
	if pod != nil {
		for _, cs := range pod.Status.ContainerStatuses {
			if cs.State.Terminated != nil && cs.State.Terminated.Message != "" {
				return cs.State.Terminated.Message
			}
		}
	}
	return ""
}

func GetOperationsResult(ctx context.Context, pod *corev1.Pod, result interface{}) error {
	if pod == nil {
		return fmt.Errorf("pod was nil")
	}
	terminationMessage := GetFirstTerminationMessage(pod)
	if terminationMessage == "" {
		return fmt.Errorf("did not find termination message for pod %q", pod.Name)
	}
	err := json.Unmarshal([]byte(terminationMessage), &result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal terminationmessage: %q : %q", terminationMessage, err)
	}
	return nil
}
