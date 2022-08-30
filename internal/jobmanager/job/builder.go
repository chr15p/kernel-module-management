package job

import (
	"fmt"

	kmmv1beta1 "github.com/rh-ecosystem-edge/kernel-module-management/api/v1beta1"
	"github.com/rh-ecosystem-edge/kernel-module-management/internal/jobmanager"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

//go:generate mockgen -source=maker.go -package=job -destination=mock_maker.go

type builder struct {
	name string
	helper build.Helper
	scheme *runtime.Scheme
}

func NewBuilder(helper build.Helper, scheme *runtime.Scheme) Job {
	return &builder{name: "Build", helper: helper, scheme: scheme}
}

func (m *builder) GetName() string {
	return m.name
}

func (m *builder) ShouldRun(mod *kmmv1beta1.Module, km *kmmv1beta1.KernelMapping) bool{
	if mod.Spec.ModuleLoader.Container.Build == nil && km.Build == nil {
		return false
	}
	return true
}


func (m *builder) PullOptions(km kmmv1beta1.KernelMapping) kmmv1beta1.PullOptions{
	return km.Build.Pull
}

func (m *builder) GetOutputImage(mod kmmv1beta1.Module, km *kmmv1beta1.KernelMapping) (string,error) {
        switch {
        case km.Sign.UnsignedImage != "":
                return km.Sign.UnsignedImage, nil
        case km.ContainerImage != "":
                return km.ContainerImage, nil
        default:
                return "",fmt.Errorf("Failed to find container image name")
        }
}


func (m *builder) MakeJob(mod kmmv1beta1.Module, km *kmmv1beta1.KernelMapping, targetKernel string) (*batchv1.Job, error) {

        buildConfig := m.helper.GetRelevantBuild(mod, *km)

        containerImage, err := m.GetOutputImage(mod,km)
        if err != nil {
                return nil, err
        }

	args := []string{"--destination", containerImage}

	buildArgs := m.helper.ApplyBuildArgOverrides(
		buildConfig.BuildArgs,
		kmmv1beta1.BuildArg{Name: "KERNEL_VERSION", Value: targetKernel},
	)

	for _, ba := range buildArgs {
		args = append(args, "--build-arg", fmt.Sprintf("%s=%s", ba.Name, ba.Value))
	}

	if buildConfig.Pull.Insecure {
		args = append(args, "--insecure-pull")
	}

	if buildConfig.Pull.InsecureSkipTLSVerify {
		args = append(args, "--skip-tls-verify-pull")
	}

	if buildConfig.Push.Insecure {
		args = append(args, "--insecure")
	}

	if buildConfig.Push.InsecureSkipTLSVerify {
		args = append(args, "--skip-tls-verify")
	}

	const dockerfileVolumeName = "dockerfile"

	dockerFileVolume := v1.Volume{
		Name: dockerfileVolumeName,
		VolumeSource: v1.VolumeSource{
			DownwardAPI: &v1.DownwardAPIVolumeSource{
				Items: []v1.DownwardAPIVolumeFile{
					{
						Path:     "Dockerfile",
						FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.annotations['Dockerfile']"},
					},
				},
			},
		},
	}

	dockerFileVolumeMount := v1.VolumeMount{
		Name:      dockerfileVolumeName,
		ReadOnly:  true,
		MountPath: "/workspace",
	}

	volumes := []v1.Volume{dockerFileVolume}
	volumeMounts := []v1.VolumeMount{dockerFileVolumeMount}
	if mod.Spec.ImageRepoSecret != nil {
		volumes = append(volumes, m.makeImagePullSecretVolume(mod.Spec.ImageRepoSecret))
		volumeMounts = append(volumeMounts, m.makeImagePullSecretVolumeMount(mod.Spec.ImageRepoSecret))
	}
	volumes = append(volumes, m.makeBuildSecretVolumes(buildConfig.Secrets)...)
	volumeMounts = append(volumeMounts, m.makeBuildSecretVolumeMounts(buildConfig.Secrets)...)

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: mod.Name + "-build-",
			Namespace:    mod.Namespace,
			Labels:       labels(mod, targetKernel, m.GetName()),
		},
		Spec: batchv1.JobSpec{
			Completions: pointer.Int32(1),
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{"Dockerfile": buildConfig.Dockerfile},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Args:         args,
							Name:         "kaniko",
							Image:        "gcr.io/kaniko-project/executor:latest",
							VolumeMounts: volumeMounts,
						},
					},
					NodeSelector:  mod.Spec.Selector,
					RestartPolicy: v1.RestartPolicyOnFailure,
					Volumes:       volumes,
				},
			},
		},
	}

	if err := controllerutil.SetControllerReference(&mod, job, m.scheme); err != nil {
		return nil, fmt.Errorf("could not set the owner reference: %v", err)
	}

	return job, nil
}

func (m *builder) makeImagePullSecretVolume(secretRef *v1.LocalObjectReference) v1.Volume {

	if secretRef == nil {
		return v1.Volume{}
	}

	return v1.Volume{
		Name: m.volumeNameFromSecretRef(*secretRef),
		VolumeSource: v1.VolumeSource{
			Secret: &v1.SecretVolumeSource{
				SecretName: secretRef.Name,
				Items: []v1.KeyToPath{
					{
						Key:  v1.DockerConfigJsonKey,
						Path: "config.json",
					},
				},
			},
		},
	}
}

func (m *builder) makeImagePullSecretVolumeMount(secretRef *v1.LocalObjectReference) v1.VolumeMount {

	if secretRef == nil {
		return v1.VolumeMount{}
	}

	return v1.VolumeMount{
		Name:      m.volumeNameFromSecretRef(*secretRef),
		ReadOnly:  true,
		MountPath: "/kaniko/.docker",
	}
}

func (m *builder) makeBuildSecretVolumes(secretRefs []v1.LocalObjectReference) []v1.Volume {

	volumes := make([]v1.Volume, 0, len(secretRefs))

	for _, secretRef := range secretRefs {
		vol := v1.Volume{
			Name: m.volumeNameFromSecretRef(secretRef),
			VolumeSource: v1.VolumeSource{
				Secret: &v1.SecretVolumeSource{
					SecretName: secretRef.Name,
				},
			},
		}

		volumes = append(volumes, vol)
	}

	return volumes
}

func (m *builder) makeBuildSecretVolumeMounts(secretRefs []v1.LocalObjectReference) []v1.VolumeMount {

	secretVolumeMounts := make([]v1.VolumeMount, 0, len(secretRefs))

	for _, secretRef := range secretRefs {
		volMount := v1.VolumeMount{
			Name:      m.volumeNameFromSecretRef(secretRef),
			ReadOnly:  true,
			MountPath: "/run/secrets/" + secretRef.Name,
		}

		secretVolumeMounts = append(secretVolumeMounts, volMount)
	}

	return secretVolumeMounts
}

func (m *builder) volumeNameFromSecretRef(ref v1.LocalObjectReference) string {
	return "secret-" + ref.Name
}
