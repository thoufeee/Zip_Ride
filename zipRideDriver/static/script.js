/* =====================================================
   MOCK DATA SOURCES
   ===================================================== */
const mockDashboardStats = [
  { id: 'activeDrivers', label: 'Active Drivers', value: 1287, icon: 'fa-user-tie', delta: '+3.8%' },
  { id: 'ongoingRides', label: 'Ongoing Rides', value: 312, icon: 'fa-road', delta: '+1.2%' },
  { id: 'totalEarnings', label: 'Total Earnings', value: '$842K', icon: 'fa-sack-dollar', delta: '+12.4%' },
  { id: 'totalDistance', label: 'Total Distance', value: '1.2M km', icon: 'fa-map-location-dot', delta: '+5.6%' }
];

const mockDrivers = [
  { id: 1, name: 'Sophia Turner', phone: '+1 415 555 3201', vehicle: 'Tesla Model 3', status: 'active', rating: 4.9 },
  { id: 2, name: 'Ethan Walker', phone: '+1 415 555 4218', vehicle: 'Toyota Prius', status: 'offline', rating: 4.6 },
  { id: 3, name: 'Olivia Chen', phone: '+1 415 555 2251', vehicle: 'Honda Civic', status: 'active', rating: 4.8 },
  { id: 4, name: 'Noah Patel', phone: '+1 415 555 2873', vehicle: 'Ford Explorer', status: 'suspended', rating: 4.1 },
  { id: 5, name: 'Ava Williams', phone: '+1 415 555 7812', vehicle: 'Chevy Bolt', status: 'active', rating: 4.7 }
];

const mockVehicles = [
  { id: 1, model: 'Tesla Model 3', plate: 'ZPR-2031', type: 'EV Sedan', driver: 'Sophia Turner' },
  { id: 2, model: 'Toyota Prius', plate: 'HYM-9412', type: 'Hybrid', driver: 'Ethan Walker' },
  { id: 3, model: 'Ford Explorer', plate: 'NRG-5529', type: 'SUV', driver: 'Noah Patel' },
  { id: 4, model: 'Honda Civic', plate: 'KGL-7730', type: 'Compact', driver: 'Olivia Chen' }
];

const mockRides = [
  { id: 1, driver: 'Sophia Turner', rider: 'Marcus Reed', fare: '$24.60', route: 'Market St → Castro St', status: 'En Route' },
  { id: 2, driver: 'Ava Williams', rider: 'Nina Rossi', fare: '$18.10', route: 'Union Sq → Fisherman’s Wharf', status: 'Picked Up' },
  { id: 3, driver: 'Ethan Walker', rider: 'David Chen', fare: '$31.40', route: 'SOMA → SFO', status: 'Drop-off' }
];

const mockPayments = [
  { period: 'This Week', rides: 1456, gross: '$178,420', net: '$134,210' },
  { period: 'This Month', rides: 5901, gross: '$698,540', net: '$520,330' },
  { period: 'Last Month', rides: 5703, gross: '$661,210', net: '$498,750' }
];

const mockAlerts = [
  { type: 'info', message: 'Scheduled maintenance for fleet analytics tonight at 11 PM PST.' },
  { type: 'danger', message: 'Driver ID #421 flagged for repeated cancellation. Review pending.' }
];

const mockReviews = [
  { driver: 'Sophia Turner', rating: 5, feedback: 'Outstanding ride – clean car, smooth driving!' },
  { driver: 'Noah Patel', rating: 3.5, feedback: 'Friendly driver but arrived 8 mins late.' }
];

const mockRatings = [
  { driver: 'Top Rated Drivers', score: '4.94 ★', info: 'Averaged across last 30 days' },
  { driver: 'Support Resolution', score: '93%', info: 'Tickets closed within SLA' }
];

/* =====================================================
   DOM HELPERS
   ===================================================== */
const qs = (selector, parent = document) => parent.querySelector(selector);
const qsa = (selector, parent = document) => [...parent.querySelectorAll(selector)];

/* =====================================================
   NAVIGATION & SECTION HANDLING
   ===================================================== */
const navLinks = qsa('.nav-link');
const sections = qsa('.content-section');

navLinks.forEach(link => {
  link.addEventListener('click', evt => {
    evt.preventDefault();
    const target = link.dataset.section;

    navLinks.forEach(l => l.classList.toggle('active', l === link));
    sections.forEach(section => section.classList.toggle('active', section.id === target));

    if (window.innerWidth <= 992) {
      sidebar.classList.remove('open');
    }
  });
});

/* =====================================================
   SIDEBAR & PROFILE INTERACTIONS
   ===================================================== */
const sidebar = qs('#sidebar');
const sidebarToggle = qs('#sidebarToggle');
const profileTrigger = qs('#profileTrigger');
const profileDropdown = qs('#profileDropdown');
const themeToggle = qs('#themeToggle');
const appShell = qs('.app-shell');

sidebarToggle.addEventListener('click', () => {
  sidebar.classList.toggle('open');
});

profileTrigger.addEventListener('click', () => {
  profileDropdown.classList.toggle('active');
});

