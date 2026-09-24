export interface CreateInvitationTokenRequest {
    userId: string;
    email: string;
    registerUrl: string;
}

export interface CreateInvitationTokenResponse {
    token: string;
}
