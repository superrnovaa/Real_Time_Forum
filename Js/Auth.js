import * as temp from "./template.js";
import * as session from "./Session.js";
import * as app from "./app.js";

export function renderLoginTemplate(error = "") {
  const template = temp.LoginTemplate(error);
  const style = `<link href="../Style/LoginPage.css" rel="stylesheet" />`;
  document.getElementById("content").innerHTML = style + template;
  renderLoginEvents();
}

export async function renderSignUpTemplate(error = "") {
  const template = temp.SignUpTemplate(error);
  const style = `<link href="../Style/LoginPage.css" rel="stylesheet" />`;
  document.getElementById("content").innerHTML = style + template;
  renderSignUpEvents();
}

function renderLoginEvents() {
  const SignUpBtn = document.getElementById("SignUp");
  SignUpBtn.addEventListener("click", function (e) {
    e.preventDefault()
    renderSignUpTemplate();
    window.history.pushState(null, '', "/SignUp");
  });

  document
    .getElementById("Login_form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      const username = document.querySelector("#username").value;
      const password = document.querySelector("#password").value;
      const jsonFormat = JSON.stringify({ username, password });

      try {
        const res = await fetch("/", {
          method: "POST",
          body: jsonFormat,
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        if (!res.ok) {
          const errorMessage = await res.text();
          renderLoginTemplate(errorMessage);
        } else {
          const responseData = await res.json();
          session.setSessionCookie(responseData.sessionId);
          e.preventDefault()
          app.renderAppTemplate()
          window.history.pushState(null, '', "/HomePage");
        
         // window.location.pathname ="/HomePage";
        }
      } catch (error) {
        renderLoginTemplate(error);
      }
    });
}

function renderSignUpEvents() {
  const LoginBtn = document.getElementById("Login");
  LoginBtn.addEventListener("click", function (e) {
    e.preventDefault();
    renderLoginTemplate();
    window.history.pushState(null, '', "/");
  });

  document.getElementById("Signup_form")
    .addEventListener("submit", async (e) => {
      e.preventDefault();

      const email = document.querySelector("#email").value;
      const username = document.querySelector("#username").value;
      const FirstName = document.querySelector("#FirstName").value;
      const LastName = document.querySelector("#LastName").value;
      const Age = document.querySelector("#Age").value;
      const Gender = document.querySelector("#Gender").value;
      const password = document.querySelector("#password").value;
      const ConfirmPassword = document.querySelector("#ConfirmPassword").value;
      if (
        checkForEmptyFields(
          email,
          username,
          FirstName,
          LastName,
          Age,
          Gender,
          password,
          ConfirmPassword
        )
      ) {
        alert("Please fill out all fields");
        return;
      }

      const jsonFormat = JSON.stringify({
        email,
        username,
        FirstName,
        LastName,
        Age,
        Gender,
        password,
        ConfirmPassword,
      });

      try {
        const res = await fetch("/SignUp", {
          method: "POST",
          body: jsonFormat,
          headers: { "Content-Type": "application/json" },
          credentials: "include",
        });

        if (!res.ok) {
          const errorMessage = await res.text();
          renderSignUpTemplate(errorMessage);
        } else {
          const responseData = await res.json();
          session.setSessionCookie(responseData.sessionId);
          e.preventDefault()
          app.renderAppTemplate()
          window.history.pushState(null, '', "/HomePage");
          //window.location.pathname ="/HomePage";
        }
      } catch (error) {
        renderSignUpTemplate(error);
      }
    });
}

function checkForEmptyFields(
  email,
  username,
  firstName,
  lastName,
  age,
  gender,
  password,
  confirmPassword
) {
  const fields = [
    email,
    username,
    firstName,
    lastName,
    age,
    gender,
    password,
    confirmPassword,
  ];

  return fields.some((field) => {
    return field.trim() === "";
  });
}
