package payload

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/glog"

	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

// SplitOnJobs enforces the rule that any Job in the payload prevents reordering or parallelism (either before or after)
func SplitOnJobs(task *Task) bool {
	return task.Manifest.GVK == schema.GroupVersionKind{Kind: "Job", Version: "v1", Group: "batch"}
}

var reMatchPattern = regexp.MustCompile(`^0000_(\d+)_([a-zA-Z0-9]+(-[a-zA-Z0-9]+)*?)_`)

const (
	groupNumber    = 1
	groupComponent = 2
)

// ByNumberAndComponent creates parallelization for tasks whose original filenames are of the form
// 0000_NN_NAME_* - files that share 0000_NN_NAME_ are run in serial, but chunks of files that have
// the same 0000_NN but different NAME can be run in parallel. If the input is not sorted in an order
// such that 0000_NN_NAME elements are next to each other, the splitter will treat those as unsplittable
// elements.
func ByNumberAndComponent(tasks []*Task) [][]*TaskNode {
	if len(tasks) <= 1 {
		return nil
	}
	count := len(tasks)
	matches := make([][]string, 0, count)
	for i := 0; i < len(tasks); i++ {
		matches = append(matches, reMatchPattern.FindStringSubmatch(tasks[i].Manifest.OriginalFilename))
	}

	var buckets [][]*TaskNode
	var lastNode *TaskNode
	for i := 0; i < count; {
		matchBase := matches[i]
		j := i + 1
		var groups []*TaskNode
		for ; j < count; j++ {
			matchNext := matches[j]
			if matchBase == nil || matchNext == nil || matchBase[groupNumber] != matchNext[groupNumber] {
				break
			}
			if matchBase[groupComponent] != matchNext[groupComponent] {
				groups = append(groups, &TaskNode{Tasks: tasks[i:j]})
				i = j
			}
			matchBase = matchNext
		}
		if len(groups) > 0 {
			groups = append(groups, &TaskNode{Tasks: tasks[i:j]})
			i = j
			buckets = append(buckets, groups)
			lastNode = nil
			continue
		}
		if lastNode == nil {
			lastNode = &TaskNode{Tasks: append([]*Task(nil), tasks[i:j]...)}
			i = j
			buckets = append(buckets, []*TaskNode{lastNode})
			continue
		}
		lastNode.Tasks = append(lastNode.Tasks, tasks[i:j]...)
		i = j
	}
	return buckets
}

type TaskNode struct {
	In    []int
	Tasks []*Task
	Out   []int
}

func (n TaskNode) String() string {
	var arr []string
	for _, t := range n.Tasks {
		if len(t.Manifest.OriginalFilename) > 0 {
			arr = append(arr, t.Manifest.OriginalFilename)
			continue
		}
		arr = append(arr, t.Manifest.GVK.String())
	}
	return "{Tasks: " + strings.Join(arr, ", ") + "}"
}

func (n *TaskNode) replaceIn(index, with int) {
	for i, from := range n.In {
		if from == index {
			n.In[i] = with
		}
	}
}

func (n *TaskNode) replaceOut(index, with int) {
	for i, to := range n.Out {
		if to == index {
			n.Out[i] = with
		}
	}
}

func (n *TaskNode) appendOut(items ...int) {
	for _, in := range items {
		if !containsInt(n.Out, in) {
			n.Out = append(n.Out, in)
		}
	}
}

// TaskGraph provides methods for parallelizing a linear sequence
// of Tasks based on Split or Parallelize functions.
type TaskGraph struct {
	Nodes []*TaskNode
}

// NewTaskGraph creates a graph with a single node containing
// the supplied tasks.
func NewTaskGraph(tasks []*Task) *TaskGraph {
	return &TaskGraph{
		Nodes: []*TaskNode{
			{
				Tasks: tasks,
			},
		},
	}
}

func containsInt(arr []int, value int) bool {
	for _, i := range arr {
		if i == value {
			return true
		}
	}
	return false
}

func (g *TaskGraph) replaceInOf(index, with int) {
	node := g.Nodes[index]
	in := node.In
	for _, pos := range in {
		g.Nodes[pos].replaceOut(index, with)
	}
}

func (g *TaskGraph) replaceOutOf(index, with int) {
	node := g.Nodes[index]
	out := node.Out
	for _, pos := range out {
		g.Nodes[pos].replaceIn(index, with)
	}
}

