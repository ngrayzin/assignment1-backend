# eti-assignment1-user-backend

This is the first microservice of the ETI assignment 1. This backend handles everything related to users. Please run the db first before starting this.

To run the backend, run `go run main.go` in the project command prompt.

## Considerations for user backend
- Users need to complete any trips before deleting or becomning a car owner account. E.g Passengers need to complete all trips or have their trips cancelled before becoming a car owner/deleting their account OR Car owners end or cancel all their trips before deleting their account.

- Once a passenger account become a car owner account, they cannot switch back to a passenger account.

- Each user must have a unqiue email when signing up.

- To verify car owner, they must have a car plate number and driver license number. I assume they would be strings and the user will be truthful about the information they provide 😃

- Users can edit any information on the profile page EXCEPT their userID and account creation date.

- Car owner account cannot enroll for trips.
  
## Architecture diagram

![Blank diagram - Page 1](https://github.com/ngrayzin/eti-assignment1-user-backend/assets/94064635/e05165d0-c9be-4b49-95ff-57f58053d203)
