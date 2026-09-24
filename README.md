# GATOR

Hey there! Thank you for taking a look at this little project!

Gator is a blog aggregator, it keeps track of xml blog feeds and saves them in a local database for reading. 

## How to Install

The easiest way is of course using the `go install github.com/AyuorusAguilar/gator` command, but by using that method, you would need to download the db migrations in the `sql/schema/` directory of this project separately. The recomended way though, is to just copy the source code using git, just go to the directory you want to save the project at and run `git clone github.com/AyuorusAguilar/gator`. No mather which method you use, you'd need to have installed the go toolchain. 

One pretty easy way to get the go toolchain is downloading it from webi, simply run `curl -sS https://webi.sh/golang | sh; source ~/.config/envman/PATH.env`.

Also, you additionally need to install postgres and create a database for the program to use. You can install postgres from the default package manager of Debian (apt) `sudo apt install postgresql postgresql-contrib`. 

After doing that just get into the psql client and create the database `gator`. Connect using `sudo -u postgres psql` and there run `CREATE DATABASE gator;`. Also you'll need a password so create a user or just set a password for the default user `postgres` using the command `ALTER USER postgres PASSWORD 'postgres';`.

Once you have that part ready, you need to get the connection string, just replace the username and password acordingly

`postgres://<username>:<password>@localhost:5432/gator`

To set up the database you can use goose and run the migrations in the `sql/schema/` folder of this project. Simply install it using the go toolchain with `go install github.com/pressly/goose/v3/cmd/goose@latest` and then navigate to the `sql/schema/` directory and run `goose postgres <connection_string> up` just replace the connection string.
Alternatively you could just manually run the migration's queries in the `sql/schema/` folder, but I think that would be pretty annoying!.

The last requisite is to create a cofig JSON file at your home directory. It must be called `.gatorconfig.json` and should contain the following fields ("You can just copy it and paste it. Replace the connection string with yours don't worry about the empty fields")

```json
{
    "db_url": "postgres://postgres:HyperLight123@localhost:5432/gator",
    "current_user_name": "",
    "current_user_id": ""
}
```

Finaly you can run it with `go run . <your_command>` or generate a binary with `go build`

## Usage

You should first create a new user using `gator register <username>`. The command is going to log you right after, but you can switch to another registered user with `gator login <username>`.

Then you'd probably want to add some feed you like to follow. Do so with `gator addfeed <custom_feed_name> <feed_url>`. You can start fetching data from feeds you're following using the `gator add <time in a 0h0m0s format>` command, but be careful not to do requests to often or you may get banned! At last you can se the retrived posts using `gator browse <num_of_posts>`