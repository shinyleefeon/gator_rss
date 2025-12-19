# gator_rss
RSS feed aggregator built using go and postgresql for the bootdev course

Requirements:
Go
Postgres

installation:
I'm pretty sure you can navigate to the package in a terminal and run go install .
This should allow you to run gator_rss (Arguments) from the terminal

This program requires a config file json at ~
Mine looks like this :
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator_rss?sslmode=disable",
  "current_user_name": "Ri"
}


Currently supported commands:
login - sets current username
register - adds new user
reset - deletes all users
users - prints all users

addfeed - with a title and url adds an rss feed to follow
feeds - lists current feeds
follow - sets user to following a feed in the database
following - shows currently followed feeds
unfollow - removes feed from followed feeds

agg - pulls posts from feeds oldest - newest
browse - prints titles of feeds the current user follows
