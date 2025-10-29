const API_URL = "http://localhost:8080/admin/config/";
const token = localStorage.getItem("accessToken");

function hasPermission(permission) {
  const permissions = JSON.parse(localStorage.getItem("permissions") || "[]");
  return permissions.includes(permission);
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

let currentConfig = {};

const axiosAuth = axios.create({
  baseURL: API_URL,
  headers: {
    Authorization: `Bearer ${token}`,
  },
});

async function loadConfig() {
  try {
    const res = await axiosAuth.get("/");
    const data = res.data.res || {};
    currentConfig = { ...data };

    document.getElementById("site_name").value = data.site_name || "";
    document.getElementById("currency").value = data.currency || "";
    document.getElementById("currency_symbol").value = data.currency_symbol || "";

    document.getElementById("payment_gateway").value = data.payment_gateway || "";
    document.getElementById("payment_publicKey").value = data.payment_publicKey || "";
    document.getElementById("payment_secretKey").value = data.payment_secretKey || "";

    document.getElementById("contact_email").value = data.contact_email || "";
    document.getElementById("contact_phone").value = data.contact_phone || "";
    document.getElementById("contact_address").value = data.contact_address || "";
  } catch (err) {
    console.error("Error loading config:", err);
    if (err.response && err.response.status === 401) {
      alert("Session expired. Please log in again.");
      window.location.href = "signin.html";
    } else {
      alert("Failed to fetch configuration data.");
    }
  }
}


async function updateSection(section) {
  let updateData = {};

  if (section === "currency") {
    const site_name = document.getElementById("site_name").value;
    const currency = document.getElementById("currency").value;
    const currency_symbol = document.getElementById("currency_symbol").value;

    if (site_name !== currentConfig.site_name) updateData.site_name = site_name;
    if (currency !== currentConfig.currency) updateData.currency = currency;
    if (currency_symbol !== currentConfig.currency_symbol)
      updateData.currency_symbol = currency_symbol;

  } else if (section === "payment") {
    const gateway = document.getElementById("payment_gateway").value;
    const pub = document.getElementById("payment_publicKey").value;
    const sec = document.getElementById("payment_secretKey").value;

    if (gateway !== currentConfig.payment_gateway) updateData.payment_gateway = gateway;
    if (pub !== currentConfig.payment_publicKey) updateData.payment_publicKey = pub;
    if (sec !== currentConfig.payment_secretKey) updateData.payment_secretKey = sec;

  } else if (section === "contact") {
    const email = document.getElementById("contact_email").value;
    const phone = document.getElementById("contact_phone").value;
    const addr = document.getElementById("contact_address").value;

    if (email !== currentConfig.contact_email) updateData.contact_email = email;
    if (phone !== currentConfig.contact_phone) updateData.contact_phone = phone;
    if (addr !== currentConfig.contact_address) updateData.contact_address = addr;
  }

  if (Object.keys(updateData).length === 0) {
    alert("No changes detected in this section.");
    return;
  }

  try {
    await axiosAuth.put("/", updateData);
    alert("Configuration updated successfully!");
    await loadConfig();
  } catch (err) {
    console.error("Error updating config:", err);
    if (err.response && err.response.status === 401) {
      alert("Session expired. Please log in again.");
      window.location.href = "signin.html";
    } else {
      alert("Failed to update configuration.");
    }
  }
}


document.addEventListener("DOMContentLoaded", () => {
  loadConfig();
});
