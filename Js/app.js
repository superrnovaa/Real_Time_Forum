import * as temp from "./template.js";
import * as create from "./createPost.js";
import * as home from "./home.js";
import * as auth from "./Auth.js";
import * as post from "./post.js";
import * as chat from "./Chat.js";
import * as socket from "./WS.js";
import * as session from "./Session.js";

export async function renderAppTemplate() {
  socket.Init();
  const data = await GetHomeData("Data");
  if (data.SessionId == session.getSessionCookie()) {
    session.sessions(data);
  }
  AppTemplate(data);
  const path = window.location.pathname;
  const segments = path.split("/");
  if (path === "/CreatePost") {
    create.renderCreatPostTemplate();
  } else if (path === "/HomePage") {
    home.renderHomePageTemplate("Data");
  } else if (segments[1] === "Post") {
    post.renderPostTemplate();
  } else if (segments[1] === "Chat") {
    chat.renderChatTemplate();
  }
}

async function AppTemplate(data) {
  const template = temp.AppTemplate(data);
  const style = `<link href="../Style/Main.css" rel="stylesheet" />`;
  document.getElementById("content").innerHTML = style + template;
  renderAppEvents();
}

export async function GetHomeData(Type) {
  try {
    const response = await fetch("/HomePage", {
      method: "POST",
      headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({RequestType: Type,}),
  });
    const data = await response.json();
    return data;
  } catch (error) {
    console.error(error);
  }
}

function renderAppEvents() { 
  document.querySelector(".HomeBtn").addEventListener("click", (e) => {
    e.preventDefault();
    home.ResetFilters();
    home.renderHomePageTemplate("Data");
    window.history.pushState(null, '', "/HomePage");
  });
  document.querySelector(".CreateBtn").addEventListener("click", (e) => {
    e.preventDefault()
    create.renderCreatPostTemplate();
    window.history.pushState(null, '', "/CreatePost");
  });
  document.querySelector(".LogoutBtn").addEventListener("click", (e) => {
    e.preventDefault()
    auth.renderLoginTemplate();
    window.history.pushState(null, '', "/Logout");
    
  });

  const personalPhoto = document.querySelector(".PersonalPhoto");
  personalPhoto.addEventListener("click", function () {
    uploadInput.click();
  });

  const uploadInput = document.getElementById("UploadInput");
  uploadInput.addEventListener("change", function () {
    const file = this.files[0];
    const formData = new FormData();
    formData.append("PImg", file);

    fetch("/ProfileImageHandler", {
      method: "POST",
      body: formData,
    })
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to upload file");
        }
        console.log("File uploaded successfully.");
        location.reload();
      })
      .catch((error) => {
        console.error(error);
      });
  });
}

export async function SortDM(users) {
  let template = temp.DMTemplate(users);
  for(let i = 0; i < 10; i++) {
    if(document.querySelector(".partners-container")) {
      document.querySelector(".partners-container").innerHTML = template;
      renderChatEvents();
      break;
    }
    await new Promise(resolve => setTimeout(resolve, 100));
  }
}

function renderChatEvents() {
  const userBoxes = document.querySelectorAll(".partner");

  userBoxes.forEach((userBox) => {
    userBox.addEventListener("click", function (e) {
      e.preventDefault()
      // Get the value of the clicked user box
      const userId = this.id;
      const userName = this.getElementsByClassName("Username")[0].textContent;
      const PImg = this.querySelector(".PartnerImage").id
      // // Store the user value in session storage
      sessionStorage.setItem("User", userId);
      sessionStorage.setItem("partner-Username", userName);
      sessionStorage.setItem("scroll", 0);
      sessionStorage.setItem("partner-Img", PImg);
      // Redirect to the /Chat page
      //update the notification to 0
      let sender = parseInt(sessionStorage.getItem("UserId"));
      let receiver = parseInt(sessionStorage.getItem("User"));
      window.history.pushState(null, '', `/Chat?Sender=${sender}&Receiver=${receiver}`);
      //window.location.href = `/Chat?Sender=${sender}&Receiver=${receiver}`;
      chat.renderChatTemplate();
      socket.LoadMessages();
      ResetNotifications(receiver, sender);
      const notificationsContainer = this.querySelector(
        ".Notifications-container"
      );
      if (notificationsContainer) {
        notificationsContainer.style.display = "none";
      }
    });
  });
}

function ResetNotifications(sender, receiver) {
  const requestData = {
    sender_id: sender,
    receiver_id: receiver,
  };
  console.log("This is request:", requestData);

  fetch("/Notifications", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(requestData),
  }).catch((error) => console.error("Error:", error));
}
