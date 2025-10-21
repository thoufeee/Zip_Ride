const BASE_URL = "http://localhost:8080/admin";


const logoutBtn = document.getElementById("logout-btn");
const usersList = document.getElementById("users-list");
const createUserForm = document.getElementById("createUserForm");
const createUserResult = document.getElementById("createUserResult");


const accessToken = localStorage.getItem("accessToken");
const userPermissions = JSON.parse(localStorage.getItem("permissions")) || [];

if (!accessToken) window.location.href = "signin.html";


axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.clear();
      window.location.href = "signin.html";
    }
    return Promise.reject(error);
  }
);


if (logoutBtn) {
  logoutBtn.addEventListener("click", () => {
    localStorage.clear();
    window.location.href = "signin.html";
  });
}


function hasPermission(requiredPermission) {
  return userPermissions.some(
    (p) => p.toUpperCase() === requiredPermission.toUpperCase()
  );
}


if (!hasPermission("VIEW_USERS")) {
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
  throw new Error("Permission denied");
}


async function fetchUsers() {
  if (!usersList) return;

  try {
    const res = await axios.get(`${BASE_URL}/allusers`, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    const users = res.data.res || [];
    renderUsers(users);
  } catch (err) {
    console.error("Error fetching users:", err);
    usersList.innerHTML = `<p class="text-red-500 text-sm">Failed to load users.</p>`;
  }
}

function renderUsers(users) {
  usersList.innerHTML = `
    <div class="grid grid-cols-7 font-semibold text-gray-700 border-b border-gray-300 pb-2 mb-2 text-sm">
      <p>First Name</p>
      <p>Last Name</p>
      <p>Email</p>
      <p>Phone</p>
      <p>Gender</p>
      <p>Location</p>
      <p class="text-center">Actions</p>
    </div>
  `;

  users.forEach((user) => {
    const div = document.createElement("div");
    div.className =
      "grid grid-cols-7 items-center gap-2 py-3 px-2 mb-1 bg-white rounded-lg border border-gray-200 shadow-sm hover:shadow-md transition duration-200 text-sm";

    const actions = [];

    if (hasPermission("EDIT_USER")) {
      actions.push(`<button class="edit-btn bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded-md text-xs font-medium transition" data-id="${user.ID}">Edit</button>`);
    }
    if (hasPermission("DELETE_USER")) {
      actions.push(`<button class="delete-btn bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-md text-xs font-medium ml-1" data-id="${user.ID}">Delete</button>`);
    }
    if (hasPermission("BLOCK_USER") && !user.Block) {
      actions.push(`<button class="block-btn bg-yellow-400 hover:bg-yellow-500 text-white px-3 py-1 rounded-md text-xs font-medium ml-1" data-id="${user.ID}">Block</button>`);
    }
    if (hasPermission("UNBLOCK_USER") && user.Block) {
      actions.push(`<button class="unblock-btn bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded-md text-xs font-medium ml-1" data-id="${user.ID}">Unblock</button>`);
    }

    div.innerHTML = `
      <p class="text-gray-800 font-medium">${user.firstname || "-"}</p>
      <p class="text-gray-800 font-medium">${user.lastname || "-"}</p>
      <p class="text-gray-600">${user.email}</p>
      <p class="text-gray-600">${user.phone || "-"}</p>
      <p class="text-gray-600">${user.gender || "-"}</p>
      <p class="text-gray-600">${user.place || "-"}</p>
      <div class="flex justify-center space-x-1">${actions.join("")}</div>
    `;

    usersList.appendChild(div);
  });

  attachUserActions(users);
}


function attachUserActions(users) {
  document.querySelectorAll(".edit-btn").forEach((btn) => {
    btn.addEventListener("click", () => {
      const id = btn.dataset.id;
      const user = users.find((u) => u.ID == id);
      if (user) showEditModal(user);
    });
  });

  document.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      if (!confirm("Are you sure you want to delete this user?")) return;
      try {
        const res = await axios.delete(`${BASE_URL}/user/${btn.dataset.id}`, { headers: { Authorization: `Bearer ${accessToken}` } });
        alert(res.data.message || "User deleted successfully!");
        fetchUsers();
      } catch (err) {
        console.error(err);
        alert(err.response?.data?.error || "Failed to delete user.");
      }
    });
  });

  document.querySelectorAll(".block-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      try {
        const res = await axios.put(`${BASE_URL}/userblock/${btn.dataset.id}`, {}, { headers: { Authorization: `Bearer ${accessToken}` } });
        alert(res.data.message || "User blocked successfully!");
        fetchUsers();
      } catch (err) {
        console.error(err);
        alert(err.response?.data?.error || "Failed to block user.");
      }
    });
  });

  document.querySelectorAll(".unblock-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      try {
        const res = await axios.put(`${BASE_URL}/userunblock/${btn.dataset.id}`, {}, { headers: { Authorization: `Bearer ${accessToken}` } });
        alert(res.data.message || "User unblocked successfully!");
        fetchUsers();
      } catch (err) {
        console.error(err);
        alert(err.response?.data?.error || "Failed to unblock user.");
      }
    });
  });
}


