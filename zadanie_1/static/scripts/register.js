// register.js (for register.html)
document.getElementById("registerForm").onsubmit = function (event) {
  event.preventDefault();
  fetch("/register", {
    method: "POST",
    body: new FormData(event.target),
  })
    .then((response) => {
      if (response.ok) {
        console.log("Registration successful. Redirecting to login...");
        window.location.href = "/login";
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
