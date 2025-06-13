package metis

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Graph represents a graph in CSR format
type Graph struct {
	Xadj   []int32 // Index array for adjacency lists
	Adjncy []int32 // Adjacency lists (concatenated)
	Vwgt   []int32 // Vertex weights (optional)
	Adjwgt []int32 // Edge weights (optional)
}

// NewGraph creates a new graph from adjacency information
func NewGraph(xadj, adjncy []int32) *Graph {
	return &Graph{
		Xadj:   xadj,
		Adjncy: adjncy,
	}
}

// NumVertices returns the number of vertices in the graph
func (g *Graph) NumVertices() int {
	return len(g.Xadj) - 1
}

// NumEdges returns the number of edges in the graph (counting each edge once)
func (g *Graph) NumEdges() int {
	return len(g.Adjncy) / 2
}

// Degree returns the degree of vertex v
func (g *Graph) Degree(v int) int {
	return int(g.Xadj[v+1] - g.Xadj[v])
}

// Neighbors returns the neighbors of vertex v
func (g *Graph) Neighbors(v int) []int32 {
	start := g.Xadj[v]
	end := g.Xadj[v+1]
	return g.Adjncy[start:end]
}

// ConvertToMetisGraph converts a mesh to a METIS graph for partitioning
func ConvertMeshToGraph(ne, nn int32, eptr, eind []int32, dual bool, ncommon int32) (*Graph, error) {
	var xadj, adjncy []int32
	var err error

	if dual {
		xadj, adjncy, err = MeshToDual(ne, nn, eptr, eind, ncommon)
	} else {
		xadj, adjncy, err = MeshToNodal(ne, nn, eptr, eind)
	}

	if err != nil {
		return nil, err
	}

	return &Graph{
		Xadj:   xadj,
		Adjncy: adjncy,
	}, nil
}

// ReadGraphFile reads a graph in METIS format
// Format:
// Line 1: <# vertices> <# edges> [fmt] [ncon]
// Following lines: vertex adjacency lists (and optional weights)
func ReadGraphFile(r io.Reader) (*Graph, error) {
	scanner := bufio.NewScanner(r)

	// Read header
	if !scanner.Scan() {
		return nil, fmt.Errorf("empty file")
	}

	header := strings.Fields(scanner.Text())
	if len(header) < 2 {
		return nil, fmt.Errorf("invalid header: %s", scanner.Text())
	}

	nvtxs, err := strconv.Atoi(header[0])
	if err != nil {
		return nil, fmt.Errorf("invalid number of vertices: %v", err)
	}

	// nedges, err := strconv.Atoi(header[1])
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid number of edges: %v", err)
	// }

	// Parse format flags if present
	hasVertexWeights := false
	hasEdgeWeights := false
	if len(header) >= 3 {
		fmt, _ := strconv.Atoi(header[2])
		hasVertexWeights = (fmt/10)%10 == 1
		hasEdgeWeights = fmt%10 == 1
	}

	// Read vertex data
	xadj := make([]int32, nvtxs+1)
	adjncy := []int32{}
	vwgt := []int32{}
	adjwgt := []int32{}

	xadj[0] = 0
	for i := 0; i < nvtxs; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("unexpected EOF at vertex %d", i)
		}

		fields := strings.Fields(scanner.Text())
		if len(fields) == 0 {
			continue // Empty adjacency list
		}

		fieldIdx := 0

		// Read vertex weight if present
		if hasVertexWeights {
			w, err := strconv.Atoi(fields[fieldIdx])
			if err != nil {
				return nil, fmt.Errorf("invalid vertex weight at vertex %d: %v", i, err)
			}
			vwgt = append(vwgt, int32(w))
			fieldIdx++
		}

		// Read adjacency list
		for j := fieldIdx; j < len(fields); j++ {
			if hasEdgeWeights && (j-fieldIdx)%2 == 1 {
				// This is an edge weight
				w, err := strconv.Atoi(fields[j])
				if err != nil {
					return nil, fmt.Errorf("invalid edge weight at vertex %d: %v", i, err)
				}
				adjwgt = append(adjwgt, int32(w))
			} else {
				// This is a vertex
				v, err := strconv.Atoi(fields[j])
				if err != nil {
					return nil, fmt.Errorf("invalid vertex id at vertex %d: %v", i, err)
				}
				// Convert to 0-based indexing
				adjncy = append(adjncy, int32(v-1))
			}
		}

		xadj[i+1] = int32(len(adjncy))
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	g := &Graph{
		Xadj:   xadj,
		Adjncy: adjncy,
	}

	if hasVertexWeights {
		g.Vwgt = vwgt
	}
	if hasEdgeWeights {
		g.Adjwgt = adjwgt
	}

	return g, nil
}

// WritePartitioning writes partition information to a writer
func WritePartitioning(w io.Writer, part []int32) error {
	for _, p := range part {
		if _, err := fmt.Fprintf(w, "%d\n", p); err != nil {
			return err
		}
	}
	return nil
}

// CalculateEdgeCut calculates the edge cut for a given partitioning
func CalculateEdgeCut(g *Graph, part []int32) int32 {
	edgeCut := int32(0)
	nvtxs := g.NumVertices()

	for i := 0; i < nvtxs; i++ {
		for j := g.Xadj[i]; j < g.Xadj[i+1]; j++ {
			neighbor := g.Adjncy[j]
			if part[i] != part[neighbor] {
				if g.Adjwgt != nil {
					edgeCut += g.Adjwgt[j]
				} else {
					edgeCut++
				}
			}
		}
	}

	return edgeCut / 2 // Each edge counted twice
}

// CalculatePartitionBalance calculates partition balance statistics
func CalculatePartitionBalance(part []int32, vwgt []int32, nparts int32) (min, max, avg float64) {
	partWeights := make([]int64, nparts)
	nvtxs := len(part)

	for i := 0; i < nvtxs; i++ {
		weight := int64(1)
		if vwgt != nil && i < len(vwgt) {
			weight = int64(vwgt[i])
		}
		partWeights[part[i]] += weight
	}

	totalWeight := int64(0)
	minWeight := partWeights[0]
	maxWeight := partWeights[0]

	for _, w := range partWeights {
		totalWeight += w
		if w < minWeight {
			minWeight = w
		}
		if w > maxWeight {
			maxWeight = w
		}
	}

	avgWeight := float64(totalWeight) / float64(nparts)

	return float64(minWeight), float64(maxWeight), avgWeight
}
