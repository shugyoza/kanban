export interface User {
    id: string;
    username: string;
    createdAt: string;
}

export interface LoginCredentials {
    username: string;
    password?: string;
}