// Package metis provides Go bindings for the METIS graph partitioning library.
package metis

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lmetis -lm
#cgo darwin CFLAGS: -I/opt/homebrew/include -I/usr/local/include
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -lmetis

#include <metis.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// Error codes from METIS
const (
	OK          = C.METIS_OK
	ErrorInput  = C.METIS_ERROR_INPUT
	ErrorMemory = C.METIS_ERROR_MEMORY
	Error       = C.METIS_ERROR
)

// Options indices
const (
	OptionPType     = C.METIS_OPTION_PTYPE
	OptionObjType   = C.METIS_OPTION_OBJTYPE
	OptionCType     = C.METIS_OPTION_CTYPE
	OptionIPType    = C.METIS_OPTION_IPTYPE
	OptionRType     = C.METIS_OPTION_RTYPE
	OptionDBGLvl    = C.METIS_OPTION_DBGLVL
	OptionNIter     = C.METIS_OPTION_NITER
	OptionNCuts     = C.METIS_OPTION_NCUTS
	OptionSeed      = C.METIS_OPTION_SEED
	OptionNo2Hop    = C.METIS_OPTION_NO2HOP
	OptionMinConn   = C.METIS_OPTION_MINCONN
	OptionContig    = C.METIS_OPTION_CONTIG
	OptionCompress  = C.METIS_OPTION_COMPRESS
	OptionCCOrder   = C.METIS_OPTION_CCORDER
	OptionPFactor   = C.METIS_OPTION_PFACTOR
	OptionNSeps     = C.METIS_OPTION_NSEPS
	OptionUFactor   = C.METIS_OPTION_UFACTOR
	OptionNumbering = C.METIS_OPTION_NUMBERING
	OptionHelp      = C.METIS_OPTION_HELP
	OptionTPWGTS    = C.METIS_OPTION_TPWGTS
	OptionNCommon   = C.METIS_OPTION_NCOMMON
	OptionNoOutput  = C.METIS_OPTION_NOOUTPUT
	OptionBalance   = C.METIS_OPTION_BALANCE
	OptionGType     = C.METIS_OPTION_GTYPE
	OptionUBVec     = C.METIS_OPTION_UBVEC
)

// Constants
const (
	NoOptions = C.METIS_NOPTIONS
)

// Partitioning types
const (
	PTypeRB   = C.METIS_PTYPE_RB
	PTypeKway = C.METIS_PTYPE_KWAY
)

// Graph types
const (
	GTypeDual  = C.METIS_GTYPE_DUAL
	GTypeNodal = C.METIS_GTYPE_NODAL
)

// Coarsening types
const (
	CTypeRM   = C.METIS_CTYPE_RM
	CTypeSHEM = C.METIS_CTYPE_SHEM
)

// Initial partitioning types
const (
	IPTypeGrow    = C.METIS_IPTYPE_GROW
	IPTypeRandom  = C.METIS_IPTYPE_RANDOM
	IPTypeEdge    = C.METIS_IPTYPE_EDGE
	IPTypeNode    = C.METIS_IPTYPE_NODE
	IPTypeMetisRB = C.METIS_IPTYPE_METISRB
)

// Refinement types
const (
	RTypeFM        = C.METIS_RTYPE_FM
	RTypeGreedy    = C.METIS_RTYPE_GREEDY
	RTypeSep2Sided = C.METIS_RTYPE_SEP2SIDED
	RTypeSep1Sided = C.METIS_RTYPE_SEP1SIDED
)

// Objective types
const (
	ObjTypeCut  = C.METIS_OBJTYPE_CUT
	ObjTypeVol  = C.METIS_OBJTYPE_VOL
	ObjTypeNode = C.METIS_OBJTYPE_NODE
)

// Debug levels
const (
	DBGInfo       = C.METIS_DBG_INFO
	DBGTime       = C.METIS_DBG_TIME
	DBGCoarsen    = C.METIS_DBG_COARSEN
	DBGRefine     = C.METIS_DBG_REFINE
	DBGIPart      = C.METIS_DBG_IPART
	DBGMoveInfo   = C.METIS_DBG_MOVEINFO
	DBGSepInfo    = C.METIS_DBG_SEPINFO
	DBGConnInfo   = C.METIS_DBG_CONNINFO
	DBGContigInfo = C.METIS_DBG_CONTIGINFO
	DBGMemory     = C.METIS_DBG_MEMORY
)

