//go:build !solution

package treeiter

func DoInOrder[E interface {
	Left() *E
	Right() *E
}](root *E, cb func(t *E)) {
	if root == nil {
		return
	}
	DoInOrder((*root).Left(), cb)
	cb(root)
	DoInOrder((*root).Right(), cb)
}
