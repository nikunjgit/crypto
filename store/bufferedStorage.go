package store

import (
	"time"
	"fmt"
	"github.com/nikunjgit/crypto/event"
)

type DataStore interface {
	Get(keys []string) (event.Messages, error)
	Set(key string, messages event.Messages) error
}

type BufferedStore struct {
	DataStore
	flushInterval time.Duration
	memory        event.Messages
	flushTime     time.Time
	Prefix        string
	batchsize	  int
}

func NewBufferedStorage(store DataStore, flushInterval time.Duration, prefix string) *BufferedStore {
	return &BufferedStore{store, flushInterval, make(event.Messages, 0, 10), time.Now(), prefix, 100}

}

func (b *BufferedStore) Get(start time.Time, end time.Time) (event.Messages, error) {
	startBucket := b.nearestBucket(start)
	endBucket := b.nearestBucket(end)
	messageArr := make(event.Messages, 0, 10)
	for bucket := startBucket; !endBucket.Before(bucket); bucket = bucket.Add(b.flushInterval * time.Duration(b.batchsize)) {
		buckets := GenerateBuckets(bucket, b.batchsize, b.flushInterval)
		message, err := b.GetInstant(buckets)
		if err != nil {
			return nil, err
		}
		messageArr = append(messageArr, message...)
	}

	return messageArr, nil
}
func GenerateBuckets(current time.Time, batchSize int, flushInterval time.Duration) []time.Time {
	batches := make([]time.Time, 0, batchSize)
	for i := 0 ; i < batchSize ; i++ {
		batches = append(batches, current.Add(time.Duration(batchSize) * flushInterval))
	}
	return batches
}

func (b *BufferedStore) GetInstant(times []time.Time) (event.Messages, error) {
	keys  := make([]string, len(times))
	for i:=0; i < len(times); i++ {
		keys[i] = b.Key(times[i])
	}

	messages, err := b.DataStore.Get(keys)
	if err != nil {
		return nil, err
	}

	return messages, nil
}


func (b *BufferedStore) Set(message *event.Message) error {
	b.memory = append(b.memory, message)
	if time.Since(b.flushTime) < b.flushInterval {
		return nil
	}

	current := time.Now()
	truncated := b.nearestBucket(current)
	if truncated.Before(b.flushTime) {
		return nil
	}

	err := b.DataStore.Set(b.Key(truncated), b.memory)
	if err != nil {
		return err
	}
	b.flushTime = current

	return nil
}

func (b *BufferedStore) Key(current time.Time) string {
	return fmt.Sprintf("%s:%d", b.Prefix, current.Unix())
}

func (b *BufferedStore) nearestBucket(current time.Time) time.Time {
	truncated := current.Truncate(b.flushInterval)
	return truncated
}