// SetDefaultOptions initializes the options array with default values
func SetDefaultOptions(opts []int32) error {
	if len(opts) != NoOptions {
		return fmt.Errorf("options array must have %d elements", NoOptions)
	}

	ret := C.METIS_SetDefaultOptions((*C.idx_t)(unsafe.Pointer(&opts[0])))
	if ret != OK {
		return getError(ret)
	}
	return nil
}

// PartGraphRecursive partitions a graph using multilevel recursive bisection
func PartGraphRecursive(xadj, adjncy []int32, nparts int32, options []int32) ([]int32, int32, error) {
	nvtxs := int32(len(xadj) - 1)
	ncon := int32(1)
	part := make([]int32, nvtxs)
	var objval C.idx_t

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_PartGraphRecursive(
		(*C.idx_t)(unsafe.Pointer(&nvtxs)),
		(*C.idx_t)(unsafe.Pointer(&ncon)),
		(*C.idx_t)(unsafe.Pointer(&xadj[0])),
		(*C.idx_t)(unsafe.Pointer(&adjncy[0])),
		nil, nil, nil,
		(*C.idx_t)(unsafe.Pointer(&nparts)),
		nil, nil,
		opts,
		&objval,
		(*C.idx_t)(unsafe.Pointer(&part[0])),
	)

	if ret != OK {
		return nil, 0, getError(ret)
	}

	return part, int32(objval), nil
}

// PartGraphKway partitions a graph using multilevel k-way partitioning
func PartGraphKway(xadj, adjncy []int32, nparts int32, options []int32) ([]int32, int32, error) {
	nvtxs := int32(len(xadj) - 1)
	ncon := int32(1)
	part := make([]int32, nvtxs)
	var objval C.idx_t

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_PartGraphKway(
		(*C.idx_t)(unsafe.Pointer(&nvtxs)),
		(*C.idx_t)(unsafe.Pointer(&ncon)),
		(*C.idx_t)(unsafe.Pointer(&xadj[0])),
		(*C.idx_t)(unsafe.Pointer(&adjncy[0])),
		nil, nil, nil,
		(*C.idx_t)(unsafe.Pointer(&nparts)),
		nil, nil,
		opts,
		&objval,
		(*C.idx_t)(unsafe.Pointer(&part[0])),
	)

	if ret != OK {
		return nil, 0, getError(ret)
	}

	return part, int32(objval), nil
}

// PartGraphRecursiveWeighted partitions a graph with vertex and edge weights using recursive bisection
func PartGraphRecursiveWeighted(xadj, adjncy, vwgt, adjwgt []int32, nparts int32, tpwgts, ubvec []float32, options []int32) ([]int32, int32, error) {
	nvtxs := int32(len(xadj) - 1)
	ncon := int32(1)
	if vwgt != nil && len(vwgt) != int(nvtxs) {
		return nil, 0, errors.New("vwgt length must equal number of vertices")
	}
	if adjwgt != nil && len(adjwgt) != len(adjncy) {
		return nil, 0, errors.New("adjwgt length must equal adjncy length")
	}

	part := make([]int32, nvtxs)
	var objval C.idx_t

	var vwgtPtr, adjwgtPtr *C.idx_t
	if vwgt != nil {
		vwgtPtr = (*C.idx_t)(unsafe.Pointer(&vwgt[0]))
	}
	if adjwgt != nil {
		adjwgtPtr = (*C.idx_t)(unsafe.Pointer(&adjwgt[0]))
	}

	var tpwgtsPtr, ubvecPtr *C.real_t
	if tpwgts != nil {
		tpwgtsPtr = (*C.real_t)(unsafe.Pointer(&tpwgts[0]))
	}
	if ubvec != nil {
		ubvecPtr = (*C.real_t)(unsafe.Pointer(&ubvec[0]))
	}

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_PartGraphRecursive(
		(*C.idx_t)(unsafe.Pointer(&nvtxs)),
		(*C.idx_t)(unsafe.Pointer(&ncon)),
		(*C.idx_t)(unsafe.Pointer(&xadj[0])),
		(*C.idx_t)(unsafe.Pointer(&adjncy[0])),
		vwgtPtr, nil, adjwgtPtr,
		(*C.idx_t)(unsafe.Pointer(&nparts)),
		tpwgtsPtr, ubvecPtr,
		opts,
		&objval,
		(*C.idx_t)(unsafe.Pointer(&part[0])),
	)

	if ret != OK {
		return nil, 0, getError(ret)
	}

	return part, int32(objval), nil
}

