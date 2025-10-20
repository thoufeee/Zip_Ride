// ==================== GLOBAL VARIABLES ====================
let currentBookingIdToAssign = null;

const tableBody = document.getElementById('bookingTable');
const assignDriverModal = document.getElementById('assignDriverModal');
const driverNameInput = document.getElementById('driverNameInput');
const modalBookingId = document.getElementById('modalBookingId');

// ==================== LOGOUT HANDLER ====================
document.getElementById("logout-btn").addEventListener("click", () => {
    localStorage.clear();
    window.location.href = "signin.html";
});

// ==================== SAMPLE DATA ====================
const bookings = [
    {id:101,user:'John Doe',phone:'+1 (123) 456-7890',vehicle:'Car',pickupLoc:'Airport Terminal 1',dropLoc:'Grand Hyatt Hotel',pickup:'2025-10-18 10:00',status:'Pending',driver:'',price:500,payment:'Cash',commission:50},
    {id:102,user:'Jane Smith',phone:'+1 (234) 567-8901',vehicle:'Bike',pickupLoc:'Main Street Station',dropLoc:'Tech Park Office',pickup:'2025-10-19 14:00',status:'Accepted',driver:'',price:200,payment:'Online',commission:20},
    {id:103,user:'Alice Brown',phone:'+1 (345) 678-9012',vehicle:'Rickshaw',pickupLoc:'Central Market',dropLoc:'Residential Area A',pickup:'2025-10-20 08:30',status:'Driver Assigned',driver:'Suresh B.',price:300,payment:'Cash',commission:30},
    {id:104,user:'Bob Johnson',phone:'+1 (456) 789-0123',vehicle:'Car',pickupLoc:'Suburban Mall',dropLoc:'Downtown Plaza',pickup:'2025-10-21 17:00',status:'Completed',driver:'Mike R.',price:750,payment:'Online',commission:75},
    {id:105,user:'Charlie Day',phone:'+1 (567) 890-1234',vehicle:'Bike',pickupLoc:'Industrial Zone',dropLoc:'City Center',pickup:'2025-10-22 12:00',status:'Cancelled',driver:'',price:150,payment:'Cash',commission:0},
];

// ==================== STATUS UPDATE FUNCTION ====================
function updateBookingStatus(id, newStatus, driverName = '') {
    const booking = bookings.find(b => b.id === id);
    if (booking) {
        booking.status = newStatus;
        if (newStatus === 'Driver Assigned') booking.driver = driverName;
        console.log(`Booking #${id} updated to: ${newStatus}`);
        renderTable();
    }
}

