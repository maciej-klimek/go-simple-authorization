// login.js (for login.html)

alert("witam, login");
document.getElementById("loginForm").onsubmit = function (event) {
  event.preventDefault();
  fetch("/login", {
    method: "POST",
    body: new FormData(event.target),
  })
    .then((response) => {
      if (response.ok) {
        window.location.href = "/content";
      } else {
        return response.text();
      }
    })
    .then((text) => {
      if (text) {
        document.getElementById("errorMessage").innerText = text;
      }
    });
};
