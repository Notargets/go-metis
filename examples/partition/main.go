package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/notargets/go-metis"
)

func main() {
	// Command line flags
	graphFile := flag.String("graph", "", "Input graph file in METIS format")
	nparts := flag.Int("nparts", 2, "Number of partitions")
	method := flag.String("method", "kway", "Partitioning method: kway or recursive")
	outFile := flag.String("output", "partition.txt", "Output partition file")
	objective := flag.String("objective", "cut", "Objective: cut or vol")
	seed := flag.Int("seed", -1, "Random seed (-1 for default)")
	verbose := flag.Bool("verbose", false, "Verbose output")

	flag.Parse()

	if *graphFile == "" {
		fmt.Fprintf(os.Stderr, "Usage: %s -graph <file> -nparts <n> [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Read graph file
	if *verbose {
		fmt.Printf("Reading graph from %s...\n", *graphFile)
	}

	file, err := os.Open(*graphFile)
	if err != nil {
		log.Fatalf("Failed to open graph file: %v", err)
	}
	defer file.Close()

	graph, err := metis.ReadGraphFile(file)
	if err != nil {
		log.Fatalf("Failed to read graph: %v", err)
	}

	if *verbose {
		fmt.Printf("Graph loaded: %d vertices, %d edges\n",
			graph.NumVertices(), graph.NumEdges())

		// Print degree distribution
		minDeg, maxDeg, avgDeg := calculateDegreeStats(graph)
		fmt.Printf("Degree stats: min=%d, max=%d, avg=%.1f\n",
			minDeg, maxDeg, avgDeg)
	}

	// Set up options
	opts := make([]int32, metis.NoOptions)
	if err := metis.SetDefaultOptions(opts); err != nil {
		log.Fatalf("Failed to set options: %v", err)
	}

	// Set objective
	if *objective == "vol" {
		opts[metis.OptionObjType] = metis.ObjTypeVol
	} else {
		opts[metis.OptionObjType] = metis.ObjTypeCut
	}

	// Set seed if specified
	if *seed >= 0 {
		opts[metis.OptionSeed] = int32(*seed)
	}

	// Set verbosity
	if *verbose {
		opts[metis.OptionDBGLvl] = metis.DBGInfo | metis.DBGTime
	}

	// Perform partitioning
	start := time.Now()

	var part []int32
	var objval int32

	if *verbose {
		fmt.Printf("\nPartitioning graph into %d parts using %s method...\n",
			*nparts, *method)
	}

	if *method == "recursive" {
		part, objval, err = metis.PartGraphRecursive(
			graph.Xadj, graph.Adjncy, int32(*nparts), opts)
	} else {
		part, objval, err = metis.PartGraphKway(
			graph.Xadj, graph.Adjncy, int32(*nparts), opts)
	}

	if err != nil {
		log.Fatalf("Partitioning failed: %v", err)
	}

	elapsed := time.Since(start)

	// Calculate and display statistics
	fmt.Printf("\nPartitioning completed in %v\n", elapsed)
	fmt.Printf("Objective (%s): %d\n", *objective, objval)

	// Calculate partition statistics
	partStats := calculatePartitionStats(part, graph.Vwgt, int32(*nparts))
	fmt.Printf("\nPartition statistics:\n")
	for i := int32(0); i < int32(*nparts); i++ {
		fmt.Printf("  Partition %d: %d vertices", i, partStats[i].count)
		if graph.Vwgt != nil {
			fmt.Printf(" (weight: %d)", partStats[i].weight)
		}
		fmt.Printf("\n")
	}

	// Calculate balance
	min, max, avg := metis.CalculatePartitionBalance(part, graph.Vwgt, int32(*nparts))
	fmt.Printf("\nBalance: %.3f (max/avg)\n", max/avg)

	// Calculate edge cut (if using cut objective)
	if *objective == "cut" {
		edgeCut := metis.CalculateEdgeCut(graph, part)
		fmt.Printf("Edge cut: %d", edgeCut)
		if edgeCut != objval {
			fmt.Printf(" (warning: doesn't match objective value!)")
		}
		fmt.Printf("\n")
	}

	// Write output
	out, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer out.Close()

	if err := metis.WritePartitioning(out, part); err != nil {
		log.Fatalf("Failed to write partitioning: %v", err)
	}

	fmt.Printf("\nPartitioning written to %s\n", *outFile)

	// Additional output formats
	if *verbose {
		// Write statistics file
		statsFile := *outFile + ".stats"
		if err := writeStatisticsFile(statsFile, graph, part, int32(*nparts), objval, elapsed); err != nil {
			log.Printf("Warning: failed to write statistics file: %v", err)
		} else {
			fmt.Printf("Statistics written to %s\n", statsFile)
		}
	}
}

type partitionInfo struct {
	count  int
	weight int64
}

func calculatePartitionStats(part []int32, vwgt []int32, nparts int32) []partitionInfo {
	stats := make([]partitionInfo, nparts)

	for i, p := range part {
		stats[p].count++
		if vwgt != nil && i < len(vwgt) {
			stats[p].weight += int64(vwgt[i])
		} else {
			stats[p].weight++
		}
	}

	return stats
}

func calculateDegreeStats(graph *metis.Graph) (min, max int, avg float64) {
	nvtxs := graph.NumVertices()
	if nvtxs == 0 {
		return 0, 0, 0
	}

	min = int(graph.Xadj[1] - graph.Xadj[0])
	max = min
	total := 0

	for i := 0; i < nvtxs; i++ {
		degree := int(graph.Xadj[i+1] - graph.Xadj[i])
		if degree < min {
			min = degree
		}
		if degree > max {
			max = degree
		}
		total += degree
	}

	avg = float64(total) / float64(nvtxs)
	return
}

func writeStatisticsFile(filename string, graph *metis.Graph, part []int32, nparts, objval int32, elapsed time.Duration) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	fmt.Fprintf(file, "METIS Partitioning Statistics\n")
	fmt.Fprintf(file, "============================\n\n")

	fmt.Fprintf(file, "Graph Information:\n")
	fmt.Fprintf(file, "  Vertices: %d\n", graph.NumVertices())
	fmt.Fprintf(file, "  Edges: %d\n", graph.NumEdges())

	minDeg, maxDeg, avgDeg := calculateDegreeStats(graph)
	fmt.Fprintf(file, "  Degree: min=%d, max=%d, avg=%.1f\n\n", minDeg, maxDeg, avgDeg)

	fmt.Fprintf(file, "Partitioning Information:\n")
	fmt.Fprintf(file, "  Number of partitions: %d\n", nparts)
	fmt.Fprintf(file, "  Objective value: %d\n", objval)
	fmt.Fprintf(file, "  Time: %v\n\n", elapsed)

	fmt.Fprintf(file, "Partition Details:\n")
	stats := calculatePartitionStats(part, graph.Vwgt, nparts)
	for i := int32(0); i < nparts; i++ {
		fmt.Fprintf(file, "  Partition %d: %d vertices", i, stats[i].count)
		if graph.Vwgt != nil {
			fmt.Fprintf(file, ", weight=%d", stats[i].weight)
		}
		fmt.Fprintf(file, "\n")
	}

	min, max, avg := metis.CalculatePartitionBalance(part, graph.Vwgt, nparts)
	fmt.Fprintf(file, "\nBalance Information:\n")
	fmt.Fprintf(file, "  Min weight: %.0f\n", min)
	fmt.Fprintf(file, "  Max weight: %.0f\n", max)
	fmt.Fprintf(file, "  Avg weight: %.1f\n", avg)
	fmt.Fprintf(file, "  Imbalance: %.3f\n", max/avg)

	// Communication statistics
	edgeCut := metis.CalculateEdgeCut(graph, part)
	fmt.Fprintf(file, "\nCommunication:\n")
	fmt.Fprintf(file, "  Edge cut: %d\n", edgeCut)

	// Calculate communication volume (simplified)
	commVol := calculateCommunicationVolume(graph, part, nparts)
	fmt.Fprintf(file, "  Communication volume: %d\n", commVol)

	return nil
}

func calculateCommunicationVolume(graph *metis.Graph, part []int32, nparts int32) int {
	// Count unique partition pairs that communicate
	commPairs := make(map[[2]int32]bool)
	nvtxs := graph.NumVertices()

	for i := 0; i < nvtxs; i++ {
		myPart := part[i]
		neighbors := make(map[int32]bool)

		for j := graph.Xadj[i]; j < graph.Xadj[i+1]; j++ {
			neighborPart := part[graph.Adjncy[j]]
			if neighborPart != myPart {
				neighbors[neighborPart] = true
			}
		}

		// Add communication pairs
		for nPart := range neighbors {
			if myPart < nPart {
				commPairs[[2]int32{myPart, nPart}] = true
			} else {
				commPairs[[2]int32{nPart, myPart}] = true
			}
		}
	}

	return len(commPairs)
}
