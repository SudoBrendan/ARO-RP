// Code generated by github.com/jewzaam/go-cosmosdb, DO NOT EDIT.

package cosmosdb

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/ugorji/go/codec"

	pkg "github.com/Azure/ARO-RP/pkg/api"
)

type fakeOpenShiftClusterManagerConfigurationDocumentTriggerHandler func(context.Context, *pkg.OpenShiftClusterManagerConfigurationDocument) error
type fakeOpenShiftClusterManagerConfigurationDocumentQueryHandler func(OpenShiftClusterManagerConfigurationDocumentClient, *Query, *Options) OpenShiftClusterManagerConfigurationDocumentRawIterator

var _ OpenShiftClusterManagerConfigurationDocumentClient = &FakeOpenShiftClusterManagerConfigurationDocumentClient{}

// NewFakeOpenShiftClusterManagerConfigurationDocumentClient returns a FakeOpenShiftClusterManagerConfigurationDocumentClient
func NewFakeOpenShiftClusterManagerConfigurationDocumentClient(h *codec.JsonHandle) *FakeOpenShiftClusterManagerConfigurationDocumentClient {
	return &FakeOpenShiftClusterManagerConfigurationDocumentClient{
		jsonHandle: h,
		openShiftClusterManagerConfigurationDocuments: make(map[string]*pkg.OpenShiftClusterManagerConfigurationDocument),
		triggerHandlers: make(map[string]fakeOpenShiftClusterManagerConfigurationDocumentTriggerHandler),
		queryHandlers:   make(map[string]fakeOpenShiftClusterManagerConfigurationDocumentQueryHandler),
	}
}

// FakeOpenShiftClusterManagerConfigurationDocumentClient is a FakeOpenShiftClusterManagerConfigurationDocumentClient
type FakeOpenShiftClusterManagerConfigurationDocumentClient struct {
	lock                                          sync.RWMutex
	jsonHandle                                    *codec.JsonHandle
	openShiftClusterManagerConfigurationDocuments map[string]*pkg.OpenShiftClusterManagerConfigurationDocument
	triggerHandlers                               map[string]fakeOpenShiftClusterManagerConfigurationDocumentTriggerHandler
	queryHandlers                                 map[string]fakeOpenShiftClusterManagerConfigurationDocumentQueryHandler
	sorter                                        func([]*pkg.OpenShiftClusterManagerConfigurationDocument)
	etag                                          int

	// returns true if documents conflict
	conflictChecker func(*pkg.OpenShiftClusterManagerConfigurationDocument, *pkg.OpenShiftClusterManagerConfigurationDocument) bool

	// err, if not nil, is an error to return when attempting to communicate
	// with this Client
	err error
}

// SetError sets or unsets an error that will be returned on any
// FakeOpenShiftClusterManagerConfigurationDocumentClient method invocation
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) SetError(err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.err = err
}

// SetSorter sets or unsets a sorter function which will be used to sort values
// returned by List() for test stability
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) SetSorter(sorter func([]*pkg.OpenShiftClusterManagerConfigurationDocument)) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.sorter = sorter
}

// SetConflictChecker sets or unsets a function which can be used to validate
// additional unique keys in a OpenShiftClusterManagerConfigurationDocument
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) SetConflictChecker(conflictChecker func(*pkg.OpenShiftClusterManagerConfigurationDocument, *pkg.OpenShiftClusterManagerConfigurationDocument) bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.conflictChecker = conflictChecker
}

// SetTriggerHandler sets or unsets a trigger handler
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) SetTriggerHandler(triggerName string, trigger fakeOpenShiftClusterManagerConfigurationDocumentTriggerHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.triggerHandlers[triggerName] = trigger
}

// SetQueryHandler sets or unsets a query handler
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) SetQueryHandler(queryName string, query fakeOpenShiftClusterManagerConfigurationDocumentQueryHandler) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.queryHandlers[queryName] = query
}

func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) deepCopy(openShiftClusterManagerConfigurationDocument *pkg.OpenShiftClusterManagerConfigurationDocument) (*pkg.OpenShiftClusterManagerConfigurationDocument, error) {
	var b []byte
	err := codec.NewEncoderBytes(&b, c.jsonHandle).Encode(openShiftClusterManagerConfigurationDocument)
	if err != nil {
		return nil, err
	}

	openShiftClusterManagerConfigurationDocument = nil
	err = codec.NewDecoderBytes(b, c.jsonHandle).Decode(&openShiftClusterManagerConfigurationDocument)
	if err != nil {
		return nil, err
	}

	return openShiftClusterManagerConfigurationDocument, nil
}

