API Documentation (API_Documentation.md)
🟢 Public Routes (/api/driver)
POST /api/driver/register
Description: Start driver registration with core profile details.
Auth Required: No
Request Headers:

Content-Type: application/json
Request Body:

json
{
  "name": "Ravi Kumar",
  "email": "ravi@example.com",
  "phone": "+919876543210",
  "password": "SecurePass@123"
}
Response Example:

json
{
  "message": "Registration successful",
  "status": "pending_verification",
  "driver_id": "drv_123"
}
GET /api/driver/registration-status/:email
Description: Check onboarding status for a submitted email.
Auth Required: No

Response Example:

json
{
  "email": "ravi@example.com",
  "status": "pending_verification"
}
POST /api/driver/auth/send-otp
Description: Send OTP to driver’s phone/email for login.
Auth Required: No

Request Body:

json
{
  "phone": "+919876543210",
  "email": "ravi@example.com"
}
Response Example:

json
{
  "message": "OTP sent successfully"
}
POST /api/driver/auth/verify-otp
Description: Verify OTP received by driver.
Auth Required: No

Request Body:

json
{
  "phone": "+919876543210",
  "otp": "123456"
}
Response Example:

json
{
  "message": "OTP verified",
  "token": "jwt-token-value"
}
POST /api/driver/auth/login
Description: Login using OTP or credentials to obtain tokens.
Auth Required: No

Request Body:

json
{
  "email": "ravi@example.com",
  "password": "SecurePass@123"
}
Response Example:

json
{
  "access_token": "jwt-access",
  "refresh_token": "jwt-refresh",
  "driver_id": "drv_123"
}
POST /api/driver/auth/logout
Description: Invalidate existing JWT session.
Auth Required: No (token in body or header, depending on implementation)

Request Body:

json
{
  "refresh_token": "jwt-refresh"
}
Response Example:

