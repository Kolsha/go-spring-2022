//go:build !solution

package tparallel

type T struct {
	finished   chan struct{}
	barrier    chan struct{}
	isParallel bool
	parent     *T
	sub        []*T
}

func newT(parent *T) *T {
	return &T{
		finished: make(chan struct{}),
		barrier:  make(chan struct{}),
		parent:   parent,
		sub:      make([]*T, 0),
	}
}

func (t *T) Parallel() {
	if t.isParallel {
		panic("test is already parallel")
	}
	t.isParallel = true
	t.parent.sub = append(t.parent.sub, t)

	t.finished <- struct{}{}
	<-t.parent.barrier
}

func (t *T) tRunner(subtest func(t *T)) {
	subtest(t)
	if len(t.sub) > 0 {
		close(t.barrier)

		for _, sub := range t.sub {
			<-sub.finished
		}
	}
	if t.isParallel {
		t.parent.finished <- struct{}{}
	}
	t.finished <- struct{}{}

}

func (t *T) Run(subtest func(t *T)) {
	subT := newT(t)
	go subT.tRunner(subtest)
	<-subT.finished
}

func Run(topTests []func(t *T)) {
	root := newT(nil)
	for _, fn := range topTests {
		root.Run(fn)
	}
	close(root.barrier)

	if len(root.sub) > 0 {
		<-root.finished
	}
}
