package multisplay

type keyType = int
type valueType = interface{}

type nodeColor int

const (
	red nodeColor = iota
	black
)

type multisplayNode struct {
	key    keyType
	value  valueType
	left   *multisplayNode
	right  *multisplayNode
	parent *multisplayNode

	deltaRefDepth int16
	deltaMinDepth int16

	isSplayRoot bool
	color       nodeColor
}

func newMultisplayNode(key keyType, value valueType) (node *multisplayNode) {
	node = new(multisplayNode)
	node.key = key
	node.value = value
	node.left = nil
	node.right = nil
	node.parent = nil
	node.isSplayRoot = true
	node.deltaRefDepth = 1
	node.deltaMinDepth = 0
	node.color = red
	return
}

func (node *multisplayNode) findWithParents(key keyType, x, z **multisplayNode, xDepth, zDepth *int16) *multisplayNode {
	current := node

	*x = nil
	*xDepth = 0

	*z = nil
	*zDepth = 0

	depth := int16(0)

	for current != nil && current.key != key {
		depth += current.deltaRefDepth

		if key < current.key {
			*z = current
			*zDepth = depth
			current = current.left
		} else {
			*x = current
			*xDepth = depth
			current = current.right
		}
	}

	return current
}

type TreeMap struct {
	root *multisplayNode
}

func (m *TreeMap) insertWithParents(created, x, z *multisplayNode, xDepth, zDepth int16) {
	var parent *multisplayNode

	if x != nil && xDepth > zDepth {
		parent = x
		m.multiSplay(x.key)

		if z != nil {
			z.left = created
			created.parent = z
			created.deltaRefDepth = xDepth + 1 - zDepth
		} else {
			x.right = created
			created.parent = x
		}
	} else {
		parent = z
		m.multiSplay(z.key)

		if x != nil {
			x.left = created
			created.parent = x
			created.deltaRefDepth = zDepth + 1 - xDepth
		} else {
			z.left = created
			created.parent = z
		}
	}

	m.virtualRebalance(parent)
}

func New() (m *TreeMap) {
	m = new(TreeMap)
	m.root = nil
	return
}

func (m *TreeMap) Find(key keyType) *valueType {
	found := m.multiSplay(key)

	if found == nil {
		return nil
	}

	return &found.value
}

func (m *TreeMap) Insert(key keyType, value valueType) {
	if m.root == nil {
		m.root = newMultisplayNode(key, value)
		m.root.color = black
		return
	}

	var x, z *multisplayNode
	var xDepth, zDepth int16

	found := m.root.findWithParents(key, &x, &z, &xDepth, &zDepth)

	if found != nil {
		found.value = value
		m.multiSplay(found.key)
	} else {
		created := newMultisplayNode(key, value)
		m.insertWithParents(created, x, z, xDepth, zDepth)
	}
}

func (m *TreeMap) Delete(key keyType) {
	// TODO
}
