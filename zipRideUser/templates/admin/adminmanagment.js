const BASE_URL = "http://localhost:8080/admin";

// 🔹 Logout
document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.clear();
  window.location.href = "signin.html";
});

// 🔹 Utility — Get stored token & permissions
function getAccessToken() {
  return localStorage.getItem("accessToken"); 
}
function getPermissions() {
  const perms = localStorage.getItem("permissions");
  try {
    return perms ? JSON.parse(perms) : [];
  } catch {
    return [];
  }
}

// 🔹 ACL Permission Check
function hasPermission(permission) {
  const perms = getPermissions();
  return perms.includes(permission);
}

// 🔹 Protect Page Access
if (!hasPermission("VIEW_STAFFS")) {
  document.body.innerHTML = `
    <div class="flex items-center justify-center min-h-screen bg-gray-50">
      <div class="bg-white p-10 rounded-xl shadow-2xl text-center">
        <h2 class="text-2xl font-bold text-red-600 mb-3">Access Denied</h2>
        <p class="text-gray-600">You don’t have permission to view this page.</p>
        <button onclick="window.location.href='admindash.html'" class="mt-6 bg-cyan-600 hover:bg-cyan-700 text-white px-6 py-2 rounded-xl shadow-md">
          Go Back
        </button>
      </div>
    </div>`;
  throw new Error("Access Denied");
}

// 🔹 DOM Elements
const tableBody = document.getElementById("adminTable");
const modal = document.getElementById("adminModal");
const modalContent = document.getElementById("modalContent");
const saveChangesBtn = document.getElementById("saveChangesBtn");

// 🔹 Load Admins
async function loadAdmins() {
    try {
        // get access token from localStorage
       const token = localStorage.getItem("accessToken");
        if (!token) {
            alert("Unauthorized: Please login again.");
            window.location.href = "signin.html";
            return;
        }

        const res = await axios.get(`${BASE_URL}/allstaffs`, {
            headers: {
                Authorization: `Bearer ${token}`,  
            },
        });

        const admins = res.data.res || [];
        renderTable(admins);

    } catch (error) {
        console.error("Error fetching admins:", error);

        if (error.response && error.response.status === 401) {
            alert("Session expired or unauthorized. Please log in again.");
            localStorage.clear();
            window.location.href = "signin.html";
        }
    }
}

// 🔹 Render Admin Table
function renderTable(admins) {
  const statusColors = {
    Active: "bg-green-100 text-green-700 border border-green-200",
    Blocked: "bg-red-100 text-red-700 border border-red-200",
  };

  tableBody.innerHTML = "";
  admins.forEach((a) => {
    // remove placeholder permission
    const validPermissions = (a.permissions || []).filter(
      (p) => p !== "Click an available permission to add it."
    );

    const status = a.Block ? "Blocked" : "Active";

    const tr = document.createElement("tr");
    tr.className = "hover:bg-gray-50 transition duration-150";
    tr.innerHTML = `
      <td class="px-6 py-3 whitespace-nowrap text-sm font-medium text-gray-900">${a.name}</td>
      <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-500">${a.email}</td>
      <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-500">${a.phonenumber}</td>
      <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-600 font-medium">${validPermissions.join(", ")}</td>
      <td class="px-6 py-3 whitespace-nowrap text-center">
        <span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${statusColors[status] || ""}">
          ${status}
        </span>
      </td>
      <td class="px-6 py-3 whitespace-nowrap text-center space-x-2">
        ${
          hasPermission("EDIT_STAFF")
            ? `<button class="bg-cyan-600 text-white px-4 py-1.5 rounded-full text-sm hover:bg-cyan-700 transition duration-150 shadow-md editBtn">Edit</button>`
            : ""
        }
        ${
          hasPermission("DELETE_STAFF")
            ? `<button class="bg-red-500 text-white px-4 py-1.5 rounded-full text-sm hover:bg-red-600 transition duration-150 shadow-md deleteBtn">Delete</button>`
            : ""
        }
        ${
          hasPermission("BLOCK_STAFF") || hasPermission("UNBLOCK_STAFF")
            ? `<button class="bg-yellow-500 text-white px-4 py-1.5 rounded-full text-sm hover:bg-yellow-600 transition duration-150 shadow-md blockBtn">${
                status === "Active" ? "Block" : "Unblock"
              }</button>`
            : ""
        }
      </td>
    `;

    tableBody.appendChild(tr);

    // 🔸 Edit Admin
    const editBtn = tr.querySelector(".editBtn");
    if (editBtn) {
      editBtn.addEventListener("click", () => openEditModal(a));
    }

    // 🔸 Delete Admin
    const deleteBtn = tr.querySelector(".deleteBtn");
    if (deleteBtn) {
      deleteBtn.addEventListener("click", () => deleteAdmin(a.ID, a.name));
    }

    // 🔸 Block / Unblock Admin
    const blockBtn = tr.querySelector(".blockBtn");
    if (blockBtn) {
      blockBtn.addEventListener("click", () => toggleBlock(a));
    }
  });
}

