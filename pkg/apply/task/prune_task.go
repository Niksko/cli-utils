// Copyright 2020 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package task

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/resource"
	"sigs.k8s.io/cli-utils/pkg/apply/prune"
	"sigs.k8s.io/cli-utils/pkg/apply/taskrunner"
	"sigs.k8s.io/cli-utils/pkg/common"
)

// PruneTask prunes objects from the cluster
// by using the PruneOptions. The provided Objects is the
// set of resources that have just been applied.
type PruneTask struct {
	PruneOptions      *prune.PruneOptions
	Objects           []*resource.Info
	DryRunStrategy    common.DryRunStrategy
	PropagationPolicy metav1.DeletionPropagation
}

// Start creates a new goroutine that will invoke
// the Run function on the PruneOptions to update
// the cluster. It will push a TaskResult on the taskChannel
// to signal to the taskrunner that the task has completed (or failed).
func (p *PruneTask) Start(taskContext *taskrunner.TaskContext) {
	go func() {
		err := p.PruneOptions.Prune(p.Objects, taskContext.EventChannel(),
			prune.Options{
				DryRunStrategy:    p.DryRunStrategy,
				PropagationPolicy: p.PropagationPolicy,
			})
		taskContext.TaskChannel() <- taskrunner.TaskResult{
			Err: err,
		}
	}()
}

// ClearTimeout is not supported by the PruneTask.
func (p *PruneTask) ClearTimeout() {}
