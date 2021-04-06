package multisplay

func (node *multisplayNode) maintainMinDepth() {
	max := int16(0)

	if node.left != nil && !node.left.isSplayRoot {
		d := node.left.deltaMinDepth - node.left.deltaRefDepth
		if d > max {
			max = d
		}
	}

	if node.right != nil && !node.right.isSplayRoot {
		d := node.left.deltaMinDepth - node.left.deltaRefDepth
		if d > max {
			max = d
		}
	}

	node.deltaMinDepth = max
}

func maintainAuxiliary(p, q, y *multisplayNode) {
	p.isSplayRoot, q.isSplayRoot = q.isSplayRoot, p.isSplayRoot

	deltaQ := q.deltaRefDepth
	q.deltaRefDepth += p.deltaRefDepth
	p.deltaRefDepth = -deltaQ

	if y != nil {
		y.deltaRefDepth += deltaQ
	}

	p.maintainMinDepth()
	q.maintainMinDepth()
}

func maintainParents(o, p, q, y *multisplayNode) {
	if o != nil {
		if o.right == p {
			o.right = q
		} else {
			o.left = q
		}
	}

	if y != nil {
		y.parent = p
	}

	p.parent = q
	q.parent = o
}

// Rotates a node (P) to the left, as shown on the following diagram:
//
//       o              o
//       |              |
//       P              Q
//      / \            / \
//     x   Q    =>    P   z
//        / \        / \
//       y   z      x   y
//
// Upper case nodes are required to be non-null for the operation to work.
// Panics otherwise.

func (p *multisplayNode) rotateLeft() {
	o := p.parent
	q := p.right
	y := q.left

	q.left = p
	p.right = y

	maintainParents(o, p, q, y)
	maintainAuxiliary(p, q, y)
}

// Rotates a node (P) to the right, as shown on the following diagram:
//
//       o              o
//       |              |
//       P              Q
//      / \            / \
//     Q   z    =>    x   P
//    / \                / \
//   x   y              y   z
//
// Upper case nodes are required to be non-null for the operation to work.
// Panics otherwise.

func (p *multisplayNode) rotateRight() {
	o := p.parent
	q := p.left
	y := q.right

	q.right = p
	p.left = y

	maintainParents(o, p, q, y)
	maintainAuxiliary(p, q, y)
}

// Splays a node until it has a given parent.
//
// If root != nil, assumes the given node is on a subtree of that node.
// Behavior is undefined otherwise.
//
// If root is nil, this will take the given node up to the root of the
// corresponding splay tree.

func (node *multisplayNode) splay(root *multisplayNode) {
	for node.parent != root && !node.isSplayRoot {
		parent := node.parent
		grandparent := parent.parent

		if parent.right == node {
			if grandparent == root {
				parent.rotateLeft()
				return
			}

			if grandparent.right == parent {
				// Zag-Zag (node = R)
				//
				//    P                      R
				//   / \                    / \
				//  x   Q                  Q   w
				//     / \       =>       / \
				//    y   R              P   z
				//       / \            / \
				//      z   w          x   y
				grandparent.rotateLeft()
				parent.rotateLeft()
			} else {
				// Zag-Zig (node = P)
				//
				//      R
				//     / \               P
				//    Q   w            /   \
				//   / \       =>     Q     R
				//  x   P            / \   / \
				//     / \          x   y z   w
				//    y   z
				parent.rotateLeft()
				grandparent.rotateRight()
			}
		} else {
			if grandparent == root {
				parent.rotateRight()
				return
			}

			if grandparent.right == parent {
				// Zig-Zag (symmetric to Zag-Zig)
				parent.rotateRight()
				grandparent.rotateLeft()
			} else {
				// Zig-Zig (symmetric to Zag-Zag)
				grandparent.rotateRight()
				parent.rotateRight()
			}
		}
	}
}

func (m *TreeMap) maintainRoot(candidate *multisplayNode) {
	if candidate.parent == nil {
		m.root = candidate
	}
}
