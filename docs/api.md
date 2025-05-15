# BPV-Matcher API Documentation

This documentation provides details on the API endpoints for the BPV-Matcher application, which facilitates internship matching between students, coordinators, and mentors.

## Authentication

Most API endpoints require authentication. To authenticate, include a Bearer token in the Authorization header:

```
Authorization: Bearer <your_jwt_token>
```

You can obtain a token by using the login endpoint.

## Rate Limiting

To prevent abuse, the API implements rate limiting. If you exceed the limit, you will receive a `429 Too Many Requests` response.

## Error Handling

The API returns standard HTTP status codes to indicate success or failure:

- `200 OK`: Request succeeded
- `201 Created`: Resource created successfully
- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Authentication required or failed
- `403 Forbidden`: Not authorized to access the resource
- `404 Not Found`: Resource not found
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

Error responses include a JSON object with an error message:

```json
{
  "error": "Error description"
}
```

## API Endpoints

### Health Check

#### GET /health

Check if the API is running.

**Response:**
```json
{
  "status": "OK"
}
```

### Authentication

#### POST /api/auth/register

Register a new user.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "firstName": "John",
  "lastName": "Doe",
  "role": "student" // Options: "student", "coordinator", "mentor"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "student",
    "profileImage": "",
    "bio": "",
    "active": true
  }
}
```

#### POST /api/auth/login

Login to get a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "student"
  }
}
```

### User Endpoints

These endpoints are available to all authenticated users.

#### GET /api/user/profile

