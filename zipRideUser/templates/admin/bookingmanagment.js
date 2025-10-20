// Logout
document.getElementById("logout-btn").addEventListener("click", () => {
    localStorage.clear();
    window.location.href = "signin.html";
});

// Sample bookings
const bookings = [
    {
        id: 'B001', user: 'John Doe', userPhone: '1234567890',
        driver: 'Mike Ross', driverPhone: '9876543210', pickup: '123 Main St, New York', dropoff: '456 Park Ave, New York',
        vehicle: 'Car', price: 500, commission: 50, payment: 'Cash', status: 'Pending'
    },
    {
        id: 'B002', user: 'Jane Smith', userPhone: '2345678901',
        driver: 'Rachel Zane', driverPhone: '8765432109', pickup: '789 Elm St, London', dropoff: '101 Maple Rd, London',
        vehicle: 'Bike', price: 200, commission: 20, payment: 'Online', status: 'In Progress'
    },
    {
        id: 'B003', user: 'Alice Brown', userPhone: '3456789012',
        driver: 'Harvey Specter', driverPhone: '7654321098', pickup: '202 Oak St, Chicago', dropoff: '303 Pine St, Chicago',
        vehicle: 'Rickshaw', price: 300, commission: 30, payment: 'Cash', status: 'Completed'
    },
    {
        id: 'B004', user: 'Bob Johnson', userPhone: '4567890123',
        driver: 'Donna Paulsen', driverPhone: '6543210987', pickup: '99 River Rd, Miami', dropoff: '404 Beach Blvd, Miami',
        vehicle: 'Car', price: 800, commission: 80, payment: 'Online', status: 'Cancelled'
    }
];

const statusColors = {
    'Pending': 'bg-yellow-100 text-yellow-700 border border-yellow-200',
    'In Progress': 'bg-blue-100 text-blue-700 border border-blue-200',
    'Completed': 'bg-green-100 text-green-700 border border-green-200',
    'Cancelled': 'bg-red-100 text-red-700 border border-red-200'
};

const tableBody = document.getElementById('bookingTable');
const modal = document.getElementById('bookingModal');
const modalContent = document.getElementById('modalContent');
const updateStatusBtn = document.getElementById('updateStatusBtn');
const cancelBookingBtn = document.getElementById('cancelBookingBtn');

let currentBooking = null;

// Render table
const renderTable = () => {
    tableBody.innerHTML = '';
    bookings.forEach(b => {
        const tr = document.createElement('tr');
        tr.className = 'hover:bg-gray-50 transition duration-150';
        tr.innerHTML = `
            <td class="px-6 py-3 whitespace-nowrap text-sm font-medium text-gray-900">${b.id}</td>
            <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-700">${b.user}</td>
            <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-500">${b.pickup.split(',')[0]}</td>
            <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-500">${b.dropoff.split(',')[0]}</td>
            <td class="px-6 py-3 whitespace-nowrap text-sm text-gray-600">${b.vehicle}</td>
            <td class="px-6 py-3 whitespace-nowrap text-center">
                <span class="px-3 py-1 inline-flex text-xs leading-5 font-semibold rounded-full ${statusColors[b.status]}">
                    ${b.status}
                </span>
            </td>
            <td class="px-6 py-3 whitespace-nowrap text-center space-x-2">
                <button data-id="${b.id}" class="bg-cyan-600 text-white px-4 py-1.5 rounded-full text-xs hover:bg-cyan-700 transition duration-150 viewBtn shadow-md">View</button>
                ${b.status !== 'Completed' && b.status !== 'Cancelled' ? `
                <button data-id="${b.id}" class="bg-yellow-500 text-white px-4 py-1.5 rounded-full text-xs hover:bg-yellow-600 transition duration-150 statusBtn shadow-md">Status</button>
                <button data-id="${b.id}" class="bg-red-500 text-white px-4 py-1.5 rounded-full text-xs hover:bg-red-600 transition duration-150 cancelBtn shadow-md">Cancel</button>
                ` : ''}
            </td>
        `;
        tableBody.appendChild(tr);

        tr.querySelector('.viewBtn').addEventListener('click', () => {
            currentBooking = b;
            showBookingModal(b);
        });

        if (b.status !== 'Completed' && b.status !== 'Cancelled') {
            tr.querySelector('.statusBtn').addEventListener('click', () => {
                const nextStatus = b.status === 'Pending' ? 'In Progress' : b.status === 'In Progress' ? 'Completed' : b.status;
                b.status = nextStatus;
                alert(`Booking ${b.id} status updated to: ${b.status}`);
                renderTable();
            });

            tr.querySelector('.cancelBtn').addEventListener('click', () => {
                if (confirm(`Are you sure you want to cancel booking ${b.id}?`)) {
                    b.status = 'Cancelled';
                    alert(`Booking ${b.id} has been Cancelled.`);
                    renderTable();
                }
            });
        }
    });
};

