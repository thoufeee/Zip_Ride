const BASE_URL = "http://localhost:8080/admin";

const logoutBtn = document.getElementById("logout-btn");
const form = document.getElementById("createAdminForm");
const permissionsList = document.getElementById("permissions-list");
const addedPermissions = document.getElementById("added-permissions");
const result = document.getElementById("result");

const accessToken = localStorage.getItem("accessToken");
const userPermissions = JSON.parse(localStorage.getItem("permissions")) || [];

if (!accessToken) {
  window.location.href = "signin.html";
}

function hasPermission(requiredPermission) {
  if (!userPermissions || userPermissions.length === 0) return false;
  return userPermissions.some(
    (p) => p.toUpperCase() === requiredPermission.toUpperCase()
  );
}

if (!hasPermission("ADD_STAFF")) {
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


// Logout
logoutBtn.addEventListener("click", () => {
  localStorage.clear();
  window.location.href = "signin.html";
});

// Axios interceptor for 401
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

// Fetch permissions from server
async function fetchPermissions() {
  try {
    const res = await axios.get(`${BASE_URL}/allpermissions`, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    const permissions = res.data.res || [];
    if (permissions.length === 0) {
      permissionsList.innerHTML =
        `<p class="text-gray-400 text-sm italic">No permissions available.</p>`;
      return;
    }

    renderPermissions(permissions);
  } catch (err) {
    console.error("Error loading permissions:", err);
    permissionsList.innerHTML =
      `<p class="text-red-500 text-sm">Failed to load permissions.</p>`;
  }
}

// Render available permissions
function renderPermissions(permissions) {
  permissionsList.innerHTML = ""; 
  permissions.forEach((perm) => {
    const tag = document.createElement("span");
    tag.textContent = perm.name.replace(/_/g, " "); 
    tag.className =
      "cursor-pointer bg-cyan-100 text-cyan-800 text-sm px-3 py-1 rounded-full font-medium hover:bg-cyan-200 transition";
    tag.addEventListener("click", () => addPermission(perm.name));
    permissionsList.appendChild(tag);
  });
}

// Add permission to selected list
function addPermission(permissionName) {
  // Remove placeholder if exists
  const placeholder = addedPermissions.querySelector("p");
  if (placeholder) placeholder.remove();

  // Check if permission already added
  if (
    [...addedPermissions.children].some(
      (el) => el.tagName === "SPAN" && el.textContent.replace(/\s/g, "_") === permissionName
    )
  ) return;

  const tag = document.createElement("span");
  tag.textContent = permissionName.replace(/_/g, " "); 
  tag.className =
    "bg-cyan-600 text-white text-sm px-3 py-1 rounded-full font-medium cursor-pointer hover:bg-cyan-700 transition";
  tag.addEventListener("click", () => {
    tag.remove();
    showPlaceholderIfEmpty();
  });

  addedPermissions.appendChild(tag);
}

// Show placeholder if no permissions are selected
function showPlaceholderIfEmpty() {
  if (![...addedPermissions.children].some(el => el.tagName === "SPAN")) {
    addedPermissions.innerHTML = `<p class="text-gray-400 text-sm italic">
      Click an available permission to add it.
    </p>`;
  }
}

// Handle form submission
form.addEventListener("submit", async (e) => {
  e.preventDefault();
  result.textContent = "";
  result.className = "";

  const name = document.getElementById("name").value.trim();
  const email = document.getElementById("email").value.trim();
  const phone = document.getElementById("phone").value.trim();
  const password = document.getElementById("password").value.trim();

  // Only include actual SPAN tags as permissions
  const permissions = [...addedPermissions.children]
    .filter(el => el.tagName === "SPAN")
    .map(el => el.textContent.replace(/\s/g, "_"));

  if (!name || !email || !phone || !password) {
    result.textContent = "Please fill all fields.";
    result.className = "text-red-600 text-sm";
    return;
  }

  if (permissions.length === 0) {
    result.textContent = "Select at least one permission.";
    result.className = "text-red-600 text-sm";
    return;
  }

  const payload = {
    name,
    email,
    phonenumber: phone,
    password,
    extra_permissions: permissions,
  };

  try {
    const res = await axios.post(`${BASE_URL}/createstaff`, payload, {
      headers: { Authorization: `Bearer ${accessToken}` },
    });

    result.textContent = res.data.res || "Account created successfully!";
    result.className = "text-green-600 font-semibold text-center pt-4";

    form.reset();
    addedPermissions.innerHTML = `<p class="text-gray-400 text-sm italic">
      Click an available permission to add it.
    </p>`;
  } catch (err) {
    console.error("Error creating account:", err);
    result.textContent =
      err.response?.data?.err || "Failed to create account.";
    result.className = "text-red-600 text-sm";
  }
});

// Initialize
fetchPermissions();
