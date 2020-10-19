//
// Copyright (c) 2019-2020 Red Hat, Inc.
// This program and the accompanying materials are made
// available under the terms of the Eclipse Public License 2.0
// which is available at https://www.eclipse.org/legal/epl-2.0/
//
// SPDX-License-Identifier: EPL-2.0
//
// Contributors:
//   Red Hat, Inc. - initial API and implementation
//

package cluster

import (
	"context"
	"fmt"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"

	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/types"
	crclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// CleanJob cleans up a job in a given namespace. If a job didn't exist then it will throw isNotFound error
func CleanJob(client crclient.Client, name string, namespace string) error {
	job, err := getJobInNamespace(client, name, namespace)
	if err != nil {
		return err
	}
	err = deleteJob(client, job)
	if err != nil {
		return err
	}
	return nil
}

// Wait for the job to complete. Times out if the job isn't complete after $(timeout) seconds
func WaitForJobCompletion(client crclient.Client, name string, namespace string, timeout time.Duration) error {
	const interval = 1 * time.Second
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		job, err := getJobInNamespace(client, name, namespace)
		if err != nil {
			return false, err
		}

		if job.Status.Succeeded > 0 {
			return true, nil
		}
		if job.Status.Failed > 0 {
			log.Info(fmt.Sprintf("Job %s has failed when attempting to wait until its complete", job.Name))
		}
		return false, nil
	})
}

func SyncJobToCluster(
	client crclient.Client,
	ctx context.Context,
	specJob *batchv1.Job,
) error {
	// Attempt to clean up the job before re-creating it. If a job didn't exist then it will throw isNotFound error
	err := CleanJob(client, specJob.Name, specJob.Namespace)
	if err != nil && !apierrors.IsNotFound(err) {
		return err
	}
	err = client.Create(ctx, specJob)
	if err != nil {
		return err
	}
	log.Info("Created Job '" + specJob.GetName() + "'")
	return nil
}

// getJobInNamespace finds a job with a given name in a namespace
func getJobInNamespace(client crclient.Client, name string, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: name, Namespace: namespace}, job)
	if err != nil {
		return job, err
	}
	return job, nil
}

// deleteJob deletes a given job and cleans up any pods associated with it
func deleteJob(c crclient.Client, job *batchv1.Job) error {
	deleteBackground := metav1.DeletePropagationBackground
	err := c.Delete(context.TODO(), job, &crclient.DeleteOptions{PropagationPolicy: &deleteBackground})
	if err != nil {
		log.Error(err, "Error deleting job: "+job.Name)
		return err
	}
	//we don't use job.Selector here since it's autogenerated and contains only generated controller-uid unique for each job instance
	//if we use job-name which is also generated label but not used in selector - we clean up pods from old jobs
	pods, err := GetPodsBySelector(c, job.Namespace, "job-name="+job.Name)
	if err != nil {
		return err
	}
	for _, pod := range pods.Items {
		log.Info(fmt.Sprintf("Job '%s/%s' related pod %s is not removed automatically", job.Namespace, job.Name, pod.Name))
		err = DeletePod(c, job.Namespace, pod.Name)
		if err != nil {
			return err
		}
	}
	return nil
}
