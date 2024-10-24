package node

import (
	"context"
	"fmt"
	"github.com/rh-ecosystem-edge/kernel-module-management/internal/meta"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

//go:generate mockgen -source=node.go -package=node -destination=mock_node.go

type Node interface {
	IsNodeSchedulable(node *v1.Node) bool
	GetNodesListBySelector(ctx context.Context, selector map[string]string) ([]v1.Node, error)
	GetNumTargetedNodes(ctx context.Context, selector map[string]string) (int, error)
	UpdateLabels(ctx context.Context, node *v1.Node, toBeAdded, toBeRemoved []string) error
	NodeBecomeReadyAfter(node *v1.Node, checkTime metav1.Time) bool
}

type node struct {
	client client.Client
}

func NewNode(client client.Client) Node {
	return &node{
		client: client,
	}
}

func (n *node) IsNodeSchedulable(node *v1.Node) bool {
	for _, taint := range node.Spec.Taints {
		if taint.Effect == v1.TaintEffectNoSchedule {
			return false
		}
	}
	return true
}

func (n *node) GetNodesListBySelector(ctx context.Context, selector map[string]string) ([]v1.Node, error) {
	logger := log.FromContext(ctx)
	logger.V(1).Info("Listing nodes", "selector", selector)

	selectedNodes := v1.NodeList{}
	opt := client.MatchingLabels(selector)
	if err := n.client.List(ctx, &selectedNodes, opt); err != nil {
		return nil, fmt.Errorf("could not list nodes: %v", err)
	}
	nodes := make([]v1.Node, 0, len(selectedNodes.Items))

	for _, node := range selectedNodes.Items {
		if n.IsNodeSchedulable(&node) {
			nodes = append(nodes, node)
		}
	}
	return nodes, nil
}

func (n *node) GetNumTargetedNodes(ctx context.Context, selector map[string]string) (int, error) {
	targetedNode, err := n.GetNodesListBySelector(ctx, selector)
	if err != nil {
		return 0, fmt.Errorf("could not list nodes: %v", err)
	}
	return len(targetedNode), nil
}

func (n *node) UpdateLabels(ctx context.Context, node *v1.Node, toBeAdded, toBeRemoved []string) error {
	patchFrom := client.MergeFrom(node.DeepCopy())

	addLabels(node, toBeAdded)
	removeLabels(node, toBeRemoved)

	if err := n.client.Patch(ctx, node, patchFrom); err != nil {
		return fmt.Errorf("could not patch node: %v", err)
	}
	return nil
}

func (n *node) NodeBecomeReadyAfter(node *v1.Node, timestamp metav1.Time) bool {
	conds := node.Status.Conditions
	for i := 0; i < len(conds); i++ {
		c := conds[i]
		if c.Type == v1.NodeReady && c.Status == v1.ConditionTrue && timestamp.Before(&c.LastTransitionTime) {
			return true
		}
	}
	return false
}

func addLabels(node *v1.Node, labels []string) {
	for _, label := range labels {
		meta.SetLabel(
			node,
			label,
			"",
		)
	}
}

func removeLabels(node *v1.Node, labels []string) {
	for _, label := range labels {
		meta.RemoveLabel(
			node,
			label,
		)
	}
}