// Split breaks a graph node with a task that onFn returns true into
// one, two, or three separate nodes, preserving the order of tasks.
// E.g. a node with [a,b,c,d] where onFn returns true of b will result
// in a graph with [a] -> [b] -> [c,d].
func (g *TaskGraph) Split(onFn func(task *Task) bool) {
	for i := 0; i < len(g.Nodes); i++ {
		node := g.Nodes[i]
		tasks := node.Tasks
		if len(tasks) <= 1 {
			continue
		}
		for j, task := range tasks {
			if !onFn(task) {
				continue
			}

			if j > 0 {
				left := tasks[0:j]
				next := len(g.Nodes)
				nextNode := &TaskNode{
					In:    node.In,
					Tasks: left,
					Out:   []int{i},
				}
				g.Nodes = append(g.Nodes, nextNode)
				g.replaceInOf(i, next)
				node.In = []int{next}
			}

			if j < (len(tasks) - 1) {
				right := tasks[j+1:]
				next := len(g.Nodes)
				nextNode := &TaskNode{
					In:    []int{i},
					Tasks: right,
					Out:   node.Out,
				}
				g.Nodes = append(g.Nodes, nextNode)
				g.replaceOutOf(i, next)
				node.Out = []int{next}
			}

			node.Tasks = tasks[j : j+1]
			break
		}
	}
}

// Parallelize takes the given breakFn and splits any TaskNode's tasks up
// into parallel groups. If breakFn returns an empty array or a single
// array item with a single task node, that is considered a no-op.
func (g *TaskGraph) Parallelize(breakFn func([]*Task) [][]*TaskNode) {
	for i := 0; i < len(g.Nodes); i++ {
		node := g.Nodes[i]
		results := breakFn(node.Tasks)
		if len(results) == 0 || (len(results) == 1 && len(results[0]) == 1) {
			continue
		}
		node.Tasks = nil
		out := node.Out
		node.Out = nil

		// starting with the left anchor, create chains of nodes,
		// and avoid M x N in/out connections by creating spacers
		in := []int{i}
		for _, inNodes := range results {
			if len(inNodes) == 0 {
				continue
			}
			singleIn, singleOut := len(in) == 1, len(inNodes) == 1

			switch {
			case singleIn && singleOut, singleIn, singleOut:
				in = g.bulkAdd(inNodes, in)
			default:
				in = g.bulkAdd([]*TaskNode{{}}, in)
				in = g.bulkAdd(inNodes, in)
			}
		}

		// make node the left anchor and nextNode the right anchor
		if len(out) > 0 {
			next := len(g.Nodes)
			nextNode := &TaskNode{
				Tasks: nil,
				Out:   out,
			}
			g.Nodes = append(g.Nodes, nextNode)
			for _, j := range out {
				g.Nodes[j].replaceIn(i, next)
			}
			for _, j := range in {
				g.Nodes[j].Out = []int{next}
				nextNode.In = append(nextNode.In, j)
			}
		}
	}
}

func (g *TaskGraph) Roots() []int {
	var roots []int
	for i, n := range g.Nodes {
		if len(n.In) > 0 {
			continue
		}
		roots = append(roots, i)
	}
	return roots
}

func (g *TaskGraph) Tree() string {
	roots := g.Roots()
	visited := make([]int, len(g.Nodes))
	stage := 0
	var out []string
	var depth []int
	for len(roots) > 0 {
		depth = append(depth, 0)
		for _, i := range roots {
			visited[i] = 1
			if d := len(g.Nodes[i].Tasks); d > depth[len(depth)-1] {
				depth[len(depth)-1] = d
			}
			out = append(out, fmt.Sprintf("%d: %d %s in=%v out=%v", stage, i, g.Nodes[i], g.Nodes[i].In, g.Nodes[i].Out))
		}
		roots = roots[0:0]
		for i, b := range visited {
			if b == 1 || !covers(visited, g.Nodes[i].In) {
				continue
			}
			roots = append(roots, i)
		}
		stage++
	}
	for i, b := range visited {
		if b == 1 {
			continue
		}
		out = append(out, fmt.Sprintf("unreachable: %d %s in=%v out=%v", i, g.Nodes[i], g.Nodes[i].In, g.Nodes[i].Out))
	}
	var totalDepth int
	var levels []string
	for _, d := range depth {
		levels = append(levels, strconv.Itoa(d))
		totalDepth += d
	}
	out = append(out, fmt.Sprintf("summary: depth=%d, levels=%s", totalDepth, strings.Join(levels, ",")))
	return strings.Join(out, "\n")
}

func covers(all []int, some []int) bool {
	for _, i := range some {
		if all[i] == 0 {
			return false
		}
	}
	return true
}

func (g *TaskGraph) bulkAdd(nodes []*TaskNode, inNodes []int) []int {
	from := len(g.Nodes)
	g.Nodes = append(g.Nodes, nodes...)
	to := len(g.Nodes)
	if len(inNodes) == 0 {
		toNodes := make([]int, to-from)
		for k := from; k < to; k++ {
			toNodes[k-from] = k
		}
		return toNodes
	}

	next := make([]int, to-from)
	for k := from; k < to; k++ {
		g.Nodes[k].In = append([]int(nil), inNodes...)
		next[k-from] = k
	}
	for _, k := range inNodes {
		g.Nodes[k].appendOut(next...)
	}
	return next
}

