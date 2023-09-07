from typing import Any, List

from pydantic import BaseModel, Json
from datetime import date, datetime, time, timedelta

APIs = {
    'polls': 'https://localhost:1082/polls',
    'voters': 'https://localhost:1081/voters',
    'votes': 'https://localhost:1080/votes',
}

class Link(BaseModel):
    Href: str
    
class Links(BaseModel):
    Self: Link
    Poll: Link
    Vote: Link
    Voter: Link
    Voters: Link
    Polls: Link
    Results: Link

class Meta(BaseModel):
    TotalPolls: int
    TotalVotes: int
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
    Options: List[PollOption]
    Results: List[Results]
    Links: Links
    Embedded: Json[Any]
    Meta: Meta

class VoterPoll(BaseModel):
    PollId: int
    VoterId: int
    VotedAt: datetime

class Voter(BaseModel):
    Id: int
    Name: str
    Email: str
    VoterPolls: List[VoterPoll]
    Links: Links
    Embedded: Json[Any]
    Meta: Meta

class Votes(BaseModel):
    Id: int
    PollId: int
    VoterId: int
    VoteValue: int
    Links: Links
    Embedded: Json[Any]
    Meta: Meta

