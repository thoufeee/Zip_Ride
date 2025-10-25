const BASE_URL = "http://localhost:8080"; 


const userPermissions = JSON.parse(localStorage.getItem("permissions") || "[]"); 

function hasPermission(permission) {
  return userPermissions.includes(permission);
}


if (!hasPermission("ACCESS_PRIZEPOOL")) {
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
  throw new Error("Unauthorized access to Prize Pool Management");
}

const token = localStorage.getItem("accessToken"); 
if (!token) {
  window.location.href = "signin.html"; 
}

const axiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json"
  }
});

const tableBody = document.getElementById("prizepool-table-body");
const modal = document.getElementById("modal");
const form = document.getElementById("prizepool-form");
const modalTitle = document.getElementById("modal-title");
const cancelBtn = document.getElementById("cancel-btn");
const addNewBtn = document.getElementById("add-new-btn");


addNewBtn.addEventListener("click", () => {
  modal.classList.remove("hidden");
  modal.classList.add("flex");
  modalTitle.textContent = "Add Prize Pool";
  form.reset();
  document.getElementById("prizepool-id").value = "";
});

cancelBtn.addEventListener("click", () => modal.classList.add("hidden"));


async function fetchPrizePools() {
  try {
    const res = await axiosInstance.get("/admin/pricepool/");
    renderPrizePools(res.data.res || []);
  } catch (err) {
    console.error("Error fetching prize pools:", err.response?.data || err);
    if (err.response?.status === 401) {
      alert("Unauthorized! Please login again.");
      window.location.href = "signin.html";
    }
  }
}


function renderPrizePools(pools) {
  tableBody.innerHTML = "";

  if (pools.length === 0) {
    tableBody.innerHTML = `
      <tr>
        <td colspan="5" class="p-4 text-center text-gray-500">No prize pools found</td>
      </tr>`;
    return;
  }

  pools.forEach(pool => {
    const tr = document.createElement("tr");
    tr.className = "border-b hover:bg-gray-50";

    tr.innerHTML = `
      <td class="p-3 text-left">${pool.vehicle_type}</td>
      <td class="p-3 text-center">${pool.commission}%</td>
      <td class="p-3 text-center">${pool.bonusamount || 0}</td>
      <td class="p-3 text-center">
        <span class="inline-block px-3 py-1 rounded-full text-white font-semibold text-sm ${pool.active ? 'bg-green-500' : 'bg-red-500'}">
          ${pool.active ? 'Active' : 'Inactive'}
        </span>
      </td>
      <td class="p-3 text-center">
        <div class="flex justify-center gap-2">
          <button onclick="editPrizePool('${pool.id}')"
            class="bg-yellow-400 hover:bg-yellow-500 text-white px-3 py-1 rounded font-medium text-sm shadow-md">
            Edit
          </button>
          <button onclick="deletePrizePool('${pool.id}')"
            class="bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded font-medium text-sm shadow-md">
            Delete
          </button>
          <button onclick="toggleActive('${pool.id}', ${pool.active})"
            class="bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded font-medium text-sm shadow-md">
            ${pool.active ? 'Deactivate' : 'Activate'}
          </button>
        </div>
      </td>
    `;

    tableBody.appendChild(tr);
  });
}

form.addEventListener("submit", async (e) => {
  e.preventDefault();

  const id = document.getElementById("prizepool-id").value;
  const data = {
    vehicle_type: document.getElementById("vehicle-type").value,
    commission: parseFloat(document.getElementById("commission").value),
    bonusamount: parseFloat(document.getElementById("bonus-amount").value) || 0,
  };

  try {
    if (id) {
      await axiosInstance.put(`/admin/pricepool/${id}`, data);
      alert("Prize pool updated successfully!");
    } else {
      await axiosInstance.post(`/admin/pricepool/`, data);
      alert("Prize pool added successfully!");
    }
    modal.classList.add("hidden");
    fetchPrizePools();
  } catch (err) {
    console.error("Error saving prize pool:", err.response?.data || err);
  }
});


async function editPrizePool(id) {
  try {
    const res = await axiosInstance.get("/admin/pricepool/");
    const pool = res.data.res.find(p => p.id === id);
    if (!pool) return;

    modal.classList.remove("hidden");
    modal.classList.add("flex");
    modalTitle.textContent = "Edit Prize Pool";

    document.getElementById("prizepool-id").value = pool.id;
    document.getElementById("vehicle-type").value = pool.vehicle_type;
    document.getElementById("commission").value = pool.commission;
    document.getElementById("bonus-amount").value = pool.bonusamount || 0;
  } catch (err) {
    console.error("Error editing pool:", err.response?.data || err);
  }
}

async function deletePrizePool(id) {
  if (!confirm("Are you sure you want to delete this prize pool?")) return;
  try {
    await axiosInstance.delete(`/admin/pricepool/${id}`);
    alert("Deleted successfully!");
    fetchPrizePools();
  } catch (err) {
    console.error("Error deleting:", err.response?.data || err);
  }
}

async function toggleActive(id, currentStatus) {
  try {
    await axiosInstance.put(`/admin/pricepool/status/${id}`, { active: !currentStatus });
    fetchPrizePools();
  } catch (err) {
    console.error("Error updating status:", err.response?.data || err);
  }
}

fetchPrizePools();
