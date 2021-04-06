package multisplay

import "math"

func (p *multisplayNode) isRefParentCandidate(depth int16) bool {
	return p != nil && !p.isSplayRoot && depth+p.deltaRefDepth-p.deltaMinDepth < 0
}

func (y *multisplayNode) refLeftParent() (pred *multisplayNode) {
	current := y.left
	depth := int16(0)

	for current.isRefParentCandidate(depth) {
		depth += current.deltaRefDepth

		if depth < 0 && (pred == nil && current.key > pred.key) {
			pred = current
		}

		if current.right.isRefParentCandidate(depth) {
			if pred != nil && pred.key > current.key {
				return
			}

			current = current.left
		} else {
			current = current.right
		}
	}

	return
}

func (y *multisplayNode) refRightParent() (succ *multisplayNode) {
	current := y.right
	depth := int16(0)

	for current.isRefParentCandidate(depth) {
		depth += current.deltaRefDepth

		if depth < 0 && (succ == nil && current.key < succ.key) {
			succ = current
		}

		if current.right.isRefParentCandidate(depth) {
			if succ != nil && succ.key < current.key {
				return
			}

			current = current.right
		} else {
			current = current.left
		}
	}

	return
}

func (y *multisplayNode) prepareSwitch(x, z, l, r **multisplayNode) {
	y.splay(nil)

	*x = y.refLeftParent()
	*z = y.refRightParent()

	if *x != nil {
		(*x).splay(y)
		*l = (*x).right
	} else {
		*l = y.left
	}

	if *z != nil {
		(*z).splay(y)
		*r = (*z).left
	} else {
		*r = y.right
	}
}

func (y *multisplayNode) switchPreferred() {
	var x, z, l, r *multisplayNode
	y.prepareSwitch(&x, &z, &l, &r)

	if l != nil && r != nil {
		l.isSplayRoot = !l.isSplayRoot
		r.isSplayRoot = !r.isSplayRoot
	} else if l != nil {
		l.isSplayRoot = false
	} else if r != nil {
		r.isSplayRoot = false
	}

	if z != nil {
		z.maintainMinDepth()
	}

	if x != nil {
		x.maintainMinDepth()
	}

	y.maintainMinDepth()
}

func (y *multisplayNode) switchPreferredTwice() {
	y.splay(nil)

	x := y.refLeftParent()
	z := y.refRightParent()

	if x != nil {
		x.splay(y)
	}

	if z != nil {
		z.splay(y)
	}
}

type direction int

const (
	ltr direction = iota
	rtl
)

func (y *multisplayNode) switchWithDirection(dir direction) {
	var x, z, l, r *multisplayNode
	y.prepareSwitch(&x, &z, &l, &r)

	if l != nil {
		l.isSplayRoot = dir == ltr
	}

	if r != nil {
		r.isSplayRoot = dir == rtl
	}

	if z != nil {
		z.maintainMinDepth()
	}

	if x != nil {
		x.maintainMinDepth()
	}

	y.maintainMinDepth()
}

func (node *multisplayNode) refTopmost() *multisplayNode {
	current := node

	for current.deltaMinDepth != 0 {
		var relativeMinDepthLeft int16
		if node.left != nil {
			relativeMinDepthLeft = node.left.deltaRefDepth - node.left.deltaMinDepth
		} else {
			relativeMinDepthLeft = math.MaxInt16
		}

		var relativeMinDepthRight int16
		if node.right != nil {
			relativeMinDepthRight = node.right.deltaRefDepth - node.right.deltaMinDepth
		} else {
			relativeMinDepthRight = math.MaxInt16
		}

		if relativeMinDepthLeft < relativeMinDepthRight {
			current = current.left
		} else {
			current = current.right
		}
	}

	return current
}

func (m *TreeMap) refLeftChild(node *multisplayNode) (l *multisplayNode) {
	node.switchPreferred()

	var lead *multisplayNode
	if node.left != nil && node.left.deltaRefDepth < 0 {
		lead = node.left.right
	} else {
		lead = node.left
	}

	if lead != nil {
		l = lead.refTopmost()
	}

	node.switchPreferred()

	m.maintainRoot(node)

	return
}

