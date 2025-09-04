# RSS Feed Aggregator

RSS Feed Aggregator is a microservice built in Go to add RSS feeds and download posts.

Go was used as it was designed for concurrency, allowing for efficient downloading of posts.

PostgreSQL was used as the database as it is a production ready database and is able to handle the posts downloaded.

## Installation

Create a .env file in the root of the folder with PORT and DB_URL

```dotenv
PORT=<PORT>
DB_URL=postgres://<username>:<password>@<host>:<port>/<database>?sslmode=disable
```

Build the project.

```bash
go build && ./rssagg
```

## Usage

For Authenticated endpoints, you need to add a Header in the format:
Authorization: ApiKey apikey

```bash
./rssagg

# Health check
https://localhost/v1/healthz

# Create User and get API Key
https://localhost/v1/users

# Create Resource (Authenticated)
https://localhost/v1/feeds

# Get posts (Authenticated)
https://localhost/v1/posts

# Unfollow feed (Authenticated)
https://localhost/v1/feed_follows/{feedFollowID}
```

## License

[MIT](https://choosealicense.com/licenses/mit/)