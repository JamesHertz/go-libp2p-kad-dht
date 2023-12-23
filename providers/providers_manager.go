package providers

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
	// "github.com/libp2p/go-libp2p/core/peerstore"
	// peerstoreImpl "github.com/libp2p/go-libp2p/p2p/host/peerstore"

	lru "github.com/hashicorp/golang-lru/simplelru"
	ds "github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/autobatch"
	dsq "github.com/ipfs/go-datastore/query"
	logging "github.com/ipfs/go-log"
	"github.com/jbenet/goprocess"
	goprocessctx "github.com/jbenet/goprocess/context"
	"github.com/multiformats/go-base32"
)

const (
	// ProvidersKeyPrefix is the prefix/namespace for ALL provider record
	// keys stored in the data store.
	ProvidersKeyPrefix = "/providers/"

	// ProviderAddrTTL is the TTL to keep the multi addresses of provider
	// peers around. Those addresses are returned alongside provider. After
	// it expires, the returned records will require an extra lookup, to
	// find the multiaddress associated with the returned peer id.
	ProviderAddrTTL = 24 * time.Hour
)

// ProvideValidity is the default time that a Provider Record should last on DHT
// This value is also known as Provider Record Expiration Interval.
var ProvideValidity = time.Hour * 48 // 24??
var defaultCleanupInterval = time.Hour
var lruCacheSize = 256
var batchBufferSize = 256
var log = logging.Logger("providers")

// ProviderStore represents a store that associates peers and their addresses to keys.
type ProviderStore interface {
	AddProvider(ctx context.Context, key []byte, prov peer.ID) error
	GetProviders(ctx context.Context, key []byte) ([]peer.ID, error)
	GetProvidersForPrefix(ctx context.Context, key []byte, prefixBitLength int) (map[string][]peer.ID, error)
}

// ProviderManager adds and pulls providers out of the datastore,
// caching them in between
type ProviderManager struct {
	self peer.ID
	// all non channel fields are meant to be accessed only within
	// the run method
	cache  lru.LRUCache
	dstore *autobatch.Datastore

    newprovs         chan *addProv
    getprovs         chan *getProv
    getProvsByPrefix chan *getProvByPrefix
    proc             goprocess.Process

	cleanupInterval time.Duration
}

type getProvByPrefix struct {
	ctx             context.Context
	key             []byte
	prefixBitLength int
	resp            chan map[string][]peer.ID
}

var _ ProviderStore = (*ProviderManager)(nil)

// Option is a function that sets a provider manager option.
type Option func(*ProviderManager) error

func (pm *ProviderManager) applyOptions(opts ...Option) error {
	for i, opt := range opts {
		if err := opt(pm); err != nil {
			return fmt.Errorf("provider manager option %d failed: %s", i, err)
		}
	}
	return nil
}

// CleanupInterval sets the time between GC runs.
// Defaults to 1h.
func CleanupInterval(d time.Duration) Option {
	return func(pm *ProviderManager) error {
		pm.cleanupInterval = d
		return nil
	}
}

// Cache sets the LRU cache implementation.
// Defaults to a simple LRU cache.
func Cache(c lru.LRUCache) Option {
	return func(pm *ProviderManager) error {
		pm.cache = c
		return nil
	}
}

type addProv struct {
	ctx context.Context
	key []byte
	val peer.ID
}

type getProv struct {
	ctx  context.Context
	key  []byte
	resp chan []peer.ID
}

// NewProviderManager constructor
func NewProviderManager(ctx context.Context, local peer.ID, dstore ds.Batching, opts ...Option) (*ProviderManager, error) {

	pm := new(ProviderManager)
	pm.self = local
	pm.getprovs = make(chan *getProv)
	pm.newprovs = make(chan *addProv)
	pm.getProvsByPrefix = make(chan *getProvByPrefix)
	pm.dstore = autobatch.NewAutoBatching(dstore, batchBufferSize)
	cache, err := lru.NewLRU(lruCacheSize, nil)
	if err != nil {
		return nil, err
	}
	pm.cache = cache
	pm.cleanupInterval = defaultCleanupInterval
	if err := pm.applyOptions(opts...); err != nil {
		return nil, err
	}
	pm.proc = goprocessctx.WithContext(ctx)
	pm.proc.Go(func(proc goprocess.Process) { pm.run(ctx, proc) })
	return pm, nil
}

