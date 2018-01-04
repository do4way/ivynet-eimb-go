package promises

import (
	"errors"
	"fmt"
	"testing"

	. "gopkg.in/check.V1"
)

func TestPromise(t *testing.T) {
	TestingT(t)
}

type PromiseTestSuite struct {
}

var (
	_ = Suite(&PromiseTestSuite{})
)

func (test *PromiseTestSuite) TestResolveFirst(c *C) {
	p := New()
	p.Resolve(20)
	p.Then(func(d interface{}) (interface{}, error) {
		c.Check(d, Equals, 20)
		dt, _ := d.(int)
		return dt * 2, nil
	}, func(err error) error { return nil })
}

func (test *PromiseTestSuite) TestResolveChained(c *C) {
	p := New()
	p.Resolve(20)
	p.Then(func(d interface{}) (interface{}, error) {
		c.Check(d, Equals, 20)
		dt, _ := d.(int)
		return dt * 2, nil
	}, func(err error) error { return nil },
	).Then(func(d interface{}) (interface{}, error) {
		c.Check(d, Equals, 40)
		return nil, nil
	}, func(err error) error { return nil })
}

func (test *PromiseTestSuite) TestResolveThen(c *C) {
	p := New()
	cnt := 0
	p.Then(func(d interface{}) (interface{}, error) {
		c.Assert(cnt, Equals, 1)
		return d, nil
	}, func(err error) error { return nil })

	cnt++
	p.Resolve(20)
}

func (test *PromiseTestSuite) TestRejectFirst(c *C) {
	p := New()
	p.Reject(errors.New("A dummy error"))
	p.Then(func(d interface{}) (interface{}, error) {
		c.Fail()
		return nil, nil
	}, func(err error) error {
		c.Check(err.Error(), Equals, "A dummy error")
		return nil
	})
}

func (test *PromiseTestSuite) TestRejectChained(c *C) {
	p := New()
	p.Reject(errors.New("A dummy error"))
	p.Then(func(d interface{}) (interface{}, error) {
		c.Fail()
		return nil, nil
	}, func(err error) error {
		c.Check(err.Error(), Equals, "A dummy error")
		return errors.New("Second error")
	},
	).Then(func(d interface{}) (interface{}, error) {
		c.Fail()
		return nil, nil
	}, func(err error) error {
		c.Check(err.Error(), Equals, "Second error")
		return nil
	})
}

func (test *PromiseTestSuite) TestRejectThen(c *C) {
	p := New()
	p.Then(func(d interface{}) (interface{}, error) {
		c.Fail()
		return nil, nil
	}, func(err error) error {
		c.Check(err.Error(), Equals, "A dummy error")
		return nil
	})
	p.Reject(errors.New("A dummy error"))
}

func (test *PromiseTestSuite) TestErrorInResolver(c *C) {
	p := New()
	p.Then(func(d interface{}) (interface{}, error) {
		return nil, errors.New("Error in resolver")
	}, func(err error) error {
		c.Assert(err, Not(Equals), nil)
		fmt.Println(err)
		return err
	})
	p.Resolve(20)
}