// Show modal
const showBookingModal = (b) => {
    modalContent.innerHTML = `
        <div class="grid grid-cols-2 gap-x-4 gap-y-3 text-sm">
            <p><strong class="text-cyan-700">Booking ID:</strong> <span class="font-bold text-gray-900">${b.id}</span></p>
            <p><strong class="text-cyan-700">Status:</strong> <span class="px-2 py-0.5 inline-flex text-xs font-semibold rounded-full ${statusColors[b.status]}">${b.status}</span></p>
            <p class="col-span-2"><strong class="text-cyan-700">User:</strong> ${b.user} (Tel: ${b.userPhone})</p>
            <p class="col-span-2"><strong class="text-cyan-700">Driver:</strong> ${b.driver} (Tel: ${b.driverPhone})</p>
            <p class="col-span-2"><strong class="text-cyan-700">Pickup:</strong> ${b.pickup}</p>
            <p class="col-span-2"><strong class="text-cyan-700">Drop-off:</strong> ${b.dropoff}</p>
            <p><strong class="text-cyan-700">Vehicle:</strong> ${b.vehicle}</p>
            <p><strong class="text-cyan-700">Payment Method:</strong> ${b.payment}</p>
            <p><strong class="text-cyan-700">Total Price:</strong> <span class="text-green-600 font-bold">₹${b.price}</span></p>
            <p><strong class="text-cyan-700">Admin Commission:</strong> <span class="text-red-600 font-bold">₹${b.commission}</span></p>
        </div>
    `;

    const isActionable = b.status !== 'Completed' && b.status !== 'Cancelled';
    updateStatusBtn.style.display = isActionable ? 'inline-block' : 'none';
    cancelBookingBtn.style.display = isActionable ? 'inline-block' : 'none';

    modal.classList.remove('hidden');
    modal.classList.add('flex');
};

// Modal actions
updateStatusBtn.onclick = () => {
    if (!currentBooking) return;
    const b = currentBooking;
    const nextStatus = b.status === 'Pending' ? 'In Progress' : b.status === 'In Progress' ? 'Completed' : b.status;
    b.status = nextStatus;
    alert(`Booking ${b.id} status updated to: ${b.status}`);
    modal.classList.add('hidden');
    modal.classList.remove('flex');
    renderTable();
};

cancelBookingBtn.onclick = () => {
    if (!currentBooking) return;
    if (confirm(`Are you sure you want to cancel booking ${currentBooking.id}?`)) {
        currentBooking.status = 'Cancelled';
        alert(`Booking ${currentBooking.id} has been Cancelled.`);
        modal.classList.add('hidden');
        modal.classList.remove('flex');
        renderTable();
    }
};

// Close Modal
document.getElementById('closeModal').addEventListener('click', () => {
    modal.classList.add('hidden');
    modal.classList.remove('flex');
    currentBooking = null;
});

// Initial render
renderTable();
