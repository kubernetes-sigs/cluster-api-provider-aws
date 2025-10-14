package certchain

import (
	"crypto/x509"
)

func Sort(certs []*x509.Certificate) []*x509.Certificate {
	graph := map[string][]string{}
	lookup := map[string]certAdapter{}

	for _, cert := range certs {
		node := certAdapter{cert}
		graph[node.ID()] = []string{}
		if node.Parent() != "" {
			graph[node.Parent()] = []string{}
		}
		lookup[node.ID()] = node
	}

	for _, node := range lookup {
		if node.Parent() == "" {
			continue
		}

		_, ok := lookup[node.Parent()]
		if !ok {
			graph[node.Parent()] = []string{}
		}

		if node.Parent() == node.ID() {
			// ignore self-signed
			continue
		}

		graph[node.Parent()] = append(graph[node.Parent()], node.ID())
	}

	ordered := topographicalSort(graph)

	var sorted []*x509.Certificate
	for _, i := range ordered {
		c, ok := lookup[i]
		if ok {
			sorted = append(sorted, c.Certificate)
		}
	}

	return sorted
}

type certAdapter struct {
	*x509.Certificate
}

func (c certAdapter) ID() string {
	return c.Subject.CommonName
}

func (c certAdapter) Parent() string {
	return c.Issuer.CommonName
}

func topographicalSort(g map[string][]string) []string {
	var ordered []string

	inDegree := map[string]int{}

	for n := range g {
		inDegree[n] = 0
	}

	for _, adjacent := range g {
		for _, v := range adjacent {
			inDegree[v]++
		}
	}

	var next []string
	for u, v := range inDegree {
		if v != 0 {
			continue
		}

		next = append(next, u)
	}

	for len(next) > 0 {
		u := next[0]
		next = next[1:]

		ordered = append(ordered, u)

		for _, v := range g[u] {
			inDegree[v]--

			if inDegree[v] == 0 {
				next = append(next, v)
			}
		}
	}

	return ordered
}
