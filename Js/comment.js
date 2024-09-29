import * as tem from "./template.js";
import * as post from "./post.js";

export async function Comment(Id) {
  var postId = parseInt(Id)
  let userName = document.querySelector(".Username").textContent;
  let profileImage = document
    .querySelector(".ProfilePhoto")
    .getAttribute("src")
    .split("/")[2];
  var commentId;
  let CommentMessage = document.getElementById(`CommentMessage${postId}`);
  var message = CommentMessage.value;

  if (message.length > 200) {
    alert("Comment cannot be more than 200 characters.");
    return;
  } else if (message.trim().length === 0) {
    alert("Empty Comments are not allowed");
  } else {
    fetch("/CommentHandler", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        message: message,
        postId: parseInt(postId),
        username: userName,
        profileImage: profileImage,
      }),
    })
      .then(function (response) {
        if (response.ok) {
          CommentMessage.value = ""; // Reset the input field
          // Handle the response from the backend if needed
          return response.json(); // Parse the response body as JSON
        } else {
          throw new Error("Error sending message.");
        }
      })
      .then(function (data) {
        commentId = data.CommentId;
        const comment = {
          CommentId: commentId,
          UserName: userName,
          Text: message,
          Time: new Date().toLocaleString(),
          CLike: 0,
          CDislike: 0,
          ProfileImage: profileImage,
        };
        var com = tem.commentTemplate(comment);

        var container = document.getElementById(`Comment${postId}`);
        
        // Directly append the HTML content of 'com' to 'container'
        container.innerHTML += com;

         // Reset the input field *after* appending the comment
         CommentMessage.value = ""; 
             
      document.querySelectorAll(".CLike").forEach(function(button) {

        button.addEventListener("click", function() {HandleCommentLike(this, 'Like', postId);});
      });
      document.querySelectorAll(".CDislike").forEach(function(button) {
        button.addEventListener("click",function(){HandleCommentLike(this,'Dislike', postId)});
      });

      })
      .catch(function (error) {
        console.error("Error sending message:", error);
      });
      
      const path = window.location.pathname;
      const segments = path.split("/");
      const data = await post.GetPostData(segments[2]);
     LikedComments(data.LikedComments)
      CreatedComments(data.CreatedComments)


  }
}

export function HandleCommentLike(clickedCheckbox, checkboxName, postId) {

  const checkboxValue = clickedCheckbox.value;
  const isChecked = clickedCheckbox.checked;
  if (
    checkboxName === "Dislike" &&
    document.getElementById(`CLike${checkboxValue}`).checked === true
  ) {
    clickedCheckbox.checked = false;
  } else if (
    checkboxName === "Like" &&
    document.getElementById(`CDislike${checkboxValue}`).checked === true
  ) {
    clickedCheckbox.checked = false;
  } else {
    fetch("/CommentLikeHandle", {
      method: "POST",
      body: JSON.stringify({
        RequestType: "like",
        Type: checkboxName,
        ID: parseInt(checkboxValue, 10),
        PostId: postId,
        Checked: isChecked,
      }),
    })
      .then((response) => {
        console.log(response);
      })
      .catch((error) => console.error("Error:", error));
    //Store Value
    const id = clickedCheckbox.getAttribute("id");
    const LikeNum = document.getElementById(`Num${id}`);

    // Get the current text content and parse it as an integer
    let currentValue = parseInt(LikeNum.textContent);
    if (isChecked === true) {
      // Set the updated value as the new text content of the element
      LikeNum.textContent = (currentValue + 1).toString();
    } else {
      LikeNum.textContent = (currentValue - 1).toString();
    }
    document
      .getElementById(`C${checkboxName}SVG${checkboxValue}`)
      .classList.toggle(`C${checkboxName}SVGClicked`);
  }
}

export function LikedComments(likedComments) {
  if (likedComments !== null) {
    likedComments.forEach(function (like) {
      let value = like.split("_");
      let checkbox = document.getElementById(value.join(""));
      if (checkbox !== null) {
        checkbox.checked = true;
        document
          .getElementById(`${value[0]}SVG${value[1]}`)
          .classList.add(`${value[0]}SVGClicked`);
      }
    });
  }
}

export function CreatedComments(createdComments) {
  if (createdComments !== null) {
    createdComments.forEach(function (CommentId) {
      let CommentBtnsContainer = document.getElementById(
        `CommentEditDeleteContainer${CommentId}`
      );
      if (CommentBtnsContainer !== null) {
        CommentBtnsContainer.style.display = "flex";
      }
    });
  }
}

export function handleEditButtonClick(editButton) {
  var commentId = editButton.value;
  var commentText = document.getElementById(`CommentText${commentId}`);

  var inputField = document.createElement("input");
  inputField.type = "text";
  inputField.value = commentText.textContent;
  inputField.className = "CommentEditTextInput";

  inputField.addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      setNewValue();
    }
  });

  editButton.disabled = true;

  function setNewValue() {
    var newValue = inputField.value;

    commentText.textContent = newValue;
    commentText.style.display = "inline"; // Show the updated comment text

    inputField.parentNode.removeChild(inputField); // Remove the input field

    UpdateComment(commentId, newValue);
    editButton.disabled = false;
  }

  commentText.style.display = "none"; // Hide the comment text
  commentText.parentNode.insertBefore(inputField, commentText); // Insert the input field before the comment text
  inputField.focus(); // Set foc
}

function UpdateComment(commentId, TextContent) {
  fetch("/Post/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      RequestType: "UpdateComment",
      ID: parseInt(commentId, 10),
      Text: TextContent,
    }),
  })
    .then((response) => {
      console.log(response);
      // Handle the response for the clicked checkbox
    })
    .catch((error) => console.error("Error:", error));
}

export function HandleCommentDelete(CommentId) {
  fetch("/Post/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      RequestType: "DeleteComment",
      ID: parseInt(CommentId, 10),
    }),
  })
    .then((response) => {
      console.log(response);
      // Handle the response for the clicked checkbox
    })
    .catch((error) => console.error("Error:", error));
  location.reload();
}
