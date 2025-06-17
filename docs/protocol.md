# Protocol

KVDB uses a binary message format to communicate.
The client and server here do so over TCP.

All keys/values are represented as strings.

## Message Format
Each message starts with a byte specifying the command value.
This is follwed by a byte specifying the length of the identifier of the item in the DB to operate on (*n*).
The next *n* bytes contain the identifier to operate on.

Optionally, another length byte and value can be included

**THIS LIMITS STORED VALUES TO 256 bytes**

## Commands


| Command | Value | Description                                         | Second Operand Required |
|---------|-------|-----------------------------------------------------|-------------------------|
| GET     | 0     | Fetch an item from the store                        | No                      |
| SET     | 1     | Set an item in the store. Overwrites existing items | Yes                     |
| DELETE  | 2     | Delete an item from the store                       | No                      |
