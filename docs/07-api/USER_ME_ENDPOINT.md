# `/me` Endpoint - Get Current User with Permissions

## 📋 Overview

The `/me` endpoint returns the currently authenticated user's profile along with their roles and permissions. This is essential for frontend applications to:

- Display user information
- Show/hide UI elements based on permissions
- Enable/disable features based on roles
- Implement role-based navigation

## 🔗 Endpoint Details

```
GET /api/v1/users/me
```

### Authentication

**Required**: Yes (Bearer Token)

### Headers

```
Authorization: Bearer <your_jwt_token>
```

## 📤 Response

### Success Response (200 OK)

```json
{
  "success": true,
  "message": "User profile retrieved successfully",
  "data": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "avatar_url": "https://example.com/avatar.jpg",
    "bio": "Software Developer",
    "status": "active",
    "roles": [
      {
        "id": 1,
        "name": "admin"
      },
      {
        "id": 2,
        "name": "editor"
      }
    ],
    "permissions": [
      {
        "resource": "users",
        "action": "read"
      },
      {
        "resource": "users",
        "action": "write"
      },
      {
        "resource": "novels",
        "action": "read"
      },
      {
        "resource": "novels",
        "action": "write"
      },
      {
        "resource": "chapters",
        "action": "read"
      },
      {
        "resource": "chapters",
        "action": "write"
      }
    ]
  }
}
```

### Error Responses

#### 401 Unauthorized - Missing Token

```json
{
  "success": false,
  "error": {
    "code": "AUTH_TOKEN_MISSING",
    "message": "Authentication token is required"
  }
}
```

#### 401 Unauthorized - Invalid Token

```json
{
  "success": false,
  "error": {
    "code": "AUTH_INVALID_TOKEN",
    "message": "Invalid or expired authentication token"
  }
}
```

#### 404 Not Found - User Not Found

```json
{
  "success": false,
  "error": {
    "code": "USER_NOT_FOUND",
    "message": "User not found"
  }
}
```

## 💻 Frontend Usage Examples

### JavaScript/Fetch

```javascript
async function getMe() {
  const token = localStorage.getItem("token");

  const response = await fetch("http://localhost:8080/api/v1/users/me", {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
      "Content-Type": "application/json",
    },
  });

  const data = await response.json();

  if (data.success) {
    const user = data.data;
    console.log("User:", user);
    console.log("Roles:", user.roles);
    console.log("Permissions:", user.permissions);

    // Store in state management (Redux, Vuex, etc.)
    return user;
  } else {
    console.error("Error:", data.error);
    throw new Error(data.error.message);
  }
}
```

### React Hook

```javascript
import { useState, useEffect } from "react";

function useCurrentUser() {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function fetchUser() {
      try {
        const token = localStorage.getItem("token");
        const response = await fetch("http://localhost:8080/api/v1/users/me", {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        });

        const data = await response.json();

        if (data.success) {
          setUser(data.data);
        } else {
          setError(data.error.message);
        }
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }

    fetchUser();
  }, []);

  return { user, loading, error };
}

// Usage in component
function UserProfile() {
  const { user, loading, error } = useCurrentUser();

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Welcome, {user.username}!</h1>
      <p>Email: {user.email}</p>
      <h3>Your Roles:</h3>
      <ul>
        {user.roles.map((role) => (
          <li key={role.id}>{role.name}</li>
        ))}
      </ul>
    </div>
  );
}
```

### Vue.js 3 Composition API

```javascript
import { ref, onMounted } from "vue";

export function useCurrentUser() {
  const user = ref(null);
  const loading = ref(true);
  const error = ref(null);

  const fetchUser = async () => {
    try {
      const token = localStorage.getItem("token");
      const response = await fetch("http://localhost:8080/api/v1/users/me", {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      const data = await response.json();

      if (data.success) {
        user.value = data.data;
      } else {
        error.value = data.error.message;
      }
    } catch (err) {
      error.value = err.message;
    } finally {
      loading.value = false;
    }
  };

  onMounted(fetchUser);

  return { user, loading, error, refetch: fetchUser };
}
```

### Axios

