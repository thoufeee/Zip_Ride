const BASE_URL = "http://localhost:8080/admin";

const logoutBtn = document.getElementById("logout-btn");
const usersList = document.getElementById("users-list");
const dynamicSection = document.getElementById("dynamicSection");
const addUserBtn = document.getElementById("addUserBtn");

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
    renderUsers(res.data.res || []);
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
      actions.push(
        `<button class="edit-btn bg-blue-500 hover:bg-blue-600 text-white px-3 py-1 rounded-md text-xs font-medium transition" data-id="${user.ID}">Edit</button>`
      );
    }
    if (hasPermission("DELETE_USER")) {
      actions.push(
        `<button class="delete-btn bg-red-500 hover:bg-red-600 text-white px-3 py-1 rounded-md text-xs font-medium ml-1" data-id="${user.ID}">Delete</button>`
      );
    }
    if (hasPermission("BLOCK_USER") && !user.Block) {
      actions.push(
        `<button class="block-btn bg-yellow-400 hover:bg-yellow-500 text-white px-3 py-1 rounded-md text-xs font-medium ml-1" data-id="${user.ID}">Block</button>`
      );
    }
    if (hasPermission("UNBLOCK_USER") && user.Block) {
      actions.push(
        `<button class="unblock-btn bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded-md text-xs font-medium ml-1" data-id="${user.ID}">Unblock</button>`
      );
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
      const user = users.find((u) => u.ID == btn.dataset.id);
      if (user) showEditModal(user);
    });
  });

  document.querySelectorAll(".delete-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      if (!confirm("Are you sure you want to delete this user?")) return;
      try {
        const res = await axios.delete(`${BASE_URL}/user/${btn.dataset.id}`, {
          headers: { Authorization: `Bearer ${accessToken}` },
        });
        alert(res.data.message || "User deleted successfully!");
        fetchUsers();
      } catch (err) {
        alert(err.response?.data?.error || "Failed to delete user.");
      }
    });
  });

  document.querySelectorAll(".block-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      try {
        await axios.put(`${BASE_URL}/userblock/${btn.dataset.id}`, {}, {
          headers: { Authorization: `Bearer ${accessToken}` },
        });
        fetchUsers();
      } catch (err) {
        alert(err.response?.data?.error || "Failed to block user.");
      }
    });
  });

  document.querySelectorAll(".unblock-btn").forEach((btn) => {
    btn.addEventListener("click", async () => {
      try {
        await axios.put(`${BASE_URL}/userunblock/${btn.dataset.id}`, {}, {
          headers: { Authorization: `Bearer ${accessToken}` },
        });
        fetchUsers();
      } catch (err) {
        alert(err.response?.data?.error || "Failed to unblock user.");
      }
    });
  });
}


addUserBtn.addEventListener("click", () => {
  if (!hasPermission("CREATE_USER") && !hasPermission("ADD_USER")) {
    alert("You do not have permission to add users.");
    return;
  }

  dynamicSection.innerHTML = `
    <div class="flex justify-between items-center mb-6">
      <h2 class="text-2xl font-bold text-gray-800">Create New User</h2>
      <button id="backToUsers" 
        class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-medium px-3 py-1.5 rounded-lg transition">
        ← Back
      </button>
    </div>

    <form id="createUserForm" class="grid grid-cols-1 md:grid-cols-2 gap-6 bg-white p-6 rounded-xl shadow-md border border-gray-200">
      <input type="text" id="newFirstname" placeholder="First Name" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400" required>
      <input type="text" id="newLastname" placeholder="Last Name" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400">
      <input type="email" id="newEmail" placeholder="Email" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400" required>
      <input type="text" id="newPhone" placeholder="Phone Number" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400" required>
      <select id="newGender" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400">
        <option value="">Select Gender</option>
        <option>Male</option><option>Female</option><option>Other</option>
      </select>
      <input type="text" id="newPlace" placeholder="Location" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400">
      <input type="password" id="newPassword" placeholder="Initial Password" class="border p-3 rounded-lg focus:ring-2 focus:ring-cyan-400" required>

      <div class="md:col-span-2 flex justify-center">
       <button type="submit" 
          class="bg-cyan-600 hover:bg-cyan-700 text-white font-semibold py-2.5 px-8 rounded-lg shadow-md hover:shadow-lg transition text-base">
          Create User
        </button>
      </div>
    </form>

    <p id="createUserResult" class="mt-4 text-center font-bold"></p>
  `;

  document.getElementById("backToUsers").addEventListener("click", () => location.reload());
  handleCreateUser();
});


function handleCreateUser() {
  const createUserForm = document.getElementById("createUserForm");
  const createUserResult = document.getElementById("createUserResult");

  createUserForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const newUser = {
      firstname: document.getElementById("newFirstname").value.trim(),
      lastname: document.getElementById("newLastname").value.trim(),
      email: document.getElementById("newEmail").value.trim(),
      phone: document.getElementById("newPhone").value.trim(),
      place: document.getElementById("newPlace").value.trim(),
      password: document.getElementById("newPassword").value.trim(),
      gender: document.getElementById("newGender").value,
    };

    try {
      const res = await axios.post(`${BASE_URL}/createuser`, newUser, {
        headers: { Authorization: `Bearer ${accessToken}` },
      });

      createUserResult.textContent = res.data.message || "User added successfully!";
      createUserResult.className = "text-green-600 font-bold";
      createUserForm.reset();
    } catch (err) {
      createUserResult.textContent = err.response?.data?.error || "Failed to add user.";
      createUserResult.className = "text-red-500 font-bold";
    }
  });
}

fetchUsers();
