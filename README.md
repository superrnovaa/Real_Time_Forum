# Real-Time Forum
This project aims to create an upgraded forum with additional features such as private messages, real-time actions, and improved user experience.

## Objectives

The main objectives of this project are:

1. Registration and Login
2. Creation of Posts
3. Commenting on Posts
4. Private Messages

The project will be divided into the following parts:

- SQLite for data storage
- Golang for handling data and Websockets (Backend)
- JavaScript for handling Frontend events and client Websockets
- HTML for organizing page elements
- CSS for styling page elements

Project is developed as a single-page application, meaning all page changes will be handled within the JavaScript code.

## Registration and Login

To access the upgraded forum, users must register and log in. The registration and login process should include the following features:

- Users can fill out a registration form with the following details:
  - Nickname
  - Age
  - Gender
  - First Name
  - Last Name
  - E-mail
  - Password

- Users can log in using either their nickname or their e-mail combined with the password.

- Users can log out from any page on the forum.

## Posts and Comments
Users are able to:

- Create posts with categories.

- Create comments on posts.

- View posts in a feed display.

- View comments on a post by clicking on it.

## Private Messages

Users will be able to send private messages to each other. The chat functionality includes the following:

- A section displaying online/offline users available for communication.

- Section is organized based on the last message sent. If a user has no previous messages, it should be organized alphabetically.

- The section displaying online users is visible at all times.

- Clicking on a user's name loads past messages between the two users.

- Chats between users:
  - Displays previous messages.
  - Loads the last 10 messages initially.
  - Loads an additional 10 messages when scrolled up, avoiding excessive scroll event triggering.

- Messages works in real-time, where a user receives notifications of new messages without refreshing the page.

## AUTHORS

- Fatema Alawadhi
- Tasneem Mearaj

## LICENSES

This program developed within the scope of Reboot.
