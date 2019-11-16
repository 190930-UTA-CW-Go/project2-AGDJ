User Story
JAG-D is a service that will allow one main server to set up and install applications from the apt package on many clients, through the use of API endpoints.
#
CLIENT
- [ ] 2 API endpoints to grab client information
- [ ] 1 API fucntion needs to parse data into a structure and send out the structure using json
- [ ] 1 JSON format data structures to be able to encode information
- [ ] 2 API endpoint to send configuration data to the client
- [ ] 1 HTTP server to handle API
#
SERVER
- [ ] 2 HTTP POST function to send a list of desired applications to be installed
- [ ] 3 HTTP GET function to grab information off the client
- [ ] 2 Database to store IP addresses and Administrator info
- [ ] 2 Have the same data structures as on client to be able to receive the buffered data
- [ ] 3 HTTP server with a front end.
- [ ] 1 Present information in an HTML template that parses the data. Stats Page
- [ ] 2 Users get a web page where they can install multiple applications form a list on their client 
- [ ] 1 Administrative users have the ability to add other moderators to edit client machines and give them a sign in and a password
- [ ] 1 A sign in page for administrators.
- [ ] 1 A welcome page that describes the application, registers new administrators, and adds new machines.

