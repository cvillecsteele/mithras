// Copyright (C) 2013-2015 by Jim Riecken

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
//     in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

(function(window) {
    'use strict'

    /**
     * A simple dependency graph
     */

    /**
     * Helper for creating a Depth-First-Search on
     * a set of edges.
     *
     * Detects cycles and throws an Error if one is detected.
     *
     * @param edges The set of edges to DFS through
     * @param leavesOnly Whether to only return "leaf" nodes (ones who have no edges)
     * @param result An array in which the results will be populated
     */
    var createDFS = function (edges, leavesOnly, result) {
	var currentPath = [];
	var visited = {};
	return function DFS(currentNode) {
	    visited[currentNode] = true;
	    currentPath.push(currentNode);
	    edges[currentNode].forEach(function (node) {
		if (!visited[node]) {
		    DFS(node);
		} else if (currentPath.indexOf(node) >= 0) {
		    currentPath.push(node);
		    throw new Error('Dependency Cycle Found: ' + currentPath.join(' -> '));
		}
	    });
	    currentPath.pop();
	    if ((!leavesOnly || edges[currentNode].length === 0) && result.indexOf(currentNode) === -1) {
		result.push(currentNode);
	    }
	};
    }

    /**
     * Simple Dependency Graph
     */
    var DepGraph = function DepGraph() {
	this.nodes = {};
	this.outgoingEdges = {}; // Node -> [Dependency Node]
	this.incomingEdges = {}; // Node -> [Dependant Node]
    };

    DepGraph.prototype = {
	/**
	 * Add a node to the dependency graph. If a node already exists, this method will do nothing.
	 */
	addNode:function (node) {
	    if (!this.hasNode(node)) {
		this.nodes[node] = node;
		this.outgoingEdges[node] = [];
		this.incomingEdges[node] = [];
	    }
	},
	/**
	 * Remove a node from the dependency graph. If a node does not exist, this method will do nothing.
	 */
	removeNode:function (node) {
	    if (this.hasNode(node)) {
		delete this.nodes[node];
		delete this.outgoingEdges[node];
		delete this.incomingEdges[node];
		[this.incomingEdges, this.outgoingEdges].forEach(function (edgeList) {
		    Object.keys(edgeList).forEach(function (key) {
			var idx = edgeList[key].indexOf(node);
			if (idx >= 0) {
			    edgeList[key].splice(idx, 1);
			}
		    }, this);
		});
	    }
	},
	/**
	 * Check if a node exists in the graph
	 */
	hasNode:function (node) {
	    return !!this.nodes[node];
	},
	/**
	 * Add a dependency between two nodes. If either of the nodes does not exist,
	 * an Error will be thrown.
	 */
	addDependency:function (from, to) {
	    if (!this.hasNode(from)) {
		throw new Error('Node does not exist: ' + from);
	    }
	    if (!this.hasNode(to)) {
		throw new Error('Node does not exist: ' + to);
	    }
	    if (this.outgoingEdges[from].indexOf(to) === -1) {
		this.outgoingEdges[from].push(to);
	    }
	    if (this.incomingEdges[to].indexOf(from) === -1) {
		this.incomingEdges[to].push(from);
	    }
	    return true;
	},
	/**
	 * Remove a dependency between two nodes.
	 */
	removeDependency:function (from, to) {
	    var idx;
	    if (this.hasNode(from)) {
		idx = this.outgoingEdges[from].indexOf(to);
		if (idx >= 0) {
		    this.outgoingEdges[from].splice(idx, 1);
		}
	    }

	    if (this.hasNode(to)) {
		idx = this.incomingEdges[to].indexOf(from);
		if (idx >= 0) {
		    this.incomingEdges[to].splice(idx, 1);
		}
	    }
	},
	/**
	 * Get an array containing the nodes that the specified node depends on (transitively).
	 *
	 * Throws an Error if the graph has a cycle, or the specified node does not exist.
	 *
	 * If `leavesOnly` is true, only nodes that do not depend on any other nodes will be returned
	 * in the array.
	 */
	dependenciesOf:function (node, leavesOnly) {
	    if (this.hasNode(node)) {
		var result = [];
		var DFS = createDFS(this.outgoingEdges, leavesOnly, result);
		DFS(node);
		var idx = result.indexOf(node);
		if (idx >= 0) {
		    result.splice(idx, 1);
		}
		return result;
	    }
	    else {
		throw new Error('Node does not exist: ' + node);
	    }
	},
	/**
	 * get an array containing the nodes that depend on the specified node (transitively).
	 *
	 * Throws an Error if the graph has a cycle, or the specified node does not exist.
	 *
	 * If `leavesOnly` is true, only nodes that do not have any dependants will be returned in the array.
	 */
	dependantsOf:function (node, leavesOnly) {
	    if (this.hasNode(node)) {
		var result = [];
		var DFS = createDFS(this.incomingEdges, leavesOnly, result);
		DFS(node);
		var idx = result.indexOf(node);
		if (idx >= 0) {
		    result.splice(idx, 1);
		}
		return result;
	    } else {
		throw new Error('Node does not exist: ' + node);
	    }
	},
	/**
	 * Construct the overall processing order for the dependency graph.
	 *
	 * Throws an Error if the graph has a cycle.
	 *
	 * If `leavesOnly` is true, only nodes that do not depend on any other nodes will be returned.
	 */
	overallOrder:function (leavesOnly) {
	    var self = this;
	    var result = [];
	    var keys = Object.keys(this.nodes);
	    if (keys.length === 0) {
		return result; // Empty graph
	    } else {
		// Look for cycles - we run the DFS starting at all the nodes in case there
		// are several disconnected subgraphs inside this dependency graph.
		var CycleDFS = createDFS(this.outgoingEdges, false, []);
		keys.forEach(function(n) {
		    CycleDFS(n);
		});

		var DFS = createDFS(this.outgoingEdges, leavesOnly, result);
		// Find all potential starting points (nodes with nothing depending on them) an
		// run a DFS starting at these points to get the order
		keys.filter(function (node) {
		    return self.incomingEdges[node].length === 0;
		}).forEach(function (n) {
		    DFS(n);
		});

		return result;
	    }
	},

    };

    // export
    if (typeof exports !== 'undefined') {
	exports.DepGraph = DepGraph;
    }
})(typeof window === 'undefined' ? this : window);
