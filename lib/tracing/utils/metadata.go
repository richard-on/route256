package utils

import (
	"strings"

	"google.golang.org/grpc/metadata"
)

// MetadataReaderWriter satisfies both the opentracing.TextMapReader and
// opentracing.TextMapWriter interfaces.
type MetadataReaderWriter struct {
	metadata.MD
}

func (w MetadataReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	w.MD[key] = append(w.MD[key], val)
}

func (w MetadataReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vals := range w.MD {
		for _, v := range vals {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}

	return nil
}
