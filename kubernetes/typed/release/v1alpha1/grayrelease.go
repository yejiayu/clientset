/*
Copyright 2017 caicloud authors. All rights reserved.
*/

package v1alpha1

import (
	scheme "github.com/caicloud/clientset/kubernetes/scheme"
	v1alpha1 "github.com/caicloud/clientset/pkg/apis/release/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GrayReleasesGetter has a method to return a GrayReleaseInterface.
// A group's client should implement this interface.
type GrayReleasesGetter interface {
	GrayReleases(namespace string) GrayReleaseInterface
}

// GrayReleaseInterface has methods to work with GrayRelease resources.
type GrayReleaseInterface interface {
	Create(*v1alpha1.GrayRelease) (*v1alpha1.GrayRelease, error)
	Update(*v1alpha1.GrayRelease) (*v1alpha1.GrayRelease, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.GrayRelease, error)
	List(opts v1.ListOptions) (*v1alpha1.GrayReleaseList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GrayRelease, err error)
	GrayReleaseExpansion
}

// grayReleases implements GrayReleaseInterface
type grayReleases struct {
	client rest.Interface
	ns     string
}

// newGrayReleases returns a GrayReleases
func newGrayReleases(c *ReleaseV1alpha1Client, namespace string) *grayReleases {
	return &grayReleases{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Create takes the representation of a grayRelease and creates it.  Returns the server's representation of the grayRelease, and an error, if there is any.
func (c *grayReleases) Create(grayRelease *v1alpha1.GrayRelease) (result *v1alpha1.GrayRelease, err error) {
	result = &v1alpha1.GrayRelease{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("grayreleases").
		Body(grayRelease).
		Do().
		Into(result)
	return
}

// Update takes the representation of a grayRelease and updates it. Returns the server's representation of the grayRelease, and an error, if there is any.
func (c *grayReleases) Update(grayRelease *v1alpha1.GrayRelease) (result *v1alpha1.GrayRelease, err error) {
	result = &v1alpha1.GrayRelease{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("grayreleases").
		Name(grayRelease.Name).
		Body(grayRelease).
		Do().
		Into(result)
	return
}

// Delete takes name of the grayRelease and deletes it. Returns an error if one occurs.
func (c *grayReleases) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("grayreleases").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *grayReleases) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("grayreleases").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Get takes name of the grayRelease, and returns the corresponding grayRelease object, and an error if there is any.
func (c *grayReleases) Get(name string, options v1.GetOptions) (result *v1alpha1.GrayRelease, err error) {
	result = &v1alpha1.GrayRelease{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("grayreleases").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GrayReleases that match those selectors.
func (c *grayReleases) List(opts v1.ListOptions) (result *v1alpha1.GrayReleaseList, err error) {
	result = &v1alpha1.GrayReleaseList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("grayreleases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested grayReleases.
func (c *grayReleases) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("grayreleases").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Patch applies the patch and returns the patched grayRelease.
func (c *grayReleases) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.GrayRelease, err error) {
	result = &v1alpha1.GrayRelease{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("grayreleases").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
