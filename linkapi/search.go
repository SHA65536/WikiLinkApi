package linkapi

type SearchHandler struct {
	DB *DatabaseHandler
}

func MakeSearchHandler(db *DatabaseHandler) *SearchHandler {
	return &SearchHandler{DB: db}
}

// ShortestPath finds the shortest path in the database between src and dst, execute hopEvent at every hop
func (s *SearchHandler) ShortestPath(src, dst uint32, hopEvent func(int)) ([]uint32, error) {
	var hops int
	var seen = map[uint32]struct{}{
		src: {},
	}
	var cur, next []*Node
	var final *Node
	// Breath First Search
	cur = []*Node{{Id: src}}
SearchLoop:
	for len(cur) > 0 {
		// Next layer
		next = []*Node{}
		for _, curNode := range cur {
			nxt, err := s.DB.GetOutgoing(curNode.Id)
			if err != nil {
				return nil, err
			}
			// For each link
			for _, nxtVal := range nxt {
				if _, ok := seen[nxtVal]; !ok { // If not seen
					seen[nxtVal] = struct{}{}
					if nxtVal == dst { // If found destination
						final = &Node{Id: nxtVal, Parent: curNode}
						break SearchLoop
					}
					next = append(next, &Node{Id: nxtVal, Parent: curNode})
				}
			}
		}
		cur = next
		hops++
		// Executing hop event
		hopEvent(hops)
	}
	return final.ToList(), nil
}

type Node struct {
	Id     uint32
	Parent *Node
}

// ToList converts a node to it's path
func (n *Node) ToList() []uint32 {
	var res []uint32
	for n != nil {
		res = append(res, n.Id)
		n = n.Parent
	}
	// List is reversed
	for i := 0; i < len(res)/2; i++ {
		j := len(res) - i - 1
		res[i], res[j] = res[j], res[i]
	}
	return res
}
