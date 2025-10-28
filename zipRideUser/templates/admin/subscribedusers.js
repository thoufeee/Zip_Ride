const usersContainer = document.getElementById("usersContainer");

document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
  localStorage.removeItem("permissions");
  window.location.href = "signin.html";
});

function getAuthConfig() {
  const token = localStorage.getItem("accessToken");
  if (!token) {
    window.location.href = "signin.html";
    throw new Error("Access token missing");
  }
  return { headers: { Authorization: `Bearer ${token}` } };
}

function hasPermission(permission) {
  const permissionsData = localStorage.getItem("permissions");
  if (!permissionsData) return false;

  try {
    const permissions = JSON.parse(permissionsData);
    return Array.isArray(permissions) && permissions.includes(permission);
  } catch {
    return false;
  }
}

async function fetchSubscribedUsers() {
  try {
    const response = await axios.get(
      "http://localhost:8080/admin/subscription/users",
      getAuthConfig()
    );

    const data = response.data.res || response.data.users || response.data || [];

    if (!data || (Array.isArray(data) && data.length === 0)) {
      usersContainer.innerHTML = `<p class="text-center text-gray-500">No subscribed users found.</p>`;
      return;
    }

    const users = Array.isArray(data) ? data : [data];

    usersContainer.innerHTML = `
      <div class="overflow-x-auto">
        <table class="min-w-full border border-gray-200 rounded-lg overflow-hidden shadow-sm">
          <thead class="bg-gray-100 text-gray-700">
            <tr>
              <th class="py-3 px-4 text-left font-semibold">User Name</th>
              <th class="py-3 px-4 text-left font-semibold">Plan Name</th>
              <th class="py-3 px-4 text-left font-semibold">Start Date</th>
              <th class="py-3 px-4 text-left font-semibold">End Date</th>
              <th class="py-3 px-4 text-left font-semibold">Status</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100">
            ${users
              .map(
                (u) => `
              <tr class="hover:bg-gray-50 transition">
                <td class="py-2 px-4">${u.user_name || "N/A"}</td>
                <td class="py-2 px-4">${u.plan_name || "N/A"}</td>
                <td class="py-2 px-4">${formatDate(u.start_date)}</td>
                <td class="py-2 px-4">${formatDate(u.end_date)}</td>
                <td class="py-2 px-4">
                  <span class="px-3 py-1 text-xs font-semibold rounded-full ${
                    u.status === "ACTIVE"
                      ? "bg-green-100 text-green-800"
                      : u.status === "EXPIRED"
                      ? "bg-red-100 text-red-800"
                      : "bg-yellow-100 text-yellow-800"
                  }">
                    ${u.status || "N/A"}
                  </span>
                </td>
              </tr>`
              )
              .join("")}
          </tbody>
        </table>
      </div>
    `;
  } catch (error) {
    console.error("Failed to fetch subscribed users:", error);
    usersContainer.innerHTML = `<p class="text-red-500 text-center">Failed to load subscribed users.</p>`;
  }
}

function formatDate(dateStr) {
  if (!dateStr) return "-";
  const date = new Date(dateStr);
  return date.toLocaleString("en-IN", {
    dateStyle: "medium",
    timeStyle: "short",
  });
}

document.addEventListener("DOMContentLoaded", () => {
  if (!hasPermission("ACCESS_SUBSCRIPTION")) {
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
    return; 
  }

  fetchSubscribedUsers();
});
