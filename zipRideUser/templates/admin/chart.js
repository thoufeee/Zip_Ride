
document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
  localStorage.removeItem("permissions"); 
  window.location.href = "signin.html";
});

function getPermissions() {
  try {
    
    const perms = JSON.parse(localStorage.getItem("permissions"));
    if (Array.isArray(perms)) return perms;

    const token = localStorage.getItem("accessToken");
    if (!token) return [];
    const payload = JSON.parse(atob(token.split(".")[1]));
    return payload.permissions || [];
  } catch {
    return [];
  }
}

function hasPermission(permission) {
  const perms = getPermissions();
  return perms.includes(permission);
}

if (!hasPermission("VIEW_ANALYTICS")) {
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
  throw new Error("Unauthorized access to Analytics page");
}

const PRIMARY_COLOR = '#06b6d4'; 
const SECONDARY_COLOR = '#3b82f6'; 
const ACCENT_COLOR = '#ef4444'; 

const chartOptions = {
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    legend: {
      position: 'top',
      labels: { color: '#4b5563' }
    },
    tooltip: {
      backgroundColor: 'rgba(30,41,59,0.8)',
      titleColor: '#fff',
      bodyColor: '#fff'
    }
  },
  scales: {
    y: {
      beginAtZero: true,
      ticks: { color: '#6b7280' },
      grid: { color: '#f3f4f6' }
    },
    x: {
      ticks: { color: '#6b7280' },
      grid: { display: false }
    }
  }
};

new Chart(document.getElementById('dailyBookingsChart').getContext('2d'), {
  type: 'bar',
  data: {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    datasets: [{
      label: 'Bookings',
      data: [5, 8, 6, 10, 9, 7, 12],
      backgroundColor: PRIMARY_COLOR,
      borderRadius: 4
    }]
  },
  options: chartOptions
});

new Chart(document.getElementById('weeklyBookingsChart').getContext('2d'), {
  type: 'line',
  data: {
    labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
    datasets: [{
      label: 'Bookings',
      data: [12, 19, 8, 15],
      borderColor: SECONDARY_COLOR,
      backgroundColor: 'rgba(59,130,246,0.3)',
      tension: 0.4,
      fill: true,
      pointRadius: 5,
      pointHoverRadius: 7
    }]
  },
  options: chartOptions
});

new Chart(document.getElementById('monthlyBookingsChart').getContext('2d'), {
  type: 'bar',
  data: {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    datasets: [{
      label: 'Total Bookings',
      data: [30, 45, 28, 60, 50],
      backgroundColor: PRIMARY_COLOR,
      borderRadius: 4
    }]
  },
  options: chartOptions
});

new Chart(document.getElementById('revenueChart').getContext('2d'), {
  type: 'line',
  data: {
    labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May'],
    datasets: [{
      label: 'Revenue (₹)',
      data: [1200, 1900, 900, 2500, 2200],
      borderColor: ACCENT_COLOR,
      backgroundColor: 'rgba(239,68,68,0.3)',
      tension: 0.4,
      fill: true,
      pointRadius: 5,
      pointHoverRadius: 7
    }]
  },
  options: chartOptions
});