// ==================== RENDER TABLE ====================
function renderTable() {
    tableBody.innerHTML = '';
    bookings.forEach(b => {
        const tr = document.createElement('tr');
        tr.className = "hover:bg-gray-50 transition duration-150";

        let statusClass;
        switch (b.status) {
            case 'Pending': statusClass = 'bg-yellow-100 text-yellow-700 border border-yellow-200'; break;
            case 'Accepted': statusClass = 'bg-blue-100 text-blue-700 border border-blue-200'; break;
            case 'Driver Assigned': statusClass = 'bg-purple-100 text-purple-700 border border-purple-200'; break;
            case 'Completed': statusClass = 'bg-green-100 text-green-700 border border-green-200'; break;
            case 'Cancelled': statusClass = 'bg-red-100 text-red-700 border border-red-200'; break;
            default: statusClass = 'bg-gray-100 text-gray-600 border border-gray-200';
        }

        let actionsHTML = '';
        if (b.status === 'Pending') {
            actionsHTML = `
                <button data-id="${b.id}" class="bg-green-600 text-white px-3 py-1 rounded-full text-xs hover:bg-green-700 transition duration-150 accept-btn shadow-sm font-semibold">Accept</button>
                <button data-id="${b.id}" class="bg-red-600 text-white px-3 py-1 rounded-full text-xs hover:bg-red-700 transition duration-150 cancel-btn shadow-sm font-semibold">Cancel</button>
            `;
        } else if (b.status === 'Accepted') {
            actionsHTML = `<button data-id="${b.id}" class="bg-cyan-600 text-white px-3 py-1 rounded-full text-xs hover:bg-cyan-700 transition duration-150 assign-btn shadow-sm font-semibold">Assign Driver</button>`;
        } else if (b.status === 'Driver Assigned') {
            actionsHTML = `<button data-id="${b.id}" class="bg-purple-600 text-white px-3 py-1 rounded-full text-xs hover:bg-purple-700 transition duration-150 complete-btn shadow-sm font-semibold">Complete</button>`;
        } else {
            actionsHTML = `<span class="text-gray-400 text-xs">No Action</span>`;
        }

        tr.innerHTML = `
            <td class="px-4 md:px-6 py-4 font-bold text-gray-900">#${b.id}</td>
            <td class="px-4 md:px-6 py-4 text-gray-700 whitespace-nowrap">${b.user}<br><span class="text-gray-500 text-xs">${b.phone}</span></td>
            <td class="px-4 md:px-6 py-4 text-gray-700">${b.vehicle}</td>
            <td class="px-4 md:px-6 py-4 text-gray-500 whitespace-nowrap">${b.pickupLoc.split(',')[0]} → ${b.dropLoc.split(',')[0]}</td>
            <td class="px-4 md:px-6 py-4 text-gray-700 whitespace-nowrap">${b.pickup}</td>
            <td class="px-4 md:px-6 py-4 text-center">
                <span class="px-3 py-1 inline-flex rounded-full text-xs font-semibold ${statusClass}">${b.status}</span>
            </td>
            <td class="px-4 md:px-6 py-4 text-gray-700 whitespace-nowrap"><span class="font-bold text-green-600">₹${b.price}</span></td>
            <td class="px-4 md:px-6 py-4 text-gray-700 whitespace-nowrap">${b.driver || '<span class="text-red-500">Unassigned</span>'}</td>
            <td class="px-4 md:px-6 py-4 text-center space-y-1 md:space-x-1 whitespace-nowrap flex flex-col md:flex-row items-center justify-center">
                ${actionsHTML}
            </td>
        `;
        tableBody.appendChild(tr);

        // Buttons handler
        if (b.status === 'Pending') {
            tr.querySelector('.accept-btn').addEventListener('click', () => updateBookingStatus(b.id, 'Accepted'));
            tr.querySelector('.cancel-btn').addEventListener('click', () => {
                if (confirm(`Cancel booking #${b.id}?`)) updateBookingStatus(b.id, 'Cancelled');
            });
        } else if (b.status === 'Accepted') {
            tr.querySelector('.assign-btn').addEventListener('click', () => {
                currentBookingIdToAssign = b.id;
                modalBookingId.textContent = `#${b.id}`;
                driverNameInput.value = '';
                assignDriverModal.classList.remove('hidden');
                assignDriverModal.classList.add('flex');
            });
        } else if (b.status === 'Driver Assigned') {
            tr.querySelector('.complete-btn').addEventListener('click', () => updateBookingStatus(b.id, 'Completed'));
        }
    });
}

// ==================== MODAL HANDLERS ====================
document.getElementById('confirmAssignBtn').addEventListener('click', () => {
    const driverName = driverNameInput.value.trim();
    if (!driverName) {
        alert("Please enter the driver's name.");
        driverNameInput.focus();
        return;
    }

    if (currentBookingIdToAssign) {
        updateBookingStatus(currentBookingIdToAssign, 'Driver Assigned', driverName);
        assignDriverModal.classList.add('hidden');
        assignDriverModal.classList.remove('flex');
        currentBookingIdToAssign = null;
    }
});

document.getElementById('closeAssignModal').addEventListener('click', () => {
    assignDriverModal.classList.add('hidden');
    assignDriverModal.classList.remove('flex');
    currentBookingIdToAssign = null;
});

// ==================== INITIAL RENDER ====================
renderTable();
