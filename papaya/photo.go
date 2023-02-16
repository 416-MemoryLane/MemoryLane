package papaya

import "fmt"

// Will change significantly based on implementation of the rest of the service
type Photo string

var ErrPhotoExists = fmt.Errorf("photo exists")
