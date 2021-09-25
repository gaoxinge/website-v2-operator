package website

import (
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"

	internalExtensionV2 "github.com/gaoxinge/website-v2-operator/pkg/apis/extensions.example.com/v2"
	internalVersioned "github.com/gaoxinge/website-v2-operator/pkg/client/clientset/versioned"
	internalInformers "github.com/gaoxinge/website-v2-operator/pkg/client/informers/externalversions"
)

type Controller struct {
	clientSet         kubernetes.Interface
	internalClientSet internalVersioned.Interface
}

func NewController(clientSet kubernetes.Interface, internalClientSet internalVersioned.Interface) *Controller {
	websiteController := Controller{
		clientSet:         clientSet,
		internalClientSet: internalClientSet,
	}
	return &websiteController
}

func (controller *Controller) Run(stopCh chan  struct{}) error {
	websiteInformerFactory := internalInformers.NewSharedInformerFactory(controller.internalClientSet, 30 * time.Second)
	websiteInformer := websiteInformerFactory.Extensions().V2().Websites()
	websiteInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    controller.add,
		UpdateFunc: controller.update,
		DeleteFunc: controller.delete,
	})
	websiteInformerFactory.Start(stopCh)
	if !cache.WaitForCacheSync(stopCh, websiteInformer.Informer().HasSynced) {
		return errors.New("run controller with fail to sync")
	}
	return nil
}

func (controller *Controller) add(obj interface{}) {
	website := obj.(*internalExtensionV2.Website)
	CreateDeployment(controller.clientSet, NewDeployment(website.Name, website.Namespace, website.Spec.GitRepo))
	CreateService(controller.clientSet, NewService(website.Name, website.Namespace))
}

func (controller *Controller) update(oldObj, newObj interface{}) {

}

func (controller *Controller) delete(obj interface{}) {
	website := obj.(*internalExtensionV2.Website)
	DeleteDeployment(controller.clientSet, NewDeployment(website.Name, website.Namespace, website.Spec.GitRepo))
	DeleteService(controller.clientSet, NewService(website.Name, website.Namespace))
}
