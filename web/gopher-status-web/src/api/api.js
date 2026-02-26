// src/api/api.js

// URL теперь относительные, так как работает прокси из vite.config.js
const MONITOR_URL = "/v1/monitors";
const AUTH_URL = "/v1/auth";

// Хелпер для получения заголовков с токеном
const getHeaders = () => {
    const token = localStorage.getItem('accessToken');
    return {
        "Content-Type": "application/json",
        // Если токен есть, добавляем заголовок Authorization
        ...(token ? { "Authorization": `Bearer ${token}` } : {})
    };
};

// ==========================================
//               АВТОРИЗАЦИЯ
// ==========================================

export const LoginUser = async (login, password) => {
    const response = await fetch(`${AUTH_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ login, password }),
    });

    const data = await response.json();
    if (!response.ok) throw new Error(data.message || "Failed to login");

    return data; // Возвращаем токен
};

export const RegisterUser = async (login, password) => {
    const response = await fetch(`${AUTH_URL}/users`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ login, password }),
    });

    const data = await response.json();
    if (!response.ok) throw new Error(data.message || "Failed to register");

    return data;
};

// ==========================================
//               МОНИТОРЫ
// ==========================================

export const ListMonitors = async () => {
    const response = await fetch(MONITOR_URL, {
        method: "GET",
        headers: getHeaders(),
    });

    if (response.status === 401 || response.status === 403) throw new Error("Unauthorized");
    if (!response.ok) throw new Error("Failed to fetch monitors");

    return response.json();
};

export const CreateMonitor = async (data) => {
    const response = await fetch(MONITOR_URL, {
        method: "POST",
        headers: getHeaders(),
        body: JSON.stringify(data),
    });

    if (response.status === 401 || response.status === 403) throw new Error("Unauthorized");
    if (!response.ok) throw new Error("Failed to create monitor");

    return response.json();
};

export const DeleteMonitor = async (id) => {
    const response = await fetch(`${MONITOR_URL}/${id}`, {
        method: "DELETE",
        headers: getHeaders(),
    });

    if (response.status === 401 || response.status === 403) throw new Error("Unauthorized");
    if (!response.ok) throw new Error("Failed to delete monitor");

    return response.json();
};