func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) apply(ctx context.Context, partitionkey string, openShiftClusterManagerConfigurationDocument *pkg.OpenShiftClusterManagerConfigurationDocument, options *Options, isCreate bool) (*pkg.OpenShiftClusterManagerConfigurationDocument, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.err != nil {
		return nil, c.err
	}

	openShiftClusterManagerConfigurationDocument, err := c.deepCopy(openShiftClusterManagerConfigurationDocument) // copy now because pretriggers can mutate openShiftClusterManagerConfigurationDocument
	if err != nil {
		return nil, err
	}

	if options != nil {
		err := c.processPreTriggers(ctx, openShiftClusterManagerConfigurationDocument, options)
		if err != nil {
			return nil, err
		}
	}

	existingOpenShiftClusterManagerConfigurationDocument, exists := c.openShiftClusterManagerConfigurationDocuments[openShiftClusterManagerConfigurationDocument.ID]
	if isCreate && exists {
		return nil, &Error{
			StatusCode: http.StatusConflict,
			Message:    "Entity with the specified id already exists in the system",
		}
	}
	if !isCreate {
		if !exists {
			return nil, &Error{StatusCode: http.StatusNotFound}
		}

		if openShiftClusterManagerConfigurationDocument.ETag != existingOpenShiftClusterManagerConfigurationDocument.ETag {
			return nil, &Error{StatusCode: http.StatusPreconditionFailed}
		}
	}

	if c.conflictChecker != nil {
		for _, openShiftClusterManagerConfigurationDocumentToCheck := range c.openShiftClusterManagerConfigurationDocuments {
			if c.conflictChecker(openShiftClusterManagerConfigurationDocumentToCheck, openShiftClusterManagerConfigurationDocument) {
				return nil, &Error{
					StatusCode: http.StatusConflict,
					Message:    "Entity with the specified id already exists in the system",
				}
			}
		}
	}

	openShiftClusterManagerConfigurationDocument.ETag = fmt.Sprint(c.etag)
	c.etag++

	c.openShiftClusterManagerConfigurationDocuments[openShiftClusterManagerConfigurationDocument.ID] = openShiftClusterManagerConfigurationDocument

	return c.deepCopy(openShiftClusterManagerConfigurationDocument)
}

// Create creates a OpenShiftClusterManagerConfigurationDocument in the database
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) Create(ctx context.Context, partitionkey string, openShiftClusterManagerConfigurationDocument *pkg.OpenShiftClusterManagerConfigurationDocument, options *Options) (*pkg.OpenShiftClusterManagerConfigurationDocument, error) {
	return c.apply(ctx, partitionkey, openShiftClusterManagerConfigurationDocument, options, true)
}

// Replace replaces a OpenShiftClusterManagerConfigurationDocument in the database
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) Replace(ctx context.Context, partitionkey string, openShiftClusterManagerConfigurationDocument *pkg.OpenShiftClusterManagerConfigurationDocument, options *Options) (*pkg.OpenShiftClusterManagerConfigurationDocument, error) {
	return c.apply(ctx, partitionkey, openShiftClusterManagerConfigurationDocument, options, false)
}

// List returns a OpenShiftClusterManagerConfigurationDocumentIterator to list all OpenShiftClusterManagerConfigurationDocuments in the database
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) List(*Options) OpenShiftClusterManagerConfigurationDocumentIterator {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.err != nil {
		return NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(c.err)
	}

	openShiftClusterManagerConfigurationDocuments := make([]*pkg.OpenShiftClusterManagerConfigurationDocument, 0, len(c.openShiftClusterManagerConfigurationDocuments))
	for _, openShiftClusterManagerConfigurationDocument := range c.openShiftClusterManagerConfigurationDocuments {
		openShiftClusterManagerConfigurationDocument, err := c.deepCopy(openShiftClusterManagerConfigurationDocument)
		if err != nil {
			return NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(err)
		}
		openShiftClusterManagerConfigurationDocuments = append(openShiftClusterManagerConfigurationDocuments, openShiftClusterManagerConfigurationDocument)
	}

	if c.sorter != nil {
		c.sorter(openShiftClusterManagerConfigurationDocuments)
	}

	return NewFakeOpenShiftClusterManagerConfigurationDocumentIterator(openShiftClusterManagerConfigurationDocuments, 0)
}

// ListAll lists all OpenShiftClusterManagerConfigurationDocuments in the database
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) ListAll(ctx context.Context, options *Options) (*pkg.OpenShiftClusterManagerConfigurationDocuments, error) {
	iter := c.List(options)
	return iter.Next(ctx, -1)
}

// Get gets a OpenShiftClusterManagerConfigurationDocument from the database
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) Get(ctx context.Context, partitionkey string, id string, options *Options) (*pkg.OpenShiftClusterManagerConfigurationDocument, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.err != nil {
		return nil, c.err
	}

	openShiftClusterManagerConfigurationDocument, exists := c.openShiftClusterManagerConfigurationDocuments[id]
	if !exists {
		return nil, &Error{StatusCode: http.StatusNotFound}
	}

	return c.deepCopy(openShiftClusterManagerConfigurationDocument)
}

// Delete deletes a OpenShiftClusterManagerConfigurationDocument from the database
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) Delete(ctx context.Context, partitionKey string, openShiftClusterManagerConfigurationDocument *pkg.OpenShiftClusterManagerConfigurationDocument, options *Options) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if c.err != nil {
		return c.err
	}

	_, exists := c.openShiftClusterManagerConfigurationDocuments[openShiftClusterManagerConfigurationDocument.ID]
	if !exists {
		return &Error{StatusCode: http.StatusNotFound}
	}

	delete(c.openShiftClusterManagerConfigurationDocuments, openShiftClusterManagerConfigurationDocument.ID)
	return nil
}

