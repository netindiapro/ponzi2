package iex

import (
	"context"
	"encoding/gob"
	"os"
	"path/filepath"
	"time"

	"github.com/btmura/ponzi2/internal/errors"
)

// GetQuotes gets quotes for stock symbols.
func (c *CacheClient) GetQuotes(ctx context.Context, req *GetQuotesRequest) ([]*Quote, error) {
	cacheClientVar.Add("get-quotes-requests", 1)

	symbol2Quote := map[string]*Quote{}
	var missingSymbols []string
	for _, sym := range req.Symbols {
		k := quoteCacheKey{Symbol: sym}
		v := c.quoteCache.get(k)
		if v != nil {
			symbol2Quote[sym] = v.Quote.DeepCopy()
		} else {
			missingSymbols = append(missingSymbols, sym)
		}
	}

	r := &GetQuotesRequest{Symbols: missingSymbols}
	missingQuotes, err := c.client.GetQuotes(ctx, r)
	if err != nil {
		return nil, err
	}

	for _, q := range missingQuotes {
		k := quoteCacheKey{Symbol: q.Symbol}
		v := &quoteCacheValue{
			Quote:          q.DeepCopy(),
			LastUpdateTime: now(),
		}
		symbol2Quote[q.Symbol] = q.DeepCopy()
		if err := c.quoteCache.put(k, v); err != nil {
			return nil, err
		}
	}

	if err := saveQuoteCache(c.quoteCache); err != nil {
		return nil, err
	}

	var quotes []*Quote
	for _, sym := range req.Symbols {
		q := symbol2Quote[sym]
		if q == nil {
			q = &Quote{Symbol: sym}
		}
		quotes = append(quotes, q)
	}
	return quotes, nil
}

// quoteCache caches data from the quote endpoint.
// Fields are exported for gob encoding and decoding.
type quoteCache struct {
	Data map[quoteCacheKey]*quoteCacheValue
}

type quoteCacheKey struct {
	Symbol string
}

type quoteCacheValue struct {
	Quote          *Quote
	LastUpdateTime time.Time
}

func (q *quoteCacheValue) deepCopy() *quoteCacheValue {
	copy := *q
	copy.Quote = copy.Quote.DeepCopy()
	return &copy
}

func (q *quoteCache) get(key quoteCacheKey) *quoteCacheValue {
	cacheClientVar.Add("quote-cache-gets", 1)

	v := q.Data[key]
	if v != nil {
		cacheClientVar.Add("quote-cache-hits", 1)
		return v.deepCopy()
	}
	cacheClientVar.Add("quote-cache-misses", 1)
	return nil
}

func (q *quoteCache) put(key quoteCacheKey, val *quoteCacheValue) error {
	cacheClientVar.Add("quote-cache-puts", 1)

	if !validSymbolRegexp.MatchString(key.Symbol) {
		return errors.Errorf("bad symbol: got %s, want: %v", key.Symbol, validSymbolRegexp)
	}

	if q.Data == nil {
		q.Data = map[quoteCacheKey]*quoteCacheValue{}
	}
	q.Data[key] = val.deepCopy()
	q.Data[key].LastUpdateTime = now()

	return nil
}

func loadQuoteCache() (*quoteCache, error) {
	t := now()
	defer func() {
		cacheClientVar.Set("quote-cache-load-time", time.Since(t))
	}()

	path, err := quoteCachePath()
	if err != nil {
		return nil, err
	}

	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return &quoteCache{}, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	q := &quoteCache{}
	dec := gob.NewDecoder(file)
	if err := dec.Decode(q); err != nil {
		return nil, err
	}
	return q, nil
}

func saveQuoteCache(q *quoteCache) error {
	t := now()
	defer func() {
		cacheClientVar.Set("quote-cache-save-time", time.Since(t))
	}()

	path, err := quoteCachePath()
	if err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := gob.NewEncoder(file)
	if err := enc.Encode(q); err != nil {
		return err
	}
	return nil
}

func quoteCachePath() (string, error) {
	dir, err := userCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "iex-quote-cache.gob"), nil
}