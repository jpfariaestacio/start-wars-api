# startWarsApi

This repo aims to be a simple service to consume https://swapi.co/ data and store it in a MongoDB (local or not) instance.

Any changes in the code, issues or comments are warmly welcome :)

**You should have a good grasp on using mod in Go to be able to use this repo**

There are some usefull flags you can pass while building your main file.

| Command | Description |
| --- | --- |
| db-host | sets the host to your MongoDB instance |
| db-user | sets the user to your MongoDB instance |
| db-password | passes a password to user user |
| db-port | sets the port to your MongoDB instance |
| db-name | Accesses a MongoDB Database especified by name |
| db-update | If true updates values in the MongoDB by comparing with those in https://swapi.co/ |





`There are some tests missing implementation. In a future release its going to be fixed.`
