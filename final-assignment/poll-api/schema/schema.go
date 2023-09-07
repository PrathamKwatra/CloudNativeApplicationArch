package schema

import (
	"time"
)

type Vote struct {
	Id        int   `json:"id"`
	PollId    int   `json:"pollId"`
	VoterId   int   `json:"voterId"`
	VoteValue int   `json:"voteValue"` // chosen option
	Links     Links `json:"_links"`
	Embedded  any   `json:"_embedded,omitempty"`
	Meta      Meta  `json:"_meta,omitempty"`
}

type VoterPoll struct {
	PollId  int       `json:"pollId"`
	VoteId  int       `json:"voteId"`
	VotedAt time.Time `json:"votedAt"`
}

type Voter struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Email      string      `json:"email"`
	VoterPolls []VoterPoll `json:"voterPolls"`
	Links      Links       `json:"_links"`
	Embedded   any         `json:"_embedded,omitempty"`
	Meta       Meta        `json:"_meta,omitempty"`
}

type pollOption struct {
	Id   int    `json:"id"`
	Text string `json:"text"`
}

type Results struct {
	OptionId int `json:"optionId"`
	Votes    int `json:"votes"`
}

type Poll struct {
	Id       int          `json:"id"`
	Title    string       `json:"title"`
	Question string       `json:"question"`
	Options  []pollOption `json:"options"`
	Results  []Results    `json:"results"`

	Links    Links `json:"_links"`
	Embedded any   `json:"_embedded,omitempty"`
	Meta     Meta  `json:"_meta,omitempty"`
}

type Link struct {
	Href string `json:"href"`
}

type Links struct {
	Self    Link `json:"self"`
	Poll    Link `json:"poll,omitempty"`
	Vote    Link `json:"vote,omitempty"`
	Votes   Link `json:"votes,omitempty"`
	Voter   Link `json:"voter,omitempty"`
	Voters  Link `json:"voters,omitempty"`
	Polls   Link `json:"polls,omitempty"`
	Results Link `json:"results,omitempty"`
}

type Meta struct {
	TotalPolls int       `json:"TotalPolls,omitempty"`
	TotalVotes int       `json:"TotalVotes,omitempty"`
	CreatedAt  time.Time `json:"CreatedAt,omitempty"`
	UpdatedAt  time.Time `json:"UpdatedAt,omitempty"`
}
