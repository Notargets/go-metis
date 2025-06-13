package metis

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lmetis -lm -lGKlib
#cgo darwin CFLAGS: -I/opt/homebrew/include -I/usr/local/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lmetis -lGKlib

#include <metis.h>
*/
import "C"
import "fmt"

// Version set by GitHub tag replacement
// GitHub replaces $Format:%(describe:tags=true)$ with the actual tag
var goMetisVersion = "$Format:%(describe:tags=true)$"

// Version returns the version of the installed METIS library
func Version() string {
	return fmt.Sprintf("%d.%d.%d", C.METIS_VER_MAJOR, C.METIS_VER_MINOR, C.METIS_VER_SUBMINOR)
}

// GoMetisVersion returns the version of go-metis from git tags
func GoMetisVersion() string {
	// If the version string contains "$Format", it means we're in development
	if goMetisVersion[0] == '$' {
		return "dev"
	}
	return goMetisVersion
}
