// Use DBML to define your database structure
// Docs: https://dbml.dbdiagram.io/docs

enum sex
{
f
m
}

enum goal
{
fun
profit
50_50
}

enum status
{
active
inactive
}

// Table for storing spaces
Table space {
id integer [pk]
name text [not null,unique]
}

// Table to store about user in spaces
Table user_space {
user_id integer [pk, ref: > user.id]
space_id integer [pk, ref: > space.id]
}

// Table for storing user information
table user {
id integer [not null, unique]
username text [not null, unique]
full_name text
birthday date
gender text
city varchar
socials varchar
position varchar
sex sex
photo_url varchar
interests text
goal goal
}

// Table for storing holiday status of users
table holidays_status {
id integer [ref: > user.id, pk]
status status
till_date date [not null, default: 'null']
}

// Table for storing meeting information between users
table event {
id integer [pk, increment]
date date
name varchar
description varchar
}

table event_members {
id integer [pk, increment]
event_id integer [ref: > event.id]
user_id integer [ref: > user.id]
}

// Table for storing reviews of meetings
table reviews {
id integer [pk, increment]
event_id integer [ref: > event.id]
who_id integer [ref: > user.id]
about_whom_id integer [ref: > user.id]
grade integer [not null]
}

// Table for storing info about user's blocks
table blocks {
id integer [ref: > user.id]
user_id integer [not null]
}