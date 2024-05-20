# Notes service

NotesService is a note-taking service designed to facilitate the creation and management of user notes. It provides user authentication, follows a multi-tier architecture, and offers a REST API for interacting with notes.

## Description

NotesService is a powerful tool for organizing and storing user notes. It allows users to create, edit, and delete their notes, as well as view all their notes in a convenient format. With built-in authentication, users can register, log in, and save their notes under their account.

The NotesService project is developed using a multi-tier architecture, which ensures logical separation of system components and improves scalability and maintainability. The REST API provides an easy way to interact with notes and allows developers to integrate the service into various client applications.

## API Endpoints

### Authentication

- **POST /auth/login** - Logs in a user. Requires a JSON body with the following fields: `email` and `password`.
- **POST /auth/register** - Registers a new user. Requires a JSON body with the following fields: `firstName`, `email`, `password`.
- **POST /auth/logout** - Logs out a user. Requires authentication using session.

### UserController

- **GET /users** - Retrieves information about a user with in session. Requires authentication using session.
- **PUT /users** - Updates user information. Requires authentication using session.
- **DELETE /users** - Deletes a user. Requires authentication using session.

### NoteController

- **POST /notes** - Creates a new note. Requires authentication using session.
- **GET /notes/{id}** - Retrieves information about a note with the specified ID. Requires authentication using session.
- **GET /notes** - Retrieves a list of all notes. Requires authentication using session.
- **PUT /notes/{id}** - Updates information about a note with the specified ID. Requires authentication using session.
- **DELETE /notes/{id}** - Deletes a note with the specified ID. Requires authentication using session.

Please note that all endpoints requiring authentication utilize the SessionMiddleware middleware.