```javascript
import axios from "axios";

const api = axios.create({
  baseURL: "http://localhost:8080/api/v1",
});

// Add token to all requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

async function getMe() {
  try {
    const response = await api.get("/users/me");
    return response.data.data;
  } catch (error) {
    console.error("Error fetching user:", error.response?.data);
    throw error;
  }
}
```

## 🎨 Permission-Based UI Rendering

### Check if User Has Permission

```javascript
function hasPermission(user, resource, action) {
  return user.permissions.some(
    (perm) => perm.resource === resource && perm.action === action
  );
}

// Usage
function NovelActions({ user }) {
  const canCreateNovel = hasPermission(user, "novels", "write");
  const canDeleteNovel = hasPermission(user, "novels", "delete");

  return (
    <div>
      {canCreateNovel && <button onClick={createNovel}>Create Novel</button>}
      {canDeleteNovel && <button onClick={deleteNovel}>Delete Novel</button>}
    </div>
  );
}
```

### Check if User Has Role

```javascript
function hasRole(user, roleName) {
  return user.roles.some((role) => role.name === roleName);
}

// Usage
function AdminPanel({ user }) {
  if (!hasRole(user, "admin")) {
    return <div>Access Denied</div>;
  }

  return <div>Admin Panel Content</div>;
}
```

### Permission Helper Class

```javascript
class PermissionChecker {
  constructor(user) {
    this.user = user;
    this.permissionMap = new Map();

    // Build permission map for O(1) lookups
    user.permissions.forEach((perm) => {
      const key = `${perm.resource}:${perm.action}`;
      this.permissionMap.set(key, true);
    });

    this.roleSet = new Set(user.roles.map((r) => r.name));
  }

  can(resource, action) {
    const key = `${resource}:${action}`;
    return this.permissionMap.has(key);
  }

  hasRole(roleName) {
    return this.roleSet.has(roleName);
  }

  hasAnyRole(...roleNames) {
    return roleNames.some((role) => this.roleSet.has(role));
  }

  hasAllRoles(...roleNames) {
    return roleNames.every((role) => this.roleSet.has(role));
  }

  isAdmin() {
    return this.hasRole("admin");
  }
}

// Usage
const checker = new PermissionChecker(user);

if (checker.can("novels", "write")) {
  // Show create novel button
}

if (checker.isAdmin()) {
  // Show admin features
}
```

## 🔄 State Management Examples

### Redux

```javascript
// actions.js
export const FETCH_ME_REQUEST = "FETCH_ME_REQUEST";
export const FETCH_ME_SUCCESS = "FETCH_ME_SUCCESS";
export const FETCH_ME_FAILURE = "FETCH_ME_FAILURE";

export const fetchMe = () => async (dispatch) => {
  dispatch({ type: FETCH_ME_REQUEST });

  try {
    const token = localStorage.getItem("token");
    const response = await fetch("http://localhost:8080/api/v1/users/me", {
      headers: { Authorization: `Bearer ${token}` },
    });

    const data = await response.json();

    if (data.success) {
      dispatch({ type: FETCH_ME_SUCCESS, payload: data.data });
    } else {
      dispatch({ type: FETCH_ME_FAILURE, error: data.error });
    }
  } catch (error) {
    dispatch({ type: FETCH_ME_FAILURE, error: error.message });
  }
};

// reducer.js
const initialState = {
  user: null,
  loading: false,
  error: null,
};

export default function userReducer(state = initialState, action) {
  switch (action.type) {
    case FETCH_ME_REQUEST:
      return { ...state, loading: true, error: null };
    case FETCH_ME_SUCCESS:
      return { ...state, loading: false, user: action.payload };
    case FETCH_ME_FAILURE:
      return { ...state, loading: false, error: action.error };
    default:
      return state;
  }
}

// selectors.js
export const selectUser = (state) => state.user.user;
export const selectUserPermissions = (state) =>
  state.user.user?.permissions || [];
export const selectUserRoles = (state) => state.user.user?.roles || [];

export const selectHasPermission = (resource, action) => (state) => {
  const permissions = selectUserPermissions(state);
  return permissions.some(
    (perm) => perm.resource === resource && perm.action === action
  );
};
```

### Zustand