// ChangeFeed is unimplemented
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) ChangeFeed(*Options) OpenShiftClusterManagerConfigurationDocumentIterator {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.err != nil {
		return NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(c.err)
	}

	return NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(ErrNotImplemented)
}

func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) processPreTriggers(ctx context.Context, openShiftClusterManagerConfigurationDocument *pkg.OpenShiftClusterManagerConfigurationDocument, options *Options) error {
	for _, triggerName := range options.PreTriggers {
		if triggerHandler := c.triggerHandlers[triggerName]; triggerHandler != nil {
			c.lock.Unlock()
			err := triggerHandler(ctx, openShiftClusterManagerConfigurationDocument)
			c.lock.Lock()
			if err != nil {
				return err
			}
		} else {
			return ErrNotImplemented
		}
	}

	return nil
}

// Query calls a query handler to implement database querying
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) Query(name string, query *Query, options *Options) OpenShiftClusterManagerConfigurationDocumentRawIterator {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.err != nil {
		return NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(c.err)
	}

	if queryHandler := c.queryHandlers[query.Query]; queryHandler != nil {
		c.lock.RUnlock()
		i := queryHandler(c, query, options)
		c.lock.RLock()
		return i
	}

	return NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(ErrNotImplemented)
}

// QueryAll calls a query handler to implement database querying
func (c *FakeOpenShiftClusterManagerConfigurationDocumentClient) QueryAll(ctx context.Context, partitionkey string, query *Query, options *Options) (*pkg.OpenShiftClusterManagerConfigurationDocuments, error) {
	iter := c.Query("", query, options)
	return iter.Next(ctx, -1)
}

func NewFakeOpenShiftClusterManagerConfigurationDocumentIterator(openShiftClusterManagerConfigurationDocuments []*pkg.OpenShiftClusterManagerConfigurationDocument, continuation int) OpenShiftClusterManagerConfigurationDocumentRawIterator {
	return &fakeOpenShiftClusterManagerConfigurationDocumentIterator{openShiftClusterManagerConfigurationDocuments: openShiftClusterManagerConfigurationDocuments, continuation: continuation}
}

type fakeOpenShiftClusterManagerConfigurationDocumentIterator struct {
	openShiftClusterManagerConfigurationDocuments []*pkg.OpenShiftClusterManagerConfigurationDocument
	continuation                                  int
	done                                          bool
}

func (i *fakeOpenShiftClusterManagerConfigurationDocumentIterator) NextRaw(ctx context.Context, maxItemCount int, out interface{}) error {
	return ErrNotImplemented
}

func (i *fakeOpenShiftClusterManagerConfigurationDocumentIterator) Next(ctx context.Context, maxItemCount int) (*pkg.OpenShiftClusterManagerConfigurationDocuments, error) {
	if i.done {
		return nil, nil
	}

	var openShiftClusterManagerConfigurationDocuments []*pkg.OpenShiftClusterManagerConfigurationDocument
	if maxItemCount == -1 {
		openShiftClusterManagerConfigurationDocuments = i.openShiftClusterManagerConfigurationDocuments[i.continuation:]
		i.continuation = len(i.openShiftClusterManagerConfigurationDocuments)
		i.done = true
	} else {
		max := i.continuation + maxItemCount
		if max > len(i.openShiftClusterManagerConfigurationDocuments) {
			max = len(i.openShiftClusterManagerConfigurationDocuments)
		}
		openShiftClusterManagerConfigurationDocuments = i.openShiftClusterManagerConfigurationDocuments[i.continuation:max]
		i.continuation += max
		i.done = i.Continuation() == ""
	}

	return &pkg.OpenShiftClusterManagerConfigurationDocuments{
		OpenShiftClusterManagerConfigurationDocuments: openShiftClusterManagerConfigurationDocuments,
		Count: len(openShiftClusterManagerConfigurationDocuments),
	}, nil
}

func (i *fakeOpenShiftClusterManagerConfigurationDocumentIterator) Continuation() string {
	if i.continuation >= len(i.openShiftClusterManagerConfigurationDocuments) {
		return ""
	}
	return fmt.Sprintf("%d", i.continuation)
}

// NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator returns a OpenShiftClusterManagerConfigurationDocumentRawIterator which
// whose methods return the given error
func NewFakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator(err error) OpenShiftClusterManagerConfigurationDocumentRawIterator {
	return &fakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator{err: err}
}

type fakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator struct {
	err error
}

func (i *fakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator) Next(ctx context.Context, maxItemCount int) (*pkg.OpenShiftClusterManagerConfigurationDocuments, error) {
	return nil, i.err
}

func (i *fakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator) NextRaw(context.Context, int, interface{}) error {
	return i.err
}

func (i *fakeOpenShiftClusterManagerConfigurationDocumentErroringRawIterator) Continuation() string {
	return ""
}
