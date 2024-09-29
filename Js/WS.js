import * as app from "./app.js";
import * as type from "./typing.js";

let socket;

export function Init() {

 socket = new WebSocket("ws://localhost:8082/ws");

socket.addEventListener("open", () => {
  console.log("WebSocket connection established");
  LoadMessages();
});

socket.addEventListener("message", (event) => {
  console.log("Received message:", event.data);
});

socket.addEventListener("error", (error) => {
  console.error("WebSocket error:", error);
});

socket.addEventListener("close", () => {
  console.log("WebSocket connection closed");
});
SocketEvents();
}

export function SocketEvents() {
socket.onmessage = function (event) {
  try {
    // Parse the received data as JSON
    const chatMessages = JSON.parse(event.data);
    if (chatMessages.type == "load_messages") {
      sessionStorage.setItem("partner-Username", chatMessages.Info.RUname);
      sessionStorage.setItem("partner-Img", chatMessages.Info.RPImg);
      // Update the profile image
      const partnerProfileImg = document.getElementById("partnerProfileImg");
      if (partnerProfileImg) {
        partnerProfileImg.src =
          "../ProfileImages/" + sessionStorage.getItem("partner-Img");
      }
      // Update the username
      const partnerUsername = document.getElementById("partnerUsername");
      if (partnerUsername) {
        partnerUsername.textContent =
          sessionStorage.getItem("partner-Username");
      }
      renderMessages(chatMessages.Data);
    } else if (chatMessages.type == "send_message") {
      sendMessage(chatMessages.Content);
    } else if (chatMessages.type == "typing") {
      type.Typing(chatMessages.Content);
    } else if (chatMessages.type == "activeUsers") {
      console.log("Message received", chatMessages.Content)
      app.SortDM(chatMessages.Content);
    }
  } catch (error) {
    console.error("Error handling WebSocket message:", error);
  }
};
}

function sendMessage( message) {
  const currentUrl = new URL(window.location.href);

  // Get the query parameters
  const searchParams = currentUrl.searchParams;

  // Check if the 'Sender' and 'Receiver' parameters exist
  if (searchParams.has("Sender") && searchParams.has("Receiver")) {
    const sender = searchParams.get("Sender");
    const receiver = searchParams.get("Receiver");
    if (
      message.Sender === parseInt(receiver) ||
      message.Sender === parseInt(sender)
    ) {
      const messageElement = createMessageElement(message);
      const chatBody = document.querySelector(".ChatBody");
      chatBody.appendChild(messageElement);
      if (message.Sender === parseInt(sender)) {
        chatBody.scrollTop = chatBody.scrollHeight;
      }
    } else {
      socket.send(
        JSON.stringify({
          type: "notification",
          Content: {
            Sender: message.Sender,
            Receiver: message.Receiver,
          },
        })
      );
    }
  } else {
    socket.send(
      JSON.stringify({
        type: "notification",
        Content: {
          Sender: message.Sender,
          Receiver: message.Receiver,
        },
      })
    );
  }
}

export async function renderChatEvents() {
  const message_Box = document.querySelector(".ChatBody");
  message_Box.addEventListener("scroll", function (e) {
    if (e.target.scrollTop === 0) {
      let scroll = parseInt(sessionStorage.getItem("scroll")) + 10;
      sessionStorage.setItem("scroll", scroll);
      throttle(LoadMessages(), 3000);
    }
  });

  const sendChatMessage = async () => {
    let message = document.getElementById("ChatMessage").value;
    if (message.trim() === "") {
      alert("Empty messages are not allowed!");
      return;
    }
    document.getElementById("ChatMessage").value = "";
    type.TypingOff(socket)

    socket.send(
      JSON.stringify({
        type: "send_message",
        Content: {
          Text: message,
          Sender: parseInt(sessionStorage.getItem("UserId")),
          Receiver: parseInt(sessionStorage.getItem("User")),
        },
      })
    );
  };

  document
    .querySelector(".SendChatMessage")
    .addEventListener("click", sendChatMessage);

  document.getElementById("ChatMessage").addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      e.preventDefault();
      sendChatMessage();
    }
  });

  document.querySelector(".MsgInput").addEventListener("keydown", (event) => {
    if ((event.key === "Backspace") || (event.key === "Enter")) {
      return;
    }
    console.log("Event worked")
    if (socket.readyState === socket.OPEN ) {
      throttle(type.TypingOn(socket), 1500);
    }
  });

  document.querySelector(".MsgInput").addEventListener("input", (event) => {
    if (socket.readyState === socket.OPEN) {
      if (event.target.value.trim() === "") {
        throttle(type.TypingOff(socket), 2000);
      }
    }
  });
}

