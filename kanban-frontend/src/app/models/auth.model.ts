export interface User {
    id: string;
    username: string;
    createdAt: string;
}

export interface Credentials {
    username: string;
    password?: string;
}