export interface InvitationTokenRequest {
    userId: string;
    email: string;
}

export interface InvitationTokenResponse {
    token: string;
}
