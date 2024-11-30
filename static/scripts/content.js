// content.js (for content.html)
document.getElementById("logoutButton").onclick = function () {
  fetch("/logout", {
    method: "POST",
  }).then((response) => {
    if (response.ok) {
      console.log("Logout successful");
      window.location.href = "/login"; // Redirect to login page after logout
    } else {
      console.error("Error during logout");
    }
  });
};
