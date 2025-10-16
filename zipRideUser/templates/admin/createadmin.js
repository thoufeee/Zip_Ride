// Load reusable components: navbar and footer
async function loadReusableComponents() {
  try {
    const navbarRes = await fetch("navbar.html");
    document.getElementById("navbar").innerHTML = await navbarRes.text();

    const footerRes = await fetch("footer.html");
    document.getElementById("footer").innerHTML = await footerRes.text();
  } catch (err) {
    console.error("Failed to load reusable components:", err);
  }
}

export async function initializeCreateAdmin() {
  await loadReusableComponents();

  const roleSelect = document.getElementById("role");
  const permissionsList = document.getElementById("permissions-list");
  const addedPermissions = document.getElementById("added-permissions");
  const form = document.getElementById("createAdminForm");
  const result = document.getElementById("result");

  const BASE_URL = "http://localhost:8080";
  const token = localStorage.getItem("accessToken");

  if (!token) {
    alert("Unauthorized! Please log in.");
    window.location.href = "signin.html";
    return;
  }

  const api = axios.create({
    baseURL: BASE_URL,
    headers: { 
      "Authorization": `Bearer ${token}`, 
      "Content-Type": "application/json" 
    }
  });

  // Fetch roles from backend
  async function fetchRoles() {
    try {
      const res = await api.get("/admin/allroles");
      const roles = res.data || [];
      roleSelect.innerHTML = `<option value="">-- Select Role --</option>`;
      roles.forEach(role => {
        const option = document.createElement("option");
        option.value = role.name || role.role_name;
        option.textContent = role.name || role.role_name;
        roleSelect.appendChild(option);
      });
    } catch (err) {
      console.error("Failed to fetch roles:", err);
      roleSelect.innerHTML = `<option value="">Error loading roles</option>`;
    }
  }

  // Fetch permissions from backend
  async function fetchPermissions() {
    try {
      const res = await api.get("/admin/allpermissions");
      const permissions = res.data || [];
      permissionsList.innerHTML = "";
      permissions.forEach(perm => {
        const div = document.createElement("div");
        div.classList.add("permission-item");
        div.textContent = perm.name || perm.permission_name;
        div.addEventListener("click", () => addPermission(perm.name || perm.permission_name));
        permissionsList.appendChild(div);
      });
    } catch (err) {
      console.error("Failed to fetch permissions:", err);
      permissionsList.innerHTML = "<p style='color:red;'>Failed to load permissions.</p>";
    }
  }

  // Add/remove permission locally
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

  // Submit form to backend
  form.addEventListener("submit", async e => {
    e.preventDefault();
    const selectedPermissions = [...addedPermissions.children].map(p => p.textContent);

    const payload = {
      name: form.name.value,
      email: form.email.value,
      phone: form.phone.value,
      password: form.password.value,
      role: form.role.value,
      permissions: selectedPermissions
    };

    try {
      await api.post("/admins", payload);
      result.textContent = "✅ Admin created successfully!";
      result.style.color = "green";
      form.reset();
      addedPermissions.innerHTML = `<p class="placeholder">No permissions selected yet.</p>`;
    } catch (err) {
      console.error(err);
      result.textContent = err.response?.data?.err || "❌ Failed to create admin.";
      result.style.color = "red";
    }
  });

  // Initialize page
  await Promise.all([fetchRoles(), fetchPermissions()]);
}

// Auto-run on page load
document.addEventListener("DOMContentLoaded", () => {
  initializeCreateAdmin();
});