document.addEventListener('click', evt => {
  if (!profileTrigger.contains(evt.target)) {
    profileDropdown.classList.remove('active');
  }
});

themeToggle.addEventListener('click', () => {
  const isDark = appShell.dataset.theme === 'dark';
  appShell.dataset.theme = isDark ? 'light' : 'dark';
  themeToggle.innerHTML = isDark ? '<i class="fas fa-moon"></i>' : '<i class="fas fa-sun"></i>';
});

/* =====================================================
   DASHBOARD STAT CARDS
   ===================================================== */
const dashboardStatsContainer = qs('#dashboardStats');

const renderDashboardStats = () => {
  dashboardStatsContainer.innerHTML = mockDashboardStats
    .map(({ label, value, icon, delta }) => `
      <article class="stat-card">
        <h3>${label}</h3>
        <div class="value">${value}</div>
        <div class="stat-meta">
          <span class="badge success">${delta}</span>
          <i class="fas ${icon}"></i>
        </div>
      </article>
    `)
    .join('');
};

qs('#refreshStats').addEventListener('click', () => {
  renderDashboardStats();
});

/* =====================================================
   CHART INITIALISATION (Chart.js)
   ===================================================== */
let ridesChart;
let earningsChart;
let revenueChart;

const initCharts = () => {
  const ridesCtx = qs('#ridesChart');
  const earningsCtx = qs('#earningsChart');
  const revenueCtx = qs('#revenueChart');

  ridesChart = new Chart(ridesCtx, {
    type: 'line',
    data: {
      labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
      datasets: [{
        label: 'Ride Volume',
        data: [420, 510, 488, 540, 620, 678, 712],
        borderColor: '#38bdf8',
        backgroundColor: 'rgba(56, 189, 248, 0.15)',
        tension: 0.4,
        fill: true
      }]
    },
    options: {
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: 'var(--color-text-muted)' }, grid: { display: false } },
        y: { ticks: { color: 'var(--color-text-muted)' }, grid: { color: 'rgba(148, 163, 184, 0.1)' } }
      }
    }
  });

  earningsChart = new Chart(earningsCtx, {
    type: 'bar',
    data: {
      labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
      datasets: [{
        label: 'Earnings',
        data: [180, 210, 240, 268],
        backgroundColor: ['#8b5cf6', '#6366f1', '#38bdf8', '#22c55e'],
        borderRadius: 12
      }]
    },
    options: {
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: 'var(--color-text-muted)' }, grid: { display: false } },
        y: { ticks: { color: 'var(--color-text-muted)' }, grid: { color: 'rgba(148, 163, 184, 0.12)' } }
      }
    }
  });

  revenueChart = new Chart(revenueCtx, {
    type: 'line',
    data: {
      labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul'],
      datasets: [{
        label: 'Revenue',
        data: [510, 580, 610, 645, 702, 735, 780],
        borderColor: '#6366f1',
        backgroundColor: 'rgba(99, 102, 241, 0.15)',
        tension: 0.4,
        fill: true
      }]
    },
    options: {
      plugins: { legend: { display: false } },
      scales: {
        x: { ticks: { color: 'var(--color-text-muted)' }, grid: { display: false } },
        y: { ticks: { color: 'var(--color-text-muted)' }, grid: { color: 'rgba(148, 163, 184, 0.12)' } }
      }
    }
  });
};

/* =====================================================
   DRIVER TABLE RENDER & FILTERING
   ===================================================== */
const driversTableBody = qs('#driversTableBody');

const renderDrivers = (drivers = mockDrivers) => {
  driversTableBody.innerHTML = drivers
    .map(({ name, phone, vehicle, status, rating, id }) => `
      <tr>
        <td>${name}</td>
        <td>${phone}</td>
        <td>${vehicle}</td>
        <td><span class="tag ${status}">${status}</span></td>
        <td><span class="rating"><i class="fas fa-star"></i>${rating.toFixed(1)}</span></td>
        <td>
          <div class="table-actions">
            <button class="btn ghost" data-role="edit-driver" data-id="${id}"><i class="fas fa-pen"></i>Edit</button>
            <button class="btn danger" data-role="deactivate-driver" data-id="${id}"><i class="fas fa-user-slash"></i>Deactivate</button>
          </div>
        </td>
      </tr>
    `)
    .join('');
};

const filterDrivers = () => {
  const searchTerm = qs('#driverSearch').value.toLowerCase().trim();
  const statusFilter = qs('#driverStatusFilter').value;
  const ratingFilter = qs('#driverRatingFilter').value;

  let filtered = mockDrivers.filter(driver => driver.name.toLowerCase().includes(searchTerm));

  if (statusFilter !== 'all') {
    filtered = filtered.filter(driver => driver.status === statusFilter);
  }

  if (ratingFilter !== 'all') {
    const threshold = Number(ratingFilter) / 10; // Ratings stored as 4.5 etc.
    filtered = filtered.filter(driver => driver.rating >= threshold);
  }

  renderDrivers(filtered);
};

