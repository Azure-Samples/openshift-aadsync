package client

import (
	"io/ioutil"
	"net"
	"os"

	clientset "pkg/aadsync/client/clientset/versioned"

	aadgroupsyncv1 "pkg/aadsync/apis/aad.microsoft.com/v1"
	v1 "pkg/aadsync/client/clientset/versioned/typed/aad.microsoft.com/v1"

	logrus "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rest "k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
)

// Client contains the internal AAD Group Sync Client details
type Client struct {
	Log       *logrus.Entry
	Config    *rest.Config
	Client    v1.AADGroupSyncInterface
	Namespace string
}

// NewClient creates a new AAD Group Sync Client with default incluster configuration. You need to be running
// incluster for this to be successful
func NewClient(namespace string, log *logrus.Entry) *Client {

	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}

	return NewClientForConfigAndNamespace(config, namespace, log)
}

// NewClientForConfigAndNamespace creates a new AAD Group Sync Client with the specified configuration and namespace
func NewClientForConfigAndNamespace(config *rest.Config, namespace string, log *logrus.Entry) *Client {

	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	client := &Client{
		Log:    log,
		Config: config,
		Client: clientset.AadV1().AADGroupSyncs(string(namespace)),
	}
	log.Info("Created aad group sync client")
	log.Debugf("Host: %s", client.Config.Host)

	return client
}

// NewClientForLocal creates a new AAD Group Sync Client from local copies of incluster resources. This is useful
// for testing
func NewClientForLocal(namespace string, log *logrus.Entry) *Client {

	// Found incluster at /var/run/secrets/kubernetes.io/serviceaccount/token
	tokenFile := os.Getenv("KUBERNETES_SERVICEACCOUNT_TOKENFILE")

	// Found incluster at /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
	rootCAFile := os.Getenv("KUBERNETES_SERVICEACCOUNT_ROOTCAFILE")

	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")

	token, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Fatal(err)
	}

	tlsClientConfig := rest.TLSClientConfig{}
	if _, err := certutil.NewPool(rootCAFile); err != nil {
		log.Fatalf("Expected to load root CA config from %s, but got err: %v", rootCAFile, err)
	} else {
		tlsClientConfig.CAFile = rootCAFile
	}

	config := &rest.Config{
		Host:            "https://" + net.JoinHostPort(host, port),
		TLSClientConfig: tlsClientConfig,
		BearerToken:     string(token),
		BearerTokenFile: tokenFile,
	}

	return NewClientForConfigAndNamespace(config, namespace, log)
}

// Get returns an existing aadgroupsyncs.aad.microsoft.com CRD
func (c *Client) Get(aadGroupName string) (*aadgroupsyncv1.AADGroupSync, error) {

	c.Log.Infof("Fetching aadgroupsyncs.aad.microsoft.com: %s", aadGroupName)

	aadGroup, err := c.Client.Get(aadGroupName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			c.Log.Infof("Not found - aadgroupsyncs.aad.microsoft.com: %s", aadGroupName)
			return nil, nil
		}
		return nil, err
	}

	return aadGroup, nil
}

// Create creates a new aadgroupsyncs.aad.microsoft.com CRD
func (c *Client) Create(aadGroup *aadgroupsyncv1.AADGroupSync) (*aadgroupsyncv1.AADGroupSync, error) {

	c.Log.Infof("Creating aadgroupsyncs.aad.microsoft.com: %s", aadGroup.ObjectMeta.Name)

	aadGroup, err := c.Client.Create(aadGroup)
	if err != nil {
		return nil, err
	}

	return aadGroup, nil
}

// Update updates an existing aadgroupsyncs.aad.microsoft.com CRD
func (c *Client) Update(aadGroup *aadgroupsyncv1.AADGroupSync) (*aadgroupsyncv1.AADGroupSync, error) {

	c.Log.Infof("Updating aadgroupsyncs.aad.microsoft.com: %s", aadGroup.ObjectMeta.Name)

	aadGroup, err := c.Client.Update(aadGroup)
	if err != nil {
		return nil, err
	}

	return aadGroup, nil
}

// Delete deletes an existing aadgroupsyncs.aad.microsoft.com CRD
func (c *Client) Delete(aadGroupName string) error {

	c.Log.Infof("Deleting aadgroupsyncs.aad.microsoft.com: %s", aadGroupName)

	deletePolicy := metav1.DeletePropagationForeground
	err := c.Client.Delete(aadGroupName, &metav1.DeleteOptions{PropagationPolicy: &deletePolicy})
	if err != nil {
		return err
	}

	return nil
}

// List returns a collection of existing aadgroupsyncs.aad.microsoft.com CRDs
func (c *Client) List() ([]aadgroupsyncv1.AADGroupSync, error) {

	c.Log.Infof("Fetching all aadgroupsyncs.aad.microsoft.com")

	aadGroupList, err := c.Client.List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return aadGroupList.Items, nil
}