function showEditModal(user) {
  const modal = document.createElement("div");
  modal.className = "fixed inset-0 flex items-center justify-center bg-black bg-opacity-40 backdrop-blur-sm z-50";

  modal.innerHTML = `
    <div class="bg-white rounded-xl shadow-xl p-6 w-[400px]">
      <h2 class="text-xl font-semibold mb-4 text-gray-800">Edit User</h2>
      <form id="editUserForm" class="space-y-3">
        <input type="hidden" id="editUserId" value="${user.ID}">
        <div><label class="block text-sm font-medium text-gray-600">First Name</label><input id="editFirstname" value="${user.firstname || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" /></div>
        <div><label class="block text-sm font-medium text-gray-600">Last Name</label><input id="editLastname" value="${user.lastname || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" /></div>
        <div><label class="block text-sm font-medium text-gray-600">Email</label><input id="editEmail" value="${user.email || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" /></div>
        <div><label class="block text-sm font-medium text-gray-600">Phone</label><input id="editPhone" value="${user.phone || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" /></div>
        <div><label class="block text-sm font-medium text-gray-600">Gender</label>
          <select id="editGender" class="w-full border border-gray-300 rounded-lg px-3 py-2">
            <option value="">Select</option>
            <option value="Male" ${user.gender === "Male" ? "selected" : ""}>Male</option>
            <option value="Female" ${user.gender === "Female" ? "selected" : ""}>Female</option>
            <option value="Other" ${user.gender === "Other" ? "selected" : ""}>Other</option>
          </select>
        </div>
        <div><label class="block text-sm font-medium text-gray-600">Location</label><input id="editPlace" value="${user.place || ''}" class="w-full border border-gray-300 rounded-lg px-3 py-2" /></div>
        <div class="flex justify-end space-x-2 mt-4">
          <button type="button" id="closeEditModal" class="bg-gray-400 hover:bg-gray-500 text-white px-4 py-2 rounded-lg">Cancel</button>
          <button type="submit" class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded-lg">Save</button>
        </div>
      </form>
    </div>
  `;

  document.body.appendChild(modal);

  document.getElementById("closeEditModal").addEventListener("click", () => modal.remove());

  document.getElementById("editUserForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    try {
      const updatedUser = {
        firstname: document.getElementById("editFirstname").value.trim(),
        lastname: document.getElementById("editLastname").value.trim(),
        email: document.getElementById("editEmail").value.trim(),
        phone: document.getElementById("editPhone").value.trim(),
        gender: document.getElementById("editGender").value,
        place: document.getElementById("editPlace").value.trim(),
      };

      const res = await axios.put(`${BASE_URL}/user/${user.ID}`, updatedUser, {
        headers: { Authorization: `Bearer ${accessToken}` },
      });

      alert(res.data.message || "User updated successfully!");
      modal.remove();
      fetchUsers();
    } catch (err) {
      console.error(err);
      alert(err.response?.data?.error || "Failed to update user.");
    }
  });
}


if (createUserForm) {
  createUserForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    if (!hasPermission("CREATE_USER") && !hasPermission("ADD_USER")) {
      createUserResult.textContent = "You do not have permission to create users.";
      createUserResult.className = "text-red-500 font-bold mt-4 text-center";
      return;
    }

    const newUser = {
      firstname: document.getElementById("newFirstname").value.trim(),
      lastname: document.getElementById("newLastname").value.trim(),
      email: document.getElementById("newEmail").value.trim(),
      phone: document.getElementById("newPhone").value.trim(),
      place: document.getElementById("newPlace").value.trim(),
      password: document.getElementById("newPassword").value.trim(),
      gender: document.getElementById("newGender")?.value || "",
    };

    try {
      const res = await axios.post(`${BASE_URL}/createuser`, newUser, {
        headers: { Authorization: `Bearer ${accessToken}` },
      });

      createUserResult.textContent = res.data.message || "User added successfully!";
      createUserResult.className = "text-green-600 font-bold mt-4 text-center";
      createUserForm.reset();
      fetchUsers();
    } catch (err) {
      console.error(err);
      createUserResult.textContent = err.response?.data?.error || "Failed to add user.";
      createUserResult.className = "text-red-500 font-bold mt-4 text-center";
    }
  });
}

fetchUsers();
