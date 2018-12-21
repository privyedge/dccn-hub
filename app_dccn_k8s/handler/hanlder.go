package handler

import (
	"context"
	"github.com/Ankr-network/refactor/app_dccn_k8s/proto"
)

type K8sHandler struct{}

func (p *K8sHandler) K8Task(ctx context.Context, req *k8s.K8S_K8TaskStream) error {
	return nil
}

