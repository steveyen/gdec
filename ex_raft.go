package gdec

// Invoked by candidates to gather votes.
type RaftVoteRequest struct {
	CandidateTerm int    // Candidate's term.
	CandidateAddr string // Candidate requesting vote.
	LastLogTerm   int    // Term of candidate's last log entry.
	LastLogIndex  int    // Index of candidate's last log entry.
}

type RaftVoteResponse struct {
	FromAddr    string
	Term        int  // Current term, for candidate to update itself.
	VoteGranted bool // True means candidate received vote.
}

type RaftAppendEntryRequest struct {
	LeaderTerm   int    // Leader's term.
	LeaderAddr   string // So follower can redirect clients.
	PrevLogTerm  int    // Term of log entry immediately preceding this one.
	PrevLogIndex int    // Index of log entry immediately preceding this one.
	Entry        string // Log entry to store (empty for heartbeat).
	CommitIndex  int    // Last entry known to be commited.
}

type RaftAppendEntryResponse struct {
	FromAddr    string
	Term        int  // Current term, for leader to update itself.
	Success     bool // True if follower contained entry matching PrevLogIndex/Term.
	CommitIndex int
}

type RaftEntry struct {
	Term  int    // Term when entry was received by leader.
	Index int    // Position of entry in the log.
	Entry string // Command for state machine.
}

type RaftPeer struct {
	Addr         string
	PrevLogIndex int
}

func RaftProtocolInit(d *D, prefix string) *D {
	d.DeclareChannel(prefix+"RaftVoteRequest", RaftVoteRequest{})
	d.DeclareChannel(prefix+"RaftVoteResponse", RaftVoteResponse{})
	d.DeclareChannel(prefix+"RaftAppendEntryRequest", RaftAppendEntryRequest{})
	d.DeclareChannel(prefix+"RaftAppendEntryResponse", RaftAppendEntryResponse{})
	return d
}

func RaftInit(d *D, prefix string) *D {
	d = RaftProtocolInit(d, prefix)

	// rvotes := d.Relations[prefix+"RaftVoteRequest"]
	// rvoteresponse := d.Relations[prefix+"RaftVoteResponse"]
	// rappends := d.Relations[prefix+"RaftAppendEntryRequest"]
	// rappendresponse := d.Relations[prefix+"RaftAppendEntryResponse"]

	// TODO.

	return d
}

func init() {
	RaftInit(NewD(""), "")
}
