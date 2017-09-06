/*
Copyright 2017 caicloud authors. All rights reserved.
*/

package v1alpha1

import (
	scheme "github.com/caicloud/clientset/kubernetes/scheme"
	v1alpha1 "github.com/caicloud/clientset/pkg/apis/config/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ConfigReferencesGetter has a method to return a ConfigReferenceInterface.
// A group's client should implement this interface.
type ConfigReferencesGetter interface {
	ConfigReferences(namespace string) ConfigReferenceInterface
}

// ConfigReferenceInterface has methods to work with ConfigReference resources.
type ConfigReferenceInterface interface {
	Create(*v1alpha1.ConfigReference) (*v1alpha1.ConfigReference, error)
	Update(*v1alpha1.ConfigReference) (*v1alpha1.ConfigReference, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.ConfigReference, error)
	List(opts v1.ListOptions) (*v1alpha1.ConfigReferenceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ConfigReference, err error)
	ConfigReferenceExpansion
}

// configReferences implements ConfigReferenceInterface
type configReferences struct {
	client rest.Interface
	ns     string
}

// newConfigReferences returns a ConfigReferences
func newConfigReferences(c *ConfigV1alpha1Client, namespace string) *configReferences {
	return &configReferences{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a configReference and creates it.  Returns the server's representation of the configReference, and an error, if there is any.
func (c *configReferences) Create(configReference *v1alpha1.ConfigReference) (result *v1alpha1.ConfigReference, err error) {
	result = &v1alpha1.ConfigReference{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("configreferences").
		Body(configReference).
		Do().
		Into(result)
	return
}

// Update takes the representation of a configReference and updates it. Returns the server's representation of the configReference, and an error, if there is any.
func (c *configReferences) Update(configReference *v1alpha1.ConfigReference) (result *v1alpha1.ConfigReference, err error) {
	result = &v1alpha1.ConfigReference{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("configreferences").
		Name(configReference.Name).
		Body(configReference).
		Do().
		Into(result)
	return
}

// Delete takes name of the configReference and deletes it. Returns an error if one occurs.
func (c *configReferences) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("configreferences").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *configReferences) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("configreferences").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the configReference, and returns the corresponding configReference object, and an error if there is any.
func (c *configReferences) Get(name string, options v1.GetOptions) (result *v1alpha1.ConfigReference, err error) {
	result = &v1alpha1.ConfigReference{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("configreferences").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of ConfigReferences that match those selectors.
func (c *configReferences) List(opts v1.ListOptions) (result *v1alpha1.ConfigReferenceList, err error) {
	result = &v1alpha1.ConfigReferenceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("configreferences").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested configReferences.
func (c *configReferences) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("configreferences").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched configReference.
func (c *configReferences) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.ConfigReference, err error) {
	result = &v1alpha1.ConfigReference{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("configreferences").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
