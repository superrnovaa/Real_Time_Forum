import * as temp from "./template.js";
import * as app from "./app.js";
import * as postFile from "./post.js";

export async function renderHomePageTemplate(x) {
  const template = temp.HomePageTemplate();
  const style = `<link href="../Style/Home.css" rel="stylesheet" />`;
  document.querySelector(".page").innerHTML = style + template;
  const data = await app.GetHomeData(x);
  const posts = temp.PostsTemplate(data.Posts);
  document.querySelector(".Posts").innerHTML = posts;
  RenderHomeEvents(data.Posts, data);
}

function RenderHomeEvents(posts, data) {
  if (posts){
    for (let post of posts) {
      document
        .getElementById(`post${post.PostId}`)
        .addEventListener("click", function (e) {
          e.preventDefault()
          window.history.pushState(null, '', `/Post/${post.PostId}`);
          postFile.renderPostTemplate();
          //window.location.href = `/Post/${post.PostId}`;
        });
    }
  }
  LikedPosts(data.LikedPosts);
  Filter();
}

function LikedPosts(likedPosts) {
  if (likedPosts !== null) {
    likedPosts.forEach(function (like) {
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

export function ResetFilters() {
  var checkboxes = document.querySelectorAll('input[name="categories"]');
  checkboxes.forEach(function (checkbox) {
   localStorage.setItem(checkbox.value, false);
  });
}


function Filter() {
  // Initialize the checkboxes based on the stored values
  var form = document.getElementById("Filter");
  var checkboxes = document.querySelectorAll('input[name="categories"]');
  checkboxes.forEach(function (checkbox) {
    var storedValue = localStorage.getItem(checkbox.value);
    if (storedValue === "true") {
      checkbox.checked = true;
    }
  });

  // Add event listener to checkboxes change
  checkboxes.forEach(function (checkbox) {
    checkbox.addEventListener("change", function () {
      updateStoredValues(checkboxes);
      submitForm(form, checkboxes);
    });
  });
}

function updateStoredValues(checkboxes) {
  checkboxes.forEach(function (checkbox) {
    localStorage.setItem(checkbox.value, checkbox.checked);
  });
}

function submitForm(form, checkboxes) {
  var categories = [];

  checkboxes.forEach(function (checkbox) {
    if (checkbox.checked) {
      categories.push(checkbox.value);
    }
  });

  var categoriesString = categories.join(" ");

  var hiddenInput = form.querySelector('input[name="categories"]');
  if (!hiddenInput) {
    hiddenInput = document.createElement("input");
    hiddenInput.type = "hidden";
    hiddenInput.name = "categories";
    form.appendChild(hiddenInput);
  }

  hiddenInput.value = categoriesString;
  renderHomePageTemplate(categoriesString)

  //form.submit();
}
