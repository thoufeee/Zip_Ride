
import { loadNavbarFooter, hasPermission, verifyToken, verifyAdminAccess } from "./common.js";

async function initializeCreateAdmin() {

  const roleSelect = document.getElementById("role");
  const permissionsList = document.getElementById("permissions-list");
  const addedPermissions = document.getElementById("added-permissions");
  const form = document.getElementById("createAdminForm");
  const messageContainer = document.getElementById("message-container");
  const result = document.getElementById("result");

  const BASE_URL = "http://localhost:8080";
  const token = localStorage.getItem("accessToken");

  if (!token) {
    alert("Unauthorized! Please log in.");
    window.location.href = "signin.html";
    return;
  }

if (!hasPermission("ADD_STAFF")) {
  const formParent = form.parentElement; 
  if (formParent) {
    formParent.innerHTML = "<p style='color:red; text-align:center;'>You don’t have permission to access this page.</p>";
  } else {
    form.style.display = "none";
    console.error("You don’t have permission to access this page.");
  }
  return;
}
  const api = axios.create({
    baseURL: BASE_URL,
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json"
    }
  });

 
  async function fetchRoles() {
    try {
      const res = await api.get("/admin/allroles");
      const roles = res.data.res || [];
      roleSelect.innerHTML = `<option value="">-- Select Role --</option>`;
      roles.forEach(role => {
        const name = role.name || role.role_name || role.Name;
        if (name) {
          const option = document.createElement("option");
          option.value = name;
          option.textContent = name;
          roleSelect.appendChild(option);
        }
      });
    } catch (err) {
      console.error("Failed to fetch roles:", err);
      roleSelect.innerHTML = `<option value="">Error loading roles</option>`;
    }
  }

  async function fetchPermissions() {
    try {
      const res = await api.get("/admin/allpermissions");
      const permissions = res.data.permissions || res.data.res || res.data.data || [];
      permissionsList.innerHTML = "";

      if (!permissions.length) {
        permissionsList.innerHTML = "<p style='color:red;'>No permissions found.</p>";
        return;
      }

      permissions.forEach(perm => {
        const permName = typeof perm === "string" ? perm : (perm.name || perm.permission_name);
        if (permName) {
          const div = document.createElement("div");
          div.classList.add("permission-item");
          div.textContent = permName;
          div.addEventListener("click", () => addPermission(permName));
          permissionsList.appendChild(div);
        }
      });
    } catch (err) {
      console.error("Failed to fetch permissions:", err);
      permissionsList.innerHTML = "<p style='color:red;'>Failed to load permissions.</p>";
    }
  }

 
  function addPermission(name) {
    const placeholder = addedPermissions.querySelector(".placeholder");
    if (placeholder) placeholder.remove();
    if ([...addedPermissions.children].some(p => p.textContent === name)) return;

    const tag = document.createElement("span");
    tag.classList.add("perm-tag");
    tag.textContent = name;
    tag.title = "Click to remove";
    tag.addEventListener("click", () => tag.remove());
    addedPermissions.appendChild(tag);
  }

  
  form.addEventListener("submit", async e => {
    e.preventDefault();
    const selectedPermissions = [...addedPermissions.children].map(p => p.textContent);

    const payload = {
      name: form.name.value,
      email: form.email.value,
      phonenumber: form.phone.value,
      password: form.password.value,
      role: form.role.value,
      extra_permissions: selectedPermissions
    };

    try {
      await api.post("/admin/createstaff", payload);
      result.textContent = " Admin created successfully!";
      result.style.color = "green";
      form.reset();
      addedPermissions.innerHTML = `<p class="placeholder">No permissions selected yet.</p>`;
    } catch (err) {
      console.error(err);
      result.textContent = err.response?.data?.err || " Failed to create admin.";
      result.style.color = "red";
    }
  });

  
  await Promise.all([fetchRoles(), fetchPermissions()]);
}


document.addEventListener("DOMContentLoaded", async () => {
  await loadNavbarFooter();

  const tokenExists = verifyToken();
  if (!tokenExists) return;

  initializeCreateAdmin();
});
