package ipcounter

import "errors"

var (
	ErrNoHexIPFound           = errors.New("no hex IP found in Aerospike key")
	ErrDecodingHexFailed      = errors.New("hex decode failed")
	ErrInvalidIPV4Length      = errors.New("invalid byte length for IPv4 conversion")
	ErrCannotConvertHexToIP   = errors.New("failed to convert hex to IP")
	ErrFailedToConnect        = errors.New("failed to connect to Aerospike")
	ErrInvalidKey             = errors.New("invalid key")
	ErrFailedToIncrementCount = errors.New("failed to increment count")
	ErrFailedScanAllRecords   = errors.New("failed to scan all records")
	ErrScanRecord             = errors.New("error scanning record")
	ErrSetKey                 = errors.New("failed to set key")
	ErrReadInputFile          = errors.New("failed to read input file")
	ErrProcessInputFile       = errors.New("failed to process input file")
	ErrProcessRecords         = errors.New("failed to process records")
	ErrListRecords            = errors.New("failed to list records")
	ErrWriteToStringBuilder   = errors.New("failed to write to string builder")
	ErrCreateFile             = errors.New("failed to create output file")
	ErrWriteFile              = errors.New("failed to write to output file")
)