['driverSearch', 'driverStatusFilter', 'driverRatingFilter'].forEach(id => {
  qs(`#${id}`).addEventListener('input', filterDrivers);
});

/* =====================================================
   VEHICLES TABLE RENDER
   ===================================================== */
const vehiclesTableBody = qs('#vehiclesTableBody');

const renderVehicles = () => {
  vehiclesTableBody.innerHTML = mockVehicles
    .map(({ model, plate, type, driver }) => `
      <tr>
        <td>${model}</td>
        <td>${plate}</td>
        <td>${type}</td>
        <td>${driver}</td>
        <td>
          <div class="avatar-placeholder">IMG</div>
        </td>
        <td>
          <div class="table-actions">
            <button class="btn ghost"><i class="fas fa-pen"></i>Edit</button>
            <button class="btn danger"><i class="fas fa-trash"></i>Remove</button>
          </div>
        </td>
      </tr>
    `)
    .join('');
};

/* =====================================================
   RIDES TABLE RENDER
   ===================================================== */
const ridesTableBody = qs('#ridesTableBody');

const renderRides = () => {
  ridesTableBody.innerHTML = mockRides
    .map(({ driver, rider, fare, route, status }) => `
      <tr>
        <td>${driver}</td>
        <td>${rider}</td>
        <td>${fare}</td>
        <td>${route}</td>
        <td><span class="badge">${status}</span></td>
      </tr>
    `)
    .join('');
};

qs('#refreshRides').addEventListener('click', renderRides);

/* =====================================================
   PAYMENTS TABLE RENDER
   ===================================================== */
const paymentsTableBody = qs('#paymentsTableBody');

const renderPayments = () => {
  paymentsTableBody.innerHTML = mockPayments
    .map(({ period, rides, gross, net }) => `
      <tr>
        <td>${period}</td>
        <td>${rides}</td>
        <td>${gross}</td>
        <td>${net}</td>
      </tr>
    `)
    .join('');
};

qs('#exportPayments').addEventListener('click', () => {
  alert('Report exported. (Mock action)');
});

/* =====================================================
   NOTIFICATIONS & REVIEWS RENDER
   ===================================================== */
const systemAlertsEl = qs('#systemAlerts');
const recentReviewsEl = qs('#recentReviews');
const driverRatingsEl = qs('#driverRatings');

const renderAlerts = () => {
  systemAlertsEl.innerHTML = `
    <h2>System Alerts</h2>
    ${mockAlerts.map(({ type, message }) => `<div class="alert-item ${type}">${message}</div>`).join('')}
  `;

  recentReviewsEl.innerHTML = `
    <h2>Recent Reviews</h2>
    ${mockReviews.map(({ driver, rating, feedback }) => `
      <div class="alert-item">
        <strong>${driver}</strong>
        <div class="rating"><i class="fas fa-star"></i>${rating.toFixed(1)}</div>
        <p>${feedback}</p>
      </div>
    `).join('')}
  `;

  driverRatingsEl.innerHTML = `
    <h2>Experience Metrics</h2>
    ${mockRatings.map(({ driver, score, info }) => `
      <div class="alert-item success">
        <strong>${driver}</strong>
        <p>${score}</p>
        <small>${info}</small>
      </div>
    `).join('')}
  `;
};

/* =====================================================
   MODAL HANDLING HELPERS
   ===================================================== */
const driverModal = qs('#driverModal');
const vehicleModal = qs('#vehicleModal');
const modalBackdrop = qs('#modalBackdrop');

const openModal = modal => {
  modal.classList.add('active');
  modalBackdrop.classList.add('active');
  modal.setAttribute('aria-hidden', 'false');
};

const closeModal = modal => {
  modal.classList.remove('active');
  modalBackdrop.classList.remove('active');
  modal.setAttribute('aria-hidden', 'true');
};

qs('#addDriverBtn').addEventListener('click', () => openModal(driverModal));
qs('#driverModalClose').addEventListener('click', () => closeModal(driverModal));
qs('#driverModalCancel').addEventListener('click', () => closeModal(driverModal));

qs('#addVehicleBtn').addEventListener('click', () => openModal(vehicleModal));
qs('#vehicleModalClose').addEventListener('click', () => closeModal(vehicleModal));
qs('#vehicleModalCancel').addEventListener('click', () => closeModal(vehicleModal));

modalBackdrop.addEventListener('click', () => {
  closeModal(driverModal);
  closeModal(vehicleModal);
});

qs('#driverForm').addEventListener('submit', evt => {
  evt.preventDefault();
  alert('Driver saved (mock).');
  closeModal(driverModal);
});

qs('#vehicleForm').addEventListener('submit', evt => {
  evt.preventDefault();
  alert('Vehicle saved (mock).');
  closeModal(vehicleModal);
});

/* =====================================================
   INITIALISATION
   ===================================================== */
window.addEventListener('DOMContentLoaded', () => {
  renderDashboardStats();
  initCharts();
  renderDrivers();
  renderVehicles();
  renderRides();
  renderPayments();
  renderAlerts();
});
