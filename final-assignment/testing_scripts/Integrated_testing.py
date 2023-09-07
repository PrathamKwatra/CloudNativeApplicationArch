import requests
import json
import time
from jsonTypes import *

# Requester
def request(url, method, data=None):
    match method:
        case "GET":
            return requests.get(url)
        case "POST":
            return requests.post(url, json=data)
        case "PUT":
            return requests.put(url, json=data)
        case "DELETE":
            return requests.delete(url)
        case _:
            raise Exception("Invalid method")


# Voter tests
# 1. Create a new voter with a valid voterID
# 2. Create a new voter with an invalid voterID
# 3. Create a new voter with a valid voterID and then delete it
class VoterTests:
    def __init__(self, url):
        self.url = url
    
    def test1(self):
        voter = Voter(
            Id=1,
            Name="Test",
            Email=""
        )
        url = self.url + "/" + str(voter.Id)
        response = request(url, "POST", voter.dict())
        if response.status_code != 200:
            raise Exception("Test 1 failed")

    def test2(self):
        voter = Voter(
            Id=2,
            Name="Test",
            Email=""
        )
        url = self.url + "/" + str(voter.Id)
        response = request(url, "POST", voter.dict())
        if response.status_code != 200:
            raise Exception("Test 2 failed")
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 2 failed")
    
    def test3(self):
        voter = Voter(
            Id=3,
            Name="Test",
            Email=""
        )
        url = self.url + "/" + str(voter.Id)
        response = request(url, "POST", voter.dict())
        if response.status_code != 200:
            raise Exception("Test 3 failed")
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 3 failed")
        response = request(url, "GET")
        if response.status_code != 404:
            raise Exception("Test 3 failed")
    
    def cleanup(self):
        # delete voter 1
        url = self.url + "/1"
        response = request(self.url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed")
        
# Poll tests
# 1. Create a new poll with a valid pollID
# 2. Create a new poll with an invalid pollID
# 3. Create a new poll with a valid pollID and then delete it
# 4. Create a new poll with a valid pollID and then update it

# Vote tests
# 1. Create a new vote with a valid voteID
# 2. Create a new vote with an invalid voteID
# 3. Create a new vote with a valid voteID and then delete it
# 4. Create a new vote with a valid voteID and then update it

# Integrated tests
# 1. Add a new vote, ask for poll results, delete the vote, ask for poll results
# 2. Add a new vote, ask for poll results, update the vote, ask for poll results, delete the vote, ask for poll results
# 3. Add a new vote, lookup voter, look at all votes created by the voter


def main():
    # run voter tests
    # voterTests = VoterTests(APIs['voters'])
    # voterTests.test1()
    # voterTests.test2()
    # voterTests.test3()
    # voterTests.cleanup()

    # 

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
            CreatedAt=time.time(),
            UpdatedAt=time.time()
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
            CreatedAt=time.time(),
            UpdatedAt=time.time()
        )
    )
    
    voterDict = voter.dict()
    print(voterDict)
    print(poll)
    print(vote)