// PartGraphKwayWeighted partitions a graph with vertex and edge weights using k-way partitioning
func PartGraphKwayWeighted(xadj, adjncy, vwgt, adjwgt []int32, nparts int32, tpwgts, ubvec []float32, options []int32) ([]int32, int32, error) {
	nvtxs := int32(len(xadj) - 1)
	ncon := int32(1)
	if vwgt != nil && len(vwgt) != int(nvtxs) {
		return nil, 0, errors.New("vwgt length must equal number of vertices")
	}
	if adjwgt != nil && len(adjwgt) != len(adjncy) {
		return nil, 0, errors.New("adjwgt length must equal adjncy length")
	}

	part := make([]int32, nvtxs)
	var objval C.idx_t

	var vwgtPtr, adjwgtPtr *C.idx_t
	if vwgt != nil {
		vwgtPtr = (*C.idx_t)(unsafe.Pointer(&vwgt[0]))
	}
	if adjwgt != nil {
		adjwgtPtr = (*C.idx_t)(unsafe.Pointer(&adjwgt[0]))
	}

	var tpwgtsPtr, ubvecPtr *C.real_t
	if tpwgts != nil {
		tpwgtsPtr = (*C.real_t)(unsafe.Pointer(&tpwgts[0]))
	}
	if ubvec != nil {
		ubvecPtr = (*C.real_t)(unsafe.Pointer(&ubvec[0]))
	}

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_PartGraphKway(
		(*C.idx_t)(unsafe.Pointer(&nvtxs)),
		(*C.idx_t)(unsafe.Pointer(&ncon)),
		(*C.idx_t)(unsafe.Pointer(&xadj[0])),
		(*C.idx_t)(unsafe.Pointer(&adjncy[0])),
		vwgtPtr, nil, adjwgtPtr,
		(*C.idx_t)(unsafe.Pointer(&nparts)),
		tpwgtsPtr, ubvecPtr,
		opts,
		&objval,
		(*C.idx_t)(unsafe.Pointer(&part[0])),
	)

	if ret != OK {
		return nil, 0, getError(ret)
	}

	return part, int32(objval), nil
}

// MeshToDual converts a mesh to its dual graph
func MeshToDual(ne, nn int32, eptr, eind []int32, ncommon int32) ([]int32, []int32, error) {
	var xadj, adjncy *C.idx_t
	var numflag C.idx_t = 0 // C-style numbering

	ret := C.METIS_MeshToDual(
		(*C.idx_t)(unsafe.Pointer(&ne)),
		(*C.idx_t)(unsafe.Pointer(&nn)),
		(*C.idx_t)(unsafe.Pointer(&eptr[0])),
		(*C.idx_t)(unsafe.Pointer(&eind[0])),
		(*C.idx_t)(unsafe.Pointer(&ncommon)),
		&numflag,
		&xadj,
		&adjncy,
	)

	if ret != OK {
		return nil, nil, getError(ret)
	}

	// Convert C arrays to Go slices
	xadjSlice := make([]int32, ne+1)
	for i := 0; i < int(ne+1); i++ {
		xadjSlice[i] = int32(*(*C.idx_t)(unsafe.Pointer(uintptr(unsafe.Pointer(xadj)) + uintptr(i)*unsafe.Sizeof(C.idx_t(0)))))
	}

	// Get size of adjncy array from xadj[ne]
	adjSize := xadjSlice[ne]
	adjncySlice := make([]int32, adjSize)
	for i := 0; i < int(adjSize); i++ {
		adjncySlice[i] = int32(*(*C.idx_t)(unsafe.Pointer(uintptr(unsafe.Pointer(adjncy)) + uintptr(i)*unsafe.Sizeof(C.idx_t(0)))))
	}

	// Free the memory allocated by METIS
	C.METIS_Free(unsafe.Pointer(xadj))
	C.METIS_Free(unsafe.Pointer(adjncy))

	return xadjSlice, adjncySlice, nil
}

