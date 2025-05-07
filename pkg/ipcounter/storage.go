package ipcounter

import (
	"context"
	"encoding/hex"
	"net"
	"regexp"
	"strings"

	as "github.com/aerospike/aerospike-client-go/v6"
)

// AeroSpike defines datastore using AeroSpike client.
type AeroSpike struct {
	client *as.Client
}

// NewAeroSpike creates a new datastore creating new connection with aerospike.
func NewAeroSpike(host string, port int) (*AeroSpike, error) {
	c, err := as.NewClient(host, port)
	if err != nil {
		return nil, ErrFailedToConnect
	}

	return &AeroSpike{client: c}, nil
}

// Increment increments the count of the given key in the specified namespace and set.
func (a *AeroSpike) Increment(_ context.Context, namespace, set, key string, count int) error {
	id, err := as.NewKey(namespace, set, key)
	if err != nil {
		return ErrInvalidKey
	}

	pol := as.NewWritePolicy(0, 0)
	bin := as.NewBin("count", count)
	if _, err = a.client.Operate(pol, id, as.AddOp(bin)); err != nil {
		return ErrFailedToIncrementCount
	}

	return nil
}

// List retrieves all records from the specified namespace and set.
func (a *AeroSpike) List(_ context.Context, namespace, set string) ([]Record, error) {
	policy := as.NewScanPolicy()
	store, err := a.client.ScanAll(policy, namespace, set)
	if err != nil {
		return nil, ErrFailedScanAllRecords
	}
	defer store.Close()

	//nolint:prealloc
	var records []Record
	for item := range store.Results() {
		if item.Err != nil {
			return nil, ErrScanRecord
		}

		input := item.Record.Key.String()
		key, err := extractAeroKey(input)
		if err != nil {
			return []Record{}, err
		}

		ip, err := convertHexToIPV4String(key)
		if err != nil {
			return nil, err
		}

		cnt, ok := item.Record.Bins["count"].(int)
		if !ok {
			continue
		}

		records = append(records, Record{Key: ip, Count: cnt})
	}

	return records, nil
}

// Close closes the connection to the AeroSpike client.
func (a *AeroSpike) Close() {
	a.client.Close()
}

// ExtractAeroKey extracts the hex key from the Aerospike key string.
func extractAeroKey(input string) (string, error) {
	const length = 2
	re := regexp.MustCompile(`::((?:[0-9a-f]{2}\s?)+)`)
	matches := re.FindStringSubmatch(input)
	if len(matches) < length {
		return "", ErrNoHexIPFound
	}

	return matches[1], nil
}

// ConvertHexToIPV4String converts a hex string to an IPv4 address string.
func convertHexToIPV4String(input string) (string, error) {
	const length = 16
	s := strings.ReplaceAll(input, " ", "")
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return "", ErrDecodingHexFailed
	}

	if len(bytes) < length {
		return "", ErrInvalidIPV4Length
	}

	ip := net.IP(bytes[12:16]).String()
	if ip == "" {
		return "", ErrCannotConvertHexToIP
	}

	return ip, nil
}
