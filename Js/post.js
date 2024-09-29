import * as temp from "./template.js";
import * as comm from "./comment.js";

export async function renderPostTemplate() {
  const container = `<div class="SinglePostContainer"></div>`;
  const path = window.location.pathname;
  const segments = path.split("/");
 const PostId = segments[2]
  const data = await GetPostData(segments[2]);
  const post = temp.PostsTemplate([data.Post]);
  const comments = temp.CommentsTemplate(data.Post);
  const style = `<link href="../Style/Home.css" rel="stylesheet" />`;
  document.querySelector(".page").innerHTML = style + container;
  document.querySelector(".SinglePostContainer").innerHTML = post + comments;
  RenderPostEvents(PostId, data);
}

function RenderPostEvents(postId, data) {
  document.getElementById(`Like${postId}`).addEventListener("click",function(){HandleLike(this,'Like')})
  document.getElementById(`Dislike${postId}`).addEventListener("click",function(){HandleLike(this,'Dislike')})
  LikedPosts(data.IsLiked,data.IsDisliked,postId)

  comm.LikedComments(data.LikedComments)
  comm.CreatedComments(data.CreatedComments)

  document.querySelector(".SendSVG").addEventListener("click", function () {

    comm.Comment(postId);
  });
  document.getElementById(`CommentMessage${postId}`).addEventListener("keydown", function (event) {
    if (event.key === "Enter") {
      comm.Comment(postId);
    }
  });

  document.querySelectorAll(".CLike").forEach(function(button) {
    button.addEventListener("click", function() {comm.HandleCommentLike(this, 'Like', parseInt(postId));});
  });
  document.querySelectorAll(".CDislike").forEach(function(button) {
    button.addEventListener("click",function(){comm.HandleCommentLike(this,'Dislike', parseInt(postId))});
  });

 
  /*
  document.querySelectorAll(".CommentEditBtn").forEach(function(button) {
    button.addEventListener("click",function(){comm.handleEditButtonClick(this)});
  });
  document.querySelectorAll(".CommentDeleteBtn").forEach(function(button) {
    button.addEventListener("click",function(){comm.HandleCommentDelete(this.value)});
  });

  if (data.IsCreated) {
    showDeleteButton(postId);
  }
    */
}

export async function GetPostData(PostId) {
  try {
    const response = await fetch(`/Post/${PostId}`, {
      method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({RequestType: "Data",}),
  });
    const data = await response.json();
    return data;
  } catch (error) {
    console.error(error);
  }
}

function HandleLike(clickedCheckbox, checkboxName) {
    const checkboxValue = clickedCheckbox.value;
    const isChecked = clickedCheckbox.checked;
    if (checkboxName === 'Dislike' && document.getElementById(`Like${checkboxValue}`).checked === true) {
      clickedCheckbox.checked = false
    } else if (checkboxName === 'Like' && document.getElementById(`Dislike${checkboxValue}`).checked === true) {
      clickedCheckbox.checked = false
    } else {
      fetch("/Post/", {
        method: "POST",
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          RequestType: "like",
          Type: checkboxName,
          ID: parseInt(checkboxValue, 10),
          Checked: isChecked,
        }),
      })
        .then((response) => {
          console.log(response);
          // Handle the response for the clicked checkbox
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
      document.getElementById(`${checkboxName}SVG${checkboxValue}`).classList.toggle(`${checkboxName}SVGClicked`)
    }
}

function LikedPosts(IsLiked, IsDisliked,postId) {
  if (IsLiked) {
    document.getElementById(`Like${postId}`).checked = true
    document.getElementById(`LikeSVG${postId}`).classList.add(`LikeSVGClicked`)
  }
  if (IsDisliked){
    document.getElementById(`Dislike${postId}`).checked = true
    document.getElementById(`DislikeSVG${postId}`).classList.add("DislikeSVGClicked")
  }
}



function showDeleteButton(PostId) {
  let DeleteBtn = document.getElementById(`DeleteBtn${PostId}`);
  if (DeleteBtn !== null) {
    DeleteBtn.style.display = "flex";
  }
  DeleteBtn.addEventListener("click", function () {
    HandleDelete(PostId);
  });
}

function HandleDelete(PostId) {
  fetch("/Post/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({
      RequestType: "delete",
      ID: PostId,
    }),
  })
    .then((response) => {
      console.log(response);
      // Handle the response for the clicked checkbox
    })
    .catch((error) => console.error("Error:", error));
  window.location.href = "/HomePage";
}

