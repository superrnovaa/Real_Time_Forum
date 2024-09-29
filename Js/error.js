let code, message;

export async function renderErrorPage( code,message ) {
  console.log("Rendering error page...");

  try {

    const template = ErrorElement(code, message);
    const style = `<link href="../Style/Main.css" rel="stylesheet" />`;
    document.querySelector(".content").innerHTML = style + template;

    document.querySelector(".back-button").addEventListener("click", function () {
      window.location.href = "/HomePage";
    });
  } catch (err) {
    console.error("Error fetching error page:", err);
    // Assuming a default error code and message
    const template = ErrorElement("Page Not Found");
    const style = `<link href="../Style/Main.css" rel="stylesheet" />`;
    document.querySelector(".content").innerHTML = style + template;
    document.querySelector(".back-button").addEventListener("click", function () {
      window.location.href = "/HomePage";
    });
  }
}

function ErrorElement(code, message) {
  return `
    <div class="Error-container">
      <h1>Error ${code}</h1>
      <p>${message}</p>
      <button class="back-button">Back</button>
    </div>
  `;
}

export async function GetError() {
  try {
    const response = await fetch("/Error");
      const { Code, Message } = await response.json();
      code = Code;
    message = Message;
    renderErrorPage(code,message)
  } catch (err) {
    console.error("Error fetching error page:", err);

  }
}