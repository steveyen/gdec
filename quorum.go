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
	qtally := d.DeclareLMax(prefix+"quorumTally")
	qreached := d.DeclareLBool(prefix + "quorumReached")

	qvotes.JoinUpdate(qvote,
		func(k *QuorumVote) *QuorumVote { return k })

	qtally.Update(
		func() int { return qvotes.Size() } )

	qreached.Update(
		func() bool { return qtally.Int() >= quorumSize })

	qresult.UpdateAsync(
		func() *QuorumResult {
			if qreached.Bool() {
				return &QuorumResult{resultAddr}
			}
			return nil
		})

	return d
}
