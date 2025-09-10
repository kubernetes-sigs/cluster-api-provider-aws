/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package helpers provides a set of utilities for testing controllers.
package helpers

import (
	"context"
	"flag"
	"fmt"
	"go/build"
	"net"
	"os"
	"path"
	"path/filepath"
	"regexp"
	goruntime "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/pkg/errors"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	logf "sigs.k8s.io/cluster-api/cmd/clusterctl/log"
	utilyaml "sigs.k8s.io/cluster-api/util/yaml"
)

const (
	mutatingWebhookKind          = "MutatingWebhookConfiguration"
	validatingWebhookKind        = "ValidatingWebhookConfiguration"
	defaultMutatingWebhookName   = "mutating-webhook-configuration"
	defaultValidatingWebhookName = "validating-webhook-configuration"
)

var (
	root                   string
	clusterAPIVersionRegex = regexp.MustCompile(`^(\W)sigs.k8s.io/cluster-api v(.+)`)
)

func init() {
	// reset flags to avoid conflicts if an imported package already called klog.InitFlags()
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	klog.InitFlags(nil)
	// additionally force all the controllers to use the Ginkgo logger.
	ctrl.SetLogger(klog.Background())
	logf.SetLogger(klog.Background())
	// add logger for ginkgo
	klog.SetOutput(ginkgo.GinkgoWriter)

	// Calculate the scheme.
	utilruntime.Must(apiextensionsv1.AddToScheme(scheme.Scheme))
	utilruntime.Must(admissionv1.AddToScheme(scheme.Scheme))
	utilruntime.Must(clusterv1.AddToScheme(scheme.Scheme))

	// Get the root of the current file to use in CRD paths.
	_, filename, _, _ := goruntime.Caller(0) //nolint:dogsled
	root = path.Join(path.Dir(filename), "..", "..")
}

type webhookConfiguration struct {
	tag              string
	relativeFilePath string
}

// TestEnvironmentConfiguration encapsulates the interim, mutable configuration of the Kubernetes local test environment.
type TestEnvironmentConfiguration struct {
	env                   *envtest.Environment
	webhookConfigurations []webhookConfiguration
}

// TestEnvironment encapsulates a Kubernetes local test environment.
type TestEnvironment struct {
	manager.Manager
	client.Client
	Config *rest.Config
	env    *envtest.Environment
	cancel context.CancelFunc
}

// Cleanup deletes all the given objects.
func (t *TestEnvironment) Cleanup(ctx context.Context, objs ...client.Object) error {
	errs := []error{}
	for _, o := range objs {
		err := t.Client.Delete(ctx, o)
		if apierrors.IsNotFound(err) {
			continue
		}
		errs = append(errs, err)
	}
	return kerrors.NewAggregate(errs)
}

// CreateNamespace creates a new namespace with a generated name.
func (t *TestEnvironment) CreateNamespace(ctx context.Context, generateName string) (*corev1.Namespace, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-", generateName),
			Labels: map[string]string{
				"testenv/original-name": generateName,
			},
		},
	}
	if err := t.Client.Create(ctx, ns); err != nil {
		return nil, err
	}

	return ns, nil
}

// NewTestEnvironmentConfiguration creates a new test environment configuration for running tests.
func NewTestEnvironmentConfiguration(crdDirectoryPaths []string) *TestEnvironmentConfiguration {
	resolvedCrdDirectoryPaths := make([]string, len(crdDirectoryPaths))

	for i, p := range crdDirectoryPaths {
		resolvedCrdDirectoryPaths[i] = path.Join(root, p)
	}

	if capiPath := getFilePathToCAPICRDs(root); capiPath != "" {
		resolvedCrdDirectoryPaths = append(resolvedCrdDirectoryPaths, capiPath)
	}

	return &TestEnvironmentConfiguration{
		env: &envtest.Environment{
			ErrorIfCRDPathMissing: true,
			CRDDirectoryPaths:     resolvedCrdDirectoryPaths,
		},
	}
}

// WithWebhookConfiguration adds the CRD webhook configuration given a named tag and file path for the webhook manifest.
func (t *TestEnvironmentConfiguration) WithWebhookConfiguration(tag string, relativeFilePath string) *TestEnvironmentConfiguration {
	t.webhookConfigurations = append(t.webhookConfigurations, webhookConfiguration{tag: tag, relativeFilePath: relativeFilePath})
	return t
}

