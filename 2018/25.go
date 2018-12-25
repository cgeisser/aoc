package main

import (
	"bufio"
	"fmt"
	"os"
)

type star struct {
	x, y, z, t int
}

func (s star) dist(s2 star) int {
	xd := s.x - s2.x
	yd := s.y - s2.y
	zd := s.z - s2.z
	td := s.t - s2.t

	if xd < 0 {
		xd = -xd
	}
	if yd < 0 {
		yd = -yd
	}
	if zd < 0 {
		zd = -zd
	}
	if td < 0 {
		td = -td
	}

	return xd + yd + zd + td
}

type constel struct {
	stars []star
	r     int
}

func main() {
	f, _ := os.Open(os.Args[1])
	input := bufio.NewScanner(f)

	allstars := make([]star, 0)
	for input.Scan() {
		newstar := star{}
		_, ok := fmt.Sscanf(input.Text(), "%d,%d,%d,%d",
			&newstar.x,
			&newstar.y,
			&newstar.z,
			&newstar.t)
		if ok != nil {
			panic(fmt.Sprintf("bad line %v %v", ok, input.Text()))
		}
		allstars = append(allstars, newstar)
	}
	result := clusterAll(allstars)
	fmt.Println("built clusters", len(result))
}

type constels []*constel

func (c *constel) cluster(s star) bool {
	for _, x := range c.stars {
		//fmt.Println("clustering: ", x, s, s.dist(x))
		if s.dist(x) <= c.r {
			return true
		}
	}
	return false
}

func clusterAll(allstars []star) constels {
	clusters := make(constels, 0)
	for _, s := range allstars {
		merge := make([]int, 0)
		for i, c := range clusters {
			if c.cluster(s) {
				merge = append(merge, i)
			}
		}
		if len(merge) == 0 {
			clusters = append(clusters, &constel{[]star{s}, 3})
		} else {
			parent := clusters[merge[0]]
			parent.stars = append(parent.stars, s)
			for _, i := range merge[1:] {
				parent.mergeStars(clusters[i])
			}
			nc := make(constels, 0, len(clusters))
			for _, c := range clusters {
				if c.stars != nil {
					nc = append(nc, c)
				}
			}
			clusters = nc
		}
	}
	return clusters
}

func deleteAt(c constels, i int) constels {
	c[i] = c[len(c)-1]
	return c[:len(c)-1]
}

func (c *constel) mergeStars(m *constel) {
	c.stars = append(c.stars, m.stars...)
	m.stars = nil
}
