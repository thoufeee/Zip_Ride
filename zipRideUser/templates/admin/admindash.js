// admindash.js

// Base API URL (update with your actual API endpoint)
const API_BASE_URL = "http://localhost:8080/api/admin";

// Handle logout
document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
  window.location.href = "signin.html";
});

// Fetch dashboard data
async function loadDashboardData() {
  const token = localStorage.getItem("accessToken");
  if (!token) {
    window.location.href = "signin.html";
    return;
  }

  try {
    const response = await axios.get(`${API_BASE_URL}/dashboard`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    const { totalUsers, totalLatestUsers, totalBookings, recentBookings } = response.data;

    document.getElementById("user-count").textContent = totalUsers || 0;
    document.getElementById("staff-count").textContent = totalLatestUsers || 0;
    document.getElementById("manager-count").textContent = totalBookings || 0;

    const bookingTable = document.getElementById("booking-history-body");
    bookingTable.innerHTML = "";

    recentBookings.forEach(b => {
      const row = `
        <tr class="hover:bg-gray-50 transition duration-150">
          <td class="px-6 py-4 font-medium text-gray-900">#${b.id}</td>
          <td class="px-6 py-4 text-gray-600">${b.userName}</td>
          <td class="px-6 py-4 text-gray-600">${b.service}</td>
          <td class="px-6 py-4">
            <span class="px-3 py-1 text-xs font-semibold rounded-full ${
              b.status === "Completed"
                ? "bg-green-100 text-green-800"
                : b.status === "Pending"
                ? "bg-yellow-100 text-yellow-800"
                : "bg-blue-100 text-blue-800"
            }">${b.status}</span>
          </td>
          <td class="px-6 py-4 text-gray-500">${b.date}</td>
        </tr>`;
      bookingTable.insertAdjacentHTML("beforeend", row);
    });

  } catch (error) {
    console.error("Failed to load dashboard:", error);
  }
}

// Run on load
document.addEventListener("DOMContentLoaded", loadDashboardData);
