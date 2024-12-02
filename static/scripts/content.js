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

// File upload functionality with CSRF token
document.getElementById("fileUploadForm").onsubmit = function (event) {
  event.preventDefault();

  // Get the CSRF token from the cookie
  const csrfToken = getCookie("csrf_token");

  // Create a FormData object
  const formData = new FormData();
  const fileInput = document.getElementById("fileInput");
  formData.append("file", fileInput.files[0]);

  // Send the file and CSRF token in the request
  fetch("/content", {
    method: "POST",
    headers: {
      "X-CSRF-Token": csrfToken, // Add CSRF token to headers
    },
    body: formData,
  })
    .then((response) => {
      if (response.ok) {
        console.log("File uploaded successfully");
        window.location.href = "/content"; // Redirect or show success message
      } else {
        console.error("Error during file upload");
      }
    })
    .catch((error) => {
      console.error("Error:", error);
    });
};

// Utility function to get a cookie by name
function getCookie(name) {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);
  if (parts.length === 2) return parts.pop().split(";").shift();
}

// Add event listener for the alert button
document.getElementById("alertButton").onclick = function () {
  alert("AAA");
};
