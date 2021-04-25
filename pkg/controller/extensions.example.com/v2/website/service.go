package website

import (
	"context"
	"log"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type Service struct {
	Name      string
	Namespace string
}

func NewService(name string, namespace string) *Service {
	return &Service{Name: name, Namespace: namespace}
}

func (service *Service) Service() *apiv1.Service {
	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      service.Name,
			Namespace: service.Namespace,
			Labels: map[string]string{
				"webserver": service.Name,
			},
		},
		Spec: apiv1.ServiceSpec{
			Type: "NodePort",
			Selector: map[string]string{
				"webserver": service.Name,
			},
			Ports: []apiv1.ServicePort{
				{
					Protocol: "TCP",
					Port:     80,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 80,
					},
				},
			},
		},
	}
}

func CreateService(client kubernetes.Interface, service *Service) {
	log.Printf("create service %s", service.Name)
	result, err := client.
		CoreV1().
		Services(service.Service().Namespace).
		Create(context.TODO(), service.Service(), metav1.CreateOptions{})
	if err != nil {
		log.Printf("create service %s with error %v", service.Name, err)
		return
	}
	log.Printf("create service %s with result %v", service.Name, result)
}

func DeleteService(client kubernetes.Interface, service *Service) {
	log.Printf("delete service %s", service.Name)
	err := client.
		CoreV1().
		Services(service.Service().Namespace).
		Delete(context.TODO(), service.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("delete service %s with error %v", service.Name, err)
		return
	}
	log.Printf("delete service %s with success", service.Name)
}
