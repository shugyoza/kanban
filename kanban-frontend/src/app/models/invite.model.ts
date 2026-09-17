export interface CreateInvitationTokenRequest {
    userId: string;
    email: string;
}

export interface CreateInvitationTokenResponse {
    token: string;
}

export interface ValidateInvitationTokenResponse {
    valid: boolean;
}
