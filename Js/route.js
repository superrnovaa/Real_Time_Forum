import * as auth from './Auth.js';
import * as app from './app.js';
import * as error from './error.js';

async function handleRoute() {
  const path = window.location.pathname;

  if (path === "/") {
    auth.renderLoginTemplate();
  } else if (path === "/SignUp") {
    auth.renderSignUpTemplate();
  } else if ((path === "/HomePage") || (path === "/Chat") || (path.startsWith("/Post/")) || (path === "/CreatePost")) {
    app.renderAppTemplate()
  } else {
    error.GetError()
  }
};

handleRoute();
