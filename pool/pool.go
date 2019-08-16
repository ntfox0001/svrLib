package pool

import (
	"errors"
)

type Factory func() (Resource, error)

type Resource interface {
	Close()
}

type Pool struct {
	resources chan Resource
	factory   Factory
}

func NewPool(factory Factory, capacity int) *Pool {
	if capacity <= 0 {
		panic(errors.New("invalid/out of range capacity"))
	}
	rp := &Pool{
		resources: make(chan Resource, capacity),
		factory:   factory,
	}
	return rp
}

func (rp *Pool) Get() (resource Resource, err error) {
	var res Resource
	ok := true
	select {
	case res, ok = <-rp.resources:
		return res, nil
	default:
	}
	if !ok {
		return nil, errors.New("ChanClosed")
	}

	res, err = rp.factory()
	if err != nil {
		rp.resources <- res
	}

	return res, err
}

func (rp *Pool) Put(res Resource) {

	select {
	case rp.resources <- res:
	default:
		res.Close()
		return
	}
}
