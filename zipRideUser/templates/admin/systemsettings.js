const BASE_URL = "http://localhost:8080/admin";

document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
  localStorage.removeItem("permissions");
  window.location.href = "signin.html";
});

function hasPermission(permissionName) {
  const permissions = JSON.parse(localStorage.getItem("permissions") || "[]");
  return permissions.includes(permissionName);
}

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

const fareList = document.getElementById("fareList");
const addFareBtn = document.getElementById("addFareBtn");

async function loadFares() {
  try {
    const token = localStorage.getItem("accessToken");
    const response = await axios.get(`${BASE_URL}/vehiclefare/`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    const fares = response.data || [];
    if (fares.length === 0) {
      fareList.innerHTML = `<p class="text-gray-500 text-center">No fare configurations found.</p>`;
      return;
    }

    fareList.innerHTML = fares
      .map(
        (fare) => `
      <div class="border border-gray-200 p-6 rounded-xl shadow-sm flex flex-col sm:flex-row sm:justify-between sm:items-center bg-gray-50 hover:bg-gray-100 transition" data-id="${fare.id}">
        <div>
          <p class="text-lg font-bold text-cyan-700 capitalize">${fare.vehicle_type}</p>
          <p class="text-sm text-gray-600 mt-1">Base Fare: ₹${fare.base_fare}</p>
          <p class="text-sm text-gray-600">Fare per Km: ₹${fare.per_km_rate}</p>
          <p class="text-sm text-gray-600">Fare per Minute: ₹${fare.per_min_rate}</p>
          <p class="text-sm text-gray-600">People Count: ${fare.people_count}</p>
        </div>
        <div class="flex space-x-3 mt-4 sm:mt-0">
          <button onclick="editFare(${fare.id})" class="bg-cyan-600 hover:bg-cyan-700 text-white py-1.5 px-4 rounded-xl text-sm font-medium">Edit</button>
          <button onclick="deleteFare(${fare.id})" class="bg-red-600 hover:bg-red-700 text-white py-1.5 px-4 rounded-xl text-sm font-medium">Delete</button>
        </div>
      </div>
    `
      )
      .join("");
  } catch (err) {
    console.error(err);
    fareList.innerHTML = `<p class="text-red-500 text-center">Failed to load fare configurations.</p>`;
  }
}


function showFareModal(fare = null) {
  const isEdit = fare !== null;

  const modal = document.createElement("div");
  modal.className = "fixed inset-0 flex items-center justify-center bg-black bg-opacity-40 backdrop-blur-sm z-50";

  modal.innerHTML = `
    <div class="bg-white rounded-xl shadow-xl p-6 w-[400px]">
      <h2 class="text-xl font-semibold mb-4 text-gray-800">${isEdit ? "Edit Fare" : "Add Fare"}</h2>
      <form id="fareModalForm" class="space-y-3">
        <input type="hidden" id="modalFareId" value="${fare?.id || ''}">
        <div>
          <label class="block text-sm font-medium text-gray-600">Vehicle Type</label>
          <input id="modalVehicleType" value="${fare?.vehicle_type || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600">Base Fare (₹)</label>
          <input id="modalBaseFare" type="number" value="${fare?.base_fare || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600">Fare per Km (₹)</label>
          <input id="modalPerKmRate" type="number" value="${fare?.per_km_rate || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" required />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600">Fare per Minute (₹)</label>
          <input id="modalPerMinRate" type="number" value="${fare?.per_min_rate || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-600">People Count</label>
          <input id="modalPeopleCount" type="number" value="${fare?.people_count || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" required />
        </div>
        <div class="flex justify-end space-x-2 mt-4">
          <button type="button" id="closeFareModal" class="bg-gray-400 hover:bg-gray-500 text-white px-4 py-2 rounded-lg">Cancel</button>
          <button type="submit" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">${isEdit ? "Save Changes" : "Add Fare"}</button>
        </div>
      </form>
    </div>
  `;

  document.body.appendChild(modal);

  document.getElementById("closeFareModal").addEventListener("click", () => modal.remove());

  document.getElementById("fareModalForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const fareData = {
      vehicle_type: document.getElementById("modalVehicleType").value.trim(),
      base_fare: parseFloat(document.getElementById("modalBaseFare").value),
      per_km_rate: parseFloat(document.getElementById("modalPerKmRate").value),
      per_min_rate: parseFloat(document.getElementById("modalPerMinRate").value) || 0,
      people_count: parseInt(document.getElementById("modalPeopleCount").value),
    };

    try {
      const token = localStorage.getItem("accessToken");
      if (isEdit) {
        await axios.put(`${BASE_URL}/vehiclefare/${fare.id}`, fareData, { headers: { Authorization: `Bearer ${token}` } });
        alert("Fare updated successfully!");
      } else {
        await axios.post(`${BASE_URL}/vehiclefare/`, fareData, { headers: { Authorization: `Bearer ${token}` } });
        alert("Fare added successfully!");
      }
      modal.remove();
      loadFares();
    } catch (err) {
      console.error(err);
      alert("Failed to save fare configuration.");
    }
  });
}

function editFare(id) {
  const fareDiv = Array.from(fareList.children).find(div => parseInt(div.dataset.id) === id);
  if (!fareDiv) return;

  const fare = {
    id,
    vehicle_type: fareDiv.querySelector('p:nth-child(1)').textContent,
    base_fare: parseFloat(fareDiv.querySelector('p:nth-child(2)').textContent.replace(/\D/g, '')),
    per_km_rate: parseFloat(fareDiv.querySelector('p:nth-child(3)').textContent.replace(/\D/g, '')),
    per_min_rate: parseFloat(fareDiv.querySelector('p:nth-child(4)').textContent.replace(/\D/g, '')),
    people_count: parseInt(fareDiv.querySelector('p:nth-child(5)').textContent.replace(/\D/g, '')),
  };

  showFareModal(fare);
}

async function deleteFare(id) {
  if (!confirm("Are you sure you want to delete this fare?")) return;

  try {
    const token = localStorage.getItem("accessToken");
    await axios.delete(`${BASE_URL}/vehiclefare/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    alert("Fare deleted successfully!");
    loadFares();
  } catch (err) {
    console.error(err);
    alert("Failed to delete fare.");
  }
}


addFareBtn.addEventListener("click", () => showFareModal());


loadFares();
