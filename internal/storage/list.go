package storage

import (
	"context"
	"strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func (b *GCSBackend) ModuleVersions(path string, ctx context.Context) ([]string, error) {
	var versions []string
	query := &storage.Query{
		Prefix: path,
	}
	it := b.Client.Bucket(b.Bucket).Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return versions, err
		}
		if attrs.ContentType != "application/zip" {
			continue
		}
		version := strings.Split(attrs.Name, "/")[3]
		if version == "" {
			continue
		}
		vs := version
		versions = append(versions, vs)
	}
	return versions, nil
}
