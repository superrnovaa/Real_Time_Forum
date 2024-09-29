import * as temp from "./template.js";

export async function renderCreatPostTemplate() {
  const template = temp.CreatePostTemplate();
  const style = `<link href="../Style/CreatePost.css" rel="stylesheet" />`;
  document.querySelector(".page").innerHTML = style + template;
  RenderCreateEvents();
}

function  RenderCreateEvents() {
  // Select the form element
  const form = document.getElementById("PostForm");

  // Add a submit event listener to the form
  form.addEventListener("submit", function (event) {
    const checkboxes = document.querySelectorAll(".customCheckBoxInput");

    // Check if at least one checkbox is checked
    const isChecked = Array.from(checkboxes).some(
      (checkbox) => checkbox.checked
    );

    // If no checkbox is checked, prevent form submission and show an error message
    if (!isChecked) {
      event.preventDefault();
      alert("Please select at least one checkbox.");
    }
  });

  const input = document.getElementById("file");
  input.addEventListener("change", (e) => handleFileInputChange(e.target));

  const DownloadButton = document.querySelector(".Download-button");
  DownloadButton.addEventListener("click", async (e) => {
    if (SubmitCheck()) {
      e.preventDefault();
      // Get the form data
  const formData = new FormData(form);

  try {

    // Send the form data to the server using fetch
    const response = await fetch('/CreatePost', {
      method: 'POST',
      body: formData,
    });

    if (response.ok) {
      console.log('Form data submitted successfully!');
    } else {
      console.error('Error submitting form data:', response.status);
    }
  } catch (error) {
    console.error('Error submitting form data:', error);
  }

  ResetState()
    }
  });
  var ClearButton = document.querySelector(".Clear");
  ClearButton.addEventListener("click", ClearImage);
  
}

function ResetState() {
  var ClearButton = document.querySelector(".Clear");
  ClearButton.click();
  var Title = document.querySelector(".Title");
  var TextArea = document.querySelector(".TextArea");
  Title.value= "";
  TextArea.value = "";

  const checkboxes = document.querySelectorAll(".customCheckBoxInput");
  checkboxes.forEach(function (checkbox) {

      checkbox.checked = false;
  });
}

function ClearImage() {

  var input = document.getElementById("file");
  input.value = "";
  var label = input.parentNode;
  var image = label.querySelector(".selected-image");
  if (image) {
    label.removeChild(image);
  }
}

export function handleFileInputChange(input) {
  var file = input.files[0];
  var allowedTypes = ["image/jpeg", "image/png", "image/svg+xml", "image/gif"];
  var maxSize = 20 * 1024 * 1024;

  if (file && allowedTypes.includes(file.type)) {
    if (file.size > maxSize) {
      alert("Please select an image file smaller than 20MB.");
      input.value = ""; // Clear the file input value
      return;
    }
    var reader = new FileReader();

    reader.onload = function (e) {
      var label = input.parentNode;
      var image = document.createElement("div");
      image.className = "selected-image";
      image.style.backgroundImage = "url(" + e.target.result + ")";
      label.appendChild(image);
    };

    reader.readAsDataURL(file);
  } else {
    // File type is not allowed, handle accordingly (e.g., display an error message)
    alert("Please select a valid image file (JPEG, PNG, SVG or GIF).");
    input.value = ""; // Clear the file input value
  }
}

function SubmitCheck() {
  // Get the values of the input and textarea
  var titleValue = document.querySelector(".Title").value;
  var contentValue = document.querySelector(".TextArea").value;

  // Check if any of the values is empty
  if (titleValue.trim() === "" || contentValue.trim() === "") {
    // Display an alert message
    alert("Please fill in all fields");
    // Prevent form submission
    return false;
  }

  // Allow form submission
  return true;
}