```javascript
import create from "zustand";

const useUserStore = create((set) => ({
  user: null,
  loading: false,
  error: null,

  fetchMe: async () => {
    set({ loading: true, error: null });

    try {
      const token = localStorage.getItem("token");
      const response = await fetch("http://localhost:8080/api/v1/users/me", {
        headers: { Authorization: `Bearer ${token}` },
      });

      const data = await response.json();

      if (data.success) {
        set({ user: data.data, loading: false });
      } else {
        set({ error: data.error, loading: false });
      }
    } catch (error) {
      set({ error: error.message, loading: false });
    }
  },

  hasPermission: (resource, action) => {
    const user = useUserStore.getState().user;
    return (
      user?.permissions.some(
        (perm) => perm.resource === resource && perm.action === action
      ) || false
    );
  },

  hasRole: (roleName) => {
    const user = useUserStore.getState().user;
    return user?.roles.some((role) => role.name === roleName) || false;
  },
}));

// Usage
function MyComponent() {
  const { user, loading, fetchMe, hasPermission } = useUserStore();

  useEffect(() => {
    fetchMe();
  }, [fetchMe]);

  if (loading) return <div>Loading...</div>;

  return (
    <div>
      {hasPermission("novels", "write") && <button>Create Novel</button>}
    </div>
  );
}
```

## 🧪 Testing with cURL

```bash
# Get your token first (from login)
TOKEN="your_jwt_token_here"

# Call /me endpoint
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq

# With fish shell
set TOKEN "your_jwt_token_here"
curl -X GET http://localhost:8080/api/v1/users/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" | jq
```

## 🔐 Security Notes

1. **Always use HTTPS** in production to protect the JWT token
2. **Store tokens securely** (httpOnly cookies are better than localStorage)
3. **Validate token expiry** on the frontend and refresh when needed
4. **Don't expose sensitive data** - the endpoint only returns safe user info
5. **Cache wisely** - permissions can change, don't cache too long

## 📱 Mobile App Usage

### React Native

```javascript
import AsyncStorage from "@react-native-async-storage/async-storage";

async function getMe() {
  try {
    const token = await AsyncStorage.getItem("token");

    const response = await fetch("http://localhost:8080/api/v1/users/me", {
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    });

    const data = await response.json();
    return data.data;
  } catch (error) {
    console.error("Error:", error);
    throw error;
  }
}
```

### Flutter

```dart
import 'package:http/http.dart' as http;
import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';

Future<Map<String, dynamic>> getMe() async {
  final prefs = await SharedPreferences.getInstance();
  final token = prefs.getString('token');

  final response = await http.get(
    Uri.parse('http://localhost:8080/api/v1/users/me'),
    headers: {
      'Authorization': 'Bearer $token',
      'Content-Type': 'application/json',
    },
  );

  if (response.statusCode == 200) {
    final data = json.decode(response.body);
    return data['data'];
  } else {
    throw Exception('Failed to load user');
  }
}
```

## 🎯 Best Practices

1. **Fetch on App Load**: Call `/me` when your app initializes to get user context
2. **Store in Global State**: Keep user data in Redux/Vuex/Context for easy access
3. **Refetch After Updates**: If user roles/permissions change, refetch `/me`
4. **Handle Errors Gracefully**: 401 errors should redirect to login
5. **Use Permission Helpers**: Create utility functions for permission checks
6. **Cache Appropriately**: Cache for session duration, clear on logout

## 🔄 Workflow

```
1. User logs in → Receives JWT token
2. App stores token (localStorage/cookies)
3. App calls /me endpoint → Gets user + roles + permissions
4. App stores user data in state management
5. UI renders based on permissions
6. On logout → Clear token and user data
```

## 📚 Related Endpoints

- `POST /api/v1/users/register` - Register new user
- `POST /api/v1/users/login` - Login and get token
- `GET /api/v1/users/:email` - Get user by email (public info only)

## 💡 Tips

- Use this endpoint to build dynamic, permission-aware UIs
- Combine with Casbin middleware on backend for complete authorization
- Frontend permissions are for UX only - always validate on backend!
- Consider using React Query or SWR for better caching and refetching
