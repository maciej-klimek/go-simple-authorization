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
  console.log("CSRF Token: ", csrfToken);
  return csrfToken;
}

function handleFileUpload(event) {
  event.preventDefault();

  const formData = new FormData(event.target);
  const csrfToken = getCsrfToken();

  fetch("/content", {
    method: "POST",
    headers: {
      "X-CSRF-Token": csrfToken,
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

document.getElementById("fileUploadForm").addEventListener("submit", handleFileUpload);

function logout() {
  const csrfToken = getCsrfToken();

  fetch("/logout", {
    method: "POST",
    headers: {
      "X-CSRF-Token": csrfToken,
    },
  })
    .then((response) => {
      if (response.ok) {
        console.log("Logged out successfully");
        window.location.href = "/login";
      } else {
        console.error("Logout failed");
      }
    })
    .catch((error) => {
      console.error("Error during logout:", error);
    });
}

document.getElementById("logoutButton").addEventListener("click", logout);

// Add CSRF Token for File Viewing
function handleFileClick(event) {
  event.preventDefault(); // Prevent the default action (i.e., direct navigation)

  const fileUrl = event.target.getAttribute("href");
  const csrfToken = getCsrfToken();

  // Create a request with the CSRF token to validate the session
  fetch(fileUrl, {
    method: "GET",
    headers: {
      "X-CSRF-Token": csrfToken, // Include CSRF token for validation
    },
  })
    .then((response) => {
      if (response.ok) {
        // If the response is OK, just let the browser handle the file directly
        window.open(fileUrl, "_blank");
      } else {
        console.error("Failed to load the file");
        return response.text();
      }
    })
    .catch((error) => {
      console.error("Error while loading the file:", error);
    });
}

// Attach event listeners to all file links
document.querySelectorAll(".file-link").forEach((link) => {
  link.addEventListener("click", handleFileClick);
});