// MeshToNodal converts a mesh to its nodal graph
func MeshToNodal(ne, nn int32, eptr, eind []int32) ([]int32, []int32, error) {
	var xadj, adjncy *C.idx_t
	var numflag C.idx_t = 0 // C-style numbering

	ret := C.METIS_MeshToNodal(
		(*C.idx_t)(unsafe.Pointer(&ne)),
		(*C.idx_t)(unsafe.Pointer(&nn)),
		(*C.idx_t)(unsafe.Pointer(&eptr[0])),
		(*C.idx_t)(unsafe.Pointer(&eind[0])),
		&numflag,
		&xadj,
		&adjncy,
	)

	if ret != OK {
		return nil, nil, getError(ret)
	}

	// Convert C arrays to Go slices
	xadjSlice := make([]int32, nn+1)
	for i := 0; i < int(nn+1); i++ {
		xadjSlice[i] = int32(*(*C.idx_t)(unsafe.Pointer(uintptr(unsafe.Pointer(xadj)) + uintptr(i)*unsafe.Sizeof(C.idx_t(0)))))
	}

	// Get size of adjncy array from xadj[nn]
	adjSize := xadjSlice[nn]
	adjncySlice := make([]int32, adjSize)
	for i := 0; i < int(adjSize); i++ {
		adjncySlice[i] = int32(*(*C.idx_t)(unsafe.Pointer(uintptr(unsafe.Pointer(adjncy)) + uintptr(i)*unsafe.Sizeof(C.idx_t(0)))))
	}

	// Free the memory allocated by METIS
	C.METIS_Free(unsafe.Pointer(xadj))
	C.METIS_Free(unsafe.Pointer(adjncy))

	return xadjSlice, adjncySlice, nil
}

// PartMeshNodal partitions a mesh using its nodal graph
func PartMeshNodal(ne, nn int32, eptr, eind []int32, vwgt, vsize []int32, nparts int32, tpwgts []float32, options []int32) (int32, []int32, []int32, error) {
	var objval C.idx_t
	epart := make([]int32, ne)
	npart := make([]int32, nn)

	var vwgtPtr, vsizePtr *C.idx_t
	if vwgt != nil {
		vwgtPtr = (*C.idx_t)(unsafe.Pointer(&vwgt[0]))
	}
	if vsize != nil {
		vsizePtr = (*C.idx_t)(unsafe.Pointer(&vsize[0]))
	}

	var tpwgtsPtr *C.real_t
	if tpwgts != nil {
		tpwgtsPtr = (*C.real_t)(unsafe.Pointer(&tpwgts[0]))
	}

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_PartMeshNodal(
		(*C.idx_t)(unsafe.Pointer(&ne)),
		(*C.idx_t)(unsafe.Pointer(&nn)),
		(*C.idx_t)(unsafe.Pointer(&eptr[0])),
		(*C.idx_t)(unsafe.Pointer(&eind[0])),
		vwgtPtr, vsizePtr,
		(*C.idx_t)(unsafe.Pointer(&nparts)),
		tpwgtsPtr,
		opts,
		&objval,
		(*C.idx_t)(unsafe.Pointer(&epart[0])),
		(*C.idx_t)(unsafe.Pointer(&npart[0])),
	)

	if ret != OK {
		return 0, nil, nil, getError(ret)
	}

	return int32(objval), epart, npart, nil
}

