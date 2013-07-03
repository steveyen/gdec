package gdec

type ShortestPathLink struct {
	From string
	To   string
	Cost int
}

type ShortestPath struct {
	From string
	To   string
	Next string
	Cost int
}

func ShortestPathInit(d *D, prefix string) *D {
	links := d.DeclareLSet(prefix+"ShortestPathLink", ShortestPathLink{})
	paths := d.DeclareLSet(prefix+"ShortestPath", ShortestPath{})

	d.Join(links, func(link *ShortestPathLink) *ShortestPath {
		return &ShortestPath{From: link.From, To: link.To, Cost: link.Cost}
	}).Into(paths)

	d.Join(links, paths, func(link *ShortestPathLink, path *ShortestPath) *ShortestPath {
		if link.To != path.From {
			return nil
		}
		return &ShortestPath{link.From, path.To, link.To, link.Cost + path.Cost}
	}).Into(paths)

	return d
}

func init() {
	ShortestPathInit(NewD(""), "")
}
