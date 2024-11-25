import axios from 'axios';
import Cookies from 'js-cookie';
import toast from 'react-hot-toast';

const apiInstance = axios.create({
    baseURL: 'http://localhost:8000', // Replace with your Golang API URL
});

// Attach the access token to every request
apiInstance.interceptors.request.use(
    (config) => {
        const accessToken = Cookies.get('accessToken');
        if (accessToken) {
            config.headers.Authorization ="Bearer "+ accessToken;
        }
        return config;
    },
    (error) => Promise.reject(error)
);


// Refresh token interceptor
apiInstance.interceptors.response.use(
    (response) => {
        if (response.status === 302 || response.status === 301) {
            const redirectURL = response.data;
            if((redirectURL as string).includes("forgot-password"))
            toast.error("You are required by admin to change your password. Redirecting...");
            setTimeout(() => {
                window.location.href = redirectURL;
            }, 2000);
        }
        return response;
    },
    async (error) => {
        if (error.response && error.response.status === 302) {
            const redirectURL = error.response.data;
            if ((redirectURL as string).includes("forgot-password"))
                toast.error("You are required by admin to change your password. Redirecting...");
            setTimeout(() => {
                window.location.href = redirectURL;
            }, 2000);
            return Promise.reject(error);
        }
        const originalRequest = error.config;

        // If access token expired
        if (error.response && error.response.status === 401 && !originalRequest._retry) {
            originalRequest._retry = true;

            try {
                // Attempt to refresh the token
                const response = await axios.post(
                    'http://localhost:8000/api/v1/auth/refresh-token',
                    {}, // Body is empty because the refresh token is in cookies
                    { withCredentials: true }
                );

                // Save the new access token
                Cookies.set('accessToken', response.data.accessToken, {
                    expires: 1 / 24,
                    // sameSite: "strict",
                    // secure: false
                });

                // Retry the original request
                originalRequest.headers.Authorization = response.data.accessToken;
                return apiInstance(originalRequest);
            } catch (refreshError) {
                // Logout the user if token refreshing fails
                Cookies.remove('accessToken');
                alert("Session has expired")
                window.location.href = '/auth/signin'; // Redirect to login page
            }
        }

        return Promise.reject(error);
    }
);

export interface APIResponse<T> {
    success: boolean
    message: string
    data: {
        data: T
    }
}

export default apiInstance;