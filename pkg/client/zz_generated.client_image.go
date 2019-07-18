/*
	Note: This file is autogenerated! Do not edit it manually!
	Edit client_image_template.go instead, and run
	hack/generate-client.sh afterwards.
*/

package client

import (
	log "github.com/sirupsen/logrus"
	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	meta "github.com/weaveworks/ignite/pkg/apis/meta/v1alpha1"
	"github.com/weaveworks/ignite/pkg/storage"
	"github.com/weaveworks/ignite/pkg/storage/filterer"
)

// ImageClient is an interface for accessing Image-specific API objects
type ImageClient interface {
	// Get returns the Image matching given UID from the storage
	Get(meta.UID) (*api.Image, error)
	// Set saves the given Image into persistent storage
	Set(*api.Image) error
	// Find returns the Image matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	Find(filter filterer.BaseFilter) (*api.Image, error)
	// FindAll returns multiple Images matching the given filter, filters can
	// match e.g. the Object's Name, UID or a specific property
	FindAll(filter filterer.BaseFilter) ([]*api.Image, error)
	// Delete deletes the Image with the given UID from the storage
	Delete(uid meta.UID) error
	// List returns a list of all Images available
	List() ([]*api.Image, error)
}

// Images returns the ImageClient for the Client instance
func (c *Client) Images() ImageClient {
	if c.imageClient == nil {
		c.imageClient = newImageClient(c.storage)
	}

	return c.imageClient
}

// imageClient is a struct implementing the ImageClient interface
// It uses a shared storage instance passed from the Client together with its own Filterer
type imageClient struct {
	storage  storage.Storage
	filterer *filterer.Filterer
}

// newImageClient builds the imageClient struct using the storage implementation and a new Filterer
func newImageClient(s storage.Storage) ImageClient {
	return &imageClient{
		storage:  s,
		filterer: filterer.NewFilterer(s),
	}
}

// Find returns a single Image based on the given Filter
func (c *imageClient) Find(filter filterer.BaseFilter) (*api.Image, error) {
	object, err := c.filterer.Find(api.KindImage, filter)
	if err != nil {
		return nil, err
	}

	return object.(*api.Image), nil
}

// FindAll returns multiple Images based on the given Filter
func (c *imageClient) FindAll(filter filterer.BaseFilter) ([]*api.Image, error) {
	matches, err := c.filterer.FindAll(api.KindImage, filter)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Image, 0, len(matches))
	for _, item := range matches {
		results = append(results, item.(*api.Image))
	}

	return results, nil
}

// Get returns the Image matching given UID from the storage
func (c *imageClient) Get(uid meta.UID) (*api.Image, error) {
	log.Debugf("Client.Get; UID: %q, Kind: %s", uid, api.KindImage)
	object, err := c.storage.GetByID(api.KindImage, uid)
	if err != nil {
		return nil, err
	}

	return object.(*api.Image), nil
}

// Set saves the given Image into the persistent storage
func (c *imageClient) Set(image *api.Image) error {
	log.Debugf("Client.Set; UID: %q, Kind: %s", image.GetUID(), image.GetKind())
	return c.storage.Set(image)
}

// Delete deletes the Image from the storage
func (c *imageClient) Delete(uid meta.UID) error {
	log.Debugf("Client.Delete; UID: %q, Kind: %s", uid, api.KindImage)
	return c.storage.Delete(api.KindImage, uid)
}

// List returns a list of all Images available
func (c *imageClient) List() ([]*api.Image, error) {
	list, err := c.storage.List(api.KindImage)
	if err != nil {
		return nil, err
	}

	results := make([]*api.Image, 0, len(list))
	for _, item := range list {
		results = append(results, item.(*api.Image))
	}

	return results, nil
}