// 🔹 Edit Admin Modal
function openEditModal(admin) {
  modalContent.innerHTML = `
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
        <input type="text" id="adminName" class="border border-gray-300 rounded-xl px-4 py-2 w-full" value="${admin.name}">
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
        <input type="email" id="adminEmail" class="border border-gray-300 rounded-xl px-4 py-2 w-full" value="${admin.email}" readonly>
      </div>
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
        <input type="text" id="adminPhone" class="border border-gray-300 rounded-xl px-4 py-2 w-full" value="${admin.phonenumber}">
      </div>
      <div class="pt-2">
        <p class="block text-sm font-extrabold text-cyan-700 mb-2 uppercase">Access Permissions</p>
        <div class="flex flex-wrap gap-x-4 gap-y-2 text-sm">
          ${["VIEW_STAFFS", "ADD_STAFF", "EDIT_STAFF", "DELETE_STAFF", "BLOCK_STAFF", "UNBLOCK_STAFF"]
            .map(
              (perm) => `
              <label class="inline-flex items-center">
                <input type="checkbox" class="permCheckbox" value="${perm}">
                <span class="ml-2 text-gray-700">${perm.replaceAll("_", " ")}</span>
              </label>`
            )
            .join("")}
        </div>
      </div>
    </div>
  `;

  modal.querySelectorAll(".permCheckbox").forEach((cb) => {
    if (admin.permissions.includes(cb.value)) cb.checked = true;
  });

  modal.classList.remove("hidden");
  modal.classList.add("flex");

  saveChangesBtn.onclick = async () => {
    const updated = {
      name: document.getElementById("adminName").value,
      phonenumber: document.getElementById("adminPhone").value,
      permissions: Array.from(
        modal.querySelectorAll(".permCheckbox:checked")
      ).map((cb) => cb.value),
    };

    try {
      const token = getAccessToken();
      await axios.put(`${BASE_URL}/staffupdate/${admin.ID}`, updated, {
        headers: { Authorization: `Bearer ${token}` },
      });
      alert("Admin updated successfully!");
      modal.classList.add("hidden");
      loadAdmins();
    } catch (err) {
      console.error(err);
      alert("Failed to update admin.");
    }
  };
}

// 🔹 Delete Admin
async function deleteAdmin(id, name) {
  if (!confirm(`Are you sure you want to delete ${name}?`)) return;
  try {
    const token = getAccessToken();
    await axios.delete(`${BASE_URL}/deleteadmin/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    alert("Admin deleted successfully!");
    loadAdmins();
  } catch (err) {
    console.error(err);
    alert("Failed to delete admin.");
  }
}

// 🔹 Block / Unblock Admin
async function toggleBlock(admin) {
  const endpoint = admin.Block
    ? `${BASE_URL}/staffunblock/${admin.ID}` // currently blocked → unblock
    : `${BASE_URL}/staffblock/${admin.ID}`;  // currently active → block

  try {
    const token = getAccessToken();
    const res = await axios.put(endpoint, {}, { headers: { Authorization: `Bearer ${token}` } });

    alert(res.data.res || res.data.message || "Status updated successfully");

    // Update local status
    admin.Block = !admin.Block;

    // Reload table
    loadAdmins();

  } catch (err) {
    console.error(err);
    alert("Failed to update admin status.");
  }
}



// 🔹 Close Modal
document.getElementById("closeModal").addEventListener("click", () => {
  modal.classList.add("hidden");
  modal.classList.remove("flex");
});

// 🔹 Initial Load
loadAdmins();
