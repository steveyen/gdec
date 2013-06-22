package gdec

type QuorumVote struct {
	VoterId string
}

type QuorumResult struct {
	Addr string
}

func QuorumInit(d *D, prefix string,
	quorumSize int, resultAddr string) *D {
	qvote := d.DeclareChannel(prefix+"QuorumVote", QuorumVote{})
	qresult := d.DeclareChannel(prefix+"QuorumResult", QuorumResult{})

	qvotes := d.DeclareLSet(prefix+"quorumVotes", QuorumVote{})
	qtally := d.DeclareLMax(prefix + "quorumTally")
	qreached := d.DeclareLBool(prefix + "quorumReached")

	d.Join(qvote).
		Into(qvotes)

	d.Join(func() int { return qvotes.Size() }).
		Into(qtally)

	d.Join(func() bool { return qtally.Int() >= quorumSize }).
		Into(qreached)

	d.Join(func() *QuorumResult {
		if qreached.Bool() {
			return &QuorumResult{resultAddr}
		}
		return nil
	}).IntoAsync(qresult)

	return d
}
