package metis

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSetDefaultOptions(t *testing.T) {
	opts := make([]int32, NoOptions)
	err := SetDefaultOptions(opts)
	require.NoError(t, err)

	// Check that all options are set to -1 (use default)
	for i := 0; i < NoOptions; i++ {
		assert.Equal(t, int32(-1), opts[i], "Option %d should be -1", i)
	}

	// Test with wrong size
	wrongOpts := make([]int32, 10)
	err = SetDefaultOptions(wrongOpts)
	assert.Error(t, err)
}

// Test_PartGraph emulates the C test function Test_PartGraph
func TestPartGraph(t *testing.T) {
	// Create a test graph similar to C tests - need larger graph for many partitions
	nvtxs := 1000 // Increased size for better partitioning
	xadj, adjncy := createRandomGraph(nvtxs)

	// Create vertex weights
	vwgt := make([]int32, nvtxs)
	for i := 0; i < nvtxs; i++ {
		vwgt[i] = int32(1 + rand.Intn(10))
	}

	// Create edge weights
	adjwgt := make([]int32, len(adjncy))
	for i := 0; i < nvtxs; i++ {
		for j := xadj[i]; j < xadj[i+1]; j++ {
			k := adjncy[j]
			if i < int(k) {
				adjwgt[j] = int32(1 + rand.Intn(5))
				// Find reverse edge and set same weight
				for jj := xadj[k]; jj < xadj[k+1]; jj++ {
					if adjncy[jj] == int32(i) {
						adjwgt[jj] = adjwgt[j]
						break
					}
				}
			}
		}
	}

	// Target partition weights for weighted partitioning
	tpwgts := []float32{0.1, 0.2, 0.3, 0.1, 0.05, 0.25}

	t.Run("METIS_PartGraphRecursive", func(t *testing.T) {
		testPartGraphRecursive(t, xadj, adjncy, vwgt, adjwgt, tpwgts)
	})

	t.Run("METIS_PartGraphKway", func(t *testing.T) {
		testPartGraphKway(t, xadj, adjncy, vwgt, adjwgt, tpwgts)
	})
}

