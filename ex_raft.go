package gdec

// Invoked by candidates to gather votes.
type RaftVoteRequest struct {
	To           string
	From         string // Candidate requesting vote.
	Term         int    // Candidate's term.
	LastLogTerm  int    // Term of candidate's last log entry.
	LastLogIndex int    // Index of candidate's last log entry.
}

type RaftVoteResponse struct {
	To      string
	From    string
	Term    int  // Current term, for candidate to update itself.
	Granted bool // True means candidate received vote.
}

type RaftAppendEntryRequest struct {
	To           string
	LeaderTerm   int    // Leader's term.
	LeaderAddr   string // So follower can redirect clients.
	PrevLogTerm  int    // Term of log entry immediately preceding this one.
	PrevLogIndex int    // Index of log entry immediately preceding this one.
	Entry        string // Log entry to store (empty for heartbeat).
	CommitIndex  int    // Last entry known to be commited.
}

type RaftAppendEntryResponse struct {
	To          string
	From        string
	Term        int  // Current term, for leader to update itself.
	Success     bool // True if had entry matching PrevLogIndex/Term.
	CommitIndex int
}

type RaftEntry struct {
	Term  int    // Term when entry was received by leader.
	Index int    // Position of entry in the log.
	Entry string // Command for state machine.
}

const (
	// Lowest bits of a state are the 'kind' of a state,
	// where ordering matters for LMax precedence.
	state_FOLLOWER  = 0
	state_LEADER    = 1
	state_CANDIDATE = 2
	state_STEP_DOWN = 3 // Must be largest for LMax precedence.
	state_SAME      = 0 // Used to denote no change to state kind.

	state_KIND_MASK    = 0x0000000f
	state_VERSION_MASK = 0xfffffff0
	state_VERSION_NEXT = 0x00000010
)

func stateKind(s int) int        { return s & state_KIND_MASK }
func stateVersion(s int) int     { return s & state_VERSION_MASK }
func stateVersionNext(s int) int { return stateVersion(s) + state_VERSION_NEXT }

func RaftProtocolInit(d *D, prefix string) *D {
	d.DeclareChannel(prefix+"RaftVoteRequest", RaftVoteRequest{})
	d.DeclareChannel(prefix+"RaftVoteResponse", RaftVoteResponse{})
	d.DeclareChannel(prefix+"RaftAppendEntryRequest", RaftAppendEntryRequest{})
	d.DeclareChannel(prefix+"RaftAppendEntryResponse", RaftAppendEntryResponse{})
	return d
}

func RaftInit(d *D, prefix string) *D {
	d = RaftProtocolInit(d, prefix)

	rvote := d.Relations[prefix+"RaftVoteRequest"]
	rvoter := d.Relations[prefix+"RaftVoteResponse"]

	// rappends := d.Relations[prefix+"RaftAppendEntryRequest"]
	// rappendresponses := d.Relations[prefix+"RaftAppendEntryResponse"]

	// members := d.DeclareLSet(prefix + "raftMember", "addrString")
	// votedFor := d.DeclareLSet(prefix + "raftVotedFor", "addrString")
	// votedForInCurrentTerm := d.DeclareLSet(prefix + "raftVotedForInCurrentTerm", "addrString")
	// votedForInCurrentTick := d.DeclareLSet(prefix + "raftVotedForInCurrentTick", "addrString")

	currentTerm := d.DeclareLMax(prefix + "raftCurrentTerm")
	currentState := d.DeclareLMax(prefix + "raftCurrentState")

	nextTerm := Scratch(d.DeclareLMax(prefix + "raftNextTerm"))
	nextState := Scratch(d.DeclareLMax(prefix + "raftNextState"))

	d.Join(currentTerm).
		IntoAsync(nextTerm)
	d.Join(rvote, func(r *RaftVoteRequest) int { return r.Term }).
		IntoAsync(nextTerm)

	d.Join(rvote, currentTerm, currentState, func(r *RaftVoteRequest, term, state *int) int {
		if r.Term > *term {
			return stateKind(*state) + state_STEP_DOWN
		}
		return state_SAME
	}).Into(nextState)

	d.Join(currentState, nextState, func(curr, next *int) int {
		if *next == state_STEP_DOWN {
			return stateVersionNext(*curr)
		}
		return state_SAME
	}).IntoAsync(currentState)

	d.Join(rvote, currentTerm, func(r *RaftVoteRequest, myCurrentTerm *int) *RaftVoteResponse {
		if r.Term < *myCurrentTerm {
			return &RaftVoteResponse{
				To:      r.From,
				From:    r.To,
				Term:    *myCurrentTerm,
				Granted: false,
			}
		}
		return nil // TODO.
	}).IntoAsync(rvoter)

	// TODO.

	return d
}

func init() {
	RaftInit(NewD(""), "")
}
