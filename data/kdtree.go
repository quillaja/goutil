package data

import (
	"math"
	"sort"
)

// a node in the kdtree
type kdnode struct {
	axis  int
	data  Interface
	left  *kdnode
	right *kdnode
}

// KDTree implements SpacialTree using the kd-tree data structure.
type KDTree struct {
	root       *kdnode
	dimensions int
	items      []Interface
}

// NewKDTree creates an empty tree with the capacity to hold the
// number of dimensions specified.
func NewKDTree(dimensions int) *KDTree {
	return &KDTree{
		root:       nil,
		dimensions: dimensions,
		items:      nil,
	}
}

// Items returns a slice of the items held in the tree.
func (t *KDTree) Items() []Interface {
	return t.items
}

// Dimensions returns the number of dimensions the tree uses.
func (t *KDTree) Dimensions() int {
	return t.dimensions
}

// Len returns the number of items in the tree.
func (t *KDTree) Len() int {
	return len(t.items)
}

// Build will build (or rebuild) the tree with the given items.
func (t *KDTree) Build(items []Interface) {
	// check that all items have correct
	// number of dimensions (avoid index out of bounds)
	var sum int
	for i := 0; i < len(items); i++ {
		sum += len(items[i].Location())
	}
	if sum != t.dimensions*len(items) {
		panic("at least one element in 'items' does not have the expected number of dimensions")
	}

	// if t.root == nil {
	t.items = items
	t.root = buildTree(t.items, 0, t.dimensions)
	// }
}

// does actual tree build
func buildTree(items []Interface, depth, dims int) (node *kdnode) {
	if len(items) == 0 {
		return nil
	}

	// ascending sort items by axis
	axis := depth % dims // 0=x, 1=y, 2=z (for Vec3)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Location()[axis] < items[j].Location()[axis]
	})

	// create node
	median := len(items) / 2
	node = &kdnode{
		data: items[median],
		axis: axis}

	node.left = buildTree(items[:median], depth+1, dims)
	node.right = buildTree(items[median+1:], depth+1, dims)

	return
}

// QueryPoint returns true if the item is found in the tree.
func (t *KDTree) QueryPoint(item Interface) bool {
	return dfsPoint(t.root, item)
}

// used in QueryPoint()
func dfsPoint(node *kdnode, item Interface) (found bool) {
	// 1. check current node
	if node == nil {
		return false
	}
	if node.data == item {
		return true
	}

	// 2. if not it, compare item to node's item to
	// determine which branch to follow. If node and item
	// are equal on the axis, have to check both branches.
	nodeAxialVal := node.data.Location()[node.axis]
	itemAxialVal := item.Location()[node.axis]
	if itemAxialVal <= nodeAxialVal {
		found = dfsPoint(node.left, item)
	}
	if !found && itemAxialVal >= nodeAxialVal {
		found = dfsPoint(node.right, item)
	}

	return
}

// QueryRange returns all items within the n-dimensional range specified.
// Each entry to 'ranges' is a [2]float64 where [0] = MIN and [1] = MAX of range.
func (t *KDTree) QueryRange(ranges [][2]float64) []Interface {
	if len(ranges) != t.Dimensions() {
		panic("incorrect number of dimensions in 'ranges'")
	}
	found := make([]Interface, 0, t.Len()/4) // starting cap 25% of size
	dfsRange(t.root, ranges, &found)
	return found
}

// used in QueryRange()
func dfsRange(node *kdnode, ranges [][2]float64, found *[]Interface) {
	// using DFS
	// 1. check node for nil, then check each of the node's
	// n-dimensional values are within the corresponding range.
	// 1.1 if so, add node.data to return slice
	if node == nil {
		return
	}
	inSearchRange := true
	itemLoc := node.data.Location()
	for axis, r := range ranges {
		if !(r[0] <= itemLoc[axis] && itemLoc[axis] <= r[1]) {
			inSearchRange = false
			break
		}
	}
	if inSearchRange {
		*found = append(*found, node.data)
	}

	// 2. determine which branch(s) to go down.
	// 2.1 if range's axial MAX is <= node's axial val, go left only
	// 2.2 if range's axial MIN is >= node's axial val, go right only
	// 2.3 if the node's axial val is IN the axial range, go down both
	axialRange := ranges[node.axis]
	nodeAxialVal := node.data.Location()[node.axis]
	nodeInRange := axialRange[0] <= nodeAxialVal && nodeAxialVal <= axialRange[1]
	if nodeInRange || axialRange[1] <= nodeAxialVal {
		dfsRange(node.left, ranges, found)
	}
	if nodeInRange || axialRange[0] >= nodeAxialVal {
		dfsRange(node.right, ranges, found)
	}
}

