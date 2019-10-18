package backend

import (
	"context"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/sirupsen/logrus"

	"github.com/jim-minter/rp/pkg/api"
	"github.com/jim-minter/rp/pkg/database"
	"github.com/jim-minter/rp/pkg/queue"
)

const (
	maxWorkers      = 100
	maxDequeueCount = 5
)

type backend struct {
	baseLog *logrus.Entry
	q       queue.Queue
	db      database.OpenShiftClusters

	authorizer autorest.Authorizer

	mu       sync.Mutex
	cond     *sync.Cond
	workers  int32
	stopping atomic.Value
}

// Runnable represents a runnable object
type Runnable interface {
	Run(stop <-chan struct{})
}

// NewBackend returns a new runnable backend
func NewBackend(log *logrus.Entry, authorizer autorest.Authorizer, q queue.Queue, db database.OpenShiftClusters) Runnable {
	b := &backend{
		baseLog: log,
		q:       q,
		db:      db,

		authorizer: authorizer,
	}

	b.cond = sync.NewCond(&b.mu)
	b.stopping.Store(false)

	return b
}

func (b *backend) Run(stop <-chan struct{}) {
	t := time.NewTicker(time.Second)
	defer t.Stop()

	go func() {
		<-stop
		b.baseLog.Print("stopping")
		b.stopping.Store(true)
		b.cond.Signal()
	}()

	for {
		b.mu.Lock()
		for atomic.LoadInt32(&b.workers) == maxWorkers && !b.stopping.Load().(bool) {
			b.cond.Wait()
		}
		b.mu.Unlock()

		if b.stopping.Load().(bool) {
			break
		}

		m, err := b.q.Get()
		if err != nil || m == nil {
			if err != nil {
				b.baseLog.Error(err)
			}
			<-t.C
			continue
		}

		log := b.baseLog.WithField("resource", m.ID())
		if m.DequeueCount() == maxDequeueCount {
			log.Warnf("dequeued %d times, failing", maxDequeueCount)
			b.setTerminalState(m.ID(), api.ProvisioningStateFailed)
			m.Done(nil)
		} else {
			log.Print("dequeued")
			go b.handle(context.Background(), log, m)
		}
	}
}

func (b *backend) handle(ctx context.Context, log *logrus.Entry, m queue.Message) {
	t := time.Now()

	var err error
	defer func() { // must add a closure to delay evaluation of err
		if e := recover(); e != nil {
			err = fmt.Errorf("panic: %#v\n%s", e, string(debug.Stack()))
			log.Error(err)
		}

		m.Done(err)
	}()

	atomic.AddInt32(&b.workers, 1)

	defer func() {
		atomic.AddInt32(&b.workers, -1)
		b.cond.Signal()
	}()

	err = b.handleOne(ctx, log, m)
	log = log.WithField("durationMs", int(time.Now().Sub(t)/time.Millisecond))
	if err != nil {
		log.Error(err)
	} else {
		log.Print("done")
	}
}

func (b *backend) handleOne(ctx context.Context, log *logrus.Entry, m queue.Message) error {
	doc, err := b.db.Get(m.ID())
	if err != nil {
		return err
	}

	switch doc.OpenShiftCluster.Properties.ProvisioningState {
	case api.ProvisioningStateUpdating:
		log.Printf("updating")
		err = b.update(ctx, log, doc)
	case api.ProvisioningStateDeleting:
		log.Printf("deleting")
		err = b.delete(ctx, log, doc)
	}
	if err != nil {
		log.Error(err)
		return b.setTerminalState(m.ID(), api.ProvisioningStateFailed)
	}

	switch doc.OpenShiftCluster.Properties.ProvisioningState {
	case api.ProvisioningStateUpdating:
		return b.setTerminalState(m.ID(), api.ProvisioningStateSucceeded)

	case api.ProvisioningStateDeleting:
		return b.db.Delete(m.ID())

	default:
		return fmt.Errorf("unexpected state %q", doc.OpenShiftCluster.Properties.ProvisioningState)
	}
}

func (b *backend) setTerminalState(resourceID string, state api.ProvisioningState) error {
	_, err := b.db.Patch(resourceID, func(doc *api.OpenShiftClusterDocument) error {
		doc.OpenShiftCluster.Properties.ProvisioningState = state
		return nil
	})
	return err
}
