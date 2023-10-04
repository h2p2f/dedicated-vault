# dedicated-vault

This project is an implementation of a technical specification for a client-server solution for storing and transmitting sensitive information that is vulnerable to compromise, including credit card information, login credentials, arbitrary text data, and binary files.

To encode and decode user data, a user's passphrase is used, which is not stored either on the server or on the client in any form.

The server solution is built on a `MongoDB` database and client connections are made through `GRPC`. Implementation features include only accepting `GRPC` server connections with trusted `TLS` parameters and user authorization verification through `JWT` tokens. User passwords are stored in hashed form on the server side, with no possibility of decryption of sensitive information in the event of unauthorized access to the server database.

Implementation simplifications and features for the server include loading database access parameters from a YAML file `./config/config.yaml`, with production deployments requiring them to be taken from environment variables when starting containers. Docker-compose containerization has not been implemented.

The client solution is based on locally storing user data in an encrypted `sqlite3` database. A GUI interface has been implemented with `Fyne` library to allow users to register and log in to the server, add, edit, and delete information locally and remotely in the server database, and perform full data synchronization with the remote server. A single user can have multiple clients on different devices, with a mechanism for controlling the time of the last data change on the server to maintain data currency in local databases.

Implementation simplifications and features for the client include the lack of graceful shutdown due to its unique implementation in fine, the inability to delete a user from the server, and anomalous length of GUI code that is difficult to read and refactor due to multiple callbacks in element descriptions. Distribution of the client is not intended for commercial use, with key files needing to be placed in `/tmp/dedicated-vault/crypto` on Unix systems.

In Windows systems, information storage is similar to Unix systems, but in the `C:\Users\Public\` folder.

Client distributions are prepared automatically by the `cmd/client/compile.sh` script, in which you can specify the software version number, build date and specify the path to the cryptographic keys.
Distributions are located in `cmd/client/packages`