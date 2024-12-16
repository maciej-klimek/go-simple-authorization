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
  console.log(csrfToken);
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