// PartMeshDual partitions a mesh using its dual graph
func PartMeshDual(ne, nn int32, eptr, eind []int32, vwgt, vsize []int32, ncommon, nparts int32, tpwgts []float32, options []int32) (int32, []int32, []int32, error) {
	var objval C.idx_t
	epart := make([]int32, ne)
	npart := make([]int32, nn)

	var vwgtPtr, vsizePtr *C.idx_t
	if vwgt != nil {
		vwgtPtr = (*C.idx_t)(unsafe.Pointer(&vwgt[0]))
	}
	if vsize != nil {
		vsizePtr = (*C.idx_t)(unsafe.Pointer(&vsize[0]))
	}

	var tpwgtsPtr *C.real_t
	if tpwgts != nil {
		tpwgtsPtr = (*C.real_t)(unsafe.Pointer(&tpwgts[0]))
	}

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_PartMeshDual(
		(*C.idx_t)(unsafe.Pointer(&ne)),
		(*C.idx_t)(unsafe.Pointer(&nn)),
		(*C.idx_t)(unsafe.Pointer(&eptr[0])),
		(*C.idx_t)(unsafe.Pointer(&eind[0])),
		vwgtPtr, vsizePtr,
		(*C.idx_t)(unsafe.Pointer(&ncommon)),
		(*C.idx_t)(unsafe.Pointer(&nparts)),
		tpwgtsPtr,
		opts,
		&objval,
		(*C.idx_t)(unsafe.Pointer(&epart[0])),
		(*C.idx_t)(unsafe.Pointer(&npart[0])),
	)

	if ret != OK {
		return 0, nil, nil, getError(ret)
	}

	return int32(objval), epart, npart, nil
}

// NodeND computes fill reducing ordering using nested dissection
func NodeND(xadj, adjncy, vwgt []int32, options []int32) ([]int32, []int32, error) {
	nvtxs := int32(len(xadj) - 1)
	perm := make([]int32, nvtxs)
	iperm := make([]int32, nvtxs)

	var vwgtPtr *C.idx_t
	if vwgt != nil && len(vwgt) == int(nvtxs) {
		vwgtPtr = (*C.idx_t)(unsafe.Pointer(&vwgt[0]))
	}

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_NodeND(
		(*C.idx_t)(unsafe.Pointer(&nvtxs)),
		(*C.idx_t)(unsafe.Pointer(&xadj[0])),
		(*C.idx_t)(unsafe.Pointer(&adjncy[0])),
		vwgtPtr,
		opts,
		(*C.idx_t)(unsafe.Pointer(&perm[0])),
		(*C.idx_t)(unsafe.Pointer(&iperm[0])),
	)

	if ret != OK {
		return nil, nil, getError(ret)
	}

	return perm, iperm, nil
}

// ComputeVertexSeparator computes a vertex separator from an edge separator
func ComputeVertexSeparator(xadj, adjncy, vwgt []int32, options []int32) (int32, []int32, error) {
	nvtxs := int32(len(xadj) - 1)
	part := make([]int32, nvtxs)
	var sepsize C.idx_t

	var vwgtPtr *C.idx_t
	if vwgt != nil && len(vwgt) == int(nvtxs) {
		vwgtPtr = (*C.idx_t)(unsafe.Pointer(&vwgt[0]))
	}

	var opts *C.idx_t
	if options != nil && len(options) == NoOptions {
		opts = (*C.idx_t)(unsafe.Pointer(&options[0]))
	}

	ret := C.METIS_ComputeVertexSeparator(
		(*C.idx_t)(unsafe.Pointer(&nvtxs)),
		(*C.idx_t)(unsafe.Pointer(&xadj[0])),
		(*C.idx_t)(unsafe.Pointer(&adjncy[0])),
		vwgtPtr,
		opts,
		&sepsize,
		(*C.idx_t)(unsafe.Pointer(&part[0])),
	)

	if ret != OK {
		return 0, nil, getError(ret)
	}

	return int32(sepsize), part, nil
}

// Version returns the METIS version
func Version() string {
	return fmt.Sprintf("%d.%d.%d", C.METIS_VER_MAJOR, C.METIS_VER_MINOR, C.METIS_VER_SUBMINOR)
}

// getError converts METIS error codes to Go errors
func getError(status C.int) error {
	switch status {
	case ErrorInput:
		return errors.New("METIS error: erroneous inputs and/or options")
	case ErrorMemory:
		return errors.New("METIS error: insufficient memory")
	case Error:
		return errors.New("METIS error: general error")
	default:
		return fmt.Errorf("METIS error: unknown error code %d", status)
	}
}