Get the current user's profile.

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "student",
    "profileImage": "",
    "bio": "",
    "active": true
  },
  "studentProfile": {
    "id": 1,
    "userID": 1,
    "programmingLanguages": "JavaScript, Python",
    "preferences": "Frontend development",
    "portfolioURL": "https://github.com/johndoe",
    "cv": "",
    "coverLetter": "",
    "status": "searching"
  }
}
```

Note: The response will include role-specific profile data based on the user's role.

#### PUT /api/user/profile

Update user profile information.

**Request Body:**
```json
{
  "firstName": "John",
  "lastName": "Doe",
  "profileImage": "profile.jpg",
  "bio": "Passionate student looking for an internship"
}
```

**Response:**
```json
{
  "message": "Profile updated successfully",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstName": "John",
    "lastName": "Doe",
    "role": "student",
    "profileImage": "profile.jpg",
    "bio": "Passionate student looking for an internship",
    "active": true
  }
}
```

#### POST /api/user/change-password

Change the user's password.

**Request Body:**
```json
{
  "currentPassword": "oldPassword123",
  "newPassword": "newPassword456"
}
```

**Response:**
```json
{
  "message": "Password changed successfully"
}
```

### Student Endpoints

These endpoints are only accessible to users with the "student" role.

#### GET /api/student/profile

Get the student's detailed profile.

**Response:**
```json
{
  "student": {
    "id": 1,
    "userID": 1,
    "user": {
      "id": 1,
      "email": "student@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "student",
      "profileImage": "",
      "bio": "",
      "active": true
    },
    "programmingLanguages": "JavaScript, Python",
    "preferences": "Frontend development",
    "portfolioURL": "https://github.com/johndoe",
    "cv": "uploads/cv/user_1_resume.pdf",
    "coverLetter": "uploads/cover_letters/user_1_cover.pdf",
    "status": "searching"
  }
}
```

#### PUT /api/student/profile

Update the student's detailed profile.

**Request Body:**
```json
{
  "programmingLanguages": "JavaScript, Python, Go",
  "preferences": "Backend development",
  "portfolioURL": "https://github.com/johndoe/portfolio"
}
```

**Response:**
```json
{
  "message": "Student profile updated successfully",
  "student": {
    "id": 1,
    "userID": 1,
    "programmingLanguages": "JavaScript, Python, Go",
    "preferences": "Backend development",
    "portfolioURL": "https://github.com/johndoe/portfolio",
    "cv": "uploads/cv/user_1_resume.pdf",
    "coverLetter": "uploads/cover_letters/user_1_cover.pdf",
    "status": "searching"
  }
}
```

#### POST /api/student/cv

Upload a CV document.

**Request:**
Form data with a file field named "cv".

**Response:**
```json
{
  "message": "CV uploaded successfully",
  "filePath": "uploads/cv/user_1_resume.pdf"
}
```

#### POST /api/student/cover-letter

Upload a cover letter document.

**Request:**
Form data with a file field named "coverLetter".

**Response:**
```json
{
  "message": "Cover letter uploaded successfully",
  "filePath": "uploads/cover_letters/user_1_cover.pdf"
}
```

#### PUT /api/student/status

Update the student's application status.

**Request Body:**
```json
{
  "status": "searching" // Options: "inactive", "searching", "applied", "placed"
}
```

**Response:**
```json
{
  "message": "Status updated successfully"
}
```

### Coordinator Endpoints

These endpoints are only accessible to users with the "coordinator" role.

#### GET /api/coordinator/profile

Get the coordinator's profile.

**Response:**
```json
{
  "coordinator": {
    "id": 1,
    "userID": 2,
    "user": {
      "id": 2,
      "email": "coordinator@example.com",
      "firstName": "Jane",
      "lastName": "Smith",
      "role": "coordinator",
      "profileImage": "",
      "bio": "",
      "active": true
    },
    "department": "Computer Science",
    "canConfigSystem": true,
    "canManageCompanies": true,
    "canViewStatistics": true,
    "canMonitorStudents": true
  }
}
```

#### GET /api/coordinator/companies

Get all companies.

**Response:**
```json
{
  "companies": [
    {
      "id": 1,
      "name": "Tech Solutions Inc",
      "description": "Software development company",
      "website": "https://techsolutions.com",
      "location": "New York",
      "industry": "Technology",
      "isPartner": true,
      "isPotential": false,
      "contact": "John Manager",
      "email": "contact@techsolutions.com",
      "phone": "+1234567890",
      "notes": "Looking for interns in summer"
    }
  ]
}
```

#### GET /api/coordinator/companies/:id

Get a specific company by ID.

**Response:**
```json
{
  "company": {
    "id": 1,
    "name": "Tech Solutions Inc",
    "description": "Software development company",
    "website": "https://techsolutions.com",
    "location": "New York",
    "industry": "Technology",
    "isPartner": true,
    "isPotential": false,
    "contact": "John Manager",
    "email": "contact@techsolutions.com",
    "phone": "+1234567890",
    "notes": "Looking for interns in summer"
  }
}
```

#### POST /api/coordinator/companies

Create a new company.

**Request Body:**
```json
{
  "name": "Tech Solutions Inc",
  "description": "Software development company",
  "website": "https://techsolutions.com",
  "location": "New York",
  "industry": "Technology",
  "isPartner": true,
  "isPotential": false,
  "contact": "John Manager",
  "email": "contact@techsolutions.com",
  "phone": "+1234567890",
  "notes": "Looking for interns in summer"
}
```

**Response:**
```json
{
  "message": "Company created successfully",
  "company": {
    "id": 1,
    "name": "Tech Solutions Inc",
    "description": "Software development company",
    "website": "https://techsolutions.com",
    "location": "New York",
    "industry": "Technology",
    "isPartner": true,
    "isPotential": false,
    "contact": "John Manager",
    "email": "contact@techsolutions.com",
    "phone": "+1234567890",
    "notes": "Looking for interns in summer"
  }
}
```

#### PUT /api/coordinator/companies/:id

Update a company.

**Request Body:**
Same format as POST request.

**Response:**
```json
{
  "message": "Company updated successfully",
  "company": {
    "id": 1,
    "name": "Tech Solutions Inc",
    "description": "Updated description",
    "website": "https://techsolutions.com",
    "location": "New York",
    "industry": "Technology",
    "isPartner": true,
    "isPotential": false,
    "contact": "John Manager",
    "email": "contact@techsolutions.com",
    "phone": "+1234567890",
    "notes": "Looking for interns in summer"
  }
}
```

#### PUT /api/coordinator/companies/:id/partner

Update a company's partner status.

**Request Body:**
```json
{
  "isPartner": true
}
```

**Response:**
```json
{
  "message": "Company updated to partner successfully"
}
```

#### GET /api/coordinator/statistics

Get student placement statistics.

**Response:**
```json
{
  "statistics": {
    "total": 50,
    "inactive": 5,
    "searching": 20,
    "applied": 15,
    "placed": 10
  }
}
```

#### GET /api/coordinator/students

Get all students with their profiles.

**Response:**
```json
{
  "students": [
    {
      "id": 1,
      "userID": 3,
      "user": {
        "id": 3,
        "email": "student@example.com",
        "firstName": "John",
        "lastName": "Doe",
        "role": "student",
        "profileImage": "",
        "bio": "",
        "active": true
      },
      "programmingLanguages": "JavaScript, Python",
      "preferences": "Frontend development",
      "portfolioURL": "https://github.com/johndoe",
      "cv": "uploads/cv/user_3_resume.pdf",
      "coverLetter": "uploads/cover_letters/user_3_cover.pdf",
      "status": "searching"
    }
  ]
}
```

### Mentor Endpoints

These endpoints are only accessible to users with the "mentor" role.

#### GET /api/mentor/profile

Get the mentor's profile.

**Response:**
```json
{
  "mentor": {
    "id": 1,
    "userID": 4,
    "user": {
      "id": 4,
      "email": "mentor@example.com",
      "firstName": "Robert",
      "lastName": "Johnson",
      "role": "mentor",
      "profileImage": "",
      "bio": "",
      "active": true
    },
    "department": "Software Engineering",
    "canMonitorApplications": true,
    "canTrackActivity": true,
    "canNotifyInactiveStudents": true,
    "assignedStudentCount": 10
  }
}
```

#### GET /api/mentor/students/applications

Get students who are searching or have applied for internships.

**Response:**
```json
{
  "students": [
    {
      "id": 1,
      "userID": 3,
      "user": {
        "id": 3,
        "email": "student@example.com",
        "firstName": "John",
        "lastName": "Doe",
        "role": "student",
        "profileImage": "",
        "bio": "",
        "active": true
      },
      "programmingLanguages": "JavaScript, Python",
      "preferences": "Frontend development",
      "portfolioURL": "https://github.com/johndoe",
      "cv": "uploads/cv/user_3_resume.pdf",
      "coverLetter": "uploads/cover_letters/user_3_cover.pdf",
      "status": "applied"
    }
  ]
}
```

#### GET /api/mentor/students/inactive

Get inactive students or those who haven't logged in for a month.

**Response:**
```json
{
  "students": [
    {
      "id": 2,
      "userID": 5,
      "user": {
        "id": 5,
        "email": "inactive@example.com",
        "firstName": "Mary",
        "lastName": "Wilson",
        "role": "student",
        "profileImage": "",
        "bio": "",
        "active": true
      },
      "programmingLanguages": "Java, C++",
      "preferences": "Mobile development",
      "portfolioURL": "https://github.com/marywilson",
      "cv": "",
      "coverLetter": "",
      "status": "inactive"
    }
  ]
}
```

#### GET /api/mentor/students/:id

Get details about a specific student.

**Response:**
```json
{
  "student": {
    "id": 1,
    "userID": 3,
    "user": {
      "id": 3,
      "email": "student@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "student",
      "profileImage": "",
      "bio": "",
      "active": true
    },
    "programmingLanguages": "JavaScript, Python",
    "preferences": "Frontend development",
    "portfolioURL": "https://github.com/johndoe",
    "cv": "uploads/cv/user_3_resume.pdf",
    "coverLetter": "uploads/cover_letters/user_3_cover.pdf",
    "status": "applied"
  }
}
```

#### POST /api/mentor/students/:id/notify

Send a notification to a student.

**Request Body:**
```json
{
  "message": "Please update your profile and start applying for internships."
}
```

**Response:**
```json
{
  "message": "Notification sent successfully",
  "details": {
    "studentID": 1,
    "notification": "Please update your profile and start applying for internships."
  }
}
```

#### GET /api/mentor/dashboard

Get student activity statistics for the dashboard.

**Response:**
```json
{
  "statistics": {
    "total": 50,
    "inactive": 5,
    "searching": 20,
    "applied": 15,
    "placed": 10,
    "recentlyActive": 30
  }
}
```

### Admin Endpoints

These endpoints are only accessible to users with the "admin" role.

#### GET /api/admin/users

Get all users.

**Response:**
```json
{
  "users": [
    {
      "id": 1,
      "email": "admin@example.com",
      "firstName": "Admin",
      "lastName": "User",
      "role": "admin",
      "profileImage": "",
      "bio": "",
      "active": true
    },
    {
      "id": 2,
      "email": "student@example.com",
      "firstName": "John",
      "lastName": "Doe",
      "role": "student",
      "profileImage": "",
      "bio": "",
      "active": true
    }
  ]
}
```

#### PUT /api/admin/users/:id/role

Update a user's role.

**Request Body:**
```json
{
  "role": "mentor" // Options: "student", "coordinator", "mentor"
}
```

**Response:**
```json
{
  "message": "User role updated successfully"
}
```

#### PUT /api/admin/users/:id/deactivate

Deactivate a user account.

**Response:**
```json
{
  "message": "User deactivated successfully"
}
```

#### PUT /api/admin/users/:id/reactivate

Reactivate a user account.

**Response:**
```json
{
  "message": "User reactivated successfully"
}
```

#### DELETE /api/admin/users/:id

Delete a user account.

**Response:**
```json
{
  "message": "User deleted successfully"
}
```

#### POST /api/admin/users/:id/coordinator

Add coordinator privileges to a user.

**Response:**
```json
{
  "message": "Coordinator privileges added successfully"
}
```

#### POST /api/admin/users/:id/mentor

Add mentor privileges to a user.

**Response:**
```json
{
  "message": "Mentor privileges added successfully"
}
```

#### DELETE /api/admin/users/:id/privileges

Remove special privileges (coordinator or mentor) from a user.

**Response:**
```json
{
  "message": "Special privileges removed successfully"
}
```

## Security Considerations

1. All sensitive data should be transferred over HTTPS.
2. JWT tokens expire after a configured duration (default: 8 hours).
3. Passwords are hashed using bcrypt.
4. Session timeout is enforced for all authenticated endpoints.
5. Role-based access control is enforced for all endpoints.
6. File uploads are restricted to certain types and sizes.
7. Cross-Origin Resource Sharing (CORS) is properly configured. 