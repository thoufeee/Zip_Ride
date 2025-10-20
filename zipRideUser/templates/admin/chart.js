// analytics.js

// Logout handler
document.getElementById("logout-btn").addEventListener("click", () => {
  localStorage.removeItem("accessToken");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("role");
  window.location.href = "signin.html";
});

// Color palette
const PRIMARY_COLOR = '#06b6d4'; // cyan
const SECONDARY_COLOR = '#3b82f6'; // blue
const ACCENT_COLOR = '#ef4444'; // red

// Common chart options
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

// 📊 1. Daily Bookings Chart
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

// 📈 2. Weekly Bookings Chart
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

// 📊 3. Monthly Bookings Chart
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

// 💰 4. Monthly Revenue Chart
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
