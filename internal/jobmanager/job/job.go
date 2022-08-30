package job


import (
	kmmv1beta1 "github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1"
	batchv1 "k8s.io/api/batch/v1"
)


type Job interface {
	MakeJob(mod kmmv1beta1.Module, m *kmmv1beta1.KernelMapping, targetKernel string) (*batchv1.Job, error)
	PullOptions(km kmmv1beta1.KernelMapping) kmmv1beta1.PullOptions
	ShouldRun(mod *kmmv1beta1.Module, km *kmmv1beta1.KernelMapping) bool
	GetName() string
	GetOutputImage(mod kmmv1beta1.Module, km *kmmv1beta1.KernelMapping) (string,error)
}
