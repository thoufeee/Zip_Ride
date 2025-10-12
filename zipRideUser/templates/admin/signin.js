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

        // Destructure access, refresh, and role
        const { access, refresh, role } = response.data;

        if (access && refresh && role) {
            // Store tokens and role
            localStorage.setItem("accessToken", access);
            localStorage.setItem("refreshToken", refresh);
            localStorage.setItem("role", role);

            // Redirect based on role
            if (role === "SUPER_ADMIN") {
                window.location.href = "admindash.html";
            } else if (role === "MANAGER") {
                window.location.href = "managerdash.html";
            } else if (role === "STAFF") {
                window.location.href = "staffdash.html";
            } else {
                errorMsg.textContent = "Unknown role. Contact system admin.";
            }
        } else {
            errorMsg.textContent = "Login failed: tokens or role not received.";
        }
    } catch (err) {
        console.error(err);

        if (err.response && err.response.data && err.response.data.err) {
            errorMsg.textContent = err.response.data.err;
        } else if (err.request) {
            errorMsg.textContent = "No response from server. Try again later.";
        } else {
            errorMsg.textContent = "An unexpected error occurred.";
        }
    }
});
