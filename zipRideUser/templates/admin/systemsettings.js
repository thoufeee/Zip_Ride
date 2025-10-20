const BASE_URL = "http://localhost:8080/admin"; // backend base URL

// 🔹 Logout
document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("access");
  localStorage.removeItem("refresh");
  localStorage.removeItem("role");
  localStorage.removeItem("permissions");
  window.location.href = "signin.html";
});

// 🔹 Permission check
function hasPermission(permissionName) {
  const permissions = JSON.parse(localStorage.getItem("permissions") || "[]");
  return permissions.includes(permissionName);
}

// ✅ Restrict page access
if (!hasPermission("SYSTEM_SETTINGS")) {
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
  throw new Error("Unauthorized access to System Settings");
}

// 🔹 Load Fare List
const fareList = document.getElementById("fareList");
const addFareBtn = document.getElementById("addFareBtn");

async function loadFares() {
  try {
    const token = localStorage.getItem("access");
    const response = await axios.get(`${BASE_URL}/fares`, {
      headers: { Authorization: `Bearer ${token}` }
    });

    const fares = response.data.fares || [];
    if (fares.length === 0) {
      fareList.innerHTML = `<p class="text-gray-500 text-center">No fare configurations found.</p>`;
      return;
    }

    fareList.innerHTML = fares.map(fare => `
      <div class="border border-gray-200 p-6 rounded-xl shadow-sm flex flex-col sm:flex-row sm:justify-between sm:items-center bg-gray-50 hover:bg-gray-100 transition">
        <div>
          <p class="text-lg font-bold text-cyan-700 capitalize">${fare.vehicle_type}</p>
          <p class="text-sm text-gray-600 mt-1">Base Fare: ₹${fare.base_fare}</p>
          <p class="text-sm text-gray-600">Fare per Km: ₹${fare.fare_per_km}</p>
          <p class="text-sm text-gray-600">Fare per Minute: ₹${fare.fare_per_minute}</p>
          <p class="text-sm text-gray-600">Commission: ${fare.commission}%</p>
        </div>
        <div class="flex space-x-3 mt-4 sm:mt-0">
          <button onclick="editFare('${fare.id}')" class="bg-cyan-600 hover:bg-cyan-700 text-white py-1.5 px-4 rounded-xl text-sm font-medium">Edit</button>
          <button onclick="deleteFare('${fare.id}')" class="bg-red-600 hover:bg-red-700 text-white py-1.5 px-4 rounded-xl text-sm font-medium">Delete</button>
        </div>
      </div>
    `).join('');
  } catch (err) {
    console.error(err);
    fareList.innerHTML = `<p class="text-red-500 text-center">Failed to load fare configurations.</p>`;
  }
}

// 🔹 Edit Fare
function editFare(id) {
  window.location.href = `fare_edit.html?id=${id}`;
}

// 🔹 Delete Fare
async function deleteFare(id) {
  if (!confirm("Are you sure you want to delete this fare?")) return;
  try {
    const token = localStorage.getItem("access");
    await axios.delete(`${BASE_URL}/fares/${id}`, {
      headers: { Authorization: `Bearer ${token}` }
    });
    alert("Fare deleted successfully!");
    loadFares();
  } catch (err) {
    console.error(err);
    alert("Failed to delete fare.");
  }
}

// 🔹 Add Fare Button
addFareBtn.addEventListener("click", () => {
  window.location.href = "fare_edit.html"; // open blank form for new fare
});

// 🔹 Demo Save Handlers for other forms
document.querySelectorAll("form").forEach(form => {
  form.addEventListener("submit", e => {
    e.preventDefault();
    alert(`Settings saved for ${form.id}! (Demo only - form submission prevented)`);
  });
});

// Initial load
loadFares();
