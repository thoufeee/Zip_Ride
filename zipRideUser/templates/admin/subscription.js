const subscriptionTable = document.getElementById("subscriptionTable");
const formModal = document.getElementById("formModal");
const subscriptionForm = document.getElementById("subscriptionForm");
const createBtn = document.getElementById("createBtn");
const cancelBtn = document.getElementById("cancelBtn");
const formTitle = document.getElementById("formTitle");

let editingId = null;
let allPlans = [];

function getAuthConfig() {
  const token = localStorage.getItem("accessToken");
  if (!token) {
    console.warn("No token found — redirecting to signin");
    window.location.href = "signin.html";
    throw new Error("Access token missing");
  }
  return { headers: { Authorization: `Bearer ${token}` } };
}

function hasPermission(permission) {
  const permissions = JSON.parse(localStorage.getItem("permissions") || "[]");
  return permissions.includes(permission);
}

if (!hasPermission("ACCESS_SUBSCRIPTION")) {
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
  throw new Error("Unauthorized access to Subscription Settings");
}

async function fetchSubscriptions() {
  try {
    const res = await axios.get("http://localhost:8080/admin/subscription/", getAuthConfig());
    console.log("Subscription API response:", res.data);

    let subscriptions = res.data.res || res.data;
    allPlans = subscriptions; 

    subscriptionTable.innerHTML = "";

    subscriptions.forEach(plan => {
      const tr = document.createElement("tr");
      tr.innerHTML = `
        <td class="py-2 px-4">${plan.planname}</td>
        <td class="py-2 px-4">${plan.description}</td>
        <td class="py-2 px-4">${plan.duration_days}</td>
        <td class="py-2 px-4">${plan.price}</td>
        <td class="py-2 px-4">${plan.comission_discount}</td>
        <td class="py-2 px-4 space-x-2">
          <button onclick="editPlan('${plan.id}')" class="bg-yellow-500 px-2 py-1 rounded text-white hover:bg-yellow-600">Edit</button>
          <button onclick="deletePlan('${plan.id}')" class="bg-red-500 px-2 py-1 rounded text-white hover:bg-red-600">Delete</button>
        </td>
      `;
      subscriptionTable.appendChild(tr);
    });
  } catch (error) {
    console.error("Failed to fetch subscriptions:", error);
    alert("Failed to fetch subscriptions. Check console for details.");
  }
}

createBtn.addEventListener("click", () => {
  editingId = null;
  formTitle.innerText = "Create Subscription";
  subscriptionForm.reset();
  formModal.classList.remove("hidden");
  formModal.classList.add("flex");
});

cancelBtn.addEventListener("click", () => {
  formModal.classList.add("hidden");
  formModal.classList.remove("flex");
});


subscriptionForm.addEventListener("submit", async (e) => {
  e.preventDefault();
  const planData = {
    planname: document.getElementById("planName").value,
    description: document.getElementById("description").value,
    duration_days: Number(document.getElementById("duration").value),
    price: Number(document.getElementById("price").value),
    comission_discount: Number(document.getElementById("commission").value),
  };

  try {
    if (editingId) {
      await axios.put(`http://localhost:8080/admin/subscription/${editingId}`, planData, getAuthConfig());
    } else {
      await axios.post("http://localhost:8080/admin/subscription/", planData, getAuthConfig());
    }
    formModal.classList.add("hidden");
    formModal.classList.remove("flex");
    fetchSubscriptions();
  } catch (error) {
    console.error("Failed to save subscription:", error);
    alert("Failed to save subscription. Check console for details.");
  }
});

function editPlan(id) {
  const plan = allPlans.find(p => p.id === id);
  if (!plan) return alert("Plan not found");

  editingId = id;
  formTitle.innerText = "Update Subscription";

  document.getElementById("planName").value = plan.planname;
  document.getElementById("description").value = plan.description;
  document.getElementById("duration").value = plan.duration_days;
  document.getElementById("price").value = plan.price;
  document.getElementById("commission").value = plan.comission_discount;

  formModal.classList.remove("hidden");
  formModal.classList.add("flex");
}

function deletePlan(id) {
  if (confirm("Are you sure you want to delete this plan?")) {
    axios.delete(`http://localhost:8080/admin/subscription/${id}`, getAuthConfig())
      .then(() => fetchSubscriptions())
      .catch(err => {
        console.error("Failed to delete:", err);
        alert("Failed to delete plan.");
      });
  }
}

fetchSubscriptions();
