package deps

import (
	err "errors"
)

type DFS struct {
	currentPath []*Node
	visited     VisitMap
	edges       NodeSliceMap
	leavesOnly  bool
	result      []*Node
	visitor     func(*Node)
}

type VisitMap map[*Node]bool

type NodeSliceMap map[*Node][]*Node

type Node string

type Graph struct {
	nodes    map[*Node]*Node
	incoming NodeSliceMap
	outgoing NodeSliceMap
}

func MakeGraph() *Graph {
	return &Graph{
		nodes:    make(map[*Node]*Node),
		incoming: make(map[*Node][]*Node),
		outgoing: make(map[*Node][]*Node),
	}
}

func NewNode(id string) *Node {
	n := Node(id)
	return &n
}

func (self *Graph) HasNode(node *Node) bool {
	_, ok := self.nodes[node]
	if ok {
		return true
	}
	return false
}

func (self *Graph) Add(node *Node) (*Graph, error) {
	if !self.HasNode(node) {
		self.nodes[node] = node
		self.incoming[node] = make([]*Node, 0, 0)
		self.outgoing[node] = make([]*Node, 0, 0)
		return self, nil
	}
	return self, err.New("Graph already contains node.")
}

func SliceIndex(m []*Node, predicate func(*Node) bool) int {
	for i, e := range m {
		if predicate(e) {
			return i
		}
	}
	return -1
}

func SafeCut(l []*Node, idx int) []*Node {
	l[idx] = nil // GC can nail it
	copy(l[idx:], l[idx+1:])
	return l[:len(l)-1]
}

func (self *Graph) Remove(node *Node) (*Graph, error) {
	if self.HasNode(node) {
		delete(self.nodes, node)
		delete(self.incoming, node)
		delete(self.outgoing, node)
		for _, list := range self.incoming {
			idx := SliceIndex(list,
				func(e *Node) bool {
					return e == node
				})
			if idx != -1 {
				SafeCut(list, idx)
			}
		}
		return self, nil
	}
	return self, err.New("Node not found.")
}

func (self *Graph) AddDependency(from *Node, to *Node) (*Graph, error) {
	if !self.HasNode(from) {
		return self, err.New("From node not found.")
	}
	if !self.HasNode(to) {
		return self, err.New("To node not found.")
	}
	idx := SliceIndex(self.outgoing[from],
		func(e *Node) bool { return e == to })
	if idx == -1 {
		self.outgoing[from] = append(self.outgoing[from], to)
	}
	idx = SliceIndex(self.incoming[to],
		func(e *Node) bool { return e == from })
	if idx == -1 {
		self.incoming[to] = append(self.incoming[to], from)
	}
	return self, nil
}

func (self *Graph) RemoveDependency(from *Node, to *Node) (*Graph, error) {
	hasFrom := self.HasNode(from)
	hasTo := self.HasNode(to)
	if hasFrom {
		idx := SliceIndex(self.outgoing[from],
			func(e *Node) bool { return e == to })
		if idx >= 0 {
			SafeCut(self.outgoing[from], idx)
		}
	}
	if hasTo {
		idx := SliceIndex(self.incoming[to],
			func(e *Node) bool { return e == from })
		if idx >= 0 {
			SafeCut(self.incoming[to], idx)
		}
	}
	return self, nil
}

func CreateDFS(edges NodeSliceMap, leavesOnly bool) *DFS {
	dfs := DFS{
		currentPath: []*Node{},
		visited:     VisitMap{},
		edges:       edges,
		leavesOnly:  leavesOnly,
		result:      []*Node{},
	}
	dfs.visitor = func(current *Node) {
		dfs.visited[current] = true
		dfs.currentPath = append(dfs.currentPath, current)
		for _, n := range dfs.edges[current] {
			_, found := dfs.visited[n]
			if !found {
				dfs.visitor(n)
			} else {
				idx := SliceIndex(dfs.currentPath,
					func(e *Node) bool { return e == n })
				if idx >= 0 {
					dfs.currentPath = append(dfs.currentPath, n)
					panic("Dependency cycle!")
				}
			}
		}
		// pop
		clen := len(dfs.currentPath)
		_, dfs.currentPath = dfs.currentPath[clen-1], dfs.currentPath[:clen-1]

		i := SliceIndex(dfs.result, func(e *Node) bool { return e == current })
		if (!dfs.leavesOnly || (len(dfs.edges) == 0)) && i == -1 {
			dfs.result = append(dfs.result, current)
		}
	}
	return &dfs
}

func (self *Graph) depender(m NodeSliceMap, node *Node, leavesOnly bool) ([]*Node, error) {
	if self.HasNode(node) {
		dfs := CreateDFS(m, leavesOnly)
		dfs.visitor(node)
		idx := SliceIndex(dfs.result, func(e *Node) bool { return e == node })
		if idx >= 0 {
			dfs.result = dfs.result[:idx]
		}
		return dfs.result, nil
	}
	return nil, err.New("Node not found.")
}

func (self *Graph) DependenciesOf(node *Node, leavesOnly bool) ([]*Node, error) {
	return self.depender(self.outgoing, node, leavesOnly)
}

func (self *Graph) DependantsOf(node *Node, leavesOnly bool) ([]*Node, error) {
	return self.depender(self.incoming, node, leavesOnly)
}

/**
 * Construct the overall processing order for the dependency graph.
 *
 * Throws an Error if the graph has a cycle.
 *
 * If `leavesOnly` is true, only nodes that do not depend on any other nodes will be returned.
 **/
func (self *Graph) Order(leavesOnly bool) ([]*Node, error) {
	if len(self.nodes) == 0 {
		return []*Node{}, nil
	}

	// Look for cycles - we run the DFS starting at all the nodes in case there
	// are several disconnected subgraphs inside this dependency graph.
	cycles := CreateDFS(self.outgoing, false)
	for n := range self.outgoing {
		cycles.visitor(n)
	}

	// Find all roots and run a dfs staring at these nodes to get the order
	// dfs := CreateDFS(self.outgoing, leavesOnly)
	roots := []*Node{}
	for n := range self.outgoing {
		if len(self.incoming[n]) == 0 {
			roots = append(roots, n)
		}
	}

	dfs := CreateDFS(self.outgoing, leavesOnly)
	for _, n := range roots {
		dfs.visitor(n)
	}

	return dfs.result, nil
}
