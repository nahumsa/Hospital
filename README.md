# Hospital Management 

[![Hospital](https://circleci.com/gh/nahumsa/Hospital.svg?style=svg)](https://circleci.com/gh/nahumsa/Hospital)

-----------------------------

## Introduction of the project

The problem that this app is going to solve is to manage the rooms with sick patients in a hospital where the doctors are logged in a service which they fill out a form giving the number of the room and the number of the case. The interface will show all rooms and mark the rooms with patients with red.

Designing this app I found out it would work better if it has a login, thus I choose a NoSQL database since it is scalable and I can create the structure of the database how I please, I chose [MongoDB](https://www.mongodb.com/) because it is an open-source project and it has been growing in popularity because it is used in huge companies such as EA and Google. 

After this, I need a structure to fill out a form and assign it to a room of the patient, where the form will be filled by the login doctor, and added to a new database collection with the desired structure. 

This will complete the backend of the project, to the frontend I will use CSS in order to make the application responsive and prettier.

## Running the project

In order to run the project you should have both [docker](https://docs.docker.com/engine/install/) and [docker-compose](https://docs.docker.com/compose/install/) installed into your machine, then you simply just need to run the following command to start the application into the "http://localhost:8080/":

```
docker-compose up
```

In order to terminate the application you just need to run:

```
docker-compose down
```
