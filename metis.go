// Package metis provides Go bindings for the METIS graph partitioning library.
package metis

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lmetis -lm
#cgo darwin CFLAGS: -I/opt/homebrew/include -I/usr/local/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lmetis

#include <metis.h>
*/
import "C"

// Version returns the METIS version
func Version() string {
    // TODO: Implement
    return "5.1.0"
}
