// Copyright 2022 Namespace Labs Inc; All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package kubenaming

import (
	"crypto/sha256"
	"fmt"

	"namespacelabs.dev/go-ids"
)

func StableID(str string) string {
	h := sha256.New()
	fmt.Fprint(h, str)
	return ids.EncodeToBase32String(h.Sum(nil))
}

func StableIDN(str string, n int) string {
	if n < 0 || n > 50 {
		panic("invalid N")
	}
	return StableID(str)[:n]
}
