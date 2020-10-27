package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"time"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	utilpointer "k8s.io/utils/pointer"
)

func main() {
	// Run ping in an Alpine container
	fn := Function{
		Name:      "ping",
		Namespace: "default",
		Spec: FunctionSpec{
			Image: "alpine:3.12",
			Args:  []string{"ping", "-c", "5", "google.com"},
		},
	}

	// Run a web scrape using mocha tests against inlets.dev

	// fn := Function{
	// 	Name:      "check-inlets",
	// 	Namespace: "default",
	// 	Spec: FunctionSpec{
	// 		Image: "alexellis2/check-inlets",
	// 		Args:  strings.Split("mocha --recursive ./integration-tests", " "),
	// 	},
	// }

	kubeconfig := path.Join(os.Getenv("HOME"), ".kube/config")
	job := FunctionToJob(&fn)

	log.Println("Accepted:", fn)
	fmt.Println(job.Name)

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Error building config: %v", err)
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error building kubernetes client: %v", err)
	}

	jobRes, err := client.BatchV1().Jobs(fn.Namespace).Create(context.Background(), job, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating job: %v", err)
	}

	log.Println(jobRes.Name)

	success := false
	attempts := 300
	sleepDuration := 1 * time.Second

	for i := 0; i < attempts; i++ {
		label := labels.SelectorFromSet(labels.Set(map[string]string{"job-name": job.Name}))
		list, err := client.CoreV1().Pods(fn.Namespace).List(context.Background(),
			metav1.ListOptions{LabelSelector: label.String()})
		if err != nil {
			log.Fatalf("Error getting pods: %v", err)
		}

		breakOut := false
		for _, p := range list.Items {
			log.Println(p.Status.Phase)
			if p.Status.Phase == v1.PodFailed || p.Status.Phase == v1.PodSucceeded {
				success = true
				breakOut = true
				break
			}
		}
		if breakOut {
			break
		}
		time.Sleep(sleepDuration)
	}

	if success {
		label := labels.SelectorFromSet(labels.Set(map[string]string{"job-name": job.Name}))
		list, err := client.CoreV1().Pods(fn.Namespace).List(context.Background(),
			metav1.ListOptions{LabelSelector: label.String()})
		if err != nil {
			log.Fatalf("Error getting pods: %v", err)
		}

		for _, p := range list.Items {
			r := client.CoreV1().Pods(fn.Namespace).GetLogs(p.Name, &v1.PodLogOptions{})

			stream, err := r.Stream(context.Background())
			if err != nil {
				log.Fatalf("Error getting logs stream: %v", err)
			}
			io.Copy(os.Stdout, stream)
		}
	}

	delOpt := metav1.DeletePropagationBackground
	err = client.BatchV1().Jobs(fn.Namespace).Delete(context.Background(), job.Name, metav1.DeleteOptions{
		PropagationPolicy: &delOpt,
	})

	if err != nil {
		log.Fatalf("Error deleting job: %v", err)
	}

	log.Printf("Deleting job %s..OK.", job.Name)
}

type Function struct {
	Name      string
	Namespace string
	Spec      FunctionSpec
}

type FunctionSpec struct {
	Image string
	Args  []string
}

func FunctionToJob(fn *Function) *batchv1.Job {
	probe := &corev1.Probe{
		InitialDelaySeconds: 1,
		PeriodSeconds:       1,
		SuccessThreshold:    1,
		TimeoutSeconds:      1,
		Handler: corev1.Handler{
			TCPSocket: &corev1.TCPSocketAction{
				Port: intstr.FromInt(9001),
			},
		},
	}

	container := corev1.Container{
		Name:            "faas",
		Image:           fn.Spec.Image,
		ImagePullPolicy: corev1.PullIfNotPresent,
		ReadinessProbe:  probe,
		Ports: []corev1.ContainerPort{
			{
				Name:          "gateway",
				ContainerPort: 9001,
				Protocol:      corev1.ProtocolTCP,
			},
		},
	}

	if len(fn.Spec.Args) > 0 {
		container.Args = fn.Spec.Args
	}

	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%v-%v", fn.Name, time.Now().Unix()),
			Namespace: fn.Namespace,
		},
		TypeMeta: metav1.TypeMeta{
			Kind: "Job",
		},
		Spec: batchv1.JobSpec{
			Parallelism:  utilpointer.Int32Ptr(1),
			Completions:  utilpointer.Int32Ptr(1),
			BackoffLimit: utilpointer.Int32Ptr(1),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						container,
					},
					RestartPolicy:                 corev1.RestartPolicyNever,
					TerminationGracePeriodSeconds: utilpointer.Int64Ptr(0),
				},
			},
		},
	}

	return job
}