// Build creates a new environment spinning up a local api-server.
// This function should be called only once for each package you're running tests within,
// usually the environment is initialized in a suite_test.go file within a `BeforeSuite` ginkgo block.
func (t *TestEnvironmentConfiguration) Build() (*TestEnvironment, error) {
	mutatingWebhooks := make([]*admissionv1.MutatingWebhookConfiguration, 0, len(t.webhookConfigurations))
	validatingWebhooks := make([]*admissionv1.ValidatingWebhookConfiguration, 0, len(t.webhookConfigurations))
	for _, w := range t.webhookConfigurations {
		m, v, err := buildModifiedWebhook(w.tag, w.relativeFilePath)
		if err != nil {
			return nil, err
		}
		if m.Webhooks != nil {
			// No mutating webhook defined.
			mutatingWebhooks = append(mutatingWebhooks, &m)
		}
		if v.Webhooks != nil {
			// No validating webhook defined.
			validatingWebhooks = append(validatingWebhooks, &v)
		}
	}

	t.env.WebhookInstallOptions = envtest.WebhookInstallOptions{
		MaxTime:            20 * time.Second,
		PollInterval:       time.Second,
		ValidatingWebhooks: validatingWebhooks,
		MutatingWebhooks:   mutatingWebhooks,
		LocalServingHost:   "localhost",
	}

	if _, err := t.env.Start(); err != nil {
		panic(err)
	}

	options := manager.Options{
		Scheme: scheme.Scheme,
		Metrics: metricsserver.Options{
			BindAddress: "0",
		},
		WebhookServer: webhook.NewServer(webhook.Options{
			Host:    t.env.WebhookInstallOptions.LocalServingHost,
			Port:    t.env.WebhookInstallOptions.LocalServingPort,
			CertDir: t.env.WebhookInstallOptions.LocalServingCertDir,
		}),
	}

	mgr, err := ctrl.NewManager(t.env.Config, options)
	if err != nil {
		klog.Fatalf("Failed to start testenv manager: %v", err)
	}

	return &TestEnvironment{
		Manager: mgr,
		Client:  mgr.GetClient(),
		Config:  mgr.GetConfig(),
		env:     t.env,
	}, nil
}

func buildModifiedWebhook(tag string, relativeFilePath string) (admissionv1.MutatingWebhookConfiguration, admissionv1.ValidatingWebhookConfiguration, error) {
	var mutatingWebhook admissionv1.MutatingWebhookConfiguration
	var validatingWebhook admissionv1.ValidatingWebhookConfiguration
	data, err := os.ReadFile(filepath.Clean(filepath.Join(root, relativeFilePath)))
	if err != nil {
		return mutatingWebhook, validatingWebhook, errors.Wrap(err, "failed to read webhook configuration file")
	}
	objs, err := utilyaml.ToUnstructured(data)
	if err != nil {
		return mutatingWebhook, validatingWebhook, errors.Wrap(err, "failed to parse yaml")
	}
	for i := range objs {
		o := objs[i]
		if o.GetKind() == mutatingWebhookKind {
			// update the name in metadata
			if o.GetName() == defaultMutatingWebhookName {
				o.SetName(defaultMutatingWebhookName + "-" + tag)
				if err := scheme.Scheme.Convert(&o, &mutatingWebhook, nil); err != nil {
					klog.Fatalf("failed to convert MutatingWebhookConfiguration %s", o.GetName())
				}
			}
		}
		if o.GetKind() == validatingWebhookKind {
			// update the name in metadata
			if o.GetName() == defaultValidatingWebhookName {
				o.SetName(defaultValidatingWebhookName + "-" + tag)
				if err := scheme.Scheme.Convert(&o, &validatingWebhook, nil); err != nil {
					klog.Fatalf("failed to convert ValidatingWebhookConfiguration %s", o.GetName())
				}
			}
		}
	}
	return mutatingWebhook, validatingWebhook, nil
}

// StartManager starts the test controller against the local API server.
func (t *TestEnvironment) StartManager(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	t.cancel = cancel
	return t.Manager.Start(ctx)
}

// WaitForWebhooks will not return until the webhook port is open.
func (t *TestEnvironment) WaitForWebhooks() {
	port := t.env.WebhookInstallOptions.LocalServingPort
	klog.V(2).Infof("Waiting for webhook port %d to be open prior to running tests", port)
	timeout := 1 * time.Second
	for {
		time.Sleep(1 * time.Second)
		conn, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), timeout)
		if err != nil {
			klog.V(2).Infof("Webhook port is not ready, will retry in %v: %s", timeout, err)
			continue
		}
		conn.Close()
		klog.V(2).Info("Webhook port is now open. Continuing with tests...")
		return
	}
}

// Stop stops the test environment.
func (t *TestEnvironment) Stop() error {
	t.cancel()
	return t.env.Stop()
}

func getFilePathToCAPICRDs(root string) string {
	modBits, err := os.ReadFile(filepath.Join(root, "go.mod")) //nolint:gosec
	if err != nil {
		return ""
	}

	var clusterAPIVersion string
	for _, line := range strings.Split(string(modBits), "\n") {
		matches := clusterAPIVersionRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			clusterAPIVersion = matches[2]
		}
	}

	if clusterAPIVersion == "" {
		return ""
	}

	gopath := envOr("GOPATH", build.Default.GOPATH)
	return filepath.Join(gopath, "pkg", "mod", "sigs.k8s.io", fmt.Sprintf("cluster-api@v%s", clusterAPIVersion), "config", "crd", "bases")
}

func envOr(envKey, defaultValue string) string {
	if value, ok := os.LookupEnv(envKey); ok {
		return value
	}
	return defaultValue
}
