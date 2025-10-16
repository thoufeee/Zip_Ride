
import { loadNavbarFooter, verifyToken } from "./common.js";


async function verifyAdminAccess() {
  if (!verifyToken()) return false;

  const token = localStorage.getItem("accessToken");

  try {
    const res = await axios.get("http://localhost:8080/admin/verify", {
      headers: { Authorization: `Bearer ${token}` }
    });

    const role = res.data.role;
    const permissions = res.data.permissions || [];

    
    localStorage.setItem("role", role);
    localStorage.setItem("permissions", JSON.stringify(permissions));

    console.log("Logged-in Role:", role);

    if (["SUPER_ADMIN", "MANAGER", "STAFF"].includes(role)) {
      return true;
    } else {
      alert("Access Denied! Unauthorized Role");
      window.location.href = "signin.html";
      return false;
    }
  } catch (err) {
    console.error("Token verification failed:", err);
    alert("Session expired or invalid token");
    window.location.href = "signin.html";
    return false;
  }
}

function loadDashboardData() {
  const role = localStorage.getItem("role");

  
  const dashboardData = {
    users: 152,
    staff: 26,
    managers: 8,
    bookings: [
      { id: "BK001", user: "Arun", service: "Phone Repair", status: "Completed", date: "2025-10-11" },
      { id: "BK002", user: "Rahul", service: "Laptop Repair", status: "Pending", date: "2025-10-10" },
      { id: "BK003", user: "Priya", service: "Tablet Repair", status: "In Progress", date: "2025-10-09" },
      { id: "BK004", user: "Sana", service: "Camera Fix", status: "Cancelled", date: "2025-10-08" },
    ],
  };

  const pageTitle = document.getElementById("page-title");
  const userCountEl = document.getElementById("user-count");
  const staffCountEl = document.getElementById("staff-count");
  const managerCountEl = document.getElementById("manager-count");
  const bookingTableBody = document.getElementById("booking-history-body");

  
  document.querySelector(".summary-section").style.display = "none";
  document.querySelector(".table-section").style.display = "none";

  if (role === "SUPER_ADMIN") {
    pageTitle.textContent = "Welcome, Super Admin";
    document.querySelector(".summary-section").style.display = "flex";
    document.querySelector(".table-section").style.display = "block";

    userCountEl.textContent = dashboardData.users;
    staffCountEl.textContent = dashboardData.staff;
    managerCountEl.textContent = dashboardData.managers;


    bookingTableBody.innerHTML = "";
    dashboardData.bookings.forEach(b => {
      bookingTableBody.innerHTML += `
        <tr>
          <td>${b.id}</td>
          <td>${b.user}</td>
          <td>${b.service}</td>
          <td>${b.status}</td>
          <td>${b.date}</td>
        </tr>
      `;
    });
  } else if (role === "STAFF") {
    pageTitle.textContent = "Welcome, Staff";
    document.querySelector(".summary-section").style.display = "flex";
    document.querySelector(".table-section").style.display = "block";

    userCountEl.textContent = dashboardData.users;
    staffCountEl.parentElement.style.display = "none";     
    managerCountEl.parentElement.style.display = "none";   

    bookingTableBody.innerHTML = "";
    dashboardData.bookings.forEach(b => {
      bookingTableBody.innerHTML += `
        <tr>
          <td>${b.id}</td>
          <td>${b.user}</td>
          <td>${b.service}</td>
          <td>${b.status}</td>
          <td>${b.date}</td>
        </tr>
      `;
    });
  } else if (role === "MANAGER") {
    pageTitle.textContent = "Welcome, Manager";
    document.querySelector(".summary-section").style.display = "flex";
    document.querySelector(".table-section").style.display = "block";

    userCountEl.textContent = dashboardData.users;
    staffCountEl.parentElement.style.display = "none";
    managerCountEl.parentElement.style.display = "none";

    
    bookingTableBody.innerHTML = `
      <tr>
        <td colspan="5">Total Bookings: ${dashboardData.bookings.length}</td>
      </tr>
    `;
  }
}


document.addEventListener("DOMContentLoaded", async () => {
  await loadNavbarFooter();
  const allowed = await verifyAdminAccess();
  if (!allowed) return;

  loadDashboardData();
});
