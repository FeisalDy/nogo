# Error Handling - Before & After Comparison

## The Problem You Had

When sending invalid data to your API, you got this ugly, hard-to-read error:

```json
{
  "error": "Key: 'RegisterUserDTO.Username' Error:Field validation for 'Username' failed on the 'required' tag\nKey: 'RegisterUserDTO.Email' Error:Field validation for 'Email' failed on the 'required' tag\nKey: 'RegisterUserDTO.Password' Error:Field validation for 'Password' failed on the 'required' tag\nKey: 'RegisterUserDTO.ConfirmPassword' Error:Field validation for 'ConfirmPassword' failed on the 'required' tag"
}
```

**Problems:**

- ❌ Hard to read (all in one line with newlines)
- ❌ Uses internal struct names (`RegisterUserDTO.Username`)
- ❌ Not user-friendly
- ❌ No error code for frontend to handle
- ❌ Inconsistent format
- ❌ Hard to parse programmatically

## The Solution We Created

Now the same error looks like this:

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        {
          "field": "username",
          "message": "username is required",
          "tag": "required",
          "value": ""
        },
        {
          "field": "email",
          "message": "email is required",
          "tag": "required",
          "value": ""
        },
        {
          "field": "password",
          "message": "password is required",
          "tag": "required",
          "value": ""
        },
        {
          "field": "confirm password",
          "message": "confirm password is required",
          "tag": "required",
          "value": ""
        }
      ],
      "summary": "username: username is required; email: email is required; password: password is required; confirm password: confirm password is required"
    }
  }
}
```

**Benefits:**

- ✅ Easy to read and understand
- ✅ Human-friendly field names (`username` instead of `RegisterUserDTO.Username`)
- ✅ Clear, actionable error messages
- ✅ Unique error code (`USER007`) for frontend handling
- ✅ Consistent format across all endpoints
- ✅ Easy to parse and display in UI
- ✅ Includes both detailed field errors and a summary
- ✅ Shows the invalid values for debugging

## More Examples

### Before: Invalid Email

```json
{
  "error": "Key: 'RegisterUserDTO.Email' Error:Field validation for 'Email' failed on the 'email' tag"
}
```

### After: Invalid Email

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        {
          "field": "email",
          "message": "email must be a valid email address",
          "tag": "email",
          "value": "not-an-email"
        }
      ],
      "summary": "email: email must be a valid email address"
    }
  }
}
```

### Before: Password Too Short

```json
{
  "error": "Key: 'RegisterUserDTO.Password' Error:Field validation for 'Password' failed on the 'min' tag"
}
```

### After: Password Too Short

```json
{
  "success": false,
  "error": {
    "code": "USER007",
    "message": "Validation failed",
    "details": {
      "fields": [
        {
          "field": "password",
          "message": "password must be at least 8 characters long",
          "tag": "min",
          "value": "123"
        }
      ],
      "summary": "password: password must be at least 8 characters long"
    }
  }
}
```

### Before: User Not Found

```json
{
  "error": "user not found"
}
```

### After: User Not Found

```json
{
  "success": false,
  "error": {
    "code": "USER001",
    "message": "User not found"
  }
}
```

### Before: Internal Server Error

```json
{
  "error": "sql: no rows in result set"
}
```

### After: Internal Server Error (with details)

```json
{
  "success": false,
  "error": {
    "code": "USER003",
    "message": "Failed to create user",
    "details": {
      "reason": "duplicate key value violates unique constraint"
    }
  }
}
```

## Success Response Format

We also standardized success responses:

### Before

```json
{
  "id": "123",
  "username": "john_doe",
  "email": "john@example.com"
}
```

### After

```json
{
  "success": true,
  "data": {
    "id": "123",
    "username": "john_doe",
    "email": "john@example.com"
  },
  "message": "User created successfully"
}
```

## Frontend Integration Example

### Before (Hard to Handle)

