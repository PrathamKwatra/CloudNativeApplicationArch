from typing import Any, List, Optional

from pydantic import BaseModel, Json
from datetime import datetime
import os

APIs = {
    'polls': 'http://localhost:1082/polls',
    'voters': 'http://localhost:1081/voters',
    'votes': 'http://localhost:1080/votes',
}

class Link(BaseModel):
    Href: str
    
class Links(BaseModel):
    Self: Link
    Poll: Optional[Link]  = None
    Vote: Optional[Link]  = None
    Voter: Optional[Link]  = None
    Voters: Optional[Link]  = None
    Polls: Optional[Link]  = None
    Results: Optional[Link]  = None

class Meta(BaseModel):
    TotalPolls: int = 0
    TotalVotes: int = 0
    CreatedAt: datetime
    UpdatedAt: datetime

class PollOption(BaseModel):
    Id: int
    Text: str

class Results(BaseModel):
    OptionId: int
    Votes: int

class Poll(BaseModel):
    Id: int
    Title: str
    Question: str
    Options: Optional[List[PollOption]]
    Results: Optional[List[Results]] = None
    Links: Optional[Links] = None
    Embedded:  Any = None
    Meta: Optional[Meta] = None

class VoterPoll(BaseModel):
    PollId: int
    VoterId: int
    VotedAt: datetime

class Voter(BaseModel):
    Id: int
    Name: str
    Email: str
    VoterPolls: Optional[List[VoterPoll]] = []
    Links: Optional[Links] = None
    Embedded:  Any = None
    Meta: Optional[Meta] = None

class Votes(BaseModel):
    Id: int
    PollId: int
    VoterId: int
    VoteValue: int
    Links: Optional[Links] = None
    Embedded:  Any = None
    Meta: Optional[Meta] = None

def testModels():
        # test voter model from jsonTypes
    voter = Voter(
        Id=1,
        Name="Test",
        Email=""
    )
    
    # test poll model from jsonTypes
    poll = Poll(
        Id=1,
        Title="Test",
        Question="Test",
        Options=[PollOption(Id=1, Text="Test")],
        Results=[Results(OptionId=1, Votes=1)],
        Links=Links(
            Self=Link(Href=APIs['polls'] + "/1"),
            Voters=Link(Href=APIs['voters']),
            Votes = Link(Href=APIs['votes']),
            Polls = Link(Href=APIs['polls']),
            Results = Link(Href=APIs['polls'] + "/1/results"),
        ),
        Embedded=None,
        Meta=Meta(
            TotalVotes=1,
            CreatedAt=datetime.now(),
            UpdatedAt=datetime.now()
        )
    )

    # test vote model from jsonTypes
    vote = Votes(
        Id=1,
        PollId=1,
        VoterId=1,
        VoteValue=1,
        Links=Links(
            Self=Link(Href=APIs['votes'] + "/1"),
            Voter=Link(Href=APIs['voters'] + "/1"),
            Poll=Link(Href=APIs['polls'] + "/1"),
            Votes=Link(Href=APIs['votes']),
            Results=Link(Href=APIs['polls'] + "/1/results"),
        ),
        Embedded=None,
        Meta=Meta(
            CreatedAt=datetime.now(),
            UpdatedAt=datetime.now()
        )
    )
    
    voterDict = voter.model_dump()
    print(voterDict)
    print(poll)
    print(vote)

def envOrDefault(key, default):
    if key in os.environ:
        return os.environ[key]
    else:
        return default

def setUpAPIs():
    APIs['polls'] = envOrDefault('POLL_API', 'http://localhost:1082/polls')
    APIs['voters'] = envOrDefault('VOTER_API', 'http://localhost:1081/voters')
    APIs['votes'] = envOrDefault('VOTE_API', 'http://localhost:1080/votes')

setUpAPIs()

if __name__ == "__main__":
    testModels()
