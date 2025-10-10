document.getElementById("signinForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    const email = document.getElementById("email").value;
    const password = document.getElementById("password").value;
    const errorMsg = document.getElementById("errorMsg");

    try {
        const response = await axios.post("http://localhost:8080/signin", {
            email,
            password
        });

        // Assuming backend returns { access: "...", refresh: "..." }
        const { access, refresh } = response.data;

        // Save tokens in localStorage
        localStorage.setItem("accessToken", access);
        localStorage.setItem("refreshToken", refresh);

        // Redirect to dashboard or home page
        window.location.href = "/dashboard.html";
    } catch (err) {
        console.error(err);
        if (err.response && err.response.data) {
            errorMsg.textContent = err.response.data.err || "Invalid credentials";
        } else {
            errorMsg.textContent = "Server error. Try again later.";
        }
    }
});
