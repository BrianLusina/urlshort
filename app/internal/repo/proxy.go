package repo

import (
	"log"
	"net/rpc"
)

// ProxyStore that will forward requests to the RPC server
// replica(s) will use the proxy store and only the primary will use the UrlStore
type ProxyStore struct {
	// store is a local cache of the data. Replicas use this in order to handle Get requests
	store  *UrlStore
	client *rpc.Client
}

func NewProxyStore(addr string) *ProxyStore {
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		log.Println("Error connecting to ProxyStore: ", err)
	}
	return &ProxyStore{store: NewUrlStore(""), client: client}
}

func (ps *ProxyStore) Get(key, url *string) error {
	if err := ps.store.Get(key, url); err == nil {
		return nil
	}
	// rpc call to primary
	if err := ps.client.Call("Store.Get", key, url); err != nil {
		return err
	}
	ps.store.Set(key, url) // update local cache
	return nil
}

func (ps *ProxyStore) Put(longUrl, key *string) error {
	// rpc call to master:
	if err := ps.client.Call("Store.Put", longUrl, key); err != nil {
		return err
	}
	ps.store.Set(key, longUrl) // update local cache
	return nil
}
