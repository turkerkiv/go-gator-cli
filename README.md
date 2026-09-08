# Gator
This is a RSS feed management cli program in go. You can 
  - register,
  - login,
  - list users,
  - add RSS feed,
  - list all feeds,
  - select one of them to follow or unfollow,
  - list following feeds,
  - and get following feeds' details with browse command.
  - Also scheduled agg command will fetch feeds into db in the background, one feed at a time which is least recently fetched.

## Installation (linux)
1. First you need postgresql and go installed in your computer.
  2. Follow this link for go installation: https://go.dev/doc/install
3. Open your terminal.
4. Create database for gator in postgresql:
```
psql postgres
CREATE DATABASE gator;
\c gator
ALTER USER postgres PASSWORD '<your_db_password>';
\q
```
4. Install goose for database migration:
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```
4. Install this project into your computer:
```
git clone https://github.com/turkerkiv/go-gator-cli.git
cd go-gator-cli
go install
cd sql/schema
goose postgres://<db_username>:<db_password>@localhost:5432/gator up
```
3. Then:
```
nano ~/.gatorconfig.json 
```
4. Then put your database connection string inside it. It looks similar to this:
```
{"db_url":"postgres://<db_username>:<db_password>@localhost:5432/gator?sslmode=disable","current_user_name":""}
```
5. Now you can use the program with commands. Examples:
```
go-gator-cli register <username> (registers that username into db)
go-gator-cli login <username> (logins to that users)
go-gator-cli users (lists all users)
go-gator-cli addfeed <url_of_feed> (adds feed to db)
go-gator-cli feeds (lists all feeds)
go-gator-cli follow <url_of_feed> (follow a feed as logged in user)
go-gator-cli unfollow <url_of_feed> (unfollow a feed)
go-gator-cli following (list all following feeds of logged in user)
go-gator-cli browse <optional_count> (browse following feed posts)
go-gator-cli agg <duration> (run this in seperate terminal. It fetches all added feeds to database in order, once in given duration (10s, 30s, 1m, 10m))
go-gator-cli reset (Don't use it. It deletes all user records)
```
