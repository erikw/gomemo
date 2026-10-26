package notes

// Clone returns an independent copy of n.
func (n *Note) Clone() *Note {
	clone := *n
	return &clone
}