json
{
  "message": "Logged out successfully"
}
🔒 Authenticated Driver Routes (/api/driver/*)
All endpoints require header Authorization: Bearer {{token}}.

POST /api/driver/onboarding
Description: Complete onboarding steps such as documents and personal info.
Auth Required: Yes

Request Body:

json
{
  "driver_id": "drv_123",
  "documents_submitted": true,
  "address": "123, MG Road, Bengaluru"
}
Response Example:

json
{
  "message": "Onboarding completed",
  "status": "approved"
}
GET /api/driver/:driverId/profile
Description: Fetch driver profile summary.
Auth Required: Yes

Response Example:

json
{
  "driver_id": "drv_123",
  "name": "Ravi Kumar",
  "status": "active",
  "rating": 4.9
}
PATCH /api/driver/:driverId/status
Description: Update availability (online/offline) status.
Auth Required: Yes

Request Body:

json
{
  "status": "online"
}
Response Example:

json
{
  "driver_id": "drv_123",
  "status": "online",
  "updated_at": "2025-10-21T14:00:00Z"
}
GET /api/driver/:driverId/earnings/summary
Description: Get aggregated earnings stats.
Auth Required: Yes

Response Example:

json
{
  "total_earnings": 45230.75,
  "completed_rides": 320,
  "pending_withdrawals": 2
}
GET /api/driver/:driverId/earnings/trend
Description: Earnings trend data for charts (daily/weekly).
Auth Required: Yes

Response Example:

json
{
  "series": [
    { "date": "2025-10-18", "amount": 1230.5 },
    { "date": "2025-10-19", "amount": 980.0 }
  ]
}
POST /api/driver/:driverId/withdraw
Description: Submit payout/withdrawal request.
Auth Required: Yes

Request Body:

json
{
  "amount": 1500,
  "method": "bank_transfer"
}
Response Example:

json
{
  "message": "Withdrawal request submitted",
  "request_id": "wd_987"
}
GET /api/driver/:driverId/rides/summary
Description: Fetch ride statistics for the driver.
Auth Required: Yes

Response Example:

json
{
  "completed": 320,
  "cancelled": 12,
  "ongoing": 1
}
POST /api/driver/:driverId/location
Description: Update live GPS coordinates.
Auth Required: Yes

Request Body:

json
{
  "latitude": 12.9716,
  "longitude": 77.5946
}
Response Example:

json
{
  "message": "Location updated",
  "timestamp": "2025-10-21T14:10:00Z"
}
GET /api/driver/:driverId/requests
Description: List incoming ride requests.
Auth Required: Yes

Response Example:

json
[
  {
    "ride_id": "ride_456",
    "pickup": "Koramangala",
    "dropoff": "Indiranagar",
    "fare": 230
  }
]
POST /api/driver/:driverId/rides/:id/accept
Description: Accept a ride request.
Auth Required: Yes

Response Example:

json
{
  "ride_id": "ride_456",
  "status": "accepted"
}
POST /api/driver/:driverId/rides/:id/cancel
Description: Cancel a previously accepted ride.
Auth Required: Yes

Request Body:

json
{
  "reason": "Vehicle breakdown"
}
Response Example:

json
{
  "ride_id": "ride_456",
  "status": "cancelled"
}
GET /api/driver/:driverId/vehicles
Description: List vehicles associated with the driver.
Auth Required: Yes

Response Example:

json
[
  {
    "vehicle_id": "veh_101",
    "model": "Honda Amaze",
    "registration": "KA03AB1234",
    "status": "verified"
  }
]
POST /api/driver/:driverId/vehicles
Description: Add a new vehicle.
Auth Required: Yes

Request Body:

json
{
  "model": "Toyota Etios",
  "registration": "KA05CD6789",
  "year": 2022
}
Response Example:

json
{
  "vehicle_id": "veh_102",
  "status": "pending_verification"
}
PUT /api/driver/:driverId/vehicles/:id
Description: Update vehicle information.
Auth Required: Yes

Request Body:

json
{
  "model": "Toyota Etios",
  "registration": "KA05CD6789"
}
Response Example:

json
{
  "vehicle_id": "veh_102",
  "message": "Vehicle updated"
}
DELETE /api/driver/:driverId/vehicles/:id
Description: Remove a vehicle from the driver profile.
Auth Required: Yes

Response Example:

json
{
  "vehicle_id": "veh_102",
  "message": "Vehicle removed"
}
GET /api/driver/:driverId/documents
Description: List documents uploaded by the driver.
Auth Required: Yes

Response Example:

json
[
  {
    "document_id": "doc_201",
    "type": "License",
    "status": "approved"
  }
]
POST /api/driver/:driverId/documents
Description: Upload a new document (metadata).
Auth Required: Yes

Request Body:

json
{
  "type": "Vehicle Insurance",
  "file_url": "https://files.example.com/ins_123.pdf"
}
Response Example:

json
{
  "document_id": "doc_202",
  "status": "pending_review"
}
PATCH /api/driver/:driverId/documents/:id/status
Description: Update document verification status (internal use).
Auth Required: Yes

Request Body:

json
{
  "status": "approved"
}
Response Example:

json
{
  "document_id": "doc_202",
  "status": "approved"
}
🆘 Help Center Routes (/api/help)
GET /api/help/faqs
Description: Fetch public FAQ list.
Auth Required: No

Response Example:

json
[
  {
    "question": "How do I reset my password?",
    "answer": "Use the forgot password option in the app."
  }
]
Authenticated help center endpoints require header Authorization: Bearer {{token}}.

POST /api/help/ticket
Description: Submit a support ticket.
Auth Required: Yes

Request Body:

json
{
  "subject": "Payment issue",
  "description": "Payment not received for ride ride_456"
}
Response Example:

json
{
  "ticket_id": "ticket_789",
  "status": "open"
}
POST /api/help/report
Description: Report rider misconduct or incident.
Auth Required: Yes

Request Body:

json
{
  "ride_id": "ride_456",
  "description": "Rider was abusive"
}
Response Example:

json
{
  "report_id": "report_321",
  "status": "submitted"
}
POST /api/help/chat/start
Description: Start a live chat session with support.
Auth Required: Yes

Response Example:

json
{
  "chat_id": "chat_555",
  "status": "active"
}
🌐 WebSocket Routes
GET ws/driver/:driverId
Description: Establish WebSocket connection for live ride updates.
Auth Required: Yes (Authorization: Bearer {{token}} in headers or query, depending on client setup)

Upgrade Request Example:

GET {{base_url}}/ws/driver/drv_123
Sec-WebSocket-Version: 13
Sec-WebSocket-Key: <random>
Authorization: Bearer {{token}}
🧑‍💼 Admin Panel Routes (SSR)
These routes are session-protected via cookies (set by login). They render HTML for the admin console. Include Cookie: admin_session={{value}} in Postman when needed.

Public Admin Panel Endpoints
GET /admin/panel/
Description: Redirect to dashboard when logged in.
Auth Required: Session cookie

Response Example: Redirect 302 → /admin/panel/dashboard

GET /admin/panel/login
Description: Render admin login page.
Auth Required: No

Response Example: HTML content

POST /admin/panel/login
Description: Authenticate admin user and set session cookie.
Auth Required: No

Request Body (form or JSON):

json
{
  "email": "admin@example.com",
  "password": "AdminPass@123"
}
Response Example: Redirect 302 → /admin/panel/dashboard

GET /admin/panel/logout
Description: Clear admin session and redirect to login.
Auth Required: Session cookie

Response Example: Redirect 302 → /admin/panel/login

Authenticated Admin Pages
All require valid admin session cookie.

GET /admin/panel/dashboard
Description: View platform analytics dashboard.
Auth Required: Yes

Response Example: HTML dashboard view.

GET /admin/panel/drivers
Description: List all drivers with filters.
Auth Required: Yes

Response Example: HTML table of drivers.

GET /admin/panel/drivers/pending
Description: View drivers awaiting approval.
Auth Required: Yes

Response Example: HTML list of pending drivers.

GET /admin/panel/driver/:id
Description: Detailed driver profile page.
Auth Required: Yes

POST /admin/panel/driver/:id/approve
Description: Approve driver account.
Auth Required: Yes (ACL admin:drivers:approve)

Response Example: Redirect back to driver page.

POST /admin/panel/driver/:id/reject
Description: Reject driver application.
Auth Required: Yes

POST /admin/panel/driver/:id/suspend
Description: Suspend driver from platform.
Auth Required: Yes

GET /admin/panel/roles
Description: View roles and permissions.
Auth Required: Yes

POST /admin/panel/roles
Description: Create a role.
Auth Required: Yes (ACL admin:roles:edit)

Request Body (form/JSON):

json
{
  "name": "Support",
  "permissions": ["admin:help:view", "admin:help:reply"]
}
GET /admin/panel/admins
Description: List admin users.
Auth Required: Yes

POST /admin/panel/admins
Description: Create admin user.
Auth Required: Yes

Request Body:

json
{
  "name": "Support Admin",
  "email": "support@example.com",
  "role": "Support"
}
GET /admin/panel/vehicles
Description: Manage vehicles (verification flow).
Auth Required: Yes

POST /admin/panel/vehicles/:id/verify
Description: Verify vehicle documents.
Auth Required: Yes (ACL admin:vehicles:verify)

GET /admin/panel/rides
Description: Monitor ride activity.
Auth Required: Yes

POST /admin/panel/rides/:id/cancel
Description: Cancel problematic ride.
Auth Required: Yes

GET /admin/panel/earnings
Description: View driver earnings overview.
Auth Required: Yes

GET /admin/panel/withdrawals
Description: Review withdrawal requests.
Auth Required: Yes

POST /admin/panel/withdrawals/:id/approve
Description: Approve withdrawal.
Auth Required: Yes (ACL admin:withdrawals:approve)

POST /admin/panel/withdrawals/:id/reject
Description: Reject withdrawal.
Auth Required: Yes

GET /admin/panel/help
Description: Manage support tickets/issues.
Auth Required: Yes

POST /admin/panel/help/:id/reply
Description: Respond to a help ticket.
Auth Required: Yes (ACL admin:help:reply)

Request Body:

json
{
  "message": "We have credited your account."
}
POST /admin/panel/help/:id/close
Description: Close a help ticket.
Auth Required: Yes (ACL admin:help:close)

GET /admin/panel/settings
Description: Admin settings page.
Auth Required: Yes

POST /admin/panel/settings/password
Description: Change admin password.
Auth Required: Yes

Request Body:

json
{
  "current_password": "OldPass@123",
  "new_password": "NewSecure@456",
  "confirm_password": "NewSecure@456"
}
🟠 Health Check
GET /health
Description: Service availability check.
Auth Required: No

Response Example:

json
{
  "status": "zipride-driver-service running"
}