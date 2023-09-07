## Final Project

# Description
This project aims to containerize three different APIs that interact with each other
and allow them to scale with size. They are all created in golang and utilize gin with
redis to store and interact with voting data.

# Build
Each API has a docker file, and there is also a docker compose file. In an ideal scenario
to build and run the project: ```.\start.sh``` should be enough.

# Testing
There is a python test script that tests some basic and integrated tests in the API. Such
as multiple voters, invalid Id, non-existent poll id or voter id, updating polls and more.
To run this file, use the ```.\testing.sh```
This requires python3.7+ and activates python using ```python3 [command]```

While there is a dockerfile present for the python test, it currently has some issues on certain devices.
It works on most devices!
```
docker compose --profile test up
```

There is also some more testing features available. If needed, a sample db creator command is available to
create a sample db. Though, this requires some extra work, as it would require commenting the ```main()```
on line 556 and uncommenting the line 557.

To view the data in redis, uncomment line 9 and 16 in ```compose.yaml```

# Limitations
The DELETE commands sent on voter or poll does not search for related votes. Hence the votes
are not deleted if the poll/voter is deleted. This may cause a problem if a voter/poll is 
deleted first, as there is a check to make sure a vote cannot be deleted if the voter/poll has
been deleted. All links that can be possibly related are part of the response, while omitempty has
been set in the go schema's, for a particular that I don't know about, go is not omitting them from
the response. Lastly, there is a counter in polls and votes, that updates the counts based on the
vote, this count is not protected by any *lock* hence, they maybe slightly inaccurate.
