
const BASE_URL = "http://localhost:8080/admin";

const logoutBtn = document.getElementById("logout-btn");
if (logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    localStorage.clear();
    window.location.href = "signin.html";
  });
}

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

function hasPermission(permission) {
  const perms = getPermissions();
  return perms.some(p => String(p).toUpperCase() === String(permission).toUpperCase());
}

if (!hasPermission("VIEW_STAFFS")) {
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

const tableBody = document.getElementById("adminTable");
const modal = document.getElementById("adminModal");
const modalContent = document.getElementById("modalContent");
const saveChangesBtn = document.getElementById("saveChangesBtn");
const closeModalBtn = document.getElementById("closeModal");


if (!tableBody) console.warn("adminTable element not found.");
if (!modal || !modalContent) console.warn("Modal elements not found.");


async function loadAdmins() {
  try {
    const token = getAccessToken();
    if (!token) {
      alert("Unauthorized: Please login again.");
      window.location.href = "signin.html";
      return;
    }

    const res = await axios.get(`${BASE_URL}/allstaffs`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    const admins = res.data.res || [];
    renderTable(admins);

  } catch (error) {
    console.error("Error fetching admins:", error);
    if (error.response && error.response.status === 401) {
      alert("Session expired or unauthorized. Please log in again.");
      localStorage.clear();
      window.location.href = "signin.html";
    } else {
      alert("Failed to load admins. See console for details.");
    }
  }
}


function renderTable(admins) {
  const statusColors = {
    Active: "bg-green-100 text-green-700 border border-green-200",
    Blocked: "bg-red-100 text-red-700 border border-red-200",
  };

  tableBody.innerHTML = "";

  admins.forEach((a) => {
    const validPermissions = Array.isArray(a.permissions) ? a.permissions.filter(
      (p) => p && p !== "Click an available permission to add it."
    ) : [];

    const status = a.Block ? "Blocked" : "Active";

    const tr = document.createElement("tr");
    tr.className = "hover:bg-gray-50 transition duration-150";
    tr.innerHTML = `
      <td class="px-6 py-3 whitespace-nowrap text-sm font-medium text-gray-900">${escapeHtml(a.name)}</td>
      <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-500">${escapeHtml(a.email)}</td>
      <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-500">${escapeHtml(a.phonenumber || '')}</td>
      <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-600 font-medium">${escapeHtml(validPermissions.join(", "))}</td>
      <td class="px-6 py-3 whitespace-nowrap text-center">
        <span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${statusColors[status] || ""}">
          ${status}
        </span>
      </td>
      <td class="px-6 py-3 whitespace-nowrap text-center space-x-2">
        ${hasPermission("EDIT_STAFF") ? `<button class="bg-cyan-600 text-white px-4 py-1.5 rounded-full text-sm hover:bg-cyan-700 transition duration-150 shadow-md editBtn">Edit</button>` : ""}
        ${hasPermission("DELETE_STAFF") ? `<button class="bg-red-500 text-white px-4 py-1.5 rounded-full text-sm hover:bg-red-600 transition duration-150 shadow-md deleteBtn">Delete</button>` : ""}
        ${(hasPermission("BLOCK_STAFF") || hasPermission("UNBLOCK_STAFF")) ? `<button class="bg-yellow-500 text-white px-4 py-1.5 rounded-full text-sm hover:bg-yellow-600 transition duration-150 shadow-md blockBtn">${status === "Active" ? "Block" : "Unblock"}</button>` : ""}
      </td>
    `;

    tableBody.appendChild(tr);

    const editBtn = tr.querySelector(".editBtn");
    if (editBtn) {
      editBtn.addEventListener("click", () => openEditModal(a));
    }

    const deleteBtn = tr.querySelector(".deleteBtn");
    if (deleteBtn) {
      deleteBtn.addEventListener("click", () => deleteAdmin(a.ID, a.name));
    }

    const blockBtn = tr.querySelector(".blockBtn");
    if (blockBtn) {
      blockBtn.addEventListener("click", () => toggleBlock(a));
    }
  });
}

function escapeHtml(text) {
  if (!text && text !== 0) return "";
  return String(text)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

async function openEditModal(admin) {
  try {
    const token = getAccessToken();
    if (!token) {
      alert("Unauthorized");
      return;
    }

    const permRes = await axios.get(`${BASE_URL}/allpermissions`, {
      headers: { Authorization: `Bearer ${token}` },
    });

    const allPermissions = (permRes.data.res || []).map(p => p.name || []);

    modal.classList.remove("hidden");
    modal.classList.add("flex");

    modalContent.innerHTML = `
      <div class="relative bg-white rounded-xl shadow-2xl w-full max-w-3xl max-h-[90vh] overflow-y-auto p-6">
        <button id="closeModalBtn" class="absolute top-3 right-3 text-gray-500 hover:text-gray-700 text-xl font-bold">✖</button>
        <h2 class="text-2xl font-bold mb-4 text-center text-gray-800">Edit Admin</h2>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
            <input type="text" id="adminName" class="border border-gray-300 rounded-xl px-4 py-2 w-full" value="${escapeHtml(admin.name || '')}">
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
            <input type="email" id="adminEmail" class="border border-gray-300 rounded-xl px-4 py-2 w-full" value="${escapeHtml(admin.email || '')}">
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
            <input type="text" id="adminPhone" class="border border-gray-300 rounded-xl px-4 py-2 w-full" value="${escapeHtml(admin.phonenumber || '')}">
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Password (leave empty to keep current)</label>
            <input type="password" id="adminPassword" class="border border-gray-300 rounded-xl px-4 py-2 w-full" placeholder="New password">
          </div>

          <div>
            <p class="font-semibold text-cyan-700 mb-2">Access Permissions</p>
            <div id="permList" class="flex flex-wrap gap-3 max-h-60 overflow-y-auto"></div>
          </div>

          <div class="mt-6 text-center">
            <button id="saveChangesBtnModal" class="bg-cyan-600 hover:bg-cyan-700 text-white px-6 py-2 rounded-lg shadow-md">Save Changes</button>
          </div>
        </div>
      </div>
    `;

    const closeBtn = modalContent.querySelector("#closeModalBtn");
    closeBtn.addEventListener("click", () => {
      modal.classList.add("hidden");
      modal.classList.remove("flex");
    });

    const permList = modalContent.querySelector("#permList");
    allPermissions.forEach((perm) => {
      const checked = Array.isArray(admin.permissions) && admin.permissions.includes(perm);
      const div = document.createElement("div");
      div.className = "flex items-center space-x-2";
      div.innerHTML = `
        <input type="checkbox" class="permCheckbox" value="${perm}" ${checked ? "checked" : ""}>
        <label class="text-gray-700 text-sm">${perm.replaceAll("_", " ")}</label>
      `;
      permList.appendChild(div);
    });

    const saveBtn = modalContent.querySelector("#saveChangesBtnModal");
    saveBtn.addEventListener("click", async () => {
      const payload = {
        name: document.getElementById("adminName").value.trim(),
        email: document.getElementById("adminEmail").value.trim(),
        phonenumber: document.getElementById("adminPhone").value.trim(),
        password: document.getElementById("adminPassword").value.trim(),
        extra_permissions: Array.from(modalContent.querySelectorAll(".permCheckbox:checked")).map(cb => cb.value),
      };

      try {
        const res = await axios.put(`${BASE_URL}/staffupdate/${admin.ID}`, payload, {
          headers: { Authorization: `Bearer ${token}` },
        });
        alert(res.data.message || "Admin updated successfully!");
        modal.classList.add("hidden");
        modal.classList.remove("flex");
        loadAdmins();
      } catch (err) {
        console.error(err);
        alert(err.response?.data?.error || "Failed to update admin.");
      }
    });

  } catch (err) {
    console.error(err);
    alert("Failed to load permissions.");
  }
}

async function deleteAdmin(id, name) {
  if (!confirm(`Are you sure you want to delete ${name}?`)) return;
  try {
    const token = getAccessToken();
    await axios.delete(`${BASE_URL}/staffdelete/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    });
    alert("Admin deleted successfully!");
    loadAdmins();
  } catch (err) {
    console.error(err);
    alert(err.response?.data?.error || "Failed to delete admin.");
  }
}

async function toggleBlock(admin) {
  const endpoint = admin.Block ? `${BASE_URL}/staffunblock/${admin.ID}` : `${BASE_URL}/staffblock/${admin.ID}`;
  try {
    const token = getAccessToken();
    const res = await axios.put(endpoint, {}, { headers: { Authorization: `Bearer ${token}` } });
    alert(res.data.message || res.data.res || "Status updated successfully");
    admin.Block = !admin.Block;
    loadAdmins();
  } catch (err) {
    console.error(err);
    alert(err.response?.data?.error || "Failed to update admin status.");
  }
}

if (closeModalBtn) {
  closeModalBtn.addEventListener("click", () => {
    modal.classList.add("hidden");
    modal.classList.remove("flex");
  });
}

loadAdmins();
