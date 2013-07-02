package gdec

func QuorumInit(d *D, prefix string) *D {
	qvote := d.Input(d.DeclareLSet(prefix+"QuorumVote", "voterString"))
	qneeded := d.DeclareLMax(prefix + "QuorumNeeded")
	qreached := d.Output(d.DeclareLBool(prefix + "QuorumReached"))

	qtally := d.DeclareLSet(prefix+"quorumTally", "voterString")

	d.Join(qvote).Into(qtally)
	d.Join(func() bool { return qtally.Size() >= qneeded.Int() }).Into(qreached)

	return d
}

func init() {
	QuorumInit(NewD(""), "")
}
