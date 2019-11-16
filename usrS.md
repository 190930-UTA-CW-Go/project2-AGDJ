User Story
JAG-D is a service that will allow one main server to set up and install applications from the apt package on many clients, through the use of API endpoints.
#
CLIENT
- [ ] API endpoints to grab client information
- [ ] API fucntion needs to parse data into a structure and send out the structure using json
- [ ] JSON format data structures to be able to encode information
- [ ] API endpoint to send configuration data to the client
- [ ] HTTP server to handle API
#
SERVER
- [ ] HTTP POST function to send a list of desired applications to be installed
- [ ] HTTP GET function to grab information off the client
- [ ] Database to store IP addresses and Administrator info
- [ ] Have the same data structures as on client to be able to receive the buffered data
- [ ] HTTP server with a front end.
- [ ] Present information in an HTML template that parses the data. Stats Page
- [ ] Users should have a user sign in for administrators
- [ ] users get a welcome page describing the application
- [ ] users get a web page where they can install multiple applications form a list on their client 
- [ ] Administrative users have the abolity to add other moderators to edit client machines and give them a sign in and a password