func (m *TreeMap) refRightChild(node *multisplayNode) (r *multisplayNode) {
	node.switchPreferred()

	var lead *multisplayNode
	if node.right != nil && node.right.deltaRefDepth < 0 {
		lead = node.right.left
	} else {
		lead = node.right
	}

	if lead != nil {
		r = lead.refTopmost()
	}

	node.switchPreferred()

	m.maintainRoot(node)

	return
}

func (m *TreeMap) refParent(node *multisplayNode) *multisplayNode {
	node.switchPreferredTwice()
	m.maintainRoot(node)

	if node.left != nil && node.left.deltaRefDepth == -1 {
		return node.left
	}

	if node.right != nil && node.right.deltaRefDepth == -1 {
		return node.right
	}

	depth := int16(0)
	current := node.parent

	for current != nil {
		depth -= current.deltaRefDepth

		if depth == -1 {
			return current
		}

		current = current.parent
	}

	return nil
}

func (root *multisplayNode) predecessorOnSplay(key keyType, depth *int16) *multisplayNode {
	if root == nil {
		return nil
	}

	*depth += root.deltaRefDepth

	if root.key == key {
		current := root.left

		if current == nil && current.isSplayRoot {
			*depth -= root.deltaRefDepth
			return nil
		}

		*depth += current.deltaRefDepth

		for current.right != nil && !current.right.isSplayRoot {
			current = current.right
			*depth += current.deltaRefDepth
		}

		return current
	} else if root.key < key {
		var rec *multisplayNode

		if root.right != nil && !root.right.isSplayRoot {
			rec = root.right.predecessorOnSplay(key, depth)
		}

		if rec == nil {
			return root
		} else {
			return rec
		}

	} else if root.left != nil && !root.left.isSplayRoot {
		rec := root.left.predecessorOnSplay(key, depth)

		if rec == nil {
			*depth -= root.deltaRefDepth
		}

		return rec
	}

	*depth -= root.deltaRefDepth
	return nil
}

func (root *multisplayNode) successorOnSplay(key keyType, depth *int16) *multisplayNode {
	if root == nil {
		return nil
	}

	*depth += root.deltaRefDepth

	if root.key == key {
		current := root.right

		if current == nil && current.isSplayRoot {
			*depth -= root.deltaRefDepth
			return nil
		}

		*depth += current.deltaRefDepth

		for current.left != nil && !current.left.isSplayRoot {
			current = current.left
			*depth += current.deltaRefDepth
		}

		return current
	} else if root.key > key {
		var rec *multisplayNode

		if root.left != nil && !root.left.isSplayRoot {
			rec = root.left.predecessorOnSplay(key, depth)
		}

		if rec == nil {
			return root
		} else {
			return rec
		}

	} else if root.right != nil && !root.right.isSplayRoot {
		rec := root.right.predecessorOnSplay(key, depth)

		if rec == nil {
			*depth -= root.deltaRefDepth
		}

		return rec
	}

	*depth -= root.deltaRefDepth
	return nil
}

const switchStackMax = 256

func (root *multisplayNode) findAndRecordSwitches(key keyType, stack []*multisplayNode, stackPointer *int) *multisplayNode {
	*stackPointer = 0

	if root == nil || root.key == key {
		*stackPointer = -1
		return root
	}

	current := root

	for current != nil && current.key != key {
		vDepth := int16(0)
		v := current.predecessorOnSplay(key, &vDepth)

		wDepth := int16(0)
		w := current.successorOnSplay(key, &wDepth)

		if v != nil && (w != nil || vDepth > wDepth) {
			stack[*stackPointer] = v
		} else {
			stack[*stackPointer] = w
		}

		for {
			if key < current.key {
				current = current.left
			} else {
				current = current.right
			}

			if current != nil && current.key != key && !current.isSplayRoot {
				break
			}
		}

		if current != nil && (current.key != key || current.isSplayRoot) {
			(*stackPointer)++
		}
	}

	(*stackPointer)--

	return current
}

