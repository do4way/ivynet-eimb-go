package promises

//Promise ...
type Promise struct {
	data       interface{}
	isResolved bool
	err        error
	onSucc     func(interface{}) (interface{}, error)
	onErr      func(error) error
	next       *Promise
}

var endPromiseChain = &Promise{
	onSucc: func(data interface{}) (interface{}, error) { return nil, nil },
	onErr:  func(err error) error { return nil },
}

//New ...
func New() *Promise {
	return &Promise{next: endPromiseChain}
}

//Then ...
func (p *Promise) Then(onSucc func(interface{}) (interface{}, error), onErr func(error) error) *Promise {

	if p.onSucc != nil || p.onErr != nil {
		panic("Future handlers have been added. ")
	}
	p.next = New()
	defer func() {
		p.onSucc = onSucc
		p.onErr = onErr
		if p.isResolved {
			if p.data != nil {
				p.unsafeResolve(p.data)
			}
			if p.err != nil {
				p.unsafeReject(p.err)
			}
		}
	}()
	return p.next
}

//Resolve ...
func (p *Promise) Resolve(data interface{}) {

	if p.isResolved {
		panic("Future have been resolved.")
	}

	go func() {
		p.unsafeResolve(data)
	}()

}

func (p *Promise) unsafeResolve(data interface{}) {

	p.data = data
	p.isResolved = true

	if p.onSucc == nil {
		return
	}

	d, err := p.onSucc(p.data)
	if err != nil {
		p.err = err
		p.unsafeReject(err)
		return
	}

	p.next.Resolve(d)
}

//Reject ..
func (p *Promise) Reject(err error) {

	if p.isResolved {
		panic("Future have been resolved.")
	}

	go func() {
		p.unsafeReject(err)
	}()
}

func (p *Promise) unsafeReject(err error) {

	p.err = err
	p.isResolved = true
	if p.onErr != nil {
		p.next.Reject(p.onErr(p.err))
	}
}
