# BPV-Matcher

BPV-Matcher is an application designed to facilitate internship matching between students, coordinators, and mentors. The system provides role-based access control and features tailored to each user role.

## Features

### Student Features
- Profile management (programming languages, preferences, portfolio)
- CV and cover letter uploads
- Application status tracking
- Portfolio management

### Coordinator Features
- Company management (partner/potential)
- Placement statistics access
- Student activity monitoring
- System-wide configuration

### Mentor Features
- Student application monitoring
- Activity tracking dashboard
- Notifications for inactive students

### Admin Features
- User management
- Role assignment
- Account activation/deactivation

## Technology Stack

- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with GORM
- **Authentication**: JWT-based authentication
- **Security**: Bcrypt password hashing, role-based access control

## Installation

### Prerequisites
- Go 1.24+
- PostgreSQL 14+
- Git

### Setup

1. Clone the repository:
```
git clone https://github.com/your-username/BPV-Matcher.git
cd BPV-Matcher
```

2. Set up the database:
```
# Create PostgreSQL database
createdb bpvmatcher
```

3. Configure environment variables:
```
# Create .env file in the project root
cp .env.example .env
# Edit .env file with your configuration
```

4. Build and run:
```
cd backend
go build -o bpvmatcher ./cmd/backend
./bpvmatcher
```

## API Documentation

For detailed API documentation, see [API Documentation](docs/api.md).

## Security

The application implements several security measures:
- JWT authentication with expiration
- Password hashing using bcrypt
- Role-based access control
- Session timeout handling
- Secure file uploads

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contributing

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-new-feature`
3. Commit your changes: `git commit -am 'Add some feature'`
4. Push to the branch: `git push origin feature/my-new-feature`
5. Submit a pull request 