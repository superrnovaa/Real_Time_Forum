let isTyping = false;
let typingTimeout = null;

export function Typing(content) {
  const typingContainer = document.querySelector(
    `#Typing-container${content.sender}`
  );
  if (content.type === "start" && !isTyping) {
    isTyping = true;
    typingContainer.style.display = "flex";
    const element = typingContainer.querySelector("p");
    element.textContent = "";
    TextTypingEffect(element, "••••");
  } else if (content.type === "stop") {
    isTyping = false;
    typingContainer.style.display = "none";
    clearTimeout(typingTimeout);
  }

  if (content.type === "start") {
    clearTimeout(typingTimeout);
    typingTimeout = setTimeout(() => {
      isTyping = false;
      typingContainer.style.display = "none";
    }, 2000); // Set timeout for 4 seconds
  }
}
export function TypingOn(socket) {
  isTyping = true;
  socket.send(
    JSON.stringify({
      type: "typing",
      Content: {
        receiver: parseInt(sessionStorage.getItem("User")),
        sender: parseInt(sessionStorage.getItem("UserId")),
        type: "start",
      },
    })
  );
}
export function TypingOff(socket) {
  isTyping = false;
  socket.send(
    JSON.stringify({
      type: "typing",
      Content: {
        receiver: parseInt(sessionStorage.getItem("User")),
        sender: parseInt(sessionStorage.getItem("UserId")),
        type: "stop",
      },
    })
  );
}

function TextTypingEffect(element, text, i = 0) {
  element.textContent += text[i];
  // rest of typing logic
  if (i === text.length - 1) {
    element.textContent = "";
    TextTypingEffect(element, "••••");
  } else {
    if (isTyping) {
      setTimeout(() => {
        TextTypingEffect(element, text, i + 1);
      }, 300);
    }
  }
}