export function LoadMessages() {
  const currentUrl = new URL(window.location.href);

  // Get the query parameters
  const searchParams = currentUrl.searchParams;

  // Check if the 'Sender' and 'Receiver' parameters exist
  if (searchParams.has("Sender") && searchParams.has("Receiver")) {
    socket.send(
      JSON.stringify({
        type: "load_messages",
        Content: {
          sender: parseInt(searchParams.get("Sender")),
          receiver: parseInt(searchParams.get("Receiver")),
          scroll: parseInt(sessionStorage.getItem("scroll")),
        },
      })
    );
  }
}

function createMessageElement(message) {
  const messageElement = document.createElement("div");
  messageElement.classList.add("messages");
  const username = document.createElement("div");

  if (message.Sender === parseInt(sessionStorage.getItem("UserId"))) {
    messageElement.classList.add("sender");
    username.textContent = sessionStorage.getItem("Username");
  } else {
    messageElement.classList.add("receiver");
    username.textContent = sessionStorage.getItem("partner-Username");
  }

  const messageBody = document.createElement("div");
  messageBody.classList.add("message-body");
  messageBody.textContent = message.Text;

  const messageTime = document.createElement("div");
  messageTime.classList.add("message-time");
  messageTime.textContent = formatTime(message.Time);

  messageElement.appendChild(username);
  messageElement.appendChild(messageBody);
  messageElement.appendChild(messageTime);

  return messageElement;
}


function renderMessages(messages) {
  // Check if ChatBody exists immediately
  if (document.querySelector(".ChatBody")) {
    renderMessagesInternal(messages); // Render if it exists
  } else {
    // If not, wait for 2 seconds and then check again
    setTimeout(() => {
      if (document.querySelector(".ChatBody")) {
        renderMessagesInternal(messages); // Render if it exists after 2 seconds
      }
    }, 500);
  }
    // Check if the scroll position is stored in session storage
    const storedScrollPosition = parseInt(sessionStorage.getItem("scroll"));

    // If the stored scroll position is 0 or not found, scroll to the bottom
    if (isNaN(storedScrollPosition) || storedScrollPosition === 0) {
      const chatBody = document.querySelector(".ChatBody");
      chatBody.scrollTop = chatBody.scrollHeight;
    }
}


// Helper function to render messages
function renderMessagesInternal(messages) {
  const messagesContainer = document.querySelector(".ChatBody");
  messagesContainer.innerHTML = "";
  if (messages != null) {
    messages
      .sort((a, b) => a.Id - b.Id)
      .forEach((message) => {
        const messageElement = createMessageElement(message);
        messagesContainer.appendChild(messageElement);
      });
  }
}

// Function to format the time
function formatTime(timeString) {
  const date = new Date(timeString);
  const formattedTime = date.toLocaleString().replace(",", "");
  return formattedTime;
}

function throttle(func, delay) {
  let lastCall = 0;
  return function (...args) {
    const now = new Date().getTime();
    if (now - lastCall >= delay) {
      lastCall = now;
      return func.apply(this, args);
    }
  };
}
