package website

import (
	"context"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type Deployment struct {
	Name      string
	Namespace string
	GitRepo   string
}

func NewDeployment(name string, namespace string, gitRepo string) *Deployment {
	return &Deployment{
		Name:      name,
		Namespace: namespace,
		GitRepo:   gitRepo,
	}
}

func (deployment *Deployment) Deployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployment.Name,
			Namespace: deployment.Namespace,
			Labels: map[string]string{
				"webserver": deployment.Name,
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"webserver": deployment.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name: deployment.Name,
					Labels: map[string]string{
						"webserver": deployment.Name,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "main",
							Image: "nginx:alpine",
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "html",
									MountPath: "/usr/share/nginx/html",
									ReadOnly:  true,
								},
							},
							Ports: []apiv1.ContainerPort{
								{
									Protocol:      "TCP",
									ContainerPort: 80,
								},
							},
						},
						{
							Name:  "git-sync",
							Image: "openweb/git-sync",
							Env: []apiv1.EnvVar{
								{
									Name:  "GIT_SYNC_REPO",
									Value: deployment.GitRepo,
								},
								{
									Name:  "GIT_SYNC_DEST",
									Value: "/gitrepo",
								},
								{
									Name:  "GIT_SYNC_BRANCH",
									Value: "master",
								},
								{
									Name:  "GIT_SYNC_REV",
									Value: "FETCH_HEAD",
								},
								{
									Name:  "GIT_SYNC_WAIT",
									Value: "10",
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								{
									Name:      "html",
									MountPath: "/gitrepo",
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						{
							Name: "html",
							VolumeSource: apiv1.VolumeSource{
								EmptyDir: &apiv1.EmptyDirVolumeSource{
									Medium: "",
								},
							},
						},
					},
				},
			},
		},
	}
}

func CreateDeployment(client kubernetes.Interface, deployment *Deployment) {
	log.Printf("create deployment %s\n", deployment.Name)
	result, err := client.
		AppsV1().
		Deployments(deployment.Deployment().Namespace).
		Create(context.TODO(), deployment.Deployment(), metav1.CreateOptions{})
	if err != nil {
		log.Printf("create deployment %s with error %v\n", deployment.Name, err)
		return
	}
	log.Printf("create deployment %s with result %v\n", deployment.Name, result)
}

func DeleteDeployment(client kubernetes.Interface, deployment *Deployment) {
	log.Printf("delete deployment %s\n", deployment.Name)
	err := client.
		AppsV1().
		Deployments(deployment.Deployment().Namespace).
		Delete(context.TODO(), deployment.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("delete deployment %s with error %v\n", deployment.Name, err)
		return
	}
	log.Printf("delete deployment %s with success\n", deployment.Name)
}
