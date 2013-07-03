package gdec

import (
	"fmt"
)

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

// Invoked by leaders to replicate log entries.
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

type RaftVote struct {
	Term      int
	Candidate string
}

type RaftEntry struct {
	Term  int    // Term when entry was received by leader.
	Index int    // Position of entry in the log.
	Entry string // Command for state machine.
}

type RaftLogState struct {
	LastTerm        int
	LastIndex       int
	LastCommitIndex int
}

const (
	// The 'kind' of a state are in the lowest bits.
	state_FOLLOWER  = 0
	state_CANDIDATE = 1
	state_LEADER    = 2
	state_STEP_DOWN = 3 // Must be largest for LMax precedence.

	state_KIND_MASK    = 0x0000000f
	state_VERSION_MASK = 0xfffffff0 // Highest bits are version for precedence.
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

	member := d.DeclareLSet(prefix+"raftMember", "addrString")

	votedFor := d.DeclareLSet(prefix+"raftVotedFor", RaftVote{})
	votedForCurrTerm := d.Scratch(d.DeclareLSet(prefix+"raftVotedForCurrTerm", "addrString"))

	MultiTallyInit(d, prefix+"tally/")
	tallyVote := d.Relations[prefix+"tally/MultiTallyVote"].(*LSet)
	tallyNeed := d.Relations[prefix+"tally/MultiTallyNeed"].(*LMax)
	tallyDone := d.Relations[prefix+"tally/MultiTallyDone"].(*LMap)

	currTerm := d.DeclareLMax(prefix + "raftCurrTerm")
	currState := d.DeclareLMax(prefix + "raftCurrState")

	nextTerm := d.Scratch(d.DeclareLMax(prefix + "raftNextTerm"))
	nextState := d.Scratch(d.DeclareLMax(prefix + "raftNextState"))

	alarm := d.Scratch(d.DeclareLBool(prefix + "raftAlarm")) // TODO: periodic.
	// resetAlarm := d.Scratch(d.DeclareLBool(prefix + "raftResetAlarm")) // TODO: periodic.
	heartBeat := d.Scratch(d.DeclareLBool(prefix + "raftHeartBeat")) // TODO: periodic.

	logState := d.DeclareLSet(prefix+"raftLogState", RaftLogState{}) // TODO: sub-module.

	potentialCandidates := d.Scratch(d.DeclareLSet(prefix+"raftPotentialCandidates", RaftVoteRequest{}))

	d.Join(func() int { return member.Size() / 2 }).
		Into(tallyNeed)

	// Initialize our scratch next term/state.

	d.Join(currTerm).
		Into(nextTerm)

	d.Join(currState, func(currState *int) int { return stateKind(*currState) }).
		Into(nextState)

	// Incoming vote requests.

	d.Join(rvote, func(rvote *RaftVoteRequest) int { return rvote.Term }).
		Into(nextTerm)

	d.Join(rvote, currTerm, currState,
		func(rvote *RaftVoteRequest, currTerm *int, currState *int) int {
			if rvote.Term > *currTerm {
				return state_STEP_DOWN
			}
			return stateKind(*currState)
		}).Into(nextState)

	d.Join(rvote, currTerm,
		func(rvote *RaftVoteRequest, currTerm *int) *RaftVoteResponse {
			if rvote.Term < *currTerm {
				return &RaftVoteResponse{
					To:      rvote.From,
					From:    rvote.To,
					Term:    *currTerm,
					Granted: false,
				}
			}
			return nil // TODO.
		}).IntoAsync(rvoter)

	// Timeouts.

	d.Join(alarm, currTerm, currState, func(alarm *bool, t *int, s *int) {
		// Move to candidate state, with a new term, self-vote, and alarm reset.
		if *alarm && stateKind(*s) != state_LEADER {
			d.Add(nextTerm, *t+1)
			d.Add(nextState, state_CANDIDATE)
			d.Add(tallyVote, &MultiTallyVote{termToRace(*t + 1), d.Addr})
			// TODO: d.Add(resetAlarm, true)
			// TODO: Remove uncommitted logs.
			return
		}
	})

	// Send vote requests.

	d.Join(heartBeat, member, currState, currTerm, logState,
		func(h *bool, a *string, s *int, t *int, l *RaftLogState) *RaftVoteRequest {
			if stateKind(*s) == state_CANDIDATE &&
				!MultiTallyHasVoteFrom(d, prefix+"tally/", termToRace(*t), *a) {
				return &RaftVoteRequest{
					To:           *a,
					From:         d.Addr,
					Term:         *t,
					LastLogTerm:  l.LastTerm,
					LastLogIndex: l.LastIndex,
				}
			}
			return nil
		}).IntoAsync(rvote)

	// Tally votes when we're a candidate.

	d.Join(rvoter, func(r *RaftVoteResponse) int { return r.Term }).
		Into(nextTerm)

	d.Join(currTerm, currState, rvoter,
		func(currTerm *int, currState *int, r *RaftVoteResponse) int {
			// If our term is stale, step down as candidate or leader.
			if stateKind(*currState) != state_FOLLOWER && r.Term > *currTerm {
				return state_STEP_DOWN
			}
			return stateKind(*currState)
		}).Into(nextState)

	d.Join(currTerm, currState, rvoter,
		func(currTerm *int, currState *int, r *RaftVoteResponse) *MultiTallyVote {
			// Record granted vote if we're still a candidate in the same term.
			if stateKind(*currState) == state_CANDIDATE &&
				r.Term == *currTerm && r.Granted {
				return &MultiTallyVote{termToRace(r.Term), r.From}
			}
			return nil
		}).Into(tallyVote)

	d.Join(currTerm, currState,
		func(currTerm *int, currState *int) int {
			// Become leader if we won the race.
			if stateKind(*currState) == state_CANDIDATE {
				won := tallyDone.At(termToRace(*currTerm)).(*LBool)
				if won != nil && won.Bool() {
					return state_LEADER
				}
			}
			return stateKind(*currState)
		}).Into(nextState)

	// Cast votes.

	d.Join(currTerm, votedFor,
		func(currTerm *int, votedFor *RaftVote) *string {
			if *currTerm == votedFor.Term {
				return &votedFor.Candidate
			}
			return nil
		}).Into(votedForCurrTerm)

	d.Join(rvote, logState,
		func(rvote *RaftVoteRequest, logState *RaftLogState) *RaftVoteRequest {
			// Potential candidate only if candidate's log is >= our log.
			if (rvote.LastLogTerm > logState.LastTerm) ||
				(rvote.LastLogTerm == logState.LastTerm && rvote.LastLogIndex >= logState.LastIndex) {
				return rvote
			}
			return nil
		}).Into(potentialCandidates)

	// Incorporate next term and next state.

	d.Join(nextTerm).
		IntoAsync(currTerm)

	d.Join(nextState, currState,
		func(nextState *int, currState *int) int {
			if *nextState == state_STEP_DOWN {
				return stateVersionNext(*currState) + state_FOLLOWER
			}
			return stateVersion(*currState) + stateKind(*nextState)
		}).IntoAsync(currState)

	return d
}

func init() {
	RaftInit(NewD(""), "")
}

func termToRace(term int) string {
	return fmt.Sprintf("%d", term)
}
