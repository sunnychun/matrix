package registry

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/ironzhang/matrix/tlog"
)

type Endpoint struct {
	Service string
	Addr    string
}

type Options struct {
	TTL       int64
	Timeout   time.Duration
	Namespace string
}

func New(c *clientv3.Client, opts Options) (*Registry, error) {
	log := tlog.Std().Sugar().With("namespace", opts.Namespace)

	if opts.TTL <= 0 {
		opts.TTL = 10 // default ttl: 10s
	}

	resp, err := c.Grant(withTimeout(opts.Timeout), opts.TTL)
	if err != nil {
		log.Errorw("grant", "error", err)
		return nil, err
	}
	if _, err = c.KeepAlive(withTimeout(opts.Timeout), resp.ID); err != nil {
		log.Errorw("keep alive", "error", err)
		return nil, err
	}
	log.Debugw("grant", "leaseID", resp.ID)

	return &Registry{
		client:    c,
		leaseID:   resp.ID,
		timeout:   opts.Timeout,
		namespace: opts.Namespace,
	}, nil
}

type Registry struct {
	client    *clientv3.Client
	leaseID   clientv3.LeaseID
	timeout   time.Duration
	namespace string
}

func (r *Registry) Register(p Endpoint) error {
	log := tlog.Std().Sugar().With("namespace", r.namespace, "service", p.Service, "addr", p.Addr)
	_, err := r.client.Put(withTimeout(r.timeout), r.key(p), "1", clientv3.WithLease(r.leaseID))
	if err != nil {
		log.Errorw("put", "error", err)
		return err
	}
	log.Debug("put")
	return nil
}

func (r *Registry) Unregister(p Endpoint) error {
	log := tlog.Std().Sugar().With("namespace", r.namespace, "service", p.Service, "addr", p.Addr)
	_, err := r.client.Delete(withTimeout(r.timeout), r.key(p))
	if err != nil {
		log.Errorw("delete", "error", err)
		return err
	}
	log.Debug("delete")
	return nil
}

func (r *Registry) Close() error {
	log := tlog.Std().Sugar().With("namespace", r.namespace)
	_, err := r.client.Revoke(withTimeout(r.timeout), r.leaseID)
	if err != nil {
		log.Errorw("revoke", "error", err)
	}
	log.Debugw("revoke", "leaseID", r.leaseID)
	return nil
}

func (r *Registry) key(p Endpoint) string {
	return fmt.Sprintf("%s/%s/%s", r.namespace, p.Service, p.Addr)
}

func withTimeout(timeout time.Duration) context.Context {
	if timeout <= 0 {
		return context.Background()
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	return ctx
}
