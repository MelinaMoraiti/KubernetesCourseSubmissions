package main

import (
	"context"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"
)

var gvr = schema.GroupVersionResource{Group: "stable.dwk", Version: "v1", Resource: "dummysites"}

// A ConfigMap can hold at most 1 MiB, so leave some headroom.
const maxPageBytes = 900 * 1024

// Matches <head> or <head attr...> but not <header>.
var headRe = regexp.MustCompile(`(?i)<head(\s[^>]*)?>`)

type controller struct {
	kube     kubernetes.Interface
	informer cache.SharedIndexInformer
	queue    workqueue.TypedRateLimitingInterface[string]
	http     *http.Client
}

func main() {
	cfg, err := restConfig()
	if err != nil {
		log.Fatalf("loading kube config: %v", err)
	}
	dyn := dynamic.NewForConfigOrDie(cfg)
	kube := kubernetes.NewForConfigOrDie(cfg)

	factory := dynamicinformer.NewDynamicSharedInformerFactory(dyn, 10*time.Minute)
	informer := factory.ForResource(gvr).Informer()

	c := &controller{
		kube:     kube,
		informer: informer,
		queue:    workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[string]()),
		http:     &http.Client{Timeout: 20 * time.Second},
	}

	// Add and Update (including the periodic resync) trigger a reconcile.
	// Deletes need no handling: child objects are garbage-collected via ownerReferences.
	_, err = informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    c.enqueue,
		UpdateFunc: func(_, obj interface{}) { c.enqueue(obj) },
	})
	if err != nil {
		log.Fatalf("adding event handler: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	factory.Start(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(), informer.HasSynced) {
		log.Fatal("cache failed to sync")
	}
	log.Println("controller started, watching dummysites.stable.dwk")

	go func() {
		<-ctx.Done()
		c.queue.ShutDown()
	}()
	for c.processNext(ctx) {
	}
	log.Println("shutting down")
}

func restConfig() (*rest.Config, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}
	rules := clientcmd.NewDefaultClientConfigLoadingRules() // honours $KUBECONFIG
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{}).ClientConfig()
}

func (c *controller) enqueue(obj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		log.Printf("cannot compute key: %v", err)
		return
	}
	c.queue.Add(key)
}

func (c *controller) processNext(ctx context.Context) bool {
	key, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(key)

	if err := c.sync(ctx, key); err != nil {
		log.Printf("sync %s failed, will retry: %v", key, err)
		c.queue.AddRateLimited(key)
		return true
	}
	c.queue.Forget(key)
	return true
}

func (c *controller) sync(ctx context.Context, key string) error {
	obj, exists, err := c.informer.GetIndexer().GetByKey(key)
	if err != nil {
		return err
	}
	if !exists {
		return nil // deleted
	}
	ds, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return fmt.Errorf("unexpected object type %T", obj)
	}

	url, found, err := unstructured.NestedString(ds.Object, "spec", "website_url")
	if err != nil || !found || url == "" {
		log.Printf("%s: spec.website_url missing, skipping", key)
		return nil // retrying will not help
	}

	log.Printf("%s: fetching %s", key, url)
	page, finalURL, err := c.fetch(ctx, url)
	if err != nil {
		return fmt.Errorf("fetching %s: %w", url, err)
	}
	page = injectBase(page, finalURL)

	return c.apply(ctx, ds, page)
}

func (c *controller) fetch(ctx context.Context, url string) (body, finalURL string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("User-Agent", "dummysite-controller/1.0")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("unexpected status %s", resp.Status)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxPageBytes+1))
	if err != nil {
		return "", "", err
	}
	if len(b) > maxPageBytes {
		return "", "", fmt.Errorf("page larger than %d bytes, does not fit in a ConfigMap", maxPageBytes)
	}
	return string(b), resp.Request.URL.String(), nil
}

// injectBase adds <base href> so relative CSS, JS and image links load from the
// original site. The copy is not complete, but it renders far less broken.
func injectBase(page, url string) string {
	tag := fmt.Sprintf(`<base href="%s">`, html.EscapeString(url))
	loc := headRe.FindStringIndex(page)
	if loc == nil {
		return tag + page
	}
	return page[:loc[1]] + tag + page[loc[1]:]
}

func (c *controller) apply(ctx context.Context, ds *unstructured.Unstructured, page string) error {
	ns, name := ds.GetNamespace(), ds.GetName()
	labels := map[string]string{"app": "dummysite-" + name}

	meta := func(suffix string) metav1.ObjectMeta {
		return metav1.ObjectMeta{
			Name:      name + "-" + suffix,
			Namespace: ns,
			Labels:    labels,
			OwnerReferences: []metav1.OwnerReference{{
				APIVersion: gvr.GroupVersion().String(),
				Kind:       "DummySite",
				Name:       name,
				UID:        ds.GetUID(),
				Controller: boolPtr(true),
			}},
		}
	}

	// ConfigMap holding the HTML
	cms := c.kube.CoreV1().ConfigMaps(ns)
	cm := &corev1.ConfigMap{ObjectMeta: meta("html"), Data: map[string]string{"index.html": page}}
	if _, err := cms.Create(ctx, cm, metav1.CreateOptions{}); apierrors.IsAlreadyExists(err) {
		_, err = cms.Update(ctx, cm, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("updating configmap: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("creating configmap: %w", err)
	}

	// Deployment serving it with nginx
	deps := c.kube.AppsV1().Deployments(ns)
	dep := &appsv1.Deployment{
		ObjectMeta: meta("site"),
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{MatchLabels: labels},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: labels},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:         "nginx",
						Image:        "nginx:alpine",
						Ports:        []corev1.ContainerPort{{ContainerPort: 80}},
						VolumeMounts: []corev1.VolumeMount{{Name: "html", MountPath: "/usr/share/nginx/html"}},
					}},
					Volumes: []corev1.Volume{{
						Name: "html",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{Name: name + "-html"},
							},
						},
					}},
				},
			},
		},
	}
	if _, err := deps.Create(ctx, dep, metav1.CreateOptions{}); apierrors.IsAlreadyExists(err) {
		_, err = deps.Update(ctx, dep, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("updating deployment: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("creating deployment: %w", err)
	}

	// Service (create only; nothing in it ever changes)
	svc := &corev1.Service{
		ObjectMeta: meta("svc"),
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports:    []corev1.ServicePort{{Port: 80, TargetPort: intstr.FromInt32(80)}},
		},
	}
	if _, err := c.kube.CoreV1().Services(ns).Create(ctx, svc, metav1.CreateOptions{}); err != nil && !apierrors.IsAlreadyExists(err) {
		return fmt.Errorf("creating service: %w", err)
	}

	log.Printf("%s/%s: site ready (%d bytes)", ns, name, len(page))
	return nil
}

func int32Ptr(i int32) *int32 { return &i }
func boolPtr(b bool) *bool    { return &b }
