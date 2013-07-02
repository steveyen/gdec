package gdec

import (
	"fmt"
	"reflect"
	"strings"
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
	state_VERSION_MASK = 0xfffffff0 // Highest bits are version.
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

	// votedFor := d.DeclareLSet(prefix + "raftVotedFor", "addrString")
	// votedForInCurrTerm := d.DeclareLSet(prefix + "raftVotedForInCurrTerm", "addrString")
	// votedForInCurrTick := d.DeclareLSet(prefix + "raftVotedForInCurrTick", "addrString")

	// currVote := d.DeclareLSet(prefix+"raftCurrVote", "addrString") // My vote.
	nextVote := d.DeclareLSet(prefix+"raftNextVote", "addrString")

	tally := d.DeclareLMap(prefix + "raftTally")                  // Key: "term:addr", val: LBool.
	yesVotes := d.Scratch(d.DeclareLMap(prefix + "raftYesVotes")) // Key: "term", val: LSet.
	wonTerm := d.Scratch(d.DeclareLSet(prefix+"raftWonTerm", 0))

	currTerm := d.DeclareLMax(prefix + "raftCurrTerm")
	currState := d.DeclareLMax(prefix + "raftCurrState")

	nextTerm := d.Scratch(d.DeclareLMax(prefix + "raftNextTerm"))
	nextState := d.Scratch(d.DeclareLMax(prefix + "raftNextState"))

	alarm := d.Scratch(d.DeclareLBool(prefix + "raftAlarm"))           // TODO: periodic.
	resetAlarm := d.Scratch(d.DeclareLBool(prefix + "raftResetAlarm")) // TODO: periodic.
	heartBeat := d.Scratch(d.DeclareLBool(prefix + "raftHeartBeat"))   // TODO: periodic.

	logState := d.DeclareLSet(prefix+"raftLogState", RaftLogState{}) // TODO: sub-module.

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

	// Timeouts.

	d.Join(alarm, currTerm, currState, func(alarm *bool, currTerm *int, currState *int) {
		// Move to candidate state, with a new term, self-vote, and alarm reset.
		if *alarm && stateKind(*currState) != state_LEADER {
			d.Add(nextTerm, *currTerm+1)
			d.Add(nextState, state_CANDIDATE)
			d.Add(nextVote, d.Addr)
			d.Add(resetAlarm, true)
			return
		}
		d.Add(nextTerm, *currTerm)
		d.Add(nextState, stateKind(*currState))
		d.Add(nextVote, "")
		d.Add(resetAlarm, false)
	})

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

	// Send vote requests.

	d.Join(heartBeat, member, currState, currTerm, logState,
		func(h *bool, mAddr *string, s *int, t *int, l *RaftLogState) *RaftVoteRequest {
			if stateKind(*s) == state_CANDIDATE && !tallyHasVoteFrom(tally, *t, *mAddr) {
				return &RaftVoteRequest{
					To:           *mAddr,
					From:         d.Addr,
					Term:         *t,
					LastLogTerm:  l.LastTerm,
					LastLogIndex: l.LastIndex,
				}
			}
			return nil
		}).IntoAsync(rvote)

	// Tally votes when we're a candidate.

	d.Join(rvoter, func(rvoter *RaftVoteResponse) int { return rvoter.Term }).
		Into(nextTerm)

	d.Join(currTerm, currState, rvoter,
		func(currTerm *int, currState *int, rvoter *RaftVoteResponse) int {
			// If our term is stale, step down.
			if stateKind(*currState) != state_FOLLOWER && rvoter.Term > *currTerm {
				return state_STEP_DOWN
			}
			return stateKind(*currState)
		}).Into(nextState)

	d.Join(currTerm, currState, rvoter,
		func(currTerm *int, currState *int, rvoter *RaftVoteResponse) *LMapEntry {
			// Record the vote if we're still a candidate and in the same term.
			if stateKind(*currState) == state_CANDIDATE && rvoter.Term == *currTerm {
				granted := d.NewLBool()
				granted.DirectAdd(rvoter.Granted)
				return &LMapEntry{voteKey(rvoter.Term, rvoter.From), granted}
			}
			return nil
		}).Into(tally)

	d.Join(tally, func(e *LMapEntry) *LMapEntry {
		term := strings.Split(e.Key, ":")[0]
		if e.Val.(*LBool).Bool() {
			return &LMapEntry{term, d.NewLSet(reflect.TypeOf(0))}
		}
		return nil
	}).Into(yesVotes)

	d.Join(wonTerm, currTerm, currState,
		func(wonTerm *int, currTerm *int, currState *int) int {
			if *wonTerm == *currTerm && stateKind(*currState) == state_CANDIDATE {
				return state_LEADER
			}
			return stateKind(*currState)
		}).Into(nextState)

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

func voteKey(term int, addr string) string {
	return fmt.Sprintf("%d:%s", term, addr)
}

func tallyHasVoteFrom(tally *LMap, term int, addr string) bool {
	return tally.At(voteKey(term, addr)) != nil
}
