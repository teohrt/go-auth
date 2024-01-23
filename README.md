# Golang AWS Cognito Authentication Template

This project serves as my generic GoLang project template, originally split off from a larger side project. It is designed for kick-starting Go backends with authentication through AWS Cognito. Leveraging containerization, the template ensures portability, scalability, and consistent deployment across diverse environments, enhancing the overall development and deployment experience. The inclusion of initial database repositories and specific business logic stems from established patterns in the original side project. These components are provided as a foundation, serving as implemented patterns that can be used as a jumping-off point for various use cases in the future.

## Setup

1. Clone the repository.
2. Setup environment variables: `cp dev.env .env`
3. Install dependencies using `go mod download`.
4. Build & run the server and database using `make re-compose`.

## Endpoints

| Endpoint            | Method | Description                        |
| ------------------- | ------ | ---------------------------------- |
| `/v1/health`        | GET    | Health check endpoint.             |
| `/v1/auth/register` | POST   | User registration endpoint.        |
| `/v1/auth/login`    | POST   | User login endpoint.               |
| `/v1/auth/confirm`  | POST   | User email sent code confirmation. |

## Key Values of Using Cognito

1. **Easy Integration with AWS Services:**

   - Cognito seamlessly integrates with other AWS services, simplifying the process of securing access to various AWS resources.

2. **Scalability:**

   - Designed to handle millions of users, Cognito scales automatically as your user base grows.

3. **User Management:**

   - Built-in user management capabilities with support for various authentication mechanisms.

4. **Security Features:**

   - Robust security features to protect user data, including encryption, secure token handling, and multi-factor authentication.

5. **Authentication Flows:**

   - Flexible authentication flows, supporting user pools for app clients and identity pools for federated identity providers.

6. **Device Management:**

   - Helps manage user devices and sessions, enhancing application security.

7. **Customizable UI and Workflows:**

   - Enables customization of user interfaces and authentication workflows for a seamless user experience.

8. **Compliance and Regulations:**

   - Designed to comply with industry standards and regulations, ensuring data security and legal compliance.

9. **Analytics and Insights:**
   - Provides analytics and insights into user authentication patterns for understanding user behavior and optimizing the authentication process.