// Process returns the ProviderManager process
func (pm *ProviderManager) Process() goprocess.Process {
	return pm.proc
}

func (pm *ProviderManager) run(ctx context.Context, proc goprocess.Process) {
	var (
		gcQuery    dsq.Results
		gcQueryRes <-chan dsq.Result
		gcSkip     map[string]struct{}
		gcTime     time.Time
		gcTimer    = time.NewTimer(pm.cleanupInterval)
	)

	defer func() {
		gcTimer.Stop()
		if gcQuery != nil {
			// don't really care if this fails.
			_ = gcQuery.Close()
		}
		if err := pm.dstore.Flush(ctx); err != nil {
			log.Error("failed to flush datastore: ", err)
		}
	}()

	for {
		select {
		case np := <-pm.newprovs:
			err := pm.addProv(np.ctx, np.key, np.val)
			if err != nil {
				log.Error("error adding new providers: ", err)
				continue
			}
			if gcSkip != nil {
				// we have an gc, tell it to skip this provider
				// as we've updated it since the GC started.
				gcSkip[mkProvKeyFor(np.key, np.val)] = struct{}{}
			}
		case gp := <-pm.getprovs:
			provs, err := pm.getProvidersForKey(gp.ctx, gp.key)
			if err != nil && err != ds.ErrNotFound {
				log.Error("error reading providers: ", err)
			}

			// set the cap so the user can't append to this.
			gp.resp <- provs[0:len(provs):len(provs)]
			continue


		case gp := <-pm.getProvsByPrefix:
			provs, err := pm.getProviderSetForPrefix(gp.ctx, gp.key, gp.prefixBitLength)
			if err != nil && err != ds.ErrNotFound {
				log.Error("error reading providers: ", err)
			}

			// set the cap so the user can't append to this.
			gp.resp <- provs.keyToProviders
			continue

		case res, ok := <-gcQueryRes:
			if !ok {
				if err := gcQuery.Close(); err != nil {
					log.Error("failed to close provider GC query: ", err)
				}
				gcTimer.Reset(pm.cleanupInterval)

				// cleanup GC round
				gcQueryRes = nil
				gcSkip = nil
				gcQuery = nil
				continue
			}
			if res.Error != nil {
				log.Error("got error from GC query: ", res.Error)
				continue
			}
			if _, ok := gcSkip[res.Key]; ok {
				// We've updated this record since starting the
				// GC round, skip it.
				continue
			}

			// check expiration time
			t, err := readTimeValue(res.Value)
			switch {
			case err != nil:
				// couldn't parse the time
				log.Error("parsing providers record from disk: ", err)
				fallthrough
			case gcTime.Sub(t) > ProvideValidity:
				// or expired
				err = pm.dstore.Delete(ctx, ds.RawKey(res.Key))
				if err != nil && err != ds.ErrNotFound {
					log.Error("failed to remove provider record from disk: ", err)
				}
			}

		case gcTime = <-gcTimer.C:
			// You know the wonderful thing about caches? You can
			// drop them.
			//
			// Much faster than GCing.
			pm.cache.Purge()

			// Now, kick off a GC of the datastore.
			q, err := pm.dstore.Query(ctx, dsq.Query{
				Prefix: ProvidersKeyPrefix,
			})
			if err != nil {
				log.Error("provider record GC query failed: ", err)
				continue
			}
			gcQuery = q
			gcQueryRes = q.Next()
			gcSkip = make(map[string]struct{})
		case <-proc.Closing():
			return
		}
	}
}

