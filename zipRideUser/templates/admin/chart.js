document.addEventListener("DOMContentLoaded", async () => {
  // 1️⃣ Load Navbar
  try {
    const navbarRes = await fetch("navbar.html");
    document.getElementById("navbar").innerHTML = await navbarRes.text();

    const footerRes = await fetch("footer.html");
    document.getElementById("footer").innerHTML = await footerRes.text();
  } catch (err) {
    console.error("Failed to load reusable components:", err);
  }

  // 2️⃣ Initialize Charts
  initializeCharts();
});

function initializeCharts() {
  // Monthly Bookings Chart
  const monthlyCtx = document.getElementById("monthlyBookingsChart")?.getContext("2d");
  if (monthlyCtx) {
    new Chart(monthlyCtx, {
      type: "bar",
      data: {
        labels: ["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct"],
        datasets:[{
          label:"Monthly Bookings",
          data:[120,180,150,200,250,300,270,320,310,400],
          backgroundColor:"rgba(54,162,235,0.7)",
          borderColor:"rgba(54,162,235,1)",
          borderWidth:1
        }]
      },
      options:{
        responsive:true,
        plugins:{legend:{display:true}, title:{display:true,text:"Monthly Bookings"}},
        scales:{y:{beginAtZero:true}}
      }
    });
  }

  // Weekly Bookings Chart
  const weeklyCtx = document.getElementById("weeklyBookingsChart")?.getContext("2d");
  if (weeklyCtx) {
    new Chart(weeklyCtx, {
      type: "line",
      data: {
        labels:["Mon","Tue","Wed","Thu","Fri","Sat","Sun"],
        datasets:[{
          label:"Weekly Bookings",
          data:[20,30,25,40,35,50,45],
          borderColor:"rgba(255,159,64,1)",
          backgroundColor:"rgba(255,159,64,0.3)",
          fill:true,
          tension:0.4
        }]
      },
      options:{
        responsive:true,
        plugins:{legend:{display:true}, title:{display:true,text:"Weekly Bookings"}},
        scales:{y:{beginAtZero:true}}
      }
    });
  }

  // Daily Bookings Chart
  const dailyCtx = document.getElementById("dailyBookingsChart")?.getContext("2d");
  if (dailyCtx) {
    new Chart(dailyCtx, {
      type: "doughnut",
      data:{
        labels:["Completed","Pending","Cancelled"],
        datasets:[{
          data:[65,25,10],
          backgroundColor:["rgba(75,192,192,0.8)","rgba(255,206,86,0.8)","rgba(255,99,132,0.8)"],
          borderWidth:1
        }]
      },
      options:{
        responsive:true,
        plugins:{title:{display:true,text:"Daily Booking Status"}, legend:{position:"bottom"}}
      }
    });
  }

  // Revenue Chart
  const revenueCtx = document.getElementById("revenueChart")?.getContext("2d");
  if (revenueCtx) {
    new Chart(revenueCtx, {
      type:"bar",
      data:{
        labels:["Jan","Feb","Mar","Apr","May","Jun","Jul","Aug","Sep","Oct"],
        datasets:[{
          label:"Revenue (₹ in thousands)",
          data:[25,30,28,35,40,45,50,55,60,75],
          backgroundColor:"rgba(153,102,255,0.7)",
          borderColor:"rgba(153,102,255,1)",
          borderWidth:1
        }]
      },
      options:{
        responsive:true,
        plugins:{legend:{display:true},title:{display:true,text:"Monthly Revenue Growth"}},
        scales:{y:{beginAtZero:true}}
      }
    });
  }
}
