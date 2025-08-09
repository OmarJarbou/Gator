# 🐊 Gator - RSS Feed Aggregator

Welcome to Gator! A powerful command-line RSS feed aggregator built with Go that helps you stay updated with your favorite content sources.

## 🚀 Features

- **User Management**: Create accounts and manage multiple users
- **RSS Feed Management**: Add, follow, and unfollow RSS feeds
- **Content Aggregation**: Automatically fetch and display posts from followed feeds
- **Continuous Feed Monitoring**: Real-time RSS feed scraping with configurable intervals
- **Simple CLI Interface**: Easy-to-use commands for all operations
- **PostgreSQL Backend**: Robust data storage with SQLC

## 📋 Prerequisites

Before you begin, you'll need to have the following installed:

- **PostgreSQL**: [Download and install PostgreSQL](https://www.postgresql.org/download/)
- **Go**: [Install Go 1.24.4 or later](https://go.dev/doc/install)

## 🛠️ Installation

Install Gator using Go's package manager:

```bash
go install github.com/OmarJarbou/Gator@latest
```

**Important**: Make sure to add Go's bin directory to your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

To make this permanent, add the above line to your shell profile (`.bashrc`, `.zshrc`, etc.).

## ⚙️ Configuration Setup

Gator requires a configuration file to connect to your PostgreSQL database. Create a `.gatorconfig.json` file in your home directory:

```bash
# Create the config file
touch ~/.gatorconfig.json
```

Add the following content to the file:

```json
{
  "db_url": "postgres://username:password@localhost:5432/database_name?sslmode=disable",
  "current_user_name": ""
}
```

**Configuration fields:**
- `db_url`: Your PostgreSQL connection string
- `current_user_name`: Will be automatically set when you login (leave empty initially)

**Example PostgreSQL connection strings:**
- Local database: `postgres://username:password@localhost:5432/gator_db?sslmode=disable`
- Remote database: `postgres://username:password@your-server.com:5432/gator_db?sslmode=disable`

## 🎯 Getting Started

1. **Create your first user account:**
   ```bash
   gator register your_username
   ```

2. **Add some RSS feeds:**
   ```bash
   gator addfeeds "Tech News" "https://example.com/feed.xml"
   gator addfeeds "Blog Updates" "https://blog.example.com/rss"
   ```
   
3. **Start aggregating feeds:**
   ```bash
   gator agg 5m  # Collect new posts every 5 minutes
   ```
   
4. **Browse your aggregated content:**
   ```bash
   gator browse
   ```

## 🔄 How Aggregation Works

The `gator agg <duration>` command is the heart of Gator. It:
- Continuously monitors all RSS feeds in the system (with prioritrizing fees that never reached before, or the oldest)
- Fetches new posts at the specified time intervals
- Stores new content in the database for browsing
- Supports various time formats: `30s`, `1m`, `5m`, `1h`, `24h`
- Runs indefinitely until manually stopped (Ctrl+C)

## 📚 Available Commands

### User Management
- **`gator login <username>`** - Authenticates and sets the current user
- **`gator register <username>`** - Creates a new user account
- **`gator users`** - Lists all registered users

### Feed Management
- **`gator addfeeds <name> <url>`** - Adds a new RSS feed and automatically follows it
- **`gator feeds`** - Lists all available RSS feeds
- **`gator follow <feed_url>`** - Follows an existing RSS feed
- **`gator unfollow <feed_url>`** - Unfollows a previously followed feed
- **`gator following`** - Shows all feeds you're currently following

### Feed Aggregation
- **`gator agg <duration>`** - **🚀 CORE FEATURE**: Continuously scrapes and aggregates RSS feeds at specified intervals (e.g., `30s`, `1m`, `5m`, `1h`)

### Content Browsing
- **`gator browse [limit]`** - Displays posts from followed feeds (default: 2 posts)

### System Management
- **`gator reset`** - **⚠️ WARNING**: Clears all users, feeds, and posts from the database

## 📖 Command Examples

```bash
# Register and login
gator register john_doe
gator login john_doe

# Add popular RSS feeds
gator addfeeds "Hacker News" "https://news.ycombinator.com/rss"
gator addfeeds "Reddit Programming" "https://www.reddit.com/r/programming/.rss"

# Follow additional feeds
gator follow "https://blog.golang.org/feed.atom"

# Start aggregating feeds (collects new posts every 5 minutes)
gator agg 5m

# Start aggregating feeds (collects new posts every hour)
gator agg 1h

# Browse recent posts
gator browse 5  # Show 5 most recent posts

# Check what you're following
gator following
```

## 🗄️ Database Schema

Gator uses PostgreSQL with the following main tables:
- `users` - User accounts
- `feeds` - RSS feed sources
- `feed_follows` - User-feed relationships
- `posts` - Aggregated content from feeds

## 🔧 Development

If you want to contribute or run from source:

```bash
# Clone the repository
git clone https://github.com/OmarJarbou/Gator.git
cd Gator

# Install dependencies
go mod download

# Build the binary
go build -o gator
```

## 📝 License

This project is open source. Please check the repository for license information.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🐛 Issues

If you encounter any issues or have questions, please [open an issue](https://github.com/OmarJarbou/Gator/issues) on GitHub.

---

**Happy RSS aggregating! 🐊📰** 
