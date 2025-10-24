const usersContainer = document.getElementById("usersContainer");

document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
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

async function fetchSubscribedUsers() {
  try {
    const response = await axios.get("http://localhost:8080/admin/subscription/users", getAuthConfig());
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
  return date.toLocaleString("en-IN", { dateStyle: "medium", timeStyle: "short" });
}

document.addEventListener("DOMContentLoaded", fetchSubscribedUsers);
