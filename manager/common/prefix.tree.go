package common

type PrefixTree[T any] struct {
	children map[rune]*PrefixTree[T]
	payload  *T
}

func CreatePrefixTree[T any]() *PrefixTree[T] {
	return &PrefixTree[T]{}
}

func (t *PrefixTree[T]) Add(str string, value *T) *PrefixTree[T] {
	node := t
	for _, r := range str {
		node = node.getOrCreate(r)
	}

	node.payload = value

	return t
}

func (t *PrefixTree[T]) GetLastPayload(str string) *T {
	node := t
	var payload *T
	for _, r := range str {
		ok := false
		if node, ok = node.children[r]; !ok {
			break
		}

		if node.payload != nil {
			payload = node.payload
		}
	}

	return payload
}

func (t *PrefixTree[T]) getOrCreate(r rune) *PrefixTree[T] {
	if node, ok := t.children[r]; ok {
		return node
	}

	node := CreatePrefixTree[T]()
	t.children[r] = node

	return node
}