///// Things used in nearest neighbors ////

// used in nearest neighbor searches for best candidate(s)
type neigh struct {
	node *kdnode
	dist float64
}

// inserts and element into slice at the index.
func insertAndTrim(item *neigh, at int, s []*neigh) {
	// insert
	s = append(s, nil)
	copy(s[at+1:], s[at:])
	s[at] = item

	// remove end
	s[len(s)-1] = nil
	s = s[:len(s)-1]
}

///////////////////////////////////////

// NearestNeighbor finds the nearest neighbor to searchPt using the given
// distance metric. Returns nil if none found or if the tree's root is nil.
func (t *KDTree) NearestNeighbor(dist DistanceMetric, point ...float64) Interface {
	best := neigh{nil, math.Inf(0)}
	nnSearch(t.root, point, &best, dist)
	if best.node == nil {
		return nil
	}
	return best.node.data
}

// Does actual nearest neighbor search
func nnSearch(root *kdnode, searchPt []float64, curBest *neigh, dist DistanceMetric) {
	// if the current node is nil, just return
	if root == nil {
		return
	}

	// decide which branch to visit first, then visit it.
	// this lets search start at a leave, which should provide potentially
	// better curBests than starting at the root.
	var goDown *kdnode
	if searchPt[root.axis] <= root.data.Location()[root.axis] {
		goDown = root.left
	} else {
		goDown = root.right
	}
	nnSearch(goDown, searchPt, curBest, dist)

	// check if current node is better than current best.
	// if current best == nil/inf, set current node to best.
	if d := dist(root.data.Location(), searchPt); curBest.node == nil || d < curBest.dist {
		curBest.node = root
		curBest.dist = d
	}

	// check if points could possibly exist on the other side of the root's splitting
	// axis by checking if the distance from the searchPt to axis is less than
	// the distance to the current best.
	// searchPt-to-axis = abs(root.data.location()[axis] - seachPt[axis])
	// if search-to-axis <= curbest.dist, then go down the branch NOT taken earlier.
	searchToAxis := dist([]float64{searchPt[root.axis]}, []float64{root.data.Location()[root.axis]})
	checkBoth := searchToAxis <= curBest.dist

	// go down one not visited earlier, if required
	if goDown == root.left {
		goDown = root.right
	} else {
		goDown = root.left
	}
	if checkBoth {
		nnSearch(goDown, searchPt, curBest, dist)
	}

	return
}

// NearestNeighbors returns the nearest [0,k] neighbors to the search point.
// If fewer than k are found, the returned slice will be as long as the number
// found. Distance is determined by the given DistanceMetric.
func (t *KDTree) NearestNeighbors(dist DistanceMetric, k int, point ...float64) []Interface {
	// MUST have k+1 capacity, or the append() to the bests slice inside
	// insertAndTrim() will cause a new backing array to be allocated, and so
	// the array we want to change is NOT changed...very subtle.
	bests := make([]*neigh, k, k+1)
	knnSearch(dist, t.root, point, bests) // will alter bests
	var found []Interface
	for _, b := range bests {
		if b != nil {
			found = append(found, b.node.data)
		}
	}
	return found
}

// does actual nn search for k nodes
// curBests is a best-to-worst ORDERED list of k elements (some of which may be nil)
func knnSearch(dist DistanceMetric, root *kdnode, searchPt []float64, curBests []*neigh) {
	if root == nil {
		return
	}

	// choose and go down one branch
	var goDown *kdnode
	if searchPt[root.axis] <= root.data.Location()[root.axis] {
		goDown = root.left
	} else {
		goDown = root.right
	}
	knnSearch(dist, goDown, searchPt, curBests)

	// examine the current node
	d := dist(root.data.Location(), searchPt)
	for i := 0; i < len(curBests); i++ {
		// check each. if found a best.dist > root.dist, insert
		// to keep order and remove the worst best from the end.
		// if nil is encountered, insert.
		if curBests[i] == nil {
			curBests[i] = &neigh{root, d}
			break
		}
		if d < curBests[i].dist {
			insertAndTrim(&neigh{root, d}, i, curBests)
			break
		}
	}

	// go down other branch if necessary.
	// use similar process as nnSearch() but use worst best.
	worstBest := curBests[len(curBests)-1] // should be last
	searchToAxis := dist([]float64{searchPt[root.axis]}, []float64{root.data.Location()[root.axis]})
	checkBoth := worstBest == nil || searchToAxis < worstBest.dist

	if goDown == root.left {
		goDown = root.right
	} else {
		goDown = root.left
	}
	if checkBoth {
		knnSearch(dist, goDown, searchPt, curBests)
	}

	return
}