// AddProvider adds a provider
func (pm *ProviderManager) AddProvider(ctx context.Context, k []byte, provInfo peer.ID) error { 
	prov := &addProv{
		ctx: ctx,
		key: k,
		val: provInfo,
	}
	select {
	case pm.newprovs <- prov:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// addProv updates the cache if needed
func (pm *ProviderManager) addProv(ctx context.Context, k []byte, p peer.ID) error {
	now := time.Now()
	if provs, ok := pm.cache.Get(string(k)); ok {
		provs.(*providerSet).setVal(p, k, now)

	} // else not cached, just write through

	return writeProviderEntry(ctx, pm.dstore, k, p, now)
}

// writeProviderEntry writes the provider into the datastore
func writeProviderEntry(ctx context.Context, dstore ds.Datastore, k []byte, p peer.ID, t time.Time) error {
	dsk := mkProvKeyFor(k, p)

	buf := make([]byte, 16)
	n := binary.PutVarint(buf, t.UnixNano())

	return dstore.Put(ctx, ds.NewKey(dsk), buf[:n])
}

func mkProvKeyFor(k []byte, p peer.ID) string {
	return mkProvKey(k) + "/" + base32.RawStdEncoding.EncodeToString([]byte(p))
}

func mkProvKey(k []byte) string {
	return ProvidersKeyPrefix + hex.EncodeToString(k)
}

// GetProviders returns the set of providers for the given key.
// This method _does not_ copy the set. Do not modify it.
func (pm *ProviderManager) GetProviders(ctx context.Context, k []byte) ([]peer.ID, error) {
	gp := &getProv{
		ctx:  ctx,
		key:  k,
		resp: make(chan []peer.ID, 1), // buffered to prevent sender from blocking
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case pm.getprovs <- gp:
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case peers := <-gp.resp:
		return peers, nil
	}
}

func (pm *ProviderManager) getProvidersForKey(ctx context.Context, k []byte) ([]peer.ID, error) {
	pset, err := pm.getProviderSetForKey(ctx, k)
	if err != nil {
		return nil, err
	}
	return pset.providers, nil
}

// returns the ProviderSet if it already exists on cache, otherwise loads it from datasatore
func (pm *ProviderManager) getProviderSetForKey(ctx context.Context, k []byte) (*providerSet, error) {
	cached, ok := pm.cache.Get(string(k))
	if ok {
		return cached.(*providerSet), nil
	}

	pset, err := loadProviderSet(ctx, pm.dstore, k)
	if err != nil {
		return nil, err
	}

	if len(pset.providers) > 0 {
		pm.cache.Add(string(k), pset)
	}

	return pset, nil
}

// loads the ProviderSet out of the datastore
func loadProviderSet(ctx context.Context, dstore ds.Datastore, k []byte) (*providerSet, error) {
	res, err := dstore.Query(ctx, dsq.Query{Prefix: mkProvKey(k)})
	if err != nil {
		return nil, err
	}
	defer res.Close()

	now := time.Now()
	out := newProviderSet()
	for {
		e, ok := res.NextSync()
		if !ok {
			break
		}

		pid, decKey, t, err := handleQueryKey(ctx, dstore, e, now)
		if err != nil {
			log.Debugf("failed to handle query key: %s", err)
			continue
		}

		out.setVal(pid, decKey, t) 
	}
		
	return out, nil
}

// loads the ProviderSet out of the datastore
func loadProviderSetByPrefix(ctx context.Context, dstore ds.Datastore, k []byte, prefixBitLength int) (*providerSet, error) {
	// for prefix lookups, this already returns all providers with the prefix, so don't need to modify
	// note: we slice off the last byte since the prefix is by *bits*, so we need to manually xor and check
	// how many bits match in the final byte.
	prefixKey := mkProvKey(k[:len(k)-1])

	q := dsq.Query{
		Filters: []dsq.Filter{
			dsq.FilterKeyPrefix{Prefix: prefixKey},
		},
	}

	res, err := dstore.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	now := time.Now()
	out := newProviderSet()
	i := 0
	for {
		e, ok := res.NextSync()
		if !ok {
			break
		}

		pid, decKey, t, err := handleQueryKey(ctx, dstore, e, now)
		if err != nil {
			log.Debugf("failed to handle query key: %s", err)
			continue
		}

		// check that the last byte of the lookup key and the corresponding byte in the db key
		// match; ie that the bits set in the prefixed key match that of the db key
		// if they don't, then ignore this record
		if !prefixesMatch(k, decKey, prefixBitLength) {
			continue
		}

		out.setVal(pid, decKey, t)
		i++
	}

	return out, nil
}

func readTimeValue(data []byte) (time.Time, error) {
	nsec, n := binary.Varint(data)
	if n <= 0 {
		return time.Time{}, fmt.Errorf("failed to parse time")
	}

	return time.Unix(0, nsec), nil
}



// OTHER FUNCTIONS

// GetProvidersForPrefix returns the set of providers with the given prefix, as well
// as the full key they provide.
func (pm *ProviderManager) GetProvidersForPrefix(ctx context.Context, k []byte, prefixBitLength int) (map[string][]peer.ID, error) {
	gp := &getProvByPrefix{
		ctx:             ctx,
		key:             k,
		prefixBitLength: prefixBitLength,
		resp:            make(chan map[string][]peer.ID, 1), // buffered to prevent sender from blocking
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case pm.getProvsByPrefix <- gp:
	}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case peers := <-gp.resp:
		return peers, nil
	}
}

// returns the ProviderSet if it already exists on cache, otherwise loads it from datasatore
func (pm *ProviderManager) getProviderSetForPrefix(ctx context.Context, k []byte, prefixBitLength int) (*providerSet, error) {
	cached, ok := pm.cache.Get(string(k))
	if ok {
		return cached.(*providerSet), nil
	}

	pset, err := loadProviderSetByPrefix(ctx, pm.dstore, k, prefixBitLength)
	if err != nil {
		return nil, err
	}

	// note: prefixes are not cached, unlike full key lookups.
	// is this okay?
	return pset, nil
}


func handleQueryKey(ctx context.Context, dstore ds.Datastore, e dsq.Result, now time.Time) (peer.ID, []byte, time.Time, error) {
	if e.Error != nil {
		return "", nil, time.Time{}, e.Error
	}

	// check expiration time
	t, err := readTimeValue(e.Value)
	switch {
	case err != nil:
		// couldn't parse the time
		log.Error("parsing providers record from disk: ", err)
		fallthrough
	case now.Sub(t) > ProvideValidity:
		// or just expired
		err = dstore.Delete(ctx, ds.RawKey(e.Key))
		if err != nil && err != ds.ErrNotFound {
			return "", nil, time.Time{}, fmt.Errorf("failed to remove provider record from disk: %w", err)
		}
	}

	lix := strings.LastIndex(e.Key, "/")

	// decode peer ID
	decstr, err := base32.RawStdEncoding.DecodeString(e.Key[lix+1:])
	if err != nil {
		log.Error("hex decoding error: ", err)
		err = dstore.Delete(ctx, ds.RawKey(e.Key))
		if err != nil && err != ds.ErrNotFound {
			return "", nil, time.Time{}, fmt.Errorf("failed to remove provider record from disk: %w", err)
		}
	}

	pid := peer.ID(decstr)

	// decode provided key
	decKey, err := hex.DecodeString(e.Key[len(ProvidersKeyPrefix):lix])
	if err != nil {
		log.Error("base32 decoding error: ", err)
		err = dstore.Delete(ctx, ds.RawKey(e.Key))
		if err != nil && err != ds.ErrNotFound {
			return "", nil, time.Time{}, fmt.Errorf("failed to remove provider record from disk: %w", err)
		}
	}

	return pid, decKey, t, nil
}

// prefixesMatch returns true if the bits up to `prefixBitLength` match.
func prefixesMatch(a, b []byte, prefixBitLength int) bool {
	if prefixBitLength == 0 {
		return true
	}

	if len(a) == 0 || len(b) == 0 {
		return false
	}

	byteLen := prefixBitLength / 8
	if prefixBitLength%8 == 0 && (len(a) < byteLen || len(b) < byteLen) {
		return false
	}

	if prefixBitLength%8 != 0 && (len(a) < byteLen+1 || len(b) < byteLen+1) {
		return false
	}

	if !bytes.Equal(a[:byteLen], b[:byteLen]) {
		return false
	}

	if prefixBitLength%8 == 0 {
		return true
	}

	cb := numCommonBits(a[byteLen], b[byteLen])
	return cb >= prefixBitLength%8
}

func numCommonBits(a, b byte) int {
	// xor last byte to see how many bits in common they have
	common := a ^ b

	// left shift by 1 bit each iteration
	// once common becomes 0, return 8 - number of iterations
	i := 8
	for {
		if common == 0 {
			return i
		}

		common = common << 1
		i--
	}
}