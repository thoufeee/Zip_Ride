document.getElementById("signinForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    const email = document.getElementById("email").value.trim();
    const password = document.getElementById("password").value.trim();
    const errorMsg = document.getElementById("errorMsg");

    errorMsg.textContent = "";

    if (!email || !password) {
        errorMsg.textContent = "Please enter both email and password.";
        return;
    }

    try {
        const response = await axios.post("http://localhost:8080/signin", { email, password });

        console.log("Signin Response:", response.data);

        const access = response.data.access;
        const refresh = response.data.refresh;
        const role = response.data.role ? response.data.role.toUpperCase() : null;
        const permissions = response.data.permissions || []; 
        if (!access || !refresh) {
            errorMsg.textContent = "Login failed: tokens not received.";
            return;
        }

      
        localStorage.setItem("accessToken", access);
        localStorage.setItem("refreshToken", refresh);
        localStorage.setItem("role", role || ""); 
        localStorage.setItem("permissions", JSON.stringify(permissions));

        window.location.href = "admindash.html";

    } catch (err) {
        console.error("Signin Error:", err);

        if (err.response) {
            const data = err.response.data;
            errorMsg.textContent = data.err || data.res || "Login failed";
        } else if (err.request) {
            errorMsg.textContent = "No response from server. Try again later.";
        } else {
            errorMsg.textContent = "An unexpected error occurred.";
        }
    }
});