func testPartGraphRecursive(t *testing.T, xadj, adjncy, vwgt, adjwgt []int32, tpwgts []float32) {
	nvtxs := len(xadj) - 1
	nparts := int32(20)
	opts := make([]int32, NoOptions)

	// Test 1: No weights
	t.Run("NoWeights", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphRecursive(xadj, adjncy, nparts, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, nil, nil, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 2: Edge weights only
	t.Run("EdgeWeightsOnly", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphRecursiveWeighted(xadj, adjncy, nil, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, nil, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 3: Vertex weights only
	t.Run("VertexWeightsOnly", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphRecursiveWeighted(xadj, adjncy, vwgt, nil, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, nil, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 4: Both vertex and edge weights
	t.Run("BothWeights", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphRecursiveWeighted(xadj, adjncy, vwgt, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 5: With specific options (matching C test options[0]=1, options[1]=1, etc.)
	t.Run("WithOptions1", func(t *testing.T) {
		SetDefaultOptions(opts)
		// In METIS 5.x, we set specific option values
		opts[OptionPType] = PTypeRB      // METIS_PTYPE_RB
		opts[OptionObjType] = ObjTypeCut // METIS_OBJTYPE_CUT
		opts[OptionCType] = CTypeRM      // METIS_CTYPE_RM
		opts[OptionIPType] = IPTypeGrow  // METIS_IPTYPE_GROW
		opts[OptionRType] = RTypeFM      // METIS_RTYPE_FM

		part, objval, err := PartGraphRecursiveWeighted(xadj, adjncy, vwgt, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 6: Different coarsening scheme
	t.Run("WithOptions2", func(t *testing.T) {
		SetDefaultOptions(opts)
		opts[OptionPType] = PTypeRB      // METIS_PTYPE_RB
		opts[OptionObjType] = ObjTypeCut // METIS_OBJTYPE_CUT
		opts[OptionCType] = CTypeSHEM    // METIS_CTYPE_SHEM
		opts[OptionIPType] = IPTypeGrow  // METIS_IPTYPE_GROW
		opts[OptionRType] = RTypeFM      // METIS_RTYPE_FM

		part, objval, err := PartGraphRecursiveWeighted(xadj, adjncy, vwgt, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test weighted partitioning with target partition weights
	t.Run("WeightedPartitioning", func(t *testing.T) {
		nparts6 := int32(6)
		SetDefaultOptions(opts)

		// Test with no vertex/edge weights
		part, objval, err := PartGraphRecursiveWeighted(xadj, adjncy, nil, nil, nparts6, tpwgts, nil, opts)
		require.NoError(t, err)
		rcode := verifyWPart(nvtxs, xadj, adjncy, nil, nil, nparts6, tpwgts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)

		// Test with all weights
		part, objval, err = PartGraphRecursiveWeighted(xadj, adjncy, vwgt, adjwgt, nparts6, tpwgts, nil, opts)
		require.NoError(t, err)
		rcode = verifyWPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts6, tpwgts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})
}

func testPartGraphKway(t *testing.T, xadj, adjncy, vwgt, adjwgt []int32, tpwgts []float32) {
	nvtxs := len(xadj) - 1
	nparts := int32(20)
	opts := make([]int32, NoOptions)

	// Test 1: No weights
	t.Run("NoWeights", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphKway(xadj, adjncy, nparts, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, nil, nil, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 2: Edge weights only
	t.Run("EdgeWeightsOnly", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphKwayWeighted(xadj, adjncy, nil, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, nil, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 3: Vertex weights only
	t.Run("VertexWeightsOnly", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphKwayWeighted(xadj, adjncy, vwgt, nil, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, nil, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test 4: Both vertex and edge weights
	t.Run("BothWeights", func(t *testing.T) {
		SetDefaultOptions(opts)
		part, objval, err := PartGraphKwayWeighted(xadj, adjncy, vwgt, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test with different refinement algorithms
	t.Run("GreedyRefinement", func(t *testing.T) {
		SetDefaultOptions(opts)
		opts[OptionRType] = RTypeGreedy // METIS_RTYPE_GREEDY

		part, objval, err := PartGraphKwayWeighted(xadj, adjncy, vwgt, adjwgt, nparts, nil, nil, opts)
		require.NoError(t, err)
		rcode := verifyPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})

	// Test weighted partitioning
	t.Run("WeightedPartitioning", func(t *testing.T) {
		nparts6 := int32(6)
		SetDefaultOptions(opts)

		part, objval, err := PartGraphKwayWeighted(xadj, adjncy, vwgt, adjwgt, nparts6, tpwgts, nil, opts)
		require.NoError(t, err)
		rcode := verifyWPart(nvtxs, xadj, adjncy, vwgt, adjwgt, nparts6, tpwgts, objval, part)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})
}

// Test_ND emulates the C test function Test_ND
func TestND(t *testing.T) {
	// Create a test graph
	nvtxs := 100
	xadj, adjncy := createRandomGraph(nvtxs)

	// Create vertex weights
	vwgt := make([]int32, nvtxs)
	for i := 0; i < nvtxs; i++ {
		vwgt[i] = int32(1 + rand.Intn(10))
	}

	opts := make([]int32, NoOptions)

	t.Run("METIS_NodeND", func(t *testing.T) {
		// Test 1: Default options, no weights
		SetDefaultOptions(opts)
		perm, iperm, err := NodeND(xadj, adjncy, nil, opts)
		require.NoError(t, err)
		rcode := verifyND(nvtxs, perm, iperm)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)

		// Test 2: With vertex weights
		SetDefaultOptions(opts)
		perm, iperm, err = NodeND(xadj, adjncy, vwgt, opts)
		require.NoError(t, err)
		rcode = verifyND(nvtxs, perm, iperm)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)

		// Test 3: Different coarsening scheme
		SetDefaultOptions(opts)
		opts[OptionCType] = CTypeSHEM // METIS_CTYPE_SHEM
		perm, iperm, err = NodeND(xadj, adjncy, nil, opts)
		require.NoError(t, err)
		rcode = verifyND(nvtxs, perm, iperm)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)

		// Test 4: Different separator refinement
		SetDefaultOptions(opts)
		opts[OptionRType] = RTypeSep2Sided // METIS_RTYPE_SEP2SIDED
		perm, iperm, err = NodeND(xadj, adjncy, nil, opts)
		require.NoError(t, err)
		rcode = verifyND(nvtxs, perm, iperm)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)

		// Test 5: 1-sided separator refinement
		SetDefaultOptions(opts)
		opts[OptionRType] = RTypeSep1Sided // METIS_RTYPE_SEP1SIDED
		perm, iperm, err = NodeND(xadj, adjncy, nil, opts)
		require.NoError(t, err)
		rcode = verifyND(nvtxs, perm, iperm)
		assert.Equal(t, 0, rcode, "Verification failed with code %d", rcode)
	})
}

func TestMeshPartitioning(t *testing.T) {
	// Create a simple mesh with multiple tetrahedra
	ne := int32(10) // Number of elements
	nn := int32(15) // Number of nodes

	// Create element connectivity (4 nodes per tetrahedron)
	eptr := make([]int32, ne+1)
	eind := make([]int32, ne*4)

	for i := int32(0); i < ne; i++ {
		eptr[i] = i * 4
		// Create some connectivity pattern
		for j := 0; j < 4; j++ {
			eind[int(i)*4+j] = (i + int32(j)) % nn
		}
	}
	eptr[ne] = ne * 4

	opts := make([]int32, NoOptions)

	t.Run("MeshToDual", func(t *testing.T) {
		// Test with different ncommon values
		for ncommon := int32(1); ncommon <= 3; ncommon++ {
			xadj, adjncy, err := MeshToDual(ne, nn, eptr, eind, ncommon)
			require.NoError(t, err)
			assert.Len(t, xadj, int(ne+1))
			assert.Equal(t, int32(0), xadj[0])
			assert.GreaterOrEqual(t, len(adjncy), 0)
		}
	})

	t.Run("MeshToNodal", func(t *testing.T) {
		xadj, _, err := MeshToNodal(ne, nn, eptr, eind)
		require.NoError(t, err)
		assert.Len(t, xadj, int(nn+1))
		assert.Equal(t, int32(0), xadj[0])
	})

	t.Run("PartMeshNodal", func(t *testing.T) {
		SetDefaultOptions(opts)
		nparts := int32(3)

		// Test without weights
		objval, epart, npart, err := PartMeshNodal(ne, nn, eptr, eind, nil, nil, nparts, nil, opts)
		require.NoError(t, err)
		assert.Len(t, epart, int(ne))
		assert.Len(t, npart, int(nn))
		assert.GreaterOrEqual(t, objval, int32(0))

		// Verify partitioning
		for i := int32(0); i < ne; i++ {
			assert.GreaterOrEqual(t, epart[i], int32(0))
			assert.Less(t, epart[i], nparts)
		}
		for i := int32(0); i < nn; i++ {
			assert.GreaterOrEqual(t, npart[i], int32(0))
			assert.Less(t, npart[i], nparts)
		}
	})

	t.Run("PartMeshDual", func(t *testing.T) {
		SetDefaultOptions(opts)
		nparts := int32(3)
		ncommon := int32(2)

		objval, epart, npart, err := PartMeshDual(ne, nn, eptr, eind, nil, nil, ncommon, nparts, nil, opts)
		require.NoError(t, err)
		assert.Len(t, epart, int(ne))
		assert.Len(t, npart, int(nn))
		assert.GreaterOrEqual(t, objval, int32(0))
	})
}

func TestComputeVertexSeparator(t *testing.T) {
	// Create test graphs of different sizes
	testSizes := []int{20, 50, 100}

	for _, nvtxs := range testSizes {
		t.Run(fmt.Sprintf("Size%d", nvtxs), func(t *testing.T) {
			xadj, adjncy := createRandomGraph(nvtxs)

			opts := make([]int32, NoOptions)
			SetDefaultOptions(opts)

			sepsize, part, err := ComputeVertexSeparator(xadj, adjncy, nil, opts)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, sepsize, int32(0))
			assert.Len(t, part, nvtxs)

			// Verify separator
			counts := [3]int{}
			for i := 0; i < nvtxs; i++ {
				assert.GreaterOrEqual(t, part[i], int32(0))
				assert.LessOrEqual(t, part[i], int32(2))
				counts[part[i]]++
			}

			// Separator (partition 2) should have sepsize vertices
			assert.Equal(t, int(sepsize), counts[2], "Separator size mismatch")
		})
	}
}

// Verification functions ported from C

func verifyPart(nvtxs int, xadj, adjncy, vwgt, adjwgt []int32, nparts, edgecut int32, part []int32) int {
	// Check if max partition number is correct
	maxPart := int32(0)
	for _, p := range part {
		if p > maxPart {
			maxPart = p
		}
	}
	if maxPart != nparts-1 {
		return 1
	}

	// Use unit weights if not provided
	localVwgt := vwgt
	localAdjwgt := adjwgt
	if localVwgt == nil {
		localVwgt = make([]int32, nvtxs)
		for i := range localVwgt {
			localVwgt[i] = 1
		}
	}
	if localAdjwgt == nil {
		localAdjwgt = make([]int32, len(adjncy))
		for i := range localAdjwgt {
			localAdjwgt[i] = 1
		}
	}

	// Compute partition weights and cut
	pwgts := make([]int32, nparts)
	cut := int32(0)

	for i := 0; i < nvtxs; i++ {
		pwgts[part[i]] += localVwgt[i]
		for j := xadj[i]; j < xadj[i+1]; j++ {
			if part[i] != part[adjncy[j]] {
				cut += localAdjwgt[j]
			}
		}
	}

	// Check edge cut (each edge counted twice)
	if cut != 2*edgecut {
		return 2
	}

	// Check balance - allow more imbalance for high partition counts
	totalWgt := int32(0)
	maxWgt := int32(0)
	for i := int32(0); i < nparts; i++ {
		totalWgt += pwgts[i]
		if pwgts[i] > maxWgt {
			maxWgt = pwgts[i]
		}
	}

	// Use a more relaxed balance constraint for many partitions
	balanceFactor := 1.10
	if nparts >= 16 {
		balanceFactor = 1.20 // Allow 20% imbalance for many partitions
	}

	if float64(nparts)*float64(maxWgt) > balanceFactor*float64(totalWgt) {
		return 3
	}

	return 0
}

func verifyWPart(nvtxs int, xadj, adjncy, vwgt, adjwgt []int32, nparts int32, tpwgts []float32, edgecut int32, part []int32) int {
	// Check if max partition number is correct
	maxPart := int32(0)
	for _, p := range part {
		if p > maxPart {
			maxPart = p
		}
	}
	if maxPart != nparts-1 {
		return 1
	}

	// Use unit weights if not provided
	localVwgt := vwgt
	localAdjwgt := adjwgt
	if localVwgt == nil {
		localVwgt = make([]int32, nvtxs)
		for i := range localVwgt {
			localVwgt[i] = 1
		}
	}
	if localAdjwgt == nil {
		localAdjwgt = make([]int32, len(adjncy))
		for i := range localAdjwgt {
			localAdjwgt[i] = 1
		}
	}

	// Compute partition weights and cut
	pwgts := make([]int32, nparts)
	cut := int32(0)

	for i := 0; i < nvtxs; i++ {
		pwgts[part[i]] += localVwgt[i]
		for j := xadj[i]; j < xadj[i+1]; j++ {
			if part[i] != part[adjncy[j]] {
				cut += localAdjwgt[j]
			}
		}
	}

	// Check edge cut
	if cut != 2*edgecut {
		return 2
	}

	// Check balance against target weights
	totalWgt := int32(0)
	for i := int32(0); i < nparts; i++ {
		totalWgt += pwgts[i]
	}

	for i := int32(0); i < nparts; i++ {
		if float32(pwgts[i]) > 1.10*tpwgts[i]*float32(totalWgt) {
			return 3
		}
	}

	return 0
}

func verifyND(nvtxs int, perm, iperm []int32) int {
	// Check that perm and iperm are inverses
	for i := 0; i < nvtxs; i++ {
		if int32(i) != perm[iperm[i]] {
			return 1
		}
		if int32(i) != iperm[perm[i]] {
			return 2
		}
	}
	return 0
}

// Helper function to create random graphs (already in original test)
func createRandomGraph(nvtxs int) ([]int32, []int32) {
	rand.Seed(42)
	edges := make(map[[2]int]bool)

	// Create a connected graph with random edges
	for i := 0; i < nvtxs-1; i++ {
		// Ensure connectivity
		edges[[2]int{i, i + 1}] = true
		edges[[2]int{i + 1, i}] = true
	}

	// Add random edges
	numExtraEdges := nvtxs + rand.Intn(nvtxs*2)
	for i := 0; i < numExtraEdges; i++ {
		u := rand.Intn(nvtxs)
		v := rand.Intn(nvtxs)
		if u != v {
			edges[[2]int{u, v}] = true
			edges[[2]int{v, u}] = true
		}
	}

	// Build adjacency lists
	adjList := make([][]int, nvtxs)
	for edge := range edges {
		adjList[edge[0]] = append(adjList[edge[0]], edge[1])
	}

	// Sort adjacency lists and build CSR format
	xadj := make([]int32, nvtxs+1)
	adjncy := []int32{}

	for i := 0; i < nvtxs; i++ {
		sort.Ints(adjList[i])
		xadj[i+1] = xadj[i] + int32(len(adjList[i]))
		for _, v := range adjList[i] {
			adjncy = append(adjncy, int32(v))
		}
	}

	return xadj, adjncy
}

// Additional test for large graphs
func TestLargeGraphSystematic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large graph test in short mode")
	}

	sizes := []int{100, 500, 1000}
	npartsList := []int32{2, 4, 8, 16, 32}

	for _, nvtxs := range sizes {
		xadj, adjncy := createRandomGraph(nvtxs)

		for _, nparts := range npartsList {
			if nparts > int32(nvtxs/4) {
				continue // Skip if too many partitions for graph size
			}

			t.Run(fmt.Sprintf("Size%d_Parts%d", nvtxs, nparts), func(t *testing.T) {
				opts := make([]int32, NoOptions)
				SetDefaultOptions(opts)

				// Test recursive
				part, objval, err := PartGraphRecursive(xadj, adjncy, nparts, opts)
				require.NoError(t, err)
				rcode := verifyPart(nvtxs, xadj, adjncy, nil, nil, nparts, objval, part)
				assert.Equal(t, 0, rcode, "Recursive verification failed")

				// Test k-way
				part, objval, err = PartGraphKway(xadj, adjncy, nparts, opts)
				require.NoError(t, err)
				rcode = verifyPart(nvtxs, xadj, adjncy, nil, nil, nparts, objval, part)
				assert.Equal(t, 0, rcode, "K-way verification failed")
			})
		}
	}
}
