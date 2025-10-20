document.getElementById("signinForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const errorMsg = document.getElementById("errorMsg");

    // Clear previous error
    errorMsg.textContent = "";

    if (!email || !password) {
        errorMsg.textContent = "Please enter both email and password.";
        return;
    }

    try {
        const response = await axios.post("http://localhost:8080/signin", { email, password });

        console.log("Signin Response:", response.data); // Debug backend response

        // Destructure response safely
        const access = response.data.access;
        const refresh = response.data.refresh;
        const role = response.data.role ? response.data.role.toUpperCase() : null;

        if (!access || !refresh || !role) {
            errorMsg.textContent = "Login failed: tokens or role not received.";
            return;
        }

        // Store tokens and role
        localStorage.setItem("accessToken", access);
        localStorage.setItem("refreshToken", refresh);
        localStorage.setItem("role", role);

        // Redirect based on role
        switch (role) {
            case "SUPER_ADMIN":
                window.location.href = "admindash.html";
                return;
            case "MANAGER":
                window.location.href = "managerdash.html";
                return;
            case "STAFF":
                window.location.href = "staffdash.html";
                return;
            default:
                errorMsg.textContent = "Unknown role. Contact system admin.";
                return;
        }

    } catch (err) {
        console.error("Signin Error:", err);

        // Better error handling for various backend responses
        if (err.response) {
            // Backend returned an error
            const data = err.response.data;
            errorMsg.textContent = data.err || data.res || "Login failed";
        } else if (err.request) {
            // Request made but no response
            errorMsg.textContent = "No response from server. Try again later.";
        } else {
            // Other errors
            errorMsg.textContent = "An unexpected error occurred.";
        }
    }
});