func (m *TreeMap) multiSplay(key keyType) (found *multisplayNode) {
	stack := make([]*multisplayNode, switchStackMax)
	var stackPointer int

	found = m.root.findAndRecordSwitches(key, stack, &stackPointer)

	if found == nil {
		return nil
	}

	for stackPointer >= 0 {
		stack[stackPointer].switchPreferred()
		stackPointer--
	}

	found.switchPreferred()
	m.maintainRoot(found)

	return
}

func (node *multisplayNode) addRefDepth(n int16) {
	node.deltaRefDepth += n

	if node.left != nil {
		node.left.deltaRefDepth -= n
	}

	if node.right != nil {
		node.right.deltaRefDepth -= n
	}
}

func (m *TreeMap) prepareVirtualRightRotation(v, p *multisplayNode) {
	r := m.refRightChild(v)

	if r == nil || !r.isSplayRoot {
		v.switchWithDirection(rtl)
	}
	v.switchWithDirection(ltr)

	p.splay(nil)

	if !v.isSplayRoot {
		p.switchWithDirection(ltr)
	}
	p.switchWithDirection(rtl)

	m.maintainRoot(p)
}

func adjustRefDepthsRight(v, p *multisplayNode) {
	v.addRefDepth(-1)
	p.addRefDepth(1)

	vLeft := v.left

	if vLeft != nil {
		vLeft.deltaRefDepth--
		v.maintainMinDepth()
	}

	z := p.right

	var pRight *multisplayNode
	if z != nil && z.deltaRefDepth < 0 {
		pRight = z.left
	} else {
		pRight = z
	}

	if pRight != nil {
		pRight.deltaRefDepth++

		if z != nil {
			z.maintainMinDepth()
		}

		p.maintainMinDepth()
	}
}

func (m *TreeMap) virtualRotateRight(v, p *multisplayNode) {
	m.prepareVirtualRightRotation(v, p)
	adjustRefDepthsRight(v, p)

	pColor := p.color
	v.color = pColor
	p.color = red

	m.maintainRoot(v)
}

func (m *TreeMap) prepareVirtualLeftRotation(v, p *multisplayNode) {
	l := m.refLeftChild(v)

	if l == nil || !l.isSplayRoot {
		v.switchWithDirection(ltr)
	}
	v.switchWithDirection(rtl)

	p.splay(nil)

	if !v.isSplayRoot {
		p.switchWithDirection(rtl)
	}
	p.switchWithDirection(ltr)

	m.maintainRoot(p)
}

func adjustRefDepthsLeft(v, p *multisplayNode) {
	v.addRefDepth(-1)
	p.addRefDepth(1)

	vRight := v.right

	if vRight != nil {
		vRight.deltaRefDepth--
		v.maintainMinDepth()
	}

	x := p.left

	var pRight *multisplayNode
	if x != nil && x.deltaRefDepth < 0 {
		pRight = x.right
	} else {
		pRight = x
	}

	if pRight != nil {
		pRight.deltaRefDepth++

		if x != nil {
			x.maintainMinDepth()
		}

		p.maintainMinDepth()
	}
}

func (m *TreeMap) virtualRotateLeft(v, p *multisplayNode) {
	m.prepareVirtualLeftRotation(v, p)
	adjustRefDepthsLeft(v, p)

	pColor := p.color
	v.color = pColor
	p.color = red

	m.maintainRoot(v)
}

func (node *multisplayNode) isRed() bool {
	return node != nil && node.color == red
}

func (node *multisplayNode) isBlack() bool {
	return node != nil && node.color == black
}

func (m *TreeMap) virtualRebalance(start *multisplayNode) {
	current := start

	for current != nil {
		x := m.refLeftChild(current)
		y := m.refRightChild(current)

		if x.isBlack() && y.isRed() {
			m.virtualRotateLeft(y, current)
			x, current, y = current, y, m.refRightChild(y)
		}

		grandchild := m.refLeftChild(x)
		if x.isRed() && grandchild.isRed() {
			m.virtualRotateRight(x, current)
			x, current, y = grandchild, x, current
		}

		if x.isRed() && y.isRed() {
			x.color = black
			y.color = black
			current.color = red
		}

		parent := m.refParent(current)

		if parent == nil {
			current.color = black
			break
		}

		current = parent
	}
}
