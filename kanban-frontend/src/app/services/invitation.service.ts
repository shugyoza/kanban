import { HttpClient, HttpResponse } from '@angular/common/http';
import { inject, Service } from '@angular/core';
import { Observable } from 'rxjs';
import { CreateInvitationTokenRequest, CreateInvitationTokenResponse } from '../models/invite.model';

@Service()
export class InvitationService {
    private http = inject(HttpClient);

    public getInvitationToken(payload: CreateInvitationTokenRequest): Observable<CreateInvitationTokenResponse> {

        return this.http.post<CreateInvitationTokenResponse>(
            '/api/auth/invite-token/create',
            payload
        )
    }

    public validateInvitationToken(token: string): Observable<HttpResponse<void>> {

        return this.http.post<HttpResponse<void>>(
            '/api/auth/invite-token/validate',
            { token }
        )
    }
}
