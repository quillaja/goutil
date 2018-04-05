package data

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

// This is not the most rigorous testing.

type point [2]float64

func (p *point) Location() []float64 {
	return p[:]
}

func makeItems(n int, max float64) []Interface {
	// rand.Seed(time.Now().UnixNano())
	items := []Interface{}
	for i := 0; i < n; i++ {
		p := &point{max * rand.Float64(), max * rand.Float64()}
		items = append(items, p)
	}
	return items
}

func TestKDTree_Insert(t *testing.T) {

	const num, max = 100, 100

	tree := NewKDTree(2)
	for i := 0; i < 2; i++ {
		items := makeItems(num*(i+1), max)

		tree.Build(items)

		t.Logf("inserted %d items", num*(i+1))
		t.Logf("tree len: %d", tree.Len())
		t.Logf("tree dimens: %d", tree.Dimensions())
		t.Logf("first item: %v", tree.Items()[0].(*point))
		if tree.Len() != num*(i+1) {
			t.Fail()
		}
		if tree.Dimensions() != 2 {
			t.Fail()
		}
	}
}

func TestKDTree_QueryPoint(t *testing.T) {
	t.Log("testing that all inserted items are found in tree")
	items := makeItems(100, 100)

	tree := NewKDTree(2)
	tree.Build(items)
	t.Logf("tree len: %d", tree.Len())

	for _, item := range items {
		if !tree.QueryPoint(item) {
			t.Logf("item %v not found in tree", item)
			t.Fail()
		}
	}
}

func TestKDTree_QueryRange(t *testing.T) {
	t.Log("testing that all items returned from the range query belong,\nand that none of the others are in the result")
	testQueryRange(t, [][2]float64{{15, 35}, {15, 35}})
}

func TestKDTree_QueryRange_GiantRange(t *testing.T) {
	t.Log("same test as normal test of QueryRange, but uses whole area as search")
	testQueryRange(t, [][2]float64{{0, 50}, {0, 50}})
}

func testQueryRange(t *testing.T, searchRange [][2]float64) {
	const max = 50

	items := makeItems(100, max)
	tree := NewKDTree(2)
	tree.Build(items)
	t.Logf("tree contains: %d items", tree.Len())

	//items is now in the tree and the order and indexing is stable.
	// searchRange := [][2]float64{{max/2 - 10, max/2 + 10}, {max/2 - 10, max/2 + 10}}
	belongs := make([]bool, len(items), len(items))
	for i := 0; i < len(items); i++ {
		loc := items[i].Location()
		if searchRange[0][0] <= loc[0] && loc[0] <= searchRange[0][1] &&
			searchRange[1][0] <= loc[1] && loc[1] <= searchRange[1][1] {
			// item belongs
			belongs[i] = true
		}
	}

	found := tree.QueryRange(searchRange)
	for i := 0; i < len(items); i++ {
		for _, f := range found {
			// fail if items[i] is in found and does not belong
			if items[i] == f && !belongs[i] {
				t.Logf("item in result that does not belong: %v", items[i])
				t.Fail()
			}
		}
	}
	t.Logf("found %d items in %v", len(found), searchRange)
}

func TestKDTree_NearestNeighbor(t *testing.T) {
	t.Log("make points in one area of the graph, then manually insert one in an empty region. Test search point near that one.")

	items := []Interface{}
	for i := 0; i < 10; i++ {
		items = append(items, &point{rand.Float64()*5 + 5, rand.Float64()*5 + 5})
	}
	items = append(items, &point{1, 1})

	tree := NewKDTree(2)
	tree.Build(items)

	found := tree.NearestNeighbor(Euclidean, 0, 0) // could also use []float64...
	t.Log("search", 0, 0)
	t.Log("found", found)
	if found.Location()[0] != 1 && found.Location()[1] != 1 {
		t.Fail()
	}
	if !reflect.DeepEqual(found, bruteForceNN(Euclidean, 1, items, &point{0, 0})[0]) {
		t.Log("tree nn != bf nn, via reflect.DeepEqual")
		t.Fail()
	}

	found = tree.NearestNeighbor(Euclidean, 10, 10)
	t.Log("search", 10, 10)
	t.Log("found", found)
	if found.Location()[0] < 5 || found.Location()[1] < 5 {
		t.Fail()
	}
	if !reflect.DeepEqual(found, bruteForceNN(Euclidean, 1, items, &point{10, 10})[0]) {
		t.Log("tree nn != bf nn, via reflect.DeepEqual")
		t.Fail()
	}

}

func TestKDTree_NearestNeighbors(t *testing.T) {
	items := makeItems(50, 20)
	tree := NewKDTree(2)
	tree.Build(items)
	ks := []int{1, 2, 5, 8, 15}
	for _, k := range ks {
		found := tree.NearestNeighbors(Euclidean, k, 10, 10)
		bffound := bruteForceNN(Euclidean, k, items, &point{10, 10})
		t.Logf("search [10,10] for k=%d.", k)
		t.Logf("found %d :", len(found))
		for i, f := range found {
			t.Logf("\tkd%v\n\t\tbf%v", f.Location(), bffound[i].Location())
		}
		if len(found) != k {
			t.Log("len found != k")
			t.Fail()
		}
		if len(found) != len(bffound) {
			t.Log("len found != len bffound")
			t.Fail()
		}
		if !reflect.DeepEqual(found, bffound) {
			t.Log("tree nn != bf nn, via reflect.DeepEqual")
			t.Fail()
		}

	}
}

func bruteForceNN(dist DistanceMetric, k int, items []Interface, search Interface) (found []Interface) {
	type dp struct {
		d float64
		p Interface
	}

	dps := []dp{}
	for _, item := range items {
		dps = append(dps, dp{
			d: dist(search.Location(), item.Location()),
			p: item,
		})
	}

	sort.Slice(dps, func(i, j int) bool {
		return dps[i].d < dps[j].d
	})

	for i := 0; i < k; i++ {
		found = append(found, dps[i].p)
	}

	return
}
