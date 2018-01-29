package eimb

import (
	. "gopkg.in/check.v1"
)

type HandlerWorkersTestSuite struct {
}

var (
	_ = Suite(&HandlerWorkersTestSuite{})
)

func (s *HandlerWorkersTestSuite) TestWorkers(c *C) {
	rsts := make(chan int, 3)
	workers := NewWorkerPool(3)
	workers.Start()
	t := 0
	for i := 1; i <= 5; i++ {
		var x = i
		job := WithJob(func() {
			rsts <- x
		})
		JobQueue <- job
	}
	c.Assert(t, Equals, 0)

	for j := 1; j <= 5; j++ {
		t += <-rsts
	}
	c.Assert(t, Equals, 15)
}
