export interface CreateInvitationTokenRequest {
    userId: string;
    email: string;
}

export interface CreateInvitationTokenResponse {
    token: string;
}
