const API_BASE_URL = "http://localhost:8080"; 

// ✅ Check permission
const userPermissions = JSON.parse(localStorage.getItem("permissions") || "[]");
if (!userPermissions.includes("ACCESS_ADMINDASH")) {
  document.body.innerHTML = `
     <div class="flex flex-col items-center justify-center min-h-screen bg-gray-50 text-center">
      <h1 class="text-3xl font-bold text-red-600 mb-3">Access Denied</h1>
      <p class="text-gray-600 mb-6">You don't have permission to view this page.</p>
      <button onclick="window.location.href='admindash.html'" 
        class="bg-cyan-600 hover:bg-cyan-700 text-white px-5 py-2.5 rounded-xl font-semibold shadow-md">
        Go to Dashboard
      </button>
    </div>
  `;
  throw new Error("Unauthorized access to Admin Dashboard");
}

// ✅ Helper: fetch with Authorization header
async function fetchWithAuth(url) {
  const token = localStorage.getItem("accessToken");
  if (!token) {
    console.warn("No token found — redirecting to signin");
    window.location.href = "signin.html";
    throw new Error("Access token missing");
  }

  const response = await fetch(url, {
    headers: { Authorization: `Bearer ${token}` },
  });

  if (!response.ok) {
    throw new Error(`Request failed: ${response.status}`);
  }

  return await response.json();
}

async function loadDashboardData() {
  try {
    console.log("Loading dashboard data...");

    const [usersRes, latestUsersRes] = await Promise.all([
      fetchWithAuth(`${API_BASE_URL}/admin/alluserslength`),
      fetchWithAuth(`${API_BASE_URL}/admin/latestuserslength`),
    ]);
    
    const totalUsers = usersRes.res ?? usersRes.count ?? 0;
    const totalLatestUsers = latestUsersRes.res ?? latestUsersRes.count ?? 0;

    // Mock data for bookings
    const totalDailyBookings = 5; 
    const recentBookings = [
      {
        ID: 101,
        UserName: "Alice",
        DriverName: "Bob",
        Pickup: "Location A",
        Dropoff: "Location B",
        VehicleType: "Car",
        Status: "Completed",
        CreatedAt: new Date().toISOString(),
      },
      {
        ID: 102,
        UserName: "Charlie",
        DriverName: null,
        Pickup: "Location C",
        Dropoff: "Location D",
        VehicleType: "Bike",
        Status: "Pending",
        CreatedAt: new Date().toISOString(),
      },
    ];
    
    document.getElementById("user-count").textContent = totalUsers;
    document.getElementById("staff-count").textContent = totalLatestUsers;
    document.getElementById("manager-count").textContent = totalDailyBookings;

    const tbody = document.getElementById("booking-history-body");
    tbody.innerHTML = "";

    recentBookings.forEach((b) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td class="px-6 py-4 whitespace-nowrap">#${b.ID || "-"}</td>
        <td class="px-6 py-4">${b.UserName || "-"}</td>
        <td class="px-6 py-4">${b.DriverName || '<span class="text-red-500">Unassigned</span>'}</td>
        <td class="px-6 py-4">${b.Pickup || "-"}</td>
        <td class="px-6 py-4">${b.Dropoff || "-"}</td>
        <td class="px-6 py-4">${b.VehicleType || "-"}</td>
        <td class="px-6 py-4">
          <span class="px-3 py-1 text-xs font-semibold rounded-full ${
            b.Status === "Completed"
              ? "bg-green-100 text-green-800"
              : b.Status === "Pending"
              ? "bg-yellow-100 text-yellow-800"
              : "bg-blue-100 text-blue-800"
          }">${b.Status || "-"}</span>
        </td>
        <td class="px-6 py-4">${b.CreatedAt ? new Date(b.CreatedAt).toLocaleString() : "-"}</td>
      `;
      tbody.appendChild(row);
    });
  } catch (error) {
    console.error("Error loading dashboard:", error);
    const tbody = document.getElementById("booking-history-body");
    tbody.innerHTML =
      '<tr><td colspan="8" class="text-center py-4 text-red-500">Error loading data</td></tr>';
  }
}

// ✅ Load dashboard on page ready
window.addEventListener("DOMContentLoaded", () => {
  console.log("Admin Dashboard ready");
  loadDashboardData();
});

// ✅ Logout handler
document.getElementById("logout-btn")?.addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
  localStorage.removeItem("permissions");
  window.location.href = "signin.html";
});
