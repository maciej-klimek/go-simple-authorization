// content.js

// Function to extract CSRF token from cookies
function getCsrfToken() {
  const csrfToken =
    document.cookie
      .split("; ")
      .find((row) => row.startsWith("csrf_token="))
      ?.split("=")[1] + "=";

  if (!csrfToken) {
    console.error("CSRF token not found in cookies.");
    throw new Error("CSRF token missing.");
  }

  return csrfToken;
}

function handleFileUpload(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const csrfToken = getCsrfToken(); // Get CSRF token

  // Prepare the request with CSRF token in the header
  fetch("/content", {
    method: "POST",
    headers: {
      "X-CSRF-Token": csrfToken, // Add CSRF token to the headers
    },
    body: formData,
  })
    .then((response) => {
      if (response.ok) {
        console.log("File upload successful");
        return response.text();
      } else {
        console.error("File upload failed");
        return response.text();
      }
    })
    .then((text) => {
      console.log("Response: ", text);
    })
    .catch((error) => {
      console.error("Error during file upload:", error);
    });
}

// Attach the event listener to the file upload form submit
document.getElementById("fileUploadForm").onsubmit = handleFileUpload;

// Function for logging out
function logout() {
  const csrfToken = getCsrfToken(); // Get CSRF token for logout

  // Call the /logout handler to handle the logout on the backend
  fetch("/logout", {
    method: "POST", // or GET based on your backend implementation
    headers: {
      "X-CSRF-Token": csrfToken, // Add CSRF token to the headers
    },
  })
    .then((response) => {
      if (response.ok) {
        console.log("Logged out successfully");
        window.location.href = "/login"; // Redirect to login after successful logout
      } else {
        console.error("Logout failed");
      }
    })
    .catch((error) => {
      console.error("Error during logout:", error);
    });
}

// Attach the event listener to the logout button
document.getElementById("logoutButton").onclick = logout;
