import requests
import json
from datetime import datetime 
from jsonTypes import *
import random 
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
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 1 failed")
        response = request(url, "GET")
        # Use pydantic to validate the response
        if response.status_code != 200:
            raise Exception("Test 1 failed")
        ret = response.json()
        if ret['id'] != voter.Id or ret['name'] != voter.Name or ret['email'] != voter.Email:
            raise Exception("Test 1 failed")
        # pretty print the response
        print(json.dumps(ret, indent=4))
        
    def test2(self):
        # invalid id if id already exists, this requires test 1 to be run first
        voter = Voter(
            Id=1,
            Name="Test",
            Email=""
        )
        url = self.url + "/" + str(voter.Id)
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code == 200:
            raise Exception("Test 2 failed")
    
    def test3(self):
        voter = Voter(
            Id=3,
            Name="Test",
            Email=""
        )
        url = self.url + "/" + str(voter.Id)
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 3 failed - voter not created")
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 3 failed - voter not deleted")
        response = request(url, "GET")
        # if response.status_code != 404:
        if response.status_code == 200:
            raise Exception("Test 3 failed - voter not deleted")
    
    def cleanup(self):
        # delete voter 1
        url = self.url + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed")
        
# Poll tests
# 1. Create a new poll with a valid pollID
# 2. Create a new poll with an invalid pollID
# 3. Create a new poll with a valid pollID and then delete it
# 4. Create a new poll with a valid pollID and then update it
class PollTests:
    def __init__(self, url):
        self.url = url
    
    def test1(self):
        poll = Poll(
            Id=1,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = self.url + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 1 failed" + str(response.status_code) + " " + response.text)
        response = request(url, "GET")
        # Use pydantic to validate the response
        if response.status_code != 200:
            raise Exception("Test 1 failed" + str(response.status_code) + " " + response.text)
        ret = response.json()
        if ret['id'] != poll.Id or ret['title'] != poll.Title or ret['question'] != poll.Question:
            raise Exception("Test 1 failed")
        # pretty print the response
        print(json.dumps(ret, indent=4))

    def test2(self):
        # invalid id if id already exists, this requires test 1 to be run first
        poll = Poll(
            Id=1,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = self.url + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code == 200:
            raise Exception("Test 2 failed")

    def test3(self):
        poll = Poll(
            Id=3,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = self.url + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 3 failed - poll not created")
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 3 failed - poll not deleted")
        response = request(url, "GET")
        # if response.status_code != 404:
        if response.status_code == 200:
            raise Exception("Test 3 failed - poll not deleted")
    
    def test4(self):
        poll = Poll(
            Id=4,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = self.url + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 4 failed - poll not created")
        poll.Title = "Test2"
        poll.Options = [
            PollOption(Id=1, Text="Test2"),
            PollOption(Id=2, Text="Test2"),
            PollOption(Id=3, Text="Test2"),
            PollOption(Id=4, Text="Test2")
        ]
        response = request(url, "PUT", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 4 failed - poll not updated")
        response = request(url, "GET")
        if response.status_code != 200:
            raise Exception("Test 4 failed - poll not updated")
        ret = response.json()
        if (ret['id'] != poll.Id or ret['title'] != poll.Title or ret['question'] != poll.Question or len(ret['options']) != len(poll.Options)        
        or ret['options'][0]['id'] != poll.Options[0].Id 
        or ret['options'][0]['text'] != poll.Options[0].Text 
        or ret['options'][1]['id'] != poll.Options[1].Id 
        or ret['options'][1]['text'] != poll.Options[1].Text        
        ):
            raise Exception("Test 4 failed")
        # pretty print the response
        print(json.dumps(ret, indent=4))
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 4 failed - poll not deleted")

    def cleanup(self):
        # delete poll 1
        url = self.url + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed")

# Vote tests
# 1. Create a new vote with a valid voteID
# 2. Create a new vote with an invalid voteID
# 3. Create a new vote with an valid voteID and invalid pollID
# 4. Create a new vote with an valid voteID and invalid voterID
# 5. Create a new vote with a valid voteID and then delete it
class VoteTests:
    def __init__(self, url):
        self.url = url
    
    def startup(self):
        # create a voter
        voter = Voter(
            Id=1,
            Name="Test",
            Email=""
        )
        url = APIs['voters'] + "/" + str(voter.Id)
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")
        # create a poll
        poll = Poll(
            Id=1,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = APIs['polls'] + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")
    
    def test1(self):
        vote = Votes(
            Id=1,
            PollId=1,
            VoterId=1,
            VoteValue=1
        )
        url = self.url + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 1 failed -" + str(response.status_code) + " " + response.text)
        response = request(url, "GET")
        # Use pydantic to validate the response
        if response.status_code != 200:
            raise Exception("Test 1 failed -" + str(response.status_code) + " " + response.text)
        ret = response.json()
        if ret['id'] != vote.Id or ret['pollId'] != vote.PollId or ret['voterId'] != vote.VoterId or ret['voteValue'] != vote.VoteValue:
            raise Exception("Test 1 failed -" + str(response.status_code) + " " + response.text)
        # pretty print the response
        print(json.dumps(ret, indent=4))
    
    def test2(self):
        # invalid id if id already exists, this requires test 1 to be run first
        vote = Votes(
            Id=1,
            PollId=1,
            VoterId=1,
            VoteValue=1
        )
        url = self.url + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code == 200:
            raise Exception("Test 2 failed")
        
    def test3(self):
        # invalid poll id
        vote = Votes(
            Id=2,
            PollId=2,
            VoterId=1,
            VoteValue=1
        )
        url = self.url + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code == 200:
            raise Exception("Test 3 failed")
    
    def test4(self):
        # invalid voter id
        vote = Votes(
            Id=3,
            PollId=1,
            VoterId=2,
            VoteValue=1
        )
        url = self.url + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code == 200:
            raise Exception("Test 4 failed")
    
    def test5(self):
        vote = Votes(
            Id=5,
            PollId=1,
            VoterId=1,
            VoteValue=1
        )
        url = self.url + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 5 failed - vote not created")
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 5 failed - vote not deleted")
        response = request(url, "GET")
        if response.status_code == 200:
            raise Exception("Test 5 failed - vote not deleted")

    def cleanup(self):
        # Since there is no recursive deletion, a vote cannot be deleted if a voter or poll is deleted
        # delete vote 1
        url = self.url + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed - vote not deleted")
        # delete voter 1
        url = APIs['voters'] + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed - voter not deleted")
        # delete poll 1
        url = APIs['polls'] + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed - poll not deleted")
        

# Integrated tests
# 1. Add a new vote, ask for poll results, delete the vote, ask for poll results
# 2. Add a new vote, ask for poll result
# 3. Look for previous vote, lookup voter, look at all votes created by the voter
class IntegratedTests:
    def __init__(self, url):
        self.url = url
    
    def startup(self):
        # create a voter
        voter = Voter(
            Id=1,
            Name="Test",
            Email=""
        )
        url = APIs['voters'] + "/" + str(voter.Id)
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")
        # create a poll
        poll = Poll(
            Id=1,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = APIs['polls'] + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")
        
    def test1(self):
        vote = Votes(
            Id=1,
            PollId=1,
            VoterId=1,
            VoteValue=1
        )
        url = APIs['votes'] + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 1 failed -" + str(response.status_code) + " " + response.text)
        response = request(url, "GET")
        # use the _links to get the poll results
        if response.status_code != 200:
            raise Exception("Test 1 failed -" + str(response.status_code) + " " + response.text)
        ret = response.json()
        ret = ret['_links']['results']['href']
        response = request(ret, "GET")
        if response.status_code != 200:
            raise Exception("Test 1 failed -" + str(response.status_code) + " " + response.text)
        # pretty print the response
        print(json.dumps(response.json(), indent=4))

    def test2(self):
        vote = Votes(
            Id=2,
            PollId=1,
            VoterId=2,
            VoteValue=1
        )
        voter = Voter(
            Id=2,
            Name="Test",
            Email=""            
        )
        url = APIs['voters'] + "/" + str(voter.Id)
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")
                
               
        url = APIs['votes'] + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Test 2 failed -" + str(response.status_code) + " " + response.text)
        response = request(url, "GET")
        # use the _links to get the poll results
        if response.status_code != 200:
            raise Exception("Test 2 failed -" + str(response.status_code) + " " + response.text)
        ret = response.json()
        ret = ret['_links']['results']['href']
        response = request(ret, "GET")
        if response.status_code != 200:
            raise Exception("Test 2 failed -" + str(response.status_code) + " " + response.text)
        # expects the result to have increased by 1 vote
        ret = response.json()
        if ret['results'][1]['votes'] != 2:
            raise Exception("Test 2 failed -" + str(response.status_code) + " " + response.text)
        # pretty print the response
        print(json.dumps(response.json(), indent=4))
        
        # delete the voterId=2, pollId=2, voteId=2
        url = APIs['voters'] + "/2"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 2 failed - voter not deleted")
        url = APIs['votes'] + "/2"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Test 2 failed - vote not deleted")
        
    def test3(self):
        url = APIs['votes'] + "/1"
        response = request(url, "GET")
        # use the _links to get the voter
        if response.status_code != 200:
            raise Exception("Test 3 failed -" + str(response.status_code) + " " + response.text)
        ret = response.json()
        ret = ret['_links']['voter']['href']
        response = request(ret, "GET")
        if response.status_code != 200:
            raise Exception("Test 3 failed -" + str(response.status_code) + " " + response.text)
        # pretty print the response
        print(json.dumps(response.json(), indent=4))

    def cleanup(self):
        # delete voter 1
        url = APIs['voters'] + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed - voter not deleted")
        # delete poll 1
        url = APIs['polls'] + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed - poll not deleted")
        # delete vote 1
        url = APIs['votes'] + "/1"
        response = request(url, "DELETE")
        if response.status_code != 200:
            raise Exception("Cleanup failed - vote not deleted")
        
        
def main():
    # run voter tests
    voterTests = VoterTests(APIs['voters'])
    voterTests.test1()
    voterTests.test2()
    voterTests.test3()
    voterTests.cleanup()

    # run poll tests
    pollTests = PollTests(APIs['polls'])
    pollTests.test1()
    pollTests.test2()
    pollTests.test3()
    pollTests.test4()
    pollTests.cleanup()

    # run vote tests
    voteTests = VoteTests(APIs['votes'])
    voteTests.startup()
    voteTests.test1()
    voteTests.test2()
    voteTests.test3()
    voteTests.test4()
    voteTests.test5()
    voteTests.cleanup()
    
    # run integrated tests
    integratedTests = IntegratedTests(APIs['votes'])
    integratedTests.startup()
    integratedTests.test1()
    integratedTests.test2()
    integratedTests.test3()
    integratedTests.cleanup()
    

def makeSampleDB():
    # Add three voters
    for i in range(1, 4):
        voter = Voter(
            Id=i,
            Name="Test",
            Email=""
        )
        url = APIs['voters'] + "/" + str(voter.Id)
        response = request(url, "POST", voter.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")


    # Add three polls
    for i in range(1, 4):
        poll = Poll(
            Id=i,
            Title="Test",
            Question="Test",
            Options=[
                PollOption(Id=1, Text="Test"),
                PollOption(Id=2, Text="Test")
            ]
        )
        url = APIs['polls'] + "/" + str(poll.Id)
        response = request(url, "POST", poll.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")

    # Add 30 votes
    for i in range(1, 31):
        vote = Votes(
            Id=i,
            PollId=random.randint(1, 3),
            VoterId=random.randint(1, 3),
            VoteValue=random.randint(1, 2)
        )
        url = APIs['votes'] + "/" + str(vote.Id)
        response = request(url, "POST", vote.model_dump(mode='json'))
        if response.status_code != 200:
            raise Exception("Startup failed")
        


if __name__ == "__main__":
    main()
    # makeSampleDB()