type runTasks struct {
	index int
	tasks []*Task
}

type taskStatus struct {
	index   int
	success bool
}

func RunGraph(ctx context.Context, graph *TaskGraph, maxParallelism int, fn func(ctx context.Context, tasks []*Task) error) error {
	nestedCtx, cancelFn := context.WithCancel(ctx)
	defer cancelFn()

	// This goroutine takes nodes from the graph as they are available (their prereq has completed) and
	// sends them to workCh. It uses completeCh to know that a previously dispatched item is complete.
	completeCh := make(chan taskStatus, maxParallelism)
	defer close(completeCh)

	workCh := make(chan runTasks, maxParallelism)
	go func() {
		defer close(workCh)

		// visited tracks nodes we have not sent (0), are currently
		// waiting for completion (1), or have completed (2,3)
		const (
			nodeNotVisited int = iota
			nodeWorking
			nodeFailed
			nodeComplete
		)
		visited := make([]int, len(graph.Nodes))
		canVisit := func(node *TaskNode) bool {
			for _, previous := range node.In {
				switch visited[previous] {
				case nodeFailed, nodeWorking, nodeNotVisited:
					return false
				}
			}
			return true
		}

		remaining := len(graph.Nodes)
		var inflight int
		for {
			found := 0

			// walk the graph, filling the work queue
			for i := 0; i < len(visited); i++ {
				if visited[i] != nodeNotVisited {
					continue
				}
				if canVisit(graph.Nodes[i]) {
					select {
					case workCh <- runTasks{index: i, tasks: graph.Nodes[i].Tasks}:
						visited[i] = nodeWorking
						found++
						inflight++
					default:
						break
					}
				}
			}

			// try to empty the done channel
			for len(completeCh) > 0 {
				finished := <-completeCh
				if finished.success {
					visited[finished.index] = nodeComplete
				} else {
					visited[finished.index] = nodeFailed
				}
				remaining--
				inflight--
				found++
			}

			if found > 0 {
				continue
			}

			// no more work to hand out
			if remaining == 0 {
				glog.V(4).Infof("Graph is complete")
				return
			}

			// we walked the entire graph, there are still nodes remaining, but we're not waiting
			// for anything
			if inflight == 0 && found == 0 {
				glog.V(4).Infof("No more reachable nodes in graph, continue")
				break
			}

			// we did nothing this round, so we have to wait for more
			finished, ok := <-completeCh
			if !ok {
				// we've been aborted
				glog.V(4).Infof("Stopped graph walker due to cancel")
				return
			}
			if finished.success {
				visited[finished.index] = nodeComplete
			} else {
				visited[finished.index] = nodeFailed
			}
			remaining--
			inflight--
		}

		// take everything remaining and process in order
		var unreachable []*Task
		for i := 0; i < len(visited); i++ {
			if visited[i] == nodeNotVisited && canVisit(graph.Nodes[i]) {
				unreachable = append(unreachable, graph.Nodes[i].Tasks...)
			}
		}
		if len(unreachable) > 0 {
			sort.Slice(unreachable, func(i, j int) bool {
				a, b := unreachable[i], unreachable[j]
				return a.Index < b.Index
			})
			workCh <- runTasks{index: -1, tasks: unreachable}
			glog.V(4).Infof("Waiting for last tasks")
			<-completeCh
		}
		glog.V(4).Infof("No more work")
	}()

	errCh := make(chan error, maxParallelism)
	wg := sync.WaitGroup{}
	if maxParallelism < 1 {
		maxParallelism = 1
	}
	for i := 0; i < maxParallelism; i++ {
		wg.Add(1)
		go func(job int) {
			defer utilruntime.HandleCrash()
			defer wg.Done()
			for {
				select {
				case <-nestedCtx.Done():
					glog.V(4).Infof("Canceled worker %d", job)
					return
				case runTask, ok := <-workCh:
					if !ok {
						glog.V(4).Infof("No more work for %d", job)
						return
					}
					glog.V(4).Infof("Running %d on %d", runTask.index, job)
					err := fn(nestedCtx, runTask.tasks)
					completeCh <- taskStatus{index: runTask.index, success: err == nil}
					if err != nil {
						errCh <- err
					}
				}
			}
		}(i)
	}
	go func() {
		glog.V(4).Infof("Waiting for workers to complete")
		wg.Wait()
		glog.V(4).Infof("Workers finished")
		close(errCh)
	}()

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}
	glog.V(4).Infof("Result of work: %v", errs)
	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}
