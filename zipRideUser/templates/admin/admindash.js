// === Load Navbar & Footer ===
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

// === Load Dashboard Data ===
function loadDashboardData() {
  // Example static demo data — replace with Axios API calls later
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

  // Set summary counts
  document.getElementById("user-count").textContent = dashboardData.users;
  document.getElementById("staff-count").textContent = dashboardData.staff;
  document.getElementById("manager-count").textContent = dashboardData.managers;

  // Fill booking history table
  const tbody = document.getElementById("booking-history-body");
  tbody.innerHTML = "";
  dashboardData.bookings.forEach(b => {
    const row = `
      <tr>
        <td>${b.id}</td>
        <td>${b.user}</td>
        <td>${b.service}</td>
        <td>${b.status}</td>
        <td>${b.date}</td>
      </tr>
    `;
    tbody.innerHTML += row;
  });
}

// === Initialize Dashboard ===
document.addEventListener("DOMContentLoaded", async () => {
  await loadReusableComponents();
  loadDashboardData();
});