```javascript
// ❌ Hard to parse and display
fetch("/api/users/register", {
  method: "POST",
  body: JSON.stringify(userData),
})
  .then((res) => res.json())
  .then((data) => {
    if (data.error) {
      // What do we do with this ugly string?
      alert(data.error);
      // Shows: "Key: 'RegisterUserDTO.Username' Error:Field validation..."
    }
  });
```

### After (Easy to Handle)

```javascript
// ✅ Easy to parse and display nicely
fetch("/api/users/register", {
  method: "POST",
  body: JSON.stringify(userData),
})
  .then((res) => res.json())
  .then((data) => {
    if (!data.success) {
      // Show error code
      console.error(`Error ${data.error.code}: ${data.error.message}`);

      // Display field errors in form
      if (data.error.details?.fields) {
        data.error.details.fields.forEach((fieldError) => {
          showErrorOnField(fieldError.field, fieldError.message);
        });
      }

      // Or show summary
      alert(data.error.details?.summary || data.error.message);
    } else {
      // Success!
      console.log(data.message); // "User created successfully"
      showUser(data.data);
    }
  });
```

### React Example

```jsx
// ✅ Perfect for React forms
const [errors, setErrors] = useState({});

const handleSubmit = async (formData) => {
  const response = await fetch("/api/users/register", {
    method: "POST",
    body: JSON.stringify(formData),
  });

  const data = await response.json();

  if (!data.success) {
    // Convert field errors to error state
    const fieldErrors = {};
    data.error.details?.fields?.forEach((fieldError) => {
      fieldErrors[fieldError.field] = fieldError.message;
    });
    setErrors(fieldErrors);
  } else {
    // Success - redirect or show message
    navigate("/dashboard");
    toast.success(data.message);
  }
};

return (
  <form onSubmit={handleSubmit}>
    <input name="username" />
    {errors.username && <span className="error">{errors.username}</span>}

    <input name="email" type="email" />
    {errors.email && <span className="error">{errors.email}</span>}

    <input name="password" type="password" />
    {errors.password && <span className="error">{errors.password}</span>}
  </form>
);
```

## Error Code Based Handling

```javascript
// Handle different errors differently based on code
if (!data.success) {
  switch (data.error.code) {
    case "USER001": // User not found
      navigate("/404");
      break;
    case "USER002": // User already exists
      showMessage("This email is already registered. Try logging in.");
      break;
    case "USER007": // Validation error
      displayFieldErrors(data.error.details.fields);
      break;
    case "AUTH001": // Invalid token
      logout();
      navigate("/login");
      break;
    default:
      showMessage("An error occurred. Please try again.");
  }
}
```

## Localization Support

Error codes make localization easy:

```javascript
const errorMessages = {
  en: {
    USER001: "User not found",
    USER007: "Please check your input",
    AUTH001: "Your session has expired",
  },
  es: {
    USER001: "Usuario no encontrado",
    USER007: "Por favor revise su entrada",
    AUTH001: "Su sesión ha expirado",
  },
};

// Use the error code to get localized message
const localizedMessage = errorMessages[currentLanguage][data.error.code];
```

## Summary

| Aspect                   | Before                  | After                                 |
| ------------------------ | ----------------------- | ------------------------------------- |
| **Readability**          | Poor - technical jargon | Excellent - human-friendly            |
| **Structure**            | Inconsistent            | Standardized                          |
| **Error Codes**          | None                    | Unique codes (USER007, AUTH001, etc.) |
| **Field Details**        | No                      | Yes - per-field errors                |
| **Parsing**              | Hard                    | Easy - structured JSON                |
| **Frontend Integration** | Difficult               | Simple                                |
| **Debugging**            | Harder                  | Easier - with codes and details       |
| **User Experience**      | Bad                     | Great - clear messages                |
| **Localization**         | Impossible              | Easy - use error codes                |
| **Consistency**          | Random                  | Uniform across all endpoints          |

---

**Your API now provides professional, user-friendly error messages! 🎉**
