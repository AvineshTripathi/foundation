// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the EARLY ACCESS SOFTWARE LICENSE AGREEMENT
// available at http://github.com/namespacelabs/foundation

package tarfs

import (
	"compress/gzip"
	"context"
	"io"
	"io/fs"

	"namespacelabs.dev/foundation/internal/artifacts"
	"namespacelabs.dev/foundation/workspace/compute"
	"namespacelabs.dev/foundation/workspace/tasks"
)

func TarGunzip(contents compute.Computable[compute.ByteStream]) compute.Computable[fs.FS] {
	return compute.Map(tasks.Action("tar.fs"),
		compute.Inputs().Computable("contents", contents),
		compute.Output{},
		func(ctx context.Context, r compute.Resolved) (fs.FS, error) {
			blob := compute.GetDepValue(r, contents, "contents")
			return FS{
				TarStream: func() (io.ReadCloser, error) {
					r, err := blob.Reader()
					if err != nil {
						return nil, err
					}

					pr := artifacts.NewProgressReader(r, blob.ContentLength())
					tasks.Attachments(ctx).SetProgress(pr)

					return gzip.NewReader(pr)
				},
			}, nil
		})
}